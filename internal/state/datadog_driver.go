package state

import (
	"fmt"

	nvidiav1alpha1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1alpha1"
)

// Datadog specific: use precompiled gdrcopy image if precompiled drivers are used
func DatadogGetGDRCopyImagePath(spec *nvidiav1alpha1.NVIDIADriverSpec, pool nodePool) (string, error) {
	if spec.UsePrecompiledDrivers() {
		return spec.GDRCopy.GetPrecompiledImagePath(pool.getOS(), pool.kernel)
	}
	return spec.GDRCopy.GetImagePath(pool.getOS())
}

// DatadogGetEFAImagePath returns EFA driver image path (precompiled only)
func DatadogGetEFAImagePath(spec *nvidiav1alpha1.NVIDIADriverSpec, pool nodePool) (string, error) {
	// EFA only supports precompiled
	if !spec.UsePrecompiledDrivers() {
		return "", fmt.Errorf("EFA driver requires precompiled drivers to be enabled")
	}
	return spec.EFA.GetPrecompiledImagePath(pool.getOS(), pool.kernel)
}

// DatadogGetEFANVPeermemImagePath returns EFA NV Peermem driver image path (precompiled only)
func DatadogGetEFANVPeermemImagePath(spec *nvidiav1alpha1.NVIDIADriverSpec, pool nodePool) (string, error) {
	// EFA NV Peermem only supports precompiled
	if !spec.UsePrecompiledDrivers() {
		return "", fmt.Errorf("EFA NV Peermem driver requires precompiled drivers to be enabled")
	}
	return spec.EFANVPeermem.GetPrecompiledImagePath(pool.getOS(), pool.kernel)
}
