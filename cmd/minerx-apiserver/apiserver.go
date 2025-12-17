package main

import (
	"context"
	"os"

	_ "go.uber.org/automaxprocs/maxprocs"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apiserver/pkg/admission"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
	_ "k8s.io/component-base/logs/json/register"          // for JSON log format registration
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugins
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/klog/v2"

	"github.com/LiangNing7/minerx/cmd/minerx-apiserver/app"
	"github.com/LiangNing7/minerx/internal/apiserver/admission/plugin/minerset"

	"github.com/LiangNing7/minerx/internal/apiserver/admission/initializer"
	appsrest "github.com/LiangNing7/minerx/internal/apiserver/registry/apps/rest"
	batchrest "github.com/LiangNing7/minerx/internal/apiserver/registry/batch/rest"
	"github.com/LiangNing7/minerx/internal/pkg/config/minerprofile"
	appsv1beta1 "github.com/LiangNing7/minerx/pkg/apis/apps/v1beta1"
	batchv1beta1 "github.com/LiangNing7/minerx/pkg/apis/batch/v1beta1"
	"github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
	"github.com/LiangNing7/minerx/pkg/generated/informers"
	generatedopenapi "github.com/LiangNing7/minerx/pkg/generated/openapi"
)

func main() {
	var informerFactory informers.SharedInformerFactory

	// Please note that the following WithOptions are all required.
	command := app.NewAPIServerCommand(
		// Add custom etcd options.
		app.WithEtcdOptions("/registry/liangning7.cn", appsv1beta1.SchemeGroupVersion, batchv1beta1.SchemeGroupVersion),
		// Add custom resource storage.
		app.WithRESTStorageProviders(appsrest.RESTStorageProvider{}, batchrest.RESTStorageProvider{}),
		// Add custom dns address.
		app.WithAlternateDNS("liangning7.cn"),
		// Add custom admission plugins.
		app.WithAdmissionPlugin(minerset.PluginName, minerset.Register),
		// Add custom admission plugins initializer.
		app.WithGetOpenAPIDefinitions(generatedopenapi.GetOpenAPIDefinitions),
		app.WithAdmissionInitializers(func(c *genericapiserver.RecommendedConfig) ([]admission.PluginInitializer, error) {
			client, err := versioned.NewForConfig(c.LoopbackClientConfig)
			if err != nil {
				return nil, err
			}
			informerFactory = informers.NewSharedInformerFactory(client, c.LoopbackClientConfig.Timeout)
			// NOTICE: As we create a shared informer, we need to start it later.
			// We can usually start it by adding a PostStartHook.
			return []admission.PluginInitializer{initializer.New(informerFactory, client)}, nil
		}),
		app.WithPostStartHook(
			"start-external-informers",
			func(ctx genericapiserver.PostStartHookContext) error {
				if informerFactory != nil {
					informerFactory.Start(ctx.Done())
				}
				return nil
			}),
		app.WithPostStartHook(
			"initialize-instance-config-client",
			func(ctx genericapiserver.PostStartHookContext) error {
				client, err := versioned.NewForConfig(ctx.LoopbackClientConfig)
				if err != nil {
					return err
				}

				if err := minerprofile.Init(context.Background(), client); err != nil {
					// When returning 'NotFound' error, we should not report an error, otherwise we can not
					// create 'MinerTypesConfigMapName' configmap via minerx-apiserver
					if apierrors.IsNotFound(err) {
						return nil
					}

					klog.ErrorS(err, "Failed to init miner type cache")
					return err
				}

				return nil
			},
		),
	)

	code := cli.Run(command)
	os.Exit(code)
}
