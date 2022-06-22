# k8s-aa-basis

### 基础介绍 V1

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

#### 自定义字段

查看`v1.1`内`main.go`

#### 根据标签获取

查看`v1.2`内`main.go`

### 进阶-Ingress案例 v2

`v1`中`main.go`文件内，为了快速演示，以`json`字符串的方式保存了数据，实例上应该使用结构体对象的方式

`v2`的资源名称为`myingress`，短名称`mi`

引入代码

<table>
    <tr>
        <td>目录</td>
        <td>介绍</td>
    </tr>
    <tr>
        <td>pkg/apis</td>
        <td>通过code-generator生成的apis对象代码</td>
    </tr>
    <tr>
        <td>pkg/store</td>
        <td>数据存储的位置，一般储存与etcd，但本示例存储与内存中，重启数据将会消失</td>
    </tr>
</table>

### 将内容修改为搜索系统内的ingress

使用`Informer`返回系统内的`ingress`数据，不在用之前的内存模拟数据

部署一个角色 `kb apply -f yamls/rbac.yaml`，因为要在内部获取数据，并且设置环境变量`release=1`，因为代码内有这个环境变量获取内部角色

见 `v2.1`

### 根据myingress新增创建真实的ingress

`kb apply -f yaml/rabc/yaml` 新增了`ingress`操作权限

`kb apply -f yamls/mi.yaml` 创建了`mi`资源后，会自动增加一个`ingress`

见`v2.2`


## v3本地运行apisever

删除旧方式 `kb delete -f yamls/api.yaml && kb delete -f yamls/deploy.yaml `

运行`go run test.go`

请求`https://localhost:6443/apis/apis.jtthink.com/v1beta1/namespaces/default/myingresses/test`是无权限的

使用`kb --kubeconfig ./local_config --insecure-skip-tls-verify=true get myingresses test`

见`v3.0`