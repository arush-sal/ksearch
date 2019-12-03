package main

import (
	"flag"

	"github.com/infracloudio/ksearch/pkg/config"
	"github.com/infracloudio/ksearch/pkg/printers"
	"github.com/infracloudio/ksearch/pkg/util"
	"k8s.io/client-go/kubernetes"
)

func main() {
	resName := flag.String("name", "", "Name of the pod that you want to get.")
	namespace := flag.String("n", "", "Namespace you want that resource to be searched in.")
	kinds := flag.String("kinds", "", "List all the kinds that you want to be displayed.")

	getter := make(chan interface{})

	flag.Parse()

	cfg := config.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(cfg)

	go util.Getter(*namespace, clientset, *kinds, getter)

	for {
		select {
		case resource, ok := <-getter:
			if !ok {
				return
			}
			printers.Printer(resource, *resName)
		}
	}

}
