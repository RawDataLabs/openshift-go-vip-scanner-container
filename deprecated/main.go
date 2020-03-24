package deprecated

import (
	"fmt"
	// "log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gomodule/redigo/redis"

	_ "github.com/rawc0der/VIP_Scanner/controllers"
)

type Image struct {
	podName string
	image string
	scanned bool
	vulnerabilities []string
}

func xx() {
	//register_main_pods()
	//register_informer()
	//controllers.RunMultiPodListerController()
}

func register_informer(){

}

func register_main_pods() {
	redisConn := *RedisStart()

	c := make(chan Image)

	go ListKubePods(c)

	for rx := range c {
		fmt.Println("Recieving:")
		fmt.Println(rx)
		addToRedisHash(redisConn, rx)
		//publishNewImageScanEvent(rx)
	}

}

func getPodsFromRedis (redisConn redis.Conn) {
	//  Get pods from redis
	fmt.Println("Searching for Pods in REDIS:")
	values, err := redis.Values(redisConn.Do("KEYS", "pod:*"))
	if err != nil {
		fmt.Println(err)
		return
	}

	for len(values) > 0 {
		var pod string
		values, err = redis.Scan(values, &pod)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("found pod: "+pod)
		getRedisHashStr(redisConn, pod)

	}
}

func getRedisHashStr(redisConn redis.Conn, podName string) {
	image, err := redis.String(redisConn.Do("HGET", podName, "image"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", image)
}

func getRedisHash(redisConn redis.Conn, podName string) {
	values, err := redis.Values(redisConn.Do("HGET", podName, "image"))
	if err != nil {
		fmt.Println(err)
		return
	}

	for len(values) > 0 {
		var image string
		values, err = redis.Scan(values, &image)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("found Pod.imgage: "+image)

	}
}

func addToRedisHash(redisConn redis.Conn, rx Image) {
	redisCmd := fmt.Sprintf("HSET pod:%s image %s", rx.podName, rx.image  )
	fmt.Println(redisCmd)

	//@todo
	//redisConn.Do("HMSET", redis.Args{"pod:"+rx.podName}.AddFlat(rx)...)
	//_, err := redisConn.Do("HMSET", redis.Args{}.Add("pod:"+rx.podName).AddFlat(&rx)...)

	redisConn.Do("ZADD", "pods", 1, rx.podName)
	_, err := redisConn.Do("HMSET", "pod:"+rx.podName, "image", rx.image, "scanned", rx.scanned, "vulnerabilities", rx.vulnerabilities)
	if err != nil {
		return
	}


}

func RedisStart() (conn *redis.Conn) {
	c, err := redis.Dial("tcp", "a7e2fc632297911ea973806f8086ebf8-23313424.ca-central-1.elb.amazonaws.com:6379")
	if err != nil {
		panic(err)
	}
	_, err = c.Do("AUTH", "P4mayz4dAK")
	if err != nil {
		panic(err)
	}
	_, err2 := c.Do("SELECT", "2")
	if err2 != nil {
		panic(err)
	}

	name, err := redis.String(c.Do("GET", "name"))
	if err != nil {
		fmt.Println("name not found")
	} else {
		//Print our key if it exists
		fmt.Println("name exists: " + name)
	}
	//defer c.Close()
	return &c
}

func  ListKubePods(c chan Image ) {
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
		c <- Image{pod.Name ,  pod.Spec.Containers[0].Image, false, []string{} }
	}
	close(c)

}
