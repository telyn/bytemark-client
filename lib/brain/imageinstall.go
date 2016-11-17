package brain

// ImageInstall represents what image was most recently installed on a VM along with its root password.
// This might only be returned when creating a VM.
type ImageInstall struct {
	Distribution    string `json:"distribution"`
	FirstbootScript string `json:"firstboot_script"`
	RootPassword    string `json:"root_password"`
	PublicKeys      string `json:"ssh_public_key"`
}
