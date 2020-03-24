package clair

import (
	"log"
	"os/exec"
)

// Scanner
type Scanner struct {
	Scannable
}

// Scannable image interface
type Scannable interface {
	Scan(dockerImage string) chan error
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (*Scanner) Scan(dockerImage string) chan error {
	stdout := make(chan error, 1)
	//dockerImage := image.Image
	go func(dockerImage string, stdout chan error) {
		//quay.io/rawdatalabs/jigger-operator:v0.0.2
		cmd := exec.Command("clair-scanner", "--clair", "http://127.0.0.1:6060", "--ip", "192.168.0.102", "-r", "reports/"+dockerImage+"-scan-output.json", dockerImage)
		//cmd := exec.Command("echo",  dockerImage )
		log.Printf("Running Clair Scan command on image %s and waiting for it to finish...", dockerImage)
		err := cmd.Run()
		log.Printf("Command finished with error: %v", err)
		stdout <- err
	}(dockerImage, stdout)
	return stdout
}
