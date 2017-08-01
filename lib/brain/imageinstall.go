package brain

import (
	"fmt"
	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"io"
	"strings"
)

// ImageInstall represents what image was most recently installed on a VM along with its root password.
// This might only be returned when creating a VM.
type ImageInstall struct {
	Distribution    string `json:"distribution"`
	FirstbootScript string `json:"firstboot_script"`
	RootPassword    string `json:"root_password"`
	PublicKeys      string `json:"ssh_public_key"`
}

// PrettyPrint outputs the image install with the given level of detail.
// TODO(telyn): rewrite to use templates & prettyprint.Run
func (pp ImageInstall) PrettyPrint(wr io.Writer, detail prettyprint.DetailLevel) error {
	var output []string
	if pp.Distribution != "" {
		output = append(output, "Image: "+pp.Distribution)
	}
	if pp.PublicKeys != "" {
		var keynames []string
		for _, k := range strings.Split(pp.PublicKeys, "\n") {
			kbits := strings.SplitN(k, " ", 3)
			if len(kbits) == 3 {
				keynames = append(keynames, strings.TrimSpace(kbits[2]))
			}

		}
		output = append(output, fmt.Sprintf("%d public keys: %s", len(keynames), strings.Join(keynames, ", ")))
	}
	if pp.RootPassword != "" {
		output = append(output, "Root/Administrator password: "+pp.RootPassword)
	}
	if pp.FirstbootScript != "" {
		output = append(output, "With a firstboot script")
	}
	_, err := wr.Write([]byte(strings.Join(output, "\r\n") + "\r\n"))
	return err
}
