package main

import (
	"context"
	"os"

	types "github.com/mjpitz/credentials-operator/pkg/apis/credentials.mjpitz.com/v1alpha1"
	scheme "github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned/scheme"
	controller "github.com/mjpitz/credentials-operator/pkg/generated/controllers/credentials.mjpitz.com/v1alpha1"

	wrangercorev1 "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	kubescheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

const controllerAgentName = "credentials-operator"
const kind = "Credential"

func Register(
	ctx context.Context,
	events typedcorev1.EventInterface,
	secrets wrangercorev1.SecretController,
	credentials controller.CredentialController) {

	h := &Handler{
		secrets:          secrets,
		secretsCache:     secrets.Cache(),
		credentials:      credentials,
		credentialsCache: credentials.Cache(),
		recorder:         createEventRecorder(events),
	}

	secrets.OnChange(ctx, controllerAgentName, h.OnSecretChanged)
	credentials.OnChange(ctx, controllerAgentName, h.OnCredentialsChanged)
}

func createEventRecorder(events typedcorev1.EventInterface) record.EventRecorder {
	utilruntime.Must(scheme.AddToScheme(kubescheme.Scheme))
	logrus.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: events})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
}

type Handler struct {
	secrets          wrangercorev1.SecretClient
	secretsCache     wrangercorev1.SecretCache
	credentials      controller.CredentialController
	credentialsCache controller.CredentialCache
	recorder         record.EventRecorder
}

func (h *Handler) OnSecretChanged(key string, secret *corev1.Secret) (*corev1.Secret, error) {
	if secret == nil { // on delete from cache
		return nil, nil
	}

	ownerRef := metav1.GetControllerOf(secret)
	if ownerRef == nil {
		return nil, nil
	}

	if ownerRef.Kind != kind {
		return nil, nil
	}

	credential, err := h.credentialsCache.Get(secret.Namespace, ownerRef.Name)
	if err != nil {
		logrus.Infof("ignoring orphaned object '%s' of credential '%s'", secret.GetSelfLink(), ownerRef.Name)
		return nil, nil
	}

	h.credentials.Enqueue(credential.Namespace, credential.Name)
	return nil, nil
}

func (h *Handler) OnCredentialsChanged(key string, creds *types.Credential) (*types.Credential, error) {
	if creds == nil { // on delete from cache
		return nil, nil
	}

	last, _ := h.secretsCache.Get(creds.Namespace, creds.Name)
	if last != nil && !metav1.IsControlledBy(last, creds) {
		return nil, nil
	}

	secrets := newSecret(creds, last)
	for _, secret := range secrets {
		s, err := h.secretsCache.Get(secret.Namespace, secret.Name)

		if s == nil {
			_, err = h.secrets.Create(secret)
		} else {
			_, err = h.secrets.Update(secret)
		}

		if err != nil {
			// log?
		}
	}

	return creds, nil
}

func newSecret(credential *types.Credential, prior *corev1.Secret) []*corev1.Secret {
	if prior == nil {
		prior = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: credential.Namespace,
				Name:      credential.Name,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(credential, schema.GroupVersionKind{
						Group:   types.SchemeGroupVersion.Group,
						Version: types.SchemeGroupVersion.Version,
						Kind:    kind,
					}),
				},
			},
			StringData: map[string]string{},
		}
	}

	previous := prior.StringData
	updated := map[string]string{}

	for _, credential := range credential.Spec.Credentials {
		if prev, ok := previous[credential.Key]; ok {
			updated[credential.Key] = prev
		} else {
			updated[credential.Key] = generateValue(credential.Requirements)
		}
	}

	secrets := []*corev1.Secret{
		prior,
	}

	for _, view := range credential.Spec.Views {
		data := map[string]string{}
		for key, value := range view.StringDataTemplate {
			data[key] = os.Expand(value, func(s string) string {
				if val, ok := updated[s]; ok {
					return val
				}
				return ""
			})
		}

		secrets = append(secrets, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: credential.Namespace,
				Name:      view.SecretRef.Name,
				OwnerReferences: []metav1.OwnerReference{
					*metav1.NewControllerRef(credential, schema.GroupVersionKind{
						Group:   types.SchemeGroupVersion.Group,
						Version: types.SchemeGroupVersion.Version,
						Kind:    kind,
					}),
				},
			},
			StringData: data,
		})
	}

	return secrets
}

func generateValue(requirements types.Requirements) string {
	return ""
}
