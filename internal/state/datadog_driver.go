package state

import nvidiav1alpha1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1alpha1"

// Datadog specific: use precompiled gdrcopy image if precompiled drivers are used
func DatadogGetGDRCopyImagePath(spec *nvidiav1alpha1.NVIDIADriverSpec, pool nodePool) (string, error) {
	if spec.UsePrecompiledDrivers() {
		return spec.GDRCopy.GetPrecompiledImagePath(pool.osTag, pool.kernel)
	}
	return spec.GDRCopy.GetImagePath(pool.osTag)
}
