package state

// efaDriverSpec contains EFA driver configuration populated from environment variables
type efaDriverSpec struct {
	Enabled            bool
	ImagePath          string
	InstallerImagePath string // Populated from EFA_INSTALLER_IMAGE env var
}

// efaNVPeermemDriverSpec contains EFA NV Peermem driver configuration populated from environment variables
type efaNVPeermemDriverSpec struct {
	Enabled   bool
	ImagePath string
}
