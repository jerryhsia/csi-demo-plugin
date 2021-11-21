FROM jerry9916/centos-box-fuse:latest

RUN yum install -y util-linux e2fsprogs && yum clean all && mkdir -p /csi-demo-plugin/bin && mkdir -p /csi-demo-plugin/etc
COPY bin/csi-demo-driver /csi-demo-plugin/bin
COPY etc/* /csi-demo-plugin/etc
ENTRYPOINT ["/csi-demo-plugin/bin/csi-demo-driver"]