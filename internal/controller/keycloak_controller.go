/*
Copyright 2023.

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

package controller

import (
	"context"
	"time"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/scope"
	"appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/cloud/services/keycloak/realms"
)

// KeycloakReconciler reconciles a Keycloak object
type KeycloakReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	KeycloakUrl      string
	ReconcileTimeout time.Duration
}

// +kubebuilder:rbac:groups=appdat.appdat.io,resources=keycloaks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=appdat.appdat.io,resources=keycloaks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=appdat.appdat.io,resources=keycloaks/finalizers,verbs=update
func (r *KeycloakReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx, cancel := context.WithTimeout(ctx, r.ReconcileTimeout)
	defer cancel()

	log := log.FromContext(ctx)
	appdatKeycloak := &appdatv1alpha1.Keycloak{}
	err := r.Get(ctx, req.NamespacedName, appdatKeycloak)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("AppdatKeycloak resource not found or already deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to fetch AppdatKeycloak resource")
		return ctrl.Result{}, err
	}

	keycloak, err := GetOwnerKeycloak(ctx, r.Client, appdatKeycloak.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to get owner keycloak")
		return ctrl.Result{}, err
	}
	if keycloak == nil {
		log.Info("Keycloak Controller has not yet set OwnerRef")
		return ctrl.Result{}, err
	}

	if IsPaused(keycloak, appdatKeycloak) {
		log.Info("AppdatKeycloak of linked Keycloak is marked as paused. Won't reconcile.")
		return ctrl.Result{}, nil
	}

	keycloakScope, err := scope.NewKeycloakScope(ctx, scope.KeycloakScopeParams{
		Client:   r.Client,
		Keycloak: appdatKeycloak,
	})
	if err != nil {
		return ctrl.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	defer func() {
		if err := keycloakScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !appdatKeycloak.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, keycloakScope)
	}

	return r.reconcile(ctx, keycloakScope)
}

func (r *KeycloakReconciler) reconcile(ctx context.Context, keycloakScope *scope.KeycloakScope) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling AppdatKeycloak")

	controllerutil.AddFinalizer(keycloakScope.Keycloak, appdatv1alpha1.KeycloakFinalizer)
	if err := keycloakScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	reconcilers := []cloud.Reconciler{
		realms.New(keycloakScope),
	}

	return ctrl.Result{}, nil
}

func (r *KeycloakReconciler) reconcileDelete(ctx context.Context, keycloakScope *scope.KeycloakScope) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Delete AppdatKeycloak")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeycloakReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appdatv1alpha1.Keycloak{}).
		Complete(r)
}
