package deprecated


import (
	"fmt"
	// "log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"


)

type ImageX struct {
	podName string
	//containers[] struct{name string; image string}
	images[] string
	vulnerabilities[] string
	scanned bool
}

func init() {
	register_pods()
	//register_informer()

}
//
//func register_informer(){
//
//}

func register_pods() {

	c := make(chan ImageX)

	go ListK8sPods(c)

	for rx := range c {
		fmt.Println("Recieving:",rx)
		addToCache( rx)
	}

}



func addToCache(rx ImageX) {
	redisCmd := fmt.Sprintf("Scanning pod:%s image %s", rx.podName, rx.images  )
	fmt.Println(redisCmd)

}

func  ListK8sPods(c chan ImageX) {
	// Instantiate loader for kubeconfig file.
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Determine the Namespace referenced by the current context in the
	// kubeconfig file.
	namespace, _, err := kubeconfig.Namespace()
	if err != nil {
		panic(err)
	}

	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	restconfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	// Create a Kubernetes core/v1 client.
	coreclient, err := corev1client.NewForConfig(restconfig)
	if err != nil {
		panic(err)
	}

	// List all Pods in our current Namespace.
	pods, err := coreclient.Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)

	for _, pod := range pods.Items {
		// fmt.Printf("  %s\n", pod.Name)
		//fmt.Printf("  %+v\n", pod)
		//fmt.Println(pod.Name, pod.Spec.Containers[0].Image)
		var images = []string{}
		for _, ctn := range pod.Spec.Containers{
			//fmt.Println(pod.Name, ctn.Image)
			images = append(images, ctn.Image)
		}

		c <- ImageX{pod.Name ,  images,  []string{}, false }
	}
	close(c)

}