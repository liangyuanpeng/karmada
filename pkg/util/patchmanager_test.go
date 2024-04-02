package util

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"sigs.k8s.io/yaml"
)

func testa(t *testing.T) {
	karmadaDeschedulerDeployment := &appsv1.Deployment{}
	patchedData, err := yaml.YAMLToJSON(patchTarget.Data)
	if err != nil {
		panic(err)
	}
	patchedData, err = strategicpatch.StrategicMergePatch(
		patchedData,
		patchBytes,
		patchTarget.StrategicMergePatchObject,
	)
}
