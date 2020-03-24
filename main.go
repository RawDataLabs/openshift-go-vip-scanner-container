package main

import (
	"github.com/rawc0der/VIP_Scanner/controllers"
	//clair "github.com/rawc0der/VIP_Scanner/services"
	//v12 "k8s.io/api/core/v1"
	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"log"
)

const (
	OPENSHIFT_NAMESPACE = "h2o"
)

func main() {
	RunScanner()
}
//
//func run_test_clair_main() {
//	stdout := make(chan error)
//	image := "quay.io/rawdatalabs/jigger-operator:v0.0.2"
// 	go Cmd_Run_Clair_Scanner(image, stdout)
//	error := <- stdout
//	log.Println("Clair finished scanning image: "+image, error)
//	defer close(stdout)
//}

func RunScanner() {
	podSink := make(chan controllers.UpdateOp)

	go controllers.RunVipInformerController(podSink, OPENSHIFT_NAMESPACE)

	controllers.RunScannerController(podSink)

	defer func(){
		close(podSink)
	}()
}





