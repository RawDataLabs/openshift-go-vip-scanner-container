package controllers

import (
	"log"

	"github.com/rawc0der/VIP_Scanner/services"
	store "github.com/rawc0der/VIP_Scanner/services/image-store"
	v12 "k8s.io/api/core/v1"
)

// Create Scanner instance ( imageStore, podSinkChannel, EventHandlers )
func RunScannerController(podSink chan UpdateOp) {
	imgStore := store.NewImageStore()
	registerK8sPodChangeEventHandler(podSink,
		func(pod v12.Pod) {
			log.Printf("main/addChannel <- recieved Pod in addChannel: %s", pod.GetName())
			CheckAndScanAllImages(&pod, imgStore)
		},
		func(pod v12.Pod) {
			log.Printf("main/updateChannel <- recieved Pod in updateChannel: %s", pod.GetName())
			//scanning all images in pod
			CheckAndScanAllImages(&pod, imgStore)
		},
		func(pod v12.Pod) {
			log.Printf("main/deleteChannel <- recieved Pod in deleteChannel: %s", pod.GetName())
		})
}

// utility method to scan all images from given pod, comparing existing entry in store
func CheckAndScanAllImages(pod *v12.Pod, store *store.ImageStore) {
	//scanning all images in pod

	for _, image := range services.GetImagesFromPod(pod) {
		exists, _ := store.Exists(image)
		if exists == false {
			log.Printf("IMAGE SCANNER START SCAN on image %s", image)
			img, err := store.Create(pod.GetName(), image)
			if err != nil {
				panic(err)
			}
			stopChan := store.Scanner.Scan(img.Image)
			stdout := <-stopChan
			log.Printf("IMAGE SCANNER output: %s", stdout)
		} else {
			log.Printf("IMAGE exists bro %s", image)
		}
	}
}

// utility method to register scanner controller event handlers
func registerK8sPodChangeEventHandler(c chan UpdateOp, addEvHandler func(object v12.Pod), updateEvHandler func(object v12.Pod), deleteEvHandler func(object v12.Pod)) {
	for pod := range c {
		log.Printf("main.registerK8sPodChangeEventHandler <- recieved OPTYPE %s for Pod: %s", pod.OpType, pod.Pod.GetName())
		switch pod.OpType {
		case "add":
			addEvHandler(*pod.Pod)
		case "update":
			updateEvHandler(*pod.Pod)
		case "delete":
			deleteEvHandler(*pod.Pod)
		}

	}
}
