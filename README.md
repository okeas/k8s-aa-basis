# k8s-aa-basis

使用聚合`api`访问`apiserver`

例(查询一个`pod`并用`jq`工具解析查看)：`kb get --raw "/api/v1/pods?limit=1" | jq`

查看`api`列表：`kb get --raw "/" | jq`

> `.kube/config`文件内的必须是原版接口，不能是`rancher`之类的封装接口


