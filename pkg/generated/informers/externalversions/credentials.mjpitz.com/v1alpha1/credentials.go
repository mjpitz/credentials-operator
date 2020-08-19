// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	credentialsmjpitzcomv1alpha1 "github.com/mjpitz/credentials-operator/pkg/apis/credentials.mjpitz.com/v1alpha1"
	versioned "github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/mjpitz/credentials-operator/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/mjpitz/credentials-operator/pkg/generated/listers/credentials.mjpitz.com/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CredentialsInformer provides access to a shared informer and lister for
// Credentialses.
type CredentialsInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.CredentialsLister
}

type credentialsInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCredentialsInformer constructs a new informer for Credentials type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCredentialsInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCredentialsInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCredentialsInformer constructs a new informer for Credentials type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCredentialsInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CredentialsV1alpha1().Credentialses(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CredentialsV1alpha1().Credentialses(namespace).Watch(context.TODO(), options)
			},
		},
		&credentialsmjpitzcomv1alpha1.Credentials{},
		resyncPeriod,
		indexers,
	)
}

func (f *credentialsInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCredentialsInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *credentialsInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&credentialsmjpitzcomv1alpha1.Credentials{}, f.defaultInformer)
}

func (f *credentialsInformer) Lister() v1alpha1.CredentialsLister {
	return v1alpha1.NewCredentialsLister(f.Informer().GetIndexer())
}