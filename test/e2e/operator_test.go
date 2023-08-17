package e2e

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	"github.com/karmada-io/karmada/test/e2e/framework"
	"github.com/karmada-io/karmada/test/helper"
)

var _ = ginkgo.Describe("operator testing", func() {
	var cluster string
	var clusterClient kubernetes.Interface

	ginkgo.BeforeEach(func() {
		cluster = "host"
		clusterClient = framework.GetClusterClient(cluster)
		defaultConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
		defaultConfigFlags.Context = &karmadaContext
	})

	ginkgo.Context("Test promoting namespaced resource: deployment", func() {
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

		ginkgo.It("Test promoting a deployment from cluster member", func() {

			// Step 1,  create namespace and deployment on cluster member1
			ginkgo.By(fmt.Sprintf("Creating deployment %s with namespace %s not existed in karmada control plane", deploymentName, deploymentNamespace), func() {
				deploymentNamespaceObj := helper.NewNamespace(deploymentNamespace)
				framework.CreateNamespace(clusterClient, deploymentNamespaceObj)
				framework.CreateDeployment(clusterClient, deployment)
			})
		})

	})

})
