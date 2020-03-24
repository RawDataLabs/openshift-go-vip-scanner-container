package controllers

import (
	v12 "k8s.io/api/core/v1"
	"log"

	//corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	//v12 "k8s.io/api/core/v1"
	//"log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)


type  UpdateOp struct {
	Pod *v12.Pod
	OpType string
}

func RunVipInformerController(podSink chan UpdateOp, namespace string) {
	if namespace == "" {
		namespace = "h2o"
	}
	// Instantiate loader for kubeconfig file.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(restconfig)
	if err != nil {
		panic(err.Error())
	}
	factory := informers.NewSharedInformerFactory(clientset, 0)
	h2oNsPodInformer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	h2oNsPodInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(v1.Object)
			pod := obj.(*v12.Pod)
			if mObj.GetNamespace() == namespace {
				log.Printf("New h2o.Pod Added to Store: %s", pod.Name)
				podSink <- *&UpdateOp{
					pod,
					"add",
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(v1.Object)
			pod := obj.(*v12.Pod)
			if mObj.GetNamespace() == namespace {
				log.Printf("New h2o.Pod Deleted From Store: %s", mObj.GetName())
				podSink <- *&UpdateOp{
					pod,
					"delete",
				}
			}
		},
		UpdateFunc: func(old, obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(v1.Object)
			pod := obj.(*v12.Pod)
			if mObj.GetNamespace() == namespace {
				//log.Printf("New h2o.Pod Updates In Store: %s", getPod(mObj))
				podSink <- *&UpdateOp{
					pod,
					"update",
				}
			}
		},
	})
	h2oNsPodInformer.Run(stopper)
}
