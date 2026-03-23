package v1alpha1

import (
	"fmt"
	"os"
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

// GetEFAPrecompiledImagePath returns the precompiled EFA driver image path from environment variables.
// Reads: EFA_REPOSITORY, EFA_IMAGE, EFA_VERSION
// Format: <repository>/<image>:<major-version>-<kernel-ver>-<os-ver>
// Example: registry.ddbuild.io/images/efa-driver:v3-6.8.0-1047-aws-ubuntu22.04
func GetEFAPrecompiledImagePath(osVersion string, kernelVersion string) (string, error) {
	repository := os.Getenv("EFA_REPOSITORY")
	imageName := os.Getenv("EFA_IMAGE")
	version := os.Getenv("EFA_VERSION")

	if repository == "" || imageName == "" || version == "" {
		return "", fmt.Errorf("EFA environment variables not set (EFA_REPOSITORY, EFA_IMAGE, EFA_VERSION)")
	}

	img, err := image.ImagePath(repository, imageName, version, "")
	if err != nil {
		return "", fmt.Errorf("failed to get EFA image path: %w", err)
	}

	if strings.Contains(img, "sha256:") {
		return "", fmt.Errorf("specifying image digest is not supported when precompiled is enabled")
	}

	img = fmt.Sprintf("%s-%s-%s", img, kernelVersion, osVersion)

	_, err = ref.New(img)
	if err != nil {
		return "", fmt.Errorf("failed to parse EFA image path: %w", err)
	}

	return img, nil
}

// GetEFANVPeermemPrecompiledImagePath returns the precompiled EFA NV Peermem driver image path from environment variables.
// Reads: EFA_NV_PEERMEM_REPOSITORY, EFA_NV_PEERMEM_IMAGE, EFA_NV_PEERMEM_VERSION
// Format: <repository>/<image>:<major-version>-<kernel-ver>-<os-ver>
// Example: registry.ddbuild.io/images/efa-nv-peermem-driver:v1-6.8.0-1047-aws-ubuntu22.04
func GetEFANVPeermemPrecompiledImagePath(osVersion string, kernelVersion string) (string, error) {
	repository := os.Getenv("EFA_NV_PEERMEM_REPOSITORY")
	imageName := os.Getenv("EFA_NV_PEERMEM_IMAGE")
	version := os.Getenv("EFA_NV_PEERMEM_VERSION")

	if repository == "" || imageName == "" || version == "" {
		return "", fmt.Errorf("EFA NV Peermem environment variables not set (EFA_NV_PEERMEM_REPOSITORY, EFA_NV_PEERMEM_IMAGE, EFA_NV_PEERMEM_VERSION)")
	}

	img, err := image.ImagePath(repository, imageName, version, "")
	if err != nil {
		return "", fmt.Errorf("failed to get EFA NV Peermem image path: %w", err)
	}

	if strings.Contains(img, "sha256:") {
		return "", fmt.Errorf("specifying image digest is not supported when precompiled is enabled")
	}

	img = fmt.Sprintf("%s-%s-%s", img, kernelVersion, osVersion)

	_, err = ref.New(img)
	if err != nil {
		return "", fmt.Errorf("failed to parse EFA NV Peermem image path: %w", err)
	}

	return img, nil
}
