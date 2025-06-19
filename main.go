package main

import (
	"dev-ops-in-golang/k8s"
	"dev-ops-in-golang/router"
	"log"
)

func main() {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("failed to create k8s client: %v", err)
	}
	r := router.SetupRouter(k8sClient)
	r.Run(":8080")
}

