package state

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/NVIDIA/gpu-operator/internal/render"
)

func getYAMLStringDatadog(objs []*unstructured.Unstructured) (string, error) {
	s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme,
		scheme.Scheme, json.SerializerOptions{Yaml: true, Pretty: false, Strict: false})
	var sb strings.Builder
	for _, obj := range objs {
		var b bytes.Buffer
		err := s.Encode(obj, &b)
		if err != nil {
			return "", err
		}
		sb.WriteString(b.String())
		sb.WriteString("---\n")
	}
	return sb.String(), nil
}

func TestDriverEFA(t *testing.T) {
	const (
		testName = "driver-efa"
	)

	state, err := NewStateDriver(nil, "", nil, manifestDir)
	require.Nil(t, err)
	stateDriver, ok := state.(*stateDriver)
	require.True(t, ok)

	renderData := getMinimalDriverRenderData()

	// Configure EFA driver
	renderData.EFA = &efaDriverSpec{
		Enabled:            true,
		ImagePath:          "727006795293.dkr.ecr.us-east-1.amazonaws.com/ckmb/efa:v3-6.8.0-1046-aws-ubuntu22.04",
		InstallerImagePath: "727006795293.dkr.ecr.us-east-1.amazonaws.com/images/aws-efa-installer:1.0.0",
	}

	// Add EFA NV Peermem
	renderData.EFANVPeermem = &efaNVPeermemDriverSpec{
		Enabled:   true,
		ImagePath: "727006795293.dkr.ecr.us-east-1.amazonaws.com/ckmb/efa-nv-peermem:v1-6.8.0-1046-aws-ubuntu22.04",
	}

	// Add RDMA Core
	renderData.RDMACoreEnabled = true

	objs, err := stateDriver.renderer.RenderObjects(
		&render.TemplatingData{
			Data: renderData,
		})
	require.Nil(t, err)
	require.NotEmpty(t, objs)

	actual, err := getYAMLStringDatadog(objs)
	require.Nil(t, err)

	// Print the generated YAML for inspection
	t.Logf("Generated EFA DaemonSet:\n%s", actual)
}
