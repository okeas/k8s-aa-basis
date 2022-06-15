package store

import (
	"k8s-aa-basis/pkg/apis/myingress/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

//内存方式 --- 代替 etcd
var MemData map[string][]*v1beta1.MyIngress

func init() {
	MemData = make(map[string][]*v1beta1.MyIngress)
	//添加一个固定的测试--- 为了演示
	test := &v1beta1.MyIngress{}
	test.Name = "test"
	test.Namespace = "default"
	test.Spec.Path = "testpath"
	createMemData(test)

}

// 创建 数据
func createMemData(ingress *v1beta1.MyIngress) {
	ingress.CreationTimestamp = metav1.NewTime(time.Now())
	if _, ok := MemData[ingress.Namespace]; !ok {
		MemData[ingress.Namespace] = []*v1beta1.MyIngress{}
	}
	MemData[ingress.Namespace] = append(MemData[ingress.Namespace], ingress)
}

//根据ns 查找数据
func findByNameSpace(ns string) []*v1beta1.MyIngress {
	if list, ok := MemData[ns]; !ok {
		MemData[ns] = []*v1beta1.MyIngress{}
		return MemData[ns]
	} else {
		return list
	}
}

// 临时函数。 列出 内存数据
func ListMemData(ns string) *v1beta1.MyIngressList {
	list := v1beta1.NewMyIngressList()
	list.Items = findByNameSpace(ns)
	return list
}
