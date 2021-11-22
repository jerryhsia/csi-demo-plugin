package demo

import "fmt"

type VolumeInterface interface {
    GetName() string
    GetMountCmd(path string) string
}

const VOLUME_TYPE_OBS = "OBS"
const VOLUME_TYPE_BOS = "BOS"
const VOLUME_TYPE_OSS = "OSS"

type Volume struct {
    Name string
    Type string
}

func (v *Volume) GetMountCmd(path string) (string, string) {
    if v.Type == VOLUME_TYPE_OBS {
        return "obsfs", fmt.Sprintf("%s %s -o url=obs.cn-southwest-2.myhuaweicloud.com -o passwd_file=/csi-demo-plugin/etc/passwd-obsfs -o big_writes -o max_write=131072 -o nonempty -o use_ino -o obsfslog", v.Name, path)
    } else if v.Type == VOLUME_TYPE_BOS {
        return "bosfs", fmt.Sprintf("%s %s -o endpoint=http://bj.bcebos.com -o ak=xxx -o sk=xxx -o logfile=/csi-demo-plugin/bos.log", v.Name, path)
    } else {
        return "ossfs", fmt.Sprintf("%s %s -o passwd_file=/csi-demo-plugin/etc/passwd-ossfs -o nonempty -o url=oss-cn-chengdu.aliyuncs.com", v.Name, path)
    }
}

var allVolume []*Volume

func init()  {
    allVolume = make([]*Volume, 0)

    oss := &Volume{}
    oss.Type = VOLUME_TYPE_OSS
    oss.Name = "ossfs-jerry-test"

    bos := &Volume{}
    bos.Type = VOLUME_TYPE_BOS
    bos.Name = "bosfs-jerry-test"

    obs := &Volume{}
    obs.Type = VOLUME_TYPE_OBS
    obs.Name = "obsfs-jerry-test"

    allVolume = append(allVolume, oss)
    allVolume = append(allVolume, bos)
    allVolume = append(allVolume, obs)
}
