package main

import (
	"context"

	types "github.com/mjpitz/credentials-operator/pkg/apis/credentials.mjpitz.com/v1alpha1"
	scheme "github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned/scheme"
	controller "github.com/mjpitz/credentials-operator/pkg/generated/controllers/credentials.mjpitz.com/v1alpha1"

	wrangercorev1 "github.com/rancher/wrangler-api/pkg/generated/controllers/core/v1"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	kubescheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

const controllerAgentName = "credentials-operator"

func Register(
	ctx context.Context,
	events typedcorev1.EventInterface,
	secrets wrangercorev1.SecretController,
	credentials controller.CredentialsController) {

	h := &Handler{
		secrets:          secrets,
		secretsCache:     secrets.Cache(),
		credentials:      credentials,
		credentialsCache: credentials.Cache(),
		recorder:         createEventRecorder(events),
	}

	secrets.OnChange(ctx, "credentials-operator", h.OnSecretChanged)
	credentials.OnChange(ctx, "credentials-operator", h.OnCredentialsChanged)
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
	credentials      controller.CredentialsClient
	credentialsCache controller.CredentialsCache
	recorder         record.EventRecorder
}

func (h *Handler) OnSecretChanged(key string, secret *corev1.Secret) (*corev1.Secret, error) {
	return secret, nil
}

func (h *Handler) OnCredentialsChanged(key string, creds *types.Credentials) (*types.Credentials, error) {
	return creds, nil
}
