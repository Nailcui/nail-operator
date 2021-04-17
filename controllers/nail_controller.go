/*
Copyright 2021.

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

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	testv1alpha1 "github.com/Nailcui/nail-operator/api/v1alpha1"
)

// NailReconciler reconciles a Nail object
type NailReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.nailcui.github.io,resources=nails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.nailcui.github.io,resources=nails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.nailcui.github.io,resources=nails/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Nail object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nail", req.NamespacedName)
	fmt.Printf("namespace: %s name: %s\n", req.Namespace, req.Name)

	nail := &testv1alpha1.Nail{}
	err := r.Get(ctx, req.NamespacedName, nail)
	if err != nil {
		fmt.Printf("get error: %s\n", err)
	} else {
		fmt.Printf("new replicas: %d\n", nail.Spec.Replicas)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.Nail{}).
		Complete(r)
}
