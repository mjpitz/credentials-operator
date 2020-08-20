// Code generated by main. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned/typed/credentials.mjpitz.com/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeCredentialsV1alpha1 struct {
	*testing.Fake
}

func (c *FakeCredentialsV1alpha1) Credentials(namespace string) v1alpha1.CredentialInterface {
	return &FakeCredentials{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeCredentialsV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
