package e2e

import (
	"fmt"
	"os"

	"github.com/onsi/ginkgo/v2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	"github.com/karmada-io/karmada/test/e2e/framework"
	"github.com/karmada-io/karmada/test/helper"
)

var _ = ginkgo.Describe("operator testing", func() {
	// var clusterClient kubernetes.Interface
	homeDir := os.Getenv("HOME")
	fmt.Println("==================begin operator testing")
	kubeConfigPath := fmt.Sprintf("%s/.kube/%s.config", homeDir, "karmada")
	fmt.Println("kubeConfigPath:", kubeConfigPath)
	// clusterConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	// if err != nil {
	// 	panic(err)
	// }
	restConfig, err := framework.LoadRESTClientConfig(kubeConfigPath, "karmada-host")
	if err != nil {
		panic(err)
	}
	kubeClient, err = kubernetes.NewForConfig(restConfig)
	// operatorClient, _ := operator.NewForConfig(clusterConfig)

	ginkgo.BeforeEach(func() {
		defaultConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
		defaultConfigFlags.Context = &karmadaContext
		// kubeconfig = os.Getenv("KUBECONFIG")
		// gomega.Expect(kubeconfig).ShouldNot(gomega.BeEmpty())
		// clusterProvider = cluster.NewProvider()
		// var err error
		// restConfig, err = framework.LoadRESTClientConfig(kubeconfig, karmadaContext)
	})

	ginkgo.Context("[operator]Test namespaced resource: deployment", func() {
		var deployment *appsv1.Deployment
		var deploymentNamespace, deploymentName string

		ginkgo.BeforeEach(func() {
			deploymentNamespace = fmt.Sprintf("karmadatest-%s", rand.String(RandomStrLength))
			deploymentName = deploymentNamePrefix + rand.String(RandomStrLength)
			deployment = helper.NewDeployment(deploymentNamespace, deploymentName)
		})

		ginkgo.AfterEach(func() {
			framework.RemoveDeployment(kubeClient, deploymentNamespace, deploymentName)
			framework.RemoveNamespace(kubeClient, deploymentNamespace)
		})

		ginkgo.It("[operator]test", func() {

			ginkgo.By(fmt.Sprintf("Creating deployment %s with namespace %s ", deploymentName, deploymentNamespace), func() {
				deploymentNamespaceObj := helper.NewNamespace(deploymentNamespace)
				framework.CreateNamespace(kubeClient, deploymentNamespaceObj)
				framework.CreateDeployment(kubeClient, deployment)
				// karmada := &operatorv1alpha1.Karmada{}
				// operatorClient.OperatorV1alpha1().Karmadas("default").Create(context.TODO(), karmada, v1.CreateOptions{})
			})
		})

	})

})
