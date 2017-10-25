package main

import (
	"context"
	"flag"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	machinesv1 "k8s.io/kube-deploy/cluster-api/api/machines/v1alpha1"
)

var kubeconfig = flag.String("kubeconfig", "", "path to kubeconfig file")

func main() {
	flag.Parse()

	client, _, err := restClient()
	if err != nil {
		panic(err.Error())
	}

	run(context.Background(), client)
}

func run(ctx context.Context, client *rest.RESTClient) error {
	source := cache.NewListWatchFromClient(client, "machines", apiv1.NamespaceAll, fields.Everything())

	_, controller := cache.NewInformer(
		source,
		&machinesv1.Machine{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    onAdd,
			UpdateFunc: onUpdate,
			DeleteFunc: onDelete,
		},
	)

	controller.Run(ctx.Done())
	// unreachable
	return nil
}

func onAdd(obj interface{}) {
	machine := obj.(*machinesv1.Machine)
	fmt.Printf("object created: %s\n", machine.ObjectMeta.Name)
}

func onUpdate(oldObj, newObj interface{}) {
	oldMachine := oldObj.(*machinesv1.Machine)
	newMachine := newObj.(*machinesv1.Machine)
	fmt.Printf("object updated: %s\n", oldMachine.ObjectMeta.Name)
	fmt.Printf("  old k8s version: %s, new: %s\n", oldMachine.Spec.Versions.Kubelet, newMachine.Spec.Versions.Kubelet)
}

func onDelete(obj interface{}) {
	machine := obj.(*machinesv1.Machine)
	fmt.Printf("object deleted: %s\n", machine.ObjectMeta.Name)
}
