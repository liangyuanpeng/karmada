package kubernetes

import (
	"log"
	"os"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/patches"

	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	"sigs.k8s.io/yaml"
)

func TestAsqq(t *testing.T) {
	cmd := &CommandInitOption{}
	karmadaApiserverDeployment := cmd.makeKarmadaAPIServerDeployment()
	log.Println("container.image:", karmadaApiserverDeployment.Spec.Template.Spec.Containers[0].Image)
	// resJson, _ := json.Marshal(karmadaApiserverDeployment)
	// patchedData, err := yaml.YAMLToJSON(resJson)
	// if err != nil {
	// 	panic(err)
	// }
	patchBytes, err := os.ReadFile("/home/runner/work/lanactions/lanactions/karmada/pkg/karmadactl/cmdinit/kubernetes/patch.yaml")
	if err != nil {
		panic(err)
	}
	podYAML, err := kubeadmutil.MarshalToYaml(karmadaApiserverDeployment, v1.SchemeGroupVersion)
	patchTarget := &patches.PatchTarget{
		Name:                      karmadaApiserverDeployment.Name,
		StrategicMergePatchObject: appsv1.Deployment{},
		Data:                      podYAML,
	}
	log.Println("yaml:", string(podYAML))
	// Always convert the target data to JSON.
	patchedData, err := yaml.YAMLToJSON(patchTarget.Data)
	if err != nil {
		panic(err)
	}
	log.Println("patchedData:", string(patchedData))
	log.Println("patchBytes:", string(patchBytes))

	patchBytes, err = yaml.YAMLToJSON(patchBytes)
	if err != nil {
		panic(err)
	}

	patchedData, err = strategicpatch.StrategicMergePatch(
		patchedData,
		patchBytes,
		patchTarget.StrategicMergePatchObject,
	)
	if err != nil {
		t.Fatal("failed ptach!", err)
		panic(err)
	}
	kubeletBytes, err := yaml.JSONToYAML(patchedData)
	if err != nil {
		panic(err)
	}
	log.Println("kubeletBytes:", string(kubeletBytes))

}
