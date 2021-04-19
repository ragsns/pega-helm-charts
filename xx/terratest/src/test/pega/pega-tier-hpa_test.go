package pega

import (
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/stretchr/testify/require"
	autoscaling "k8s.io/kubernetes/pkg/apis/autoscaling"
	api "k8s.io/kubernetes/pkg/apis/core"
	"path/filepath"
	"testing"
	"strings"
	"fmt"
)



func TestPegaTierHPA(t *testing.T){
	var supportedVendors = []string{"k8s", "openshift", "eks","gke","aks","pks"}
	var supportedOperations =  []string{"deploy","install-deploy","upgrade-deploy"}

	helmChartPath, err := filepath.Abs(PegaHelmChartPath)
	require.NoError(t, err)


	for _,vendor := range supportedVendors{

		for _,operation := range supportedOperations{

			fmt.Println(vendor + "-" + operation)

			var options = &helm.Options{			
				SetValues: map[string]string{
					"global.provider":        vendor,
					"global.actions.execute": operation,
			 	},
		    }

			yamlContent := RenderTemplate(t, options, helmChartPath, []string{"templates/pega-tier-hpa.yaml"})
			VerifyPegaHPAs(t,yamlContent, options)

		}
	}


}
// VerifyPegaHPAs - Splits the HPA object from the rendered template and asserts each HPA object
func VerifyPegaHPAs(t *testing.T, yamlContent string, options *helm.Options) {
	var pegaHpaObj autoscaling.HorizontalPodAutoscaler
	hpaSlice := strings.SplitAfter(yamlContent, "85")
	for index, hpaInfo := range hpaSlice {
		if index >= 0 && index <= 1 {
			UnmarshalK8SYaml(t, hpaInfo, &pegaHpaObj)
			if index == 0 {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-web-hpa", "pega-web", "Deployment", "apps/v1"})
			} else {
				VerifyPegaHpa(t, &pegaHpaObj, hpa{"pega-batch-hpa", "pega-batch", "Deployment", "apps/v1"})
			}
		}
	}
}

// VerifyPegaHpa - Performs Pega HPA assertions with the values as provided in default values.yaml
func VerifyPegaHpa(t *testing.T, hpaObj *autoscaling.HorizontalPodAutoscaler, expectedHpa hpa) {
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Name, expectedHpa.targetRefName)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.Kind, expectedHpa.kind)
	require.Equal(t, hpaObj.Spec.ScaleTargetRef.APIVersion, expectedHpa.apiversion)
	require.Equal(t, hpaObj.Spec.Metrics[0].Resource.Name, api.ResourceName("cpu"))
	require.Equal(t, hpaObj.Spec.Metrics[1].Resource.Name, api.ResourceName("memory"))
	require.Equal(t, hpaObj.Spec.MaxReplicas, int32(5))
}

type hpa struct {
	name          string
	targetRefName string
	kind          string
	apiversion    string
}