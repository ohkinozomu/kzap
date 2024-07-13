package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	var namespace *string

	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", home+"/.kube/config", "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	namespace = flag.String("namespace", "", "name of the namespace to be deleted")

	flag.Parse()

	if *namespace == "" {
		fmt.Println("namespace is required")
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), *namespace, v1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	ns.ObjectMeta.Finalizers = nil

	nsBytes, err := json.Marshal(ns)
	if err != nil {
		panic(err.Error())
	}

	req := clientset.CoreV1().RESTClient().Put().Resource("namespaces").Name(*namespace).SubResource("finalize").Body(nsBytes)
	result := &corev1.Namespace{}
	err = req.Do(context.TODO()).Into(result)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Namespace %s has been deleted\n", *namespace)
}
