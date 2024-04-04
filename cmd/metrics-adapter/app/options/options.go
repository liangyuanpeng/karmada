/*
Copyright 2023 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"context"
	"log"

	"github.com/spf13/pflag"
	openapinamer "k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/cmd/options"
	"sigs.k8s.io/metrics-server/pkg/api"

	karmadaclientset "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	informerfactory "github.com/karmada-io/karmada/pkg/generated/informers/externalversions"
	generatedopenapi "github.com/karmada-io/karmada/pkg/generated/openapi"
	"github.com/karmada-io/karmada/pkg/metricsadapter"
	"github.com/karmada-io/karmada/pkg/sharedcli/profileflag"
	"github.com/karmada-io/karmada/pkg/version"
)

// Options contains everything necessary to create and run metrics-adapter.
type Options struct {
	CustomMetricsAdapterServerOptions *options.CustomMetricsAdapterServerOptions

	KubeConfig string

	ProfileOpts profileflag.Options
}

// NewOptions builds a default metrics-adapter options.
func NewOptions() *Options {
	o := &Options{
		CustomMetricsAdapterServerOptions: options.NewCustomMetricsAdapterServerOptions(),
	}

	return o
}

// Complete fills in fields required to have valid data.
func (o *Options) Complete() error {
	return nil
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.CustomMetricsAdapterServerOptions.AddFlags(fs)
	o.ProfileOpts.AddFlags(fs)

	fs.StringVar(&o.KubeConfig, "kubeconfig", o.KubeConfig, "Path to karmada control plane kubeconfig file.")
}

// Config returns config for the metrics-adapter server given Options
func (o *Options) Config() (*metricsadapter.MetricsServer, error) {
	restConfig, err := clientcmd.BuildConfigFromFlags("", o.KubeConfig)
	if err != nil {
		klog.Errorf("Unable to build restConfig: %v", err)
		return nil, err
	}
	log.Println("kubeconfig.opt:", restConfig.ContentType)

	karmadaClient := karmadaclientset.NewForConfigOrDie(restConfig)
	factory := informerfactory.NewSharedInformerFactory(karmadaClient, 0)
	kubeClient := kubernetes.NewForConfigOrDie(restConfig)
	kubeFactory := informers.NewSharedInformerFactory(kubeClient, 0)
	metricsController := metricsadapter.NewMetricsController(restConfig, factory, kubeFactory)
	metricsAdapter := metricsadapter.NewMetricsAdapter(metricsController, o.CustomMetricsAdapterServerOptions)
	metricsAdapter.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(api.Scheme))
	metricsAdapter.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(api.Scheme))
	metricsAdapter.OpenAPIConfig.Info.Title = "karmada-metrics-adapter"
	metricsAdapter.OpenAPIConfig.Info.Version = "1.0.0"

	server, err := metricsAdapter.Server()
	if err != nil {
		klog.Errorf("Unable to construct metrics adapter: %v", err)
		return nil, err
	}

	err = server.GenericAPIServer.AddPostStartHook("start-karmada-informers", func(context genericapiserver.PostStartHookContext) error {
		kubeFactory.Core().V1().Secrets().Informer()
		kubeFactory.Start(context.StopCh)
		kubeFactory.WaitForCacheSync(context.StopCh)
		factory.Start(context.StopCh)
		return nil
	})
	if err != nil {
		klog.Errorf("Unable to add post hook: %v", err)
		return nil, err
	}

	if err := api.Install(metricsAdapter, metricsAdapter.PodLister, metricsAdapter.NodeLister, server.GenericAPIServer, nil); err != nil {
		klog.Errorf("Unable to install resource metrics adapter: %v", err)
		return nil, err
	}

	return metricsadapter.NewMetricsServer(metricsController, metricsAdapter), nil
}

// Run runs the metrics-adapter with options. This should never exit.
func (o *Options) Run(ctx context.Context) error {
	klog.Infof("karmada-metrics-adapter version: %s", version.Get())

	profileflag.ListenAndServe(o.ProfileOpts)

	metricsServer, err := o.Config()
	if err != nil {
		return err
	}

	return metricsServer.StartServer(ctx.Done())
}
