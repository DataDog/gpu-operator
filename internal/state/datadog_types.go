package state

// efaDriverSpec contains EFA driver configuration populated from environment variables
type efaDriverSpec struct {
	Enabled              bool
	ImagePath            string // EFA driver image (precompiled, kernel-specific)
	NVPeermemImagePath   string // EFA NV Peermem driver image (precompiled, kernel-specific)
	InstallerImagePath   string // RDMA core installer image (from RDMA_CORE_INSTALLER_IMAGE)
}
