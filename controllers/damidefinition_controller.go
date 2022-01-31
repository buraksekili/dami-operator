/*
Copyright 2022.

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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/buraksekili/dami-operator/pkg"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	damigroupv1alpha1 "github.com/buraksekili/dami-operator/api/v1alpha1"
)

// DamiDefinitionReconciler reconciles a DamiDefinition object
type DamiDefinitionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Env    pkg.Env
}

//+kubebuilder:rbac:groups=damigroup.dami.io,resources=damidefinitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=damigroup.dami.io,resources=damidefinitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=damigroup.dami.io,resources=damidefinitions/finalizers,verbs=update

const damiDefinitionFinalizer = "damigroup.dami.io/finalizer"

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DamiDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	fmt.Println("-------------\n")

	l.Info("we started reconciling", "namespace", req.NamespacedName)
	damiDefinition := &damigroupv1alpha1.DamiDefinition{}
	if err := r.Get(ctx, req.NamespacedName, damiDefinition); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	l.Info("got dami definition", "meta", damiDefinition.Name, "status", damiDefinition.Status)

	// If DeletionTimestamp is not zero, then the object is being deleted.
	if !damiDefinition.ObjectMeta.DeletionTimestamp.IsZero() {
		l.Info("object is marked as deleted")
		if controllerutil.ContainsFinalizer(damiDefinition, damiDefinitionFinalizer) {
			// since finalizer exists, delete dami related dependencies here.
			// delete dami related external sources here.

			// remove finalizer here
			l.Info("removing finalizer from the object")
			controllerutil.RemoveFinalizer(damiDefinition, damiDefinitionFinalizer)
			if err := r.Update(ctx, damiDefinition); err != nil {
				return ctrl.Result{}, err
			}
		}

		// At this point, there are no finalizers existing. So, there is nothing to do.
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(damiDefinition, damiDefinitionFinalizer) {
		l.Info("object does not contain a finalizer, adding it.")
		controllerutil.AddFinalizer(damiDefinition, damiDefinitionFinalizer)
		if err := r.Update(ctx, damiDefinition); err != nil {
			l.Error(err, "failed to update damidefinition after adding finalizer")
			return ctrl.Result{}, err
		}
	}

	if damiDefinition.Status.Resp != damiDefinition.Spec.Resp {
		// make put request to dami server's update endpoint to update dami's configuration.
		endpoint := fmt.Sprintf("%s/update", r.Env.DamiURL)
		if err := makePutRequest(endpoint, &damiDefinition.Spec); err != nil {
			l.Error(err, "failed to make request to the server", "url", endpoint)
		}

		// update Status of the DamiDefinition according to desired state's spec.
		damiDefinition.Status.Resp = damiDefinition.Spec.Resp
		if err := r.Status().Update(ctx, damiDefinition); err != nil {
			l.Error(err, "failed to update damidefinition after updating status.")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DamiDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&damigroupv1alpha1.DamiDefinition{}).
		Complete(r)
}

func makePutRequest(url string, body *damigroupv1alpha1.DamiDefinitionSpec) error {
	hc := &http.Client{}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if _, err := hc.Do(req); err != nil {
		return err
	}
	return nil
}
