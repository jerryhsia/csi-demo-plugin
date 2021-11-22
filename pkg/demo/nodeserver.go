package demo

import (
	"bytes"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type nodeServer struct {
	nodeID string
}

func NewNodeServer(nodeid string) *nodeServer {
	return &nodeServer{
		nodeID: nodeid,
	}
}

func ensureDir(dir string) error {
	_, err2 := os.Stat(dir)
	if err2 != nil {
		// klog.V(4).Infof("ensureDir dir:%s not found", dir)
		if err := os.MkdirAll(dir, 0750); err != nil {
			klog.V(4).Infof("ensureDir dir:%s failed: %v", dir, err)
		}
	}

	return nil
}

// 此函数在每个节点调用一次
// 目的：将存储挂载到global全局目录StagingTargetPath
func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.V(4).Infof("NodeStageVolume: called with args %+v", *req)

	ensureDir(req.GetStagingTargetPath())

	for _, volume := range allVolume {
		volumeStagePath := filepath.Join(req.GetStagingTargetPath(), volume.Name)
		cmdStr, optionsStr := volume.GetMountCmd(volumeStagePath)
		options := strings.Split(optionsStr, " ")

		ensureDir(volumeStagePath)

		cmd := exec.Command(cmdStr, options...)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		klog.V(4).Infof("NodeStageVolume cmd:%s stdout:%s stderr:%s err:%#v", cmdStr, stdout.String(), stderr.String(), err)
	}

	return &csi.NodeStageVolumeResponse{}, nil
}

// 需保证GetStagingTargetPath
// req.GetStagingTargetPath() 目录为空，且自身不是挂载点
func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.V(4).Infof("NodeUnstageVolume: called with args %+v", *req)

	for _, volume := range allVolume {
		volumeStagePath := filepath.Join(req.GetStagingTargetPath(), volume.Name)
		cmd := exec.Command("umount", "-lf", volumeStagePath)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		klog.V(4).Infof("NodeUnstageVolume cmd:%s stdout:%s stderr:%s err:%#v", cmd, stdout.String(), stderr.String(), err)

		err = os.Remove(volumeStagePath)
		klog.V(4).Infof("NodeUnstageVolume remove dir:%s err:%#v", volumeStagePath, err)
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

// 此函数每个Pod调用一次
// 目的：将Pod的TargetPath关联至global全局目录
func (ns *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.V(4).Infof("NodePublishVolume: called with args %#v", *req)

	ensureDir(req.GetTargetPath())

	for _, volume := range allVolume {
		volumeTargetPath    := filepath.Join(req.GetTargetPath(), volume.Name)
		volumeStagingPath := filepath.Join(req.GetStagingTargetPath(), volume.Name)

		ensureDir(volumeTargetPath)

		str := fmt.Sprintf("--bind %s %s", volumeStagingPath, volumeTargetPath)
		strArr := strings.Split(str, " ")
		cmd := exec.Command("mount", strArr...)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		klog.V(4).Infof("NodeUnstageVolume type:%s stdout:%s stderr:%s err:%#v", volume.Type, stdout.String(), stderr.String(), err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

// 需保证
// req.GetTargetPath() 目录为空，且自身不是挂载点
func (ns *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.V(4).Infof("NodeUnpublishVolume: called with args %+v", *req)

	for _, volume := range allVolume {
		volumeTargetPath    := filepath.Join(req.GetTargetPath(), volume.Name)

		cmd := exec.Command("umount", "-lf", volumeTargetPath)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		klog.V(4).Infof("NodeUnpublishVolume type:%s stdout:%s stderr:%s err:%#v", volume.Type, stdout.String(), stderr.String(), err)

		// 需要删除子目录
		err = os.Remove(volumeTargetPath)
		klog.V(4).Infof("NodeUnpublishVolume remove dir:%s err: %#v", volumeTargetPath, err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetInfo 返回节点信息
func (ns *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.V(4).Infof("NodeGetInfo: called with args %+v", *req)

	return &csi.NodeGetInfoResponse{
		NodeId: ns.nodeID,
	}, nil
}

// NodeGetCapabilities 返回节点支持的功能
func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.V(4).Infof("NodeGetCapabilities: called with args %+v", *req)

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}

func (ns *nodeServer) NodeGetVolumeStats(ctx context.Context, in *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	klog.V(4).Infof("NodeGetVolumeStats: called with args %#v", in)
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *nodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	klog.V(4).Infof("NodeExpandVolume: called with args %#v", req)
	return nil, status.Error(codes.Unimplemented, "")
}
