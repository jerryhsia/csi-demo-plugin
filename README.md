# csi-demo-plugin

BOS、OSS、OBS csi插件demo，基于 [https://github.com/xzy256/csi-demo-plugin](https://github.com/xzy256/csi-demo-plugin) 修改

## 准备认证文件

```bash
# 放置obs认证文件
etc/passwd-obsfs

# 放置oss认证文件
etc/passwd-ossfs

# 修改文件权限
chmod 0600 etc/*
```

## 准备存储卷列表文件

文件名：`etc/volumes.json`，内容示例：

```json
[
  {
    "name": "bosfs-jerry-test",
    "type": "BOS",
    "cmd": "bosfs",
    "options": "-o endpoint=http://bj.bcebos.com -o ak=xxx -o sk=xxx -o logfile=/csi-demo-plugin/bos.log"
  },
  {
    "name": "ossfs-jerry-test",
    "type": "OSS",
    "cmd": "ossfs",
    "options": "-o passwd_file=/csi-demo-plugin/etc/passwd-ossfs -o nonempty -o url=oss-cn-chengdu.aliyuncs.com"
  },
  {
    "name": "obsfs-jerry-test",
    "type": "OBS",
    "cmd": "obsfs",
    "options": "-o url=obs.cn-southwest-2.myhuaweicloud.com -o passwd_file=/csi-demo-plugin/etc/passwd-obsfs -o big_writes -o max_write=131072 -o nonempty -o use_ino -o obsfslog"
  }
]
```

## 一键部署

```bash
deploy/deploy.sh
```
