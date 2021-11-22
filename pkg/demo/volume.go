package demo

import (
    "encoding/json"
    "fmt"
    "io"
    "os"
)

type VolumeInterface interface {
    GetName() string
    GetMountCmd(path string) string
}

const VOLUME_TYPE_OBS = "OBS"
const VOLUME_TYPE_BOS = "BOS"
const VOLUME_TYPE_OSS = "OSS"

type Volume struct {
    Name string `json:"name"`
    Cmd string `json:"cmd"`
    Type string `json:"type"`
    Options string `json:"options"`
}

func (v *Volume) GetMountCmd(path string) (string, string) {
    return v.Cmd, fmt.Sprintf("%s %s %s", v.Name, path, v.Options)
}

var allVolume []*Volume

func init()  {
    allVolume = make([]*Volume, 0)

    file , err1 := os.Open("etc/volumes.json")
    content, err2 := io.ReadAll(file)
    err3 := json.Unmarshal(content, &allVolume)

    fmt.Printf("InitVolumes volumes:%s err1:%#v err2:%#v err3:%#v", allVolume, err1, err2, err3)
}
