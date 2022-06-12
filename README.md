# k8s-aa-basis

### 基础介绍

使用聚合`api`访问`apiserver`

例(查询一个`pod`并用`jq`工具解析查看)：`kb get --raw "/api/v1/pods?limit=1" | jq`

> `.kube/config`文件内的必须是原版接口，不能是`rancher`之类的封装接口


查看`api`列表：`kb get --raw "/" | jq`

特殊的`/api/v1` 对应的就是 `/core/v1` 系统自带资源，而其他的资源大多数以`/apis/`开头，如 `"/apis/apps/v1/deployments"` 

### 开启功能

`kube-apiserver` 需要开启自定义聚合功能，使用`kubeadmin`安装的默认已开启，可查看`/etc/kubernetes/manifests/kube-apiserver.yaml`

二进制安装的默认没开启，可以在`systemctl`管理文件中`cat /usr/lib/systemd/system/kube-apiserver.service`新增

```
--proxy-client-key-file=/etc/k8s/certs/server-key.pem \
--proxy-client-cert-file=/etc/k8s/certs/server.pem \
--requestheader-client-ca-file=/etc/k8s/certs/ca.pem \
--requestheader-allowed-names=front-proxy-client \
--requestheader-extra-headers-prefix=X-Remote-Extra- \
--requestheader-group-headers=X-Remote-Group \
--requestheader-username-headers=X-Remote-User

# 解释
# --proxy-client-key-file= 指定私钥文件
# --proxy-client-cret-file= 客户端证书文件
# --requestheader-client-ca-file= 客户端证书文件ca证书
# --requestheader-allowed-names= 客户端证书有效名称(CN)
```

生成自定义服务端的证书：

```bash
$ openssl genrsa -out aaserver.key 2048
$ openssl req -new -key aaserver.key -out aaserver.csr -subj "/CN=front-proxy-client"
# 找一个可用的-CA 和 -CAkey文件生成 CA必须是--requestheader-client-ca-file对应的CA
$ openssl x509 -req -days 3650 -in aaserver.csr -CA /etc/k8s/certs/ca.pem -CAkey /etc/k8s/certs/ca-key.pem -CAcreateserial -out aaserver.crt
```

重启`kube-apiserver`，并将`main.go`编译与`crets`文件夹拷贝到`node01`节点

> `crets`文件夹内的是`aaserver.key`和`aaserver.crt`文件，在`main.go`文件内使用

部署服务`kb apply -f yamls/deploy.yaml`

将自定义服务加入到`aa`中：`kb apply -f yamls/api.yaml`

查看自定义`aa`部署是否成功`kb get apiservice | grep jtthink`

查看自定义`aa`服务响应是否正常`kb get --raw "/apis/apis.jtthink.com/v1beta1"`

目前代码停止于`v1.0`


