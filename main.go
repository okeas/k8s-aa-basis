package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-aa-basis/pkg/apis/myingress/v1beta1"
	"k8s-aa-basis/pkg/builders"
	"k8s-aa-basis/pkg/k8sconfig"
	"k8s-aa-basis/pkg/store"
	"k8s-aa-basis/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

//把 xx=xx,xx=xxx  解析为一个map
func parseLabelQuery(query string) map[string]string {
	m := make(map[string]string)
	if query == "" {
		return m
	}
	qs := strings.Split(query, ",")
	if len(qs) == 0 {
		return m
	}
	for _, q := range qs {
		qPair := strings.Split(q, "=")
		if len(qPair) == 2 {
			m[qPair[0]] = qPair[1]
		}
	}
	return m
}

var rootJson = `
{
  "kind":"APIResourceList",
  "apiVersion":"v1",
  "groupVersion":"apis.jtthink.com/v1beta1",
  "resources":[
     {"name":"mypods","singularName":"mypod","shortNames":["mp"],"namespaced":true,"kind":"MyPod","verbs":["get","list"]}
  ]}
`
var podsListv2 = `
{
  "kind": "MyPodList",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {},
  "items":[
    {
	  "metadata": {
        "name": "testpod1-v2",
        "namespace": "default"
       }
    },
    {
	  "metadata": {
        "name": "testpod2-v2",
        "namespace": "default"
       }
    }
   ]
}
`
var podsListv1 = `
{
  "kind": "MyPodList",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {},
  "items":[
    {
	  "metadata": {
        "name": "testpod1-v1",
        "namespace": "default"
       }
    },
    {
	  "metadata": {
        "name": "testpod2-v1",
        "namespace": "default"
       }
    }
   ]
}
`
var podDetail = `
{
  "kind": "MyPod",
  "apiVersion": "apis.jtthink.com/v1beta1",
  "metadata": {"name":"{name}","namespace":"{namespace}"},
  "spec":{"属性":"你懂的"},
  "columnDefinitions": [
        {
            "name": "Name",
            "type": "string"
        },
        {
            "name": "Created At",
            "type": "date"
        }
    ]
}
`

var (
	ROOT_URL = fmt.Sprintf("/apis/%s/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version)
	//根据NS 获取 myingress列表 kubectl get mi -n xxx
	ListByNS_URL = fmt.Sprintf("/apis/%s/%s/namespaces/:ns/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)

	ListAll_URL = fmt.Sprintf("/apis/%s/%s/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)
	// 根据NS 获取 myingress 如 kubectl get mi abc
	DetailByNS_URL = fmt.Sprintf("/apis/%s/%s/namespaces/:ns/%s/:name", v1beta1.SchemeGroupVersion.Group,
		v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)
	//新增
	PostByNS_URL = fmt.Sprintf("/apis/%s/%s/namespaces/:ns/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)
	// apply -->存在的情况下
	PatchByNS_URL = fmt.Sprintf("/apis/%s/%s/namespaces/:ns/%s/:name", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)
)

func main() {

	k8sconfig.K8sInitInformer() //启动Informer监听
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if c.Request.Method != "GET" {
			fmt.Println("METHOD", c.Request.Method, "  PATH:", c.Request.URL.Path)
		}
		defer func() {
			if e := recover(); e != nil {
				fmt.Println("myerror")
				c.AbortWithStatusJSON(400, e)
				return
			}
		}()

		c.Next()
	})
	// 根
	r.GET(ROOT_URL, func(c *gin.Context) {
		c.JSON(200, builders.ApiResourceList())
	})
	//列表  （根据ns)
	r.GET(ListByNS_URL, func(c *gin.Context) {
		list, err := store.NewClientStore().
			ListByNsOrAll(c.Param("ns"))
		if err != nil {
			status := utils.NotFoundStatus("Ingress列表不存在")
			c.AbortWithStatusJSON(404, status)
			return
		}
		c.JSON(200, utils.ConvertToTable(list))
	})

	//列表  （所有 )
	r.GET(ListAll_URL, func(c *gin.Context) {
		list, err := store.NewClientStore().
			ListByNsOrAll("") //取全部
		if err != nil {
			status := utils.NotFoundStatus("Ingress列表不存在")
			c.AbortWithStatusJSON(404, status)
			return
		}
		c.JSON(200, utils.ConvertToTable(list))
	})

	//详细 （根据ns)  好比kubectl get mi  web1
	r.GET(DetailByNS_URL, func(c *gin.Context) {
		mi, err := store.NewClientStore().GetByNs(c.Param("name"), c.Param("ns"))
		if err != nil {
			status := utils.NotFoundStatus(fmt.Sprintf("你要的Ingress:%s在%s这个命名空间没找到，去别地儿看看？",
				c.Param("name"), c.Param("ns")))
			c.AbortWithStatusJSON(404, status)
			return
		}
		c.JSON(200, mi)
	})

	//无脑新增
	r.POST(PostByNS_URL, func(c *gin.Context) {
		mi := &v1beta1.MyIngress{}
		err := c.ShouldBindJSON(mi)
		if err != nil {
			c.AbortWithStatusJSON(400, utils.ErrorStatus(400, err.Error(), metav1.StatusReasonBadRequest))
			return
		}
		//创建真实的Ingress
		err = builders.CreateIngress(mi)
		if err != nil {
			c.AbortWithStatusJSON(400,
				utils.ErrorStatus(400, err.Error(), metav1.StatusReasonBadRequest))
			return
		}
		c.JSON(200, mi)

	})

	//  如果已经存在， 执行 patch 请求
	r.PATCH(PatchByNS_URL, func(c *gin.Context) {
		apply := &v1beta1.MyIngress{} //
		err := c.ShouldBindJSON(&apply)
		if err != nil {
			c.AbortWithStatusJSON(400,
				utils.ErrorStatus(400, err.Error(),
					metav1.StatusReasonBadRequest))
			return
		}

		newMi, err := builders.PatchIngress(apply)
		if err != nil {
			c.AbortWithStatusJSON(400,
				utils.ErrorStatus(400, err.Error(), metav1.StatusReasonBadRequest))
			return
		}
		c.JSON(200, newMi)
	})
	// 8443  没有为啥
	if err := r.RunTLS(":8443",
		"certs/aaserver.crt", "certs/aaserver.key"); err != nil {
		log.Fatalln(err)
	}

	//if err:=r.Run(":8080");err!=nil{
	//	log.Fatalln(err)
	//}

}
