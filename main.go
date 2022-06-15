package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s-aa-basis/pkg/apis/myingress/v1beta1"
	"k8s-aa-basis/pkg/builders"
	"k8s-aa-basis/pkg/store"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
)

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
	ROOTURL      = fmt.Sprintf("/apis/%s/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version)
	ListByNS_URL = fmt.Sprintf("/apis/%s/%s/namespaces/:ns/%s", v1beta1.SchemeGroupVersion.Group, v1beta1.SchemeGroupVersion.Version, v1beta1.ResourceName)
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

func main() {

	r := gin.New()
	r.Use(func(c *gin.Context) {
		fmt.Println(c.Request.URL.Path)
		c.Next()
	})

	// 根
	r.GET(ROOTURL, func(c *gin.Context) {
		c.JSON(200, builders.ApiResourceList())
	})

	//列表  （根据ns)  kb get mp -l app=nginx,version=1 指定标签
	r.GET(ListByNS_URL, func(c *gin.Context) {
		//解析出query 参数(labelQuery)
		//labelQueryMap := parseLabelQuery(c.Query("labelSelector"))
		//json := ""
		//if v, ok := labelQueryMap["version"]; ok {
		//	if v == "1" {
		//		json = strings.Replace(podsListv1, "default", c.Param("ns"), -1)
		//	}
		//}
		//if json == "" {
		//	json = strings.Replace(podsListv2, "default", c.Param("ns"), -1)
		//}
		c.JSON(200, store.ListMemData(c.Param("ns")))
	})

	//列表  （所有 ) kb get mp -A
	r.GET("/apis/apis.jtthink.com/v1beta1/mypods", func(c *gin.Context) {
		json := strings.Replace(podsListv1, "default", "all", -1)
		c.JSON(200, json)
	})

	//详细 （根据ns)  kb get mp testpod1
	r.GET("/apis/apis.jtthink.com/v1beta1/namespaces/:ns/mypods/:name", func(c *gin.Context) {
		// 自定义字段
		t := metav1.Table{}
		t.Kind = "Table"
		t.APIVersion = "meta.k8s.io/v1"
		// 列
		t.ColumnDefinitions = []metav1.TableColumnDefinition{
			{Name: "name", Type: "string"},
			{Name: "命令空间", Type: "string"},
			{Name: "状态", Type: "string"},
		}

		// 内容
		t.Rows = []metav1.TableRow{
			{Cells: []interface{}{c.Param("name"), c.Param("ns"), "ready"}},
		}
		c.JSON(200, t)
	})

	//  8443  没有为啥
	if err := r.RunTLS(":8443",
		"certs/aaserver.crt", "certs/aaserver.key"); err != nil {
		log.Fatalln(err)
	}
}
