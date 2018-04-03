package brain

import (
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
)

// CreateAdminNote creates a new note on the specified thing 'spec'.
func CreateAdminNote(client lib.Client, on string, spec string, note string) (err error) {
	req, err := client.BuildRequest("POST", lib.BrainEndpoint, "/admin/notes")
	if err != nil {
		return
	}

	adminNote := brain.AdminNote{
		On:   on,
		Spec: spec,
		Note: note,
	}

	_, _, err = req.MarshalAndRun(adminNote, nil)
	return
}
