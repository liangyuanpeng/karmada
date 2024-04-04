package kubernetes

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"sigs.k8s.io/yaml"
)

func testa(t *testing.T) {
	cmd := &CommandInitOption{}
	karmadaApiserverDeployment := cmd.makeKarmadaAPIServerDeployment()
	log.Println("container.image:", karmadaApiserverDeployment.Spec.Template.Spec.Containers[0].Image)
	resJson, _ := json.Marshal(karmadaApiserverDeployment)
	patchedData, err := yaml.YAMLToJSON(resJson)
	if err != nil {
		panic(err)
	}
	patchBytes, err := os.ReadFile("/home/runner/work/lanactions/lanactions/karmada/pkg/karmadactl/cmdinit/kubernetes/patch.yaml")
	if err != nil {
		panic(err)
	}
	patchedData, err = strategicpatch.StrategicMergePatch(
		patchedData,
		patchBytes,
		patchTarget.StrategicMergePatchObject,
	)
}
