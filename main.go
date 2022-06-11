package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

var json = `
{
  "kind":"APIResourceList",
  "apiVersion":"v1",
  "groupVersion":"apis.jtthink.com/v1beta1",
  "resources":[
     {"name":"pods","singularName":"","namespaced":true,"kind":"MyPod","verbs":["get","list"]}
  ]}
`

func main() {

	r := gin.New()
	r.GET("/apis/apis.jtthink.com/v1beta1", func(c *gin.Context) {
		c.Header("content-type", "application/json")
		c.String(200, json)
	})

	//  8443  没有为啥
	if err := r.RunTLS(":8443",
		"certs/aaserver.crt", "certs/aaserver.key"); err != nil {
		log.Fatalln(err)
	}

}
