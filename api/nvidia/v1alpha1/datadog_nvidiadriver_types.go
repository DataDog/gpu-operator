package v1alpha1

import (
	"fmt"
	"strings"

	"github.com/NVIDIA/gpu-operator/internal/image"
	"github.com/regclient/regclient/types/ref"
)

// Datadog specific: be able to get precompiled gdrcopy image if precompiled drivers are used
// Complete copy/paste of `func (d *NVIDIADriverSpec) GetPrecompiledImagePath`
// GetPrecompiledImagePath returns the precompiled gdrcopy image path for a
// given os version and kernel version. Precompiled gdrcopy images follow
// the following format:
// <repository>/<image>:<gdrcopy-ver>-<kernel-ver>-<os-ver>
func (d *GDRCopySpec) GetPrecompiledImagePath(osVersion string, kernelVersion string) (string, error) {
	// We pass an empty string for the last arg, the imagePathEnvName, since
	// we do not want any environment variable in the operator container
	// to be used as the default driver image. This means that the driver
	// image must be specified in the NVIDIADriver CR spec.
	image, err := image.ImagePath(d.Repository, d.Image, d.Version, "")
	if err != nil {
		return "", fmt.Errorf("failed to get image path from crd: %w", err)
	}

	// specifying a digest in the spec is not supported when using precompiled
	if strings.Contains(image, "sha256:") {
		return "", fmt.Errorf("specifying image digest is not supported when precompiled is enabled")
	}

	// append '-<kernelVersion>-<osVersion>' to the driver tag
	image = fmt.Sprintf("%s-%s-%s", image, kernelVersion, osVersion)

	_, err = ref.New(image)
	if err != nil {
		return "", fmt.Errorf("failed to parse gdrcopy image path: %w", err)
	}

	return image, nil
}

// GetPrecompiledImagePath returns the precompiled EFA driver image path for a
// given OS version and kernel version. Precompiled EFA driver images follow
// the following format:
// <repository>/<image>:<major-version>-<kernel-ver>-<os-ver>
// Example: registry.ddbuild.io/images/efa-driver:v3-6.8.0-1047-aws-ubuntu22.04
func (e *EFASpec) GetPrecompiledImagePath(osVersion string, kernelVersion string) (string, error) {
	// Get base image path with major version (e.g., "v3")
	img, err := image.ImagePath(e.Repository, e.Image, e.Version, "")
	if err != nil {
		return "", fmt.Errorf("failed to get image path from crd: %w", err)
	}

	if strings.Contains(img, "sha256:") {
		return "", fmt.Errorf("specifying image digest is not supported when precompiled is enabled")
	}

	// Append -<kernel-ver>-<os-ver>
	img = fmt.Sprintf("%s-%s-%s", img, kernelVersion, osVersion)

	_, err = ref.New(img)
	if err != nil {
		return "", fmt.Errorf("failed to parse EFA image path: %w", err)
	}

	return img, nil
}

// GetPrecompiledImagePath returns the precompiled EFA NV Peermem driver image path for a
// given OS version and kernel version. Precompiled EFA NV Peermem driver images follow
// the following format:
// <repository>/<image>:<major-version>-<kernel-ver>-<os-ver>
// Example: registry.ddbuild.io/images/efa-nv-peermem-driver:v1-6.8.0-1047-aws-ubuntu22.04
func (e *EFANVPeermemSpec) GetPrecompiledImagePath(osVersion string, kernelVersion string) (string, error) {
	// Get base image path with major version (e.g., "v1")
	img, err := image.ImagePath(e.Repository, e.Image, e.Version, "")
	if err != nil {
		return "", fmt.Errorf("failed to get image path from crd: %w", err)
	}

	if strings.Contains(img, "sha256:") {
		return "", fmt.Errorf("specifying image digest is not supported when precompiled is enabled")
	}

	// Append -<kernel-ver>-<os-ver>
	img = fmt.Sprintf("%s-%s-%s", img, kernelVersion, osVersion)

	_, err = ref.New(img)
	if err != nil {
		return "", fmt.Errorf("failed to parse EFA NV Peermem image path: %w", err)
	}

	return img, nil
}
