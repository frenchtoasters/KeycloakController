package controller

import (
	"context"

	appdatv1alpha1 "appdat.jsc.nasa.gov/platform/controllers/mri-keycloak/api/v1alpha1"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/clientcredentials"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/cluster-api/controllers/external"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ConfigOwner provides a data interface for different config owner types.
type ConfigOwner struct {
	*unstructured.Unstructured
}

// Pointer returns the pointer of any type
func Pointer[T any](t T) *T {
	return &t
}

func GetKeycloakAccessToken(clientID, clientSecret, rootURL, realmName string) (string, error) {
	tokenURL := rootURL + "/realms/" + realmName + "/protocol/openid-connect/token"

	// Set up the client credentials config
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	// Obtain an access token using the client credentials config
	token, err := config.Token(context.Background())
	if err != nil {
		return "", err
	}

	// Return the access token
	return token.AccessToken, nil
}

// hasAnnotation returns true if the object has the specified annotation.
func hasAnnotation(o metav1.Object, annotation string) bool {
	annotations := o.GetAnnotations()
	if annotations == nil {
		return false
	}
	_, ok := annotations[annotation]
	return ok
}

// HasPaused returns true if the object has the `paused` annotation.
func HasPaused(o metav1.Object) bool {
	return hasAnnotation(o, appdatv1alpha1.PausedAnnotation)
}

// IsPaused returns true if the Keycloak is paused or the object has the `paused` annotation.
func IsPaused(keycloak *appdatv1alpha1.Keycloak, o metav1.Object) bool {
	if keycloak.Spec.Paused {
		return true
	}
	return HasPaused(o)
}

// GetOwnerByRef finds and returns the owner by looking at the object reference.
func GetOwnerByRef(ctx context.Context, c client.Client, ref *corev1.ObjectReference) (*ConfigOwner, error) {
	obj, err := external.Get(ctx, c, ref, ref.Namespace)
	if err != nil {
		return nil, err
	}
	return &ConfigOwner{obj}, nil
}

// GetOwnerKeycloak returns the Keycloak object owning the current resource.
func GetOwnerKeycloak(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*appdatv1alpha1.Keycloak, error) {
	for _, ref := range obj.GetOwnerReferences() {
		if ref.Kind != "Keycloak" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == appdatv1alpha1.GroupVersion.Group {
			return GetKeycloakByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

// GetKeycloakByName finds and return a Keycloak object using the specified params.
func GetKeycloakByName(ctx context.Context, c client.Client, namespace, name string) (*appdatv1alpha1.Keycloak, error) {
	keycloak := &appdatv1alpha1.Keycloak{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}

	if err := c.Get(ctx, key, keycloak); err != nil {
		return nil, errors.Wrapf(err, "failed to get Keycloak/%s", name)
	}

	return keycloak, nil
}
