package main

import (
	"log"
	"time"

	"github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/apis/samplecontroller/v1alpha1"
	"k8s.io/client-go/tools/cache"

	clientset "github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/generated/clientset/versioned"
	informers "github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/generated/informers/externalversions"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "config/config")
	if err != nil {
		panic(err)
	}
	exampleClient, err := clientset.NewForConfig(config)
	sharedInformers := informers.NewSharedInformerFactory(exampleClient, time.Second*2)
	informer := sharedInformers.Samplecontroller().V1alpha1().Foos().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*v1alpha1.Foo)
			log.Printf("New Foo Added: %s", mObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(*v1alpha1.Foo)
			log.Printf("Delete Foo : %s", mObj.GetName())
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	informer.Run(stopCh)
}
