// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/mjpitz/credentials-operator/pkg/apis/credentials.mjpitz.com/v1alpha1"
	"github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned/scheme"

	rest "k8s.io/client-go/rest"
)

type CredentialsV1alpha1Interface interface {
	RESTClient() rest.Interface
	CredentialsesGetter
}

// CredentialsV1alpha1Client is used to interact with features provided by the credentials.mjpitz.com group.
type CredentialsV1alpha1Client struct {
	restClient rest.Interface
}

func (c *CredentialsV1alpha1Client) Credentialses(namespace string) CredentialsInterface {
	return newCredentialses(c, namespace)
}

// NewForConfig creates a new CredentialsV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*CredentialsV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CredentialsV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new CredentialsV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CredentialsV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new CredentialsV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *CredentialsV1alpha1Client {
	return &CredentialsV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *CredentialsV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
