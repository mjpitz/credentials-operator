//go:generate go run ./cmd/codegen-cleanup/main.go
//go:generate go run ./cmd/codegen/main.go

package main

import (
	"context"
	"flag"

	"github.com/mjpitz/credentials-operator/pkg/generated/controllers/credentials.mjpitz.com"

	"github.com/rancher/wrangler-api/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/rancher/wrangler/pkg/start"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
)

var (
	masterURL      string
	kubeconfigFile string
)

func init() {
	flag.StringVar(&kubeconfigFile, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}

func main() {
	ctx := signals.SetupSignalHandler(context.Background())

	cfg, err := kubeconfig.GetNonInteractiveClientConfig(kubeconfigFile).ClientConfig()
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// Raw k8s client, used to events
	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	core := core.NewFactoryFromConfigOrDie(cfg)
	creds := credentials.NewFactoryFromConfigOrDie(cfg)

	// The typical pattern is to build all your controller/clients then just pass to each handler
	// the bare minimum of what they need.
	Register(ctx,
		kubeClient.CoreV1().Events(""),
		core.Core().V1().Secret(),
		creds.Credentials().V1alpha1().Credential())

	// Start all the controllers
	if err := start.All(ctx, 2, core, creds); err != nil {
		logrus.Fatalf("Error starting: %s", err.Error())
	}

	<-ctx.Done()
}
