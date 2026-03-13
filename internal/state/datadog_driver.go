package state

import (
	"fmt"
	"os"

	nvidiav1alpha1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1alpha1"
)

// Datadog specific: use precompiled gdrcopy image if precompiled drivers are used
func DatadogGetGDRCopyImagePath(spec *nvidiav1alpha1.NVIDIADriverSpec, pool nodePool) (string, error) {
	if spec.UsePrecompiledDrivers() {
		return spec.GDRCopy.GetPrecompiledImagePath(pool.getOS(), pool.kernel)
	}
	return spec.GDRCopy.GetImagePath(pool.getOS())
}

// DatadogGetEFASpec returns EFA driver spec populated from environment variables
func DatadogGetEFASpec(pool nodePool) (*efaDriverSpec, error) {
	// Check if EFA is enabled via environment variable
	if os.Getenv("EFA_ENABLED") != "true" {
		return nil, nil
	}

	// Get EFA driver image path from environment variables
	imagePath, err := nvidiav1alpha1.GetEFAPrecompiledImagePath(pool.getOS(), pool.kernel)
	if err != nil {
		return nil, err
	}

	nvPeermemImagePath, err := nvidiav1alpha1.GetEFANVPeermemPrecompiledImagePath(pool.getOS(), pool.kernel)
	if err != nil {
		return nil, err
	}

	// Get installer image from environment variable (set by chart)
	installerImagePath := os.Getenv("RDMA_CORE_INSTALLER_IMAGE")
	if installerImagePath == "" {
		return nil, fmt.Errorf("RDMA_CORE_INSTALLER_IMAGE environment variable must be set when EFA is enabled")
	}

	return &efaDriverSpec{
		Enabled:            true,
		ImagePath:          imagePath,
		NVPeermemImagePath: nvPeermemImagePath,
		InstallerImagePath: installerImagePath,
	}, nil
}
