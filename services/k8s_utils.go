package services

import v12 "k8s.io/api/core/v1"

func GetImagesFromPod(pod *v12.Pod) (images []string) {
	//var images = []string{}
	for _, ctn := range pod.Spec.Containers{
		//fmt.Println(pod.Name, ctn.Image)
		images = append(images, ctn.Image)
	}
	return images
}

