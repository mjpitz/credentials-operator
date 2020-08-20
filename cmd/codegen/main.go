package main

import (
	"github.com/mjpitz/credentials-operator/pkg/apis/credentials.mjpitz.com/v1alpha1"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	controllergen.Run(args.Options{
		OutputPackage: "github.com/mjpitz/credentials-operator/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"credentials.mjpitz.com": {
				Types: []interface{}{
					v1alpha1.Credential{},
				},
				GenerateTypes: true,
			},
		},
	})
}
