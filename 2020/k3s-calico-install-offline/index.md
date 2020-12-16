# 离线安装K3S with Calico

参考链接： https://www.cnblogs.com/weschen/p/12666486.html 

k3s版本：v1.16.13-k3s1

## 下载k3s离线安装文件

1. k3s (k3s binary)
2. k3s-airgap-images (k3s离线组件包)
3. install.sh ( https://raw.githubusercontent.com/rancher/k3s/master/install.sh )

## 离线安装

复制离线包到对应位置：

```shell
sudo mkdir -p /var/lib/rancher/k3s/agent/images/

#下载的离线镜像包复制，格式如下
sudo cp ./k3s-airgap-images-amd64.tar /var/lib/rancher/k3s/agent/images/
```

复制k3s bin包:

```json
#授权
chmod 755 k3s

#下载的K3S的bin包，格式如下
sudo cp ./k3s /usr/local/bin && sudo chmod 755 /usr/local/bin/k3s
```

安装其他依赖（视OS情况）：

```shell
yum install -y container-selinux selinux-policy-base
rpm -i https://rpm.rancher.io/k3s-selinux-0.1.1-rc1.el7.noarch.rpm
```

安装k3s到master节点:

```shell
chmod +x ./install.sh
# 取消flannel，使用docker，取消traefik
cat install.sh | INSTALL_K3S_SKIP_DOWNLOAD=true INSTALL_K3S_EXEC="--flannel-backend=none --docker --no-deploy=traefik" sh -
```

安装k3s到agent节点:

```shell
scp root@k3s-master:/var/lib/rancher/k3s/server/node-token ./
cat node-token # 获取token内容
# calico方式，禁用flannel，启用docker
cat install.sh | INSTALL_K3S_SKIP_DOWNLOAD=true K3S_URL=https://k3s-master:6443 K3S_TOKEN=TOKEN INSTALL_K3S_EXEC="--docker --flannel-backend=none" sh -
```

一个agent的例子：

```shell
cat ./install.sh | INSTALL_K3S_SKIP_DOWNLOAD=true K3S_URL=https://10.10.10.21:6443 K3S_TOKEN=K10637af8a502fd220ce8f58c484ae040170c69b65a379bd7707e358d7bd6457086::server:a0651f8c12f612f109c9a3f513bddb7f INSTALL_K3S_EXEC="--docker --no-flannel" sh -
```

## 删除K3S

```shell
#服务器
/usr/local/bin/k3s-uninstall.sh
#工作节点
/usr/local/bin/k3s-agent-uninstall.sh
```

## 离线安装calico

实现需要升级内核到4.16及以上，下为CentOS7的升级方式。

```shell
# https://www.howtoforge.com/tutorial/how-to-upgrade-kernel-in-centos-7-server/
yum -y update
yum -y install yum-plugin-fastestmirror
cat /etc/redhat-release
cat /etc/os-release
rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
rpm -Uvh https://www.elrepo.org/elrepo-release-7.0-3.el7.elrepo.noarch.rpm
yum repolist
yum --enablerepo=elrepo-kernel install kernel-ml
yum repolist all
sudo awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
sudo grub2-set-default 0
sudo grub2-mkconfig -o /boot/grub2/grub.cfg
sudo reboot
```

参考链接： https://docs.projectcalico.org/getting-started/kubernetes/k3s/quickstart 

需要节点指定容器后端为docker，事先拉取如下包：

* k8s.gcr.io/pause:3.1
* calico/node:v3.15.1
* calico/pod2daemon-flexvol:v3.15.1
* calico/cni:v3.15.1
* calico/kube-controllers:v3.15.1

除此之外，可能还需要k3s本身拉取以下包：

* calico/node
* calico/pod2daemon-flexvol
* calico/cni
* calico/kube-controllers
* rancher/metrics-server
* rancher/local-path-provisioner
* coredns/coredns
* traefik
* rancher/klipper-lb
* rancher/klipper-helm
* mirrorgooglecontainers/pause
* rancher/pause
* k8s.gcr.io/pause

其中`k8s.gcr.io/pause`可从`mirrorgooglecontainers/pause`下载，重新打tag：

> docker tag mirrorgooglecontainers/pause:3.1 k8s.gcr.io/pause:3.1

通过no-flannel模式启动server节点，agent节点正确加入集群后，安装calico：

```shell
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml
```

如若需要自行调整calico启动参数，将上面链接的yaml文件下载下来，并根据Calico的官方文档进行配置参数的修改。