// Code generated by main. DO NOT EDIT.

package credentials

import (
	"context"
	"time"

	clientset "github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned"
	scheme "github.com/mjpitz/credentials-operator/pkg/generated/clientset/versioned/scheme"
	informers "github.com/mjpitz/credentials-operator/pkg/generated/informers/externalversions"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/schemes"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func init() {
	scheme.AddToScheme(schemes.All)
}

type Factory struct {
	synced            bool
	informerFactory   informers.SharedInformerFactory
	clientset         clientset.Interface
	controllerManager *generic.ControllerManager
	threadiness       map[schema.GroupVersionKind]int
}

func NewFactoryFromConfigOrDie(config *rest.Config) *Factory {
	f, err := NewFactoryFromConfig(config)
	if err != nil {
		panic(err)
	}
	return f
}

func NewFactoryFromConfig(config *rest.Config) (*Factory, error) {
	return NewFactoryFromConfigWithOptions(config, nil)
}

func NewFactoryFromConfigWithNamespace(config *rest.Config, namespace string) (*Factory, error) {
	return NewFactoryFromConfigWithOptions(config, &FactoryOptions{
		Namespace: namespace,
	})
}

type FactoryOptions struct {
	Namespace string
	Resync    time.Duration
}

func NewFactoryFromConfigWithOptions(config *rest.Config, opts *FactoryOptions) (*Factory, error) {
	if opts == nil {
		opts = &FactoryOptions{}
	}

	cs, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	resync := opts.Resync
	if resync == 0 {
		resync = 2 * time.Hour
	}

	if opts.Namespace == "" {
		informerFactory := informers.NewSharedInformerFactory(cs, resync)
		return NewFactory(cs, informerFactory), nil
	}

	informerFactory := informers.NewSharedInformerFactoryWithOptions(cs, resync, informers.WithNamespace(opts.Namespace))
	return NewFactory(cs, informerFactory), nil
}

func NewFactory(clientset clientset.Interface, informerFactory informers.SharedInformerFactory) *Factory {
	return &Factory{
		threadiness:       map[schema.GroupVersionKind]int{},
		controllerManager: &generic.ControllerManager{},
		clientset:         clientset,
		informerFactory:   informerFactory,
	}
}

func (c *Factory) Controllers() map[schema.GroupVersionKind]*generic.Controller {
	return c.controllerManager.Controllers()
}

func (c *Factory) SetThreadiness(gvk schema.GroupVersionKind, threadiness int) {
	c.threadiness[gvk] = threadiness
}

func (c *Factory) Sync(ctx context.Context) error {
	c.informerFactory.Start(ctx.Done())
	c.informerFactory.WaitForCacheSync(ctx.Done())
	return nil
}

func (c *Factory) Start(ctx context.Context, defaultThreadiness int) error {
	if err := c.Sync(ctx); err != nil {
		return err
	}

	return c.controllerManager.Start(ctx, defaultThreadiness, c.threadiness)
}

func (c *Factory) Credentials() Interface {
	return New(c.controllerManager, c.informerFactory.Credentials(), c.clientset)
}