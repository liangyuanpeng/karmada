/*
Copyright 2021 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package statefulset

import (
	"context"

	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// ControllerName is the controller name that will be used when reporting events.
	ControllerName = "statefulset-controller"
)

// Controller is to sync Work.
type Controller struct {
	client.Client // used to operate Work resources.
	EventRecorder record.EventRecorder
}

// Reconcile performs a full reconciliation for the object referred to by the Request.
// The Controller will requeue the Request to be processed again if an error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (c *Controller) Reconcile(ctx context.Context, req controllerruntime.Request) (controllerruntime.Result, error) {
	klog.V(4).Infof("Namespaces sync controller reconciling %s", req.NamespacedName.String())

	return controllerruntime.Result{}, nil
}

// SetupWithManager creates a controller and register to controller manager.
func (c *Controller) SetupWithManager(mgr controllerruntime.Manager) error {
	clusterNamespaceFn := handler.MapFunc(
		func(ctx context.Context, a client.Object) []reconcile.Request {
			var requests []reconcile.Request
			namespaceList := &corev1.NamespaceList{}
			if err := c.Client.List(context.TODO(), namespaceList); err != nil {
				klog.Errorf("Failed to list namespace, error: %v", err)
				return nil
			}

			for _, namespace := range namespaceList.Items {
				requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{
					Name: namespace.Name,
				}})
			}
			return requests
		})

	clusterPredicate := builder.WithPredicates(predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return false
		},
		DeleteFunc: func(event.DeleteEvent) bool {
			return false
		},
		GenericFunc: func(event.GenericEvent) bool {
			return false
		},
	})

	return controllerruntime.NewControllerManagedBy(mgr).
		For(&workv1alpha2.ResourceKindResourceBinding{}).
		Watches(&workv1alpha2.ResourceKindResourceBinding{},
			handler.EnqueueRequestsFromMapFunc(clusterNamespaceFn),
			clusterPredicate).
		Complete(c)
}
