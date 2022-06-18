package k8sconfig

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

//全局变量

const NSFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

//POD里  体内
func K8sRestConfigInPod() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// 获取 config对象
func K8sRestConfig() *rest.Config {
	if os.Getenv("release") == "1" { //自定义环境
		log.Println("run in cluster")
		return K8sRestConfigInPod()
	}
	log.Println("run outside cluster")
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/zx/.kube/config")
	if err != nil {
		log.Fatal(err)
	}
	//config.Insecure=true
	return config
}

//初始化client-go客户端
func InitClient() *kubernetes.Clientset {
	c, err := kubernetes.NewForConfig(K8sRestConfig())
	c.RESTClient().GetRateLimiter()
	if err != nil {
		log.Fatal(err)
	}
	return c
}

var Factory informers.SharedInformerFactory

var K8sClient *kubernetes.Clientset

func K8sInitInformer() {
	K8sClient = InitClient()
	Factory = informers.NewSharedInformerFactory(K8sClient, 0)
	Factory = informers.NewSharedInformerFactory(InitClient(), 0)
	IngressInformer := Factory.Networking().V1().Ingresses() //监听Ingress
	// 暂时不写自己的 回调
	IngressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{})
	stopCh := make(chan struct{})
	Factory.Start(stopCh)
	Factory.WaitForCacheSync(stopCh)
}
