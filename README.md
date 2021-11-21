# csi-demo-plugin

BOS、OSS、OBS csi插件demo，基于 [https://github.com/xzy256/csi-demo-plugin](https://github.com/xzy256/csi-demo-plugin) 修改

## 第一步

```bash
# 放置obs认证文件
etc/passwd-obsfs

# 放置oss认证文件
etc/passwd-ossfs

# 修改文件权限
chmod 0600 etc/*
```

## 一键部署

```bash
deploy/deploy.sh
```

# csi-demo-plugin
