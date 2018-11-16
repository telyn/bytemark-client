package brain_test

import (
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainRequests "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
)

func simplePostTest(t *testing.T, url string, testBody string, runTest func(lib.Client) error) {
	testName := testutil.Name(0)
	rts := testutil.RequestTestSpec{
		Method:        "POST",
		Endpoint:      lib.BrainEndpoint,
		URL:           url,
		AssertRequest: assert.BodyString(testBody),
	}
	rts.Run(t, testName, true, func(client lib.Client) {

		err := runTest(client)
		if err != nil {
			t.Errorf("%s errored: %s", testName, err.Error())
		}
	})
}

func TestPostCreateVMDefault(t *testing.T) {
	// TODO(tom): shorten and fix, too messy & remove non required

	simplePostTest(t, "/vm_defaults",
		`{"name":"vmd-name","public":true,"server_settings":{` +
		`"vm_default":{"cores":1,"name":"vm"},` +
		`"disc":[{"storage_grade":"sata","size":1024,` +
		`"backup_schedules":[{"start_at":"","interval_seconds":604800,"capacity":1}]}],` +
		`"reimage":{"distribution":"image","firstboot_script":"script","root_password":"","ssh_public_key":""}}}`,

		func(client lib.Client) error {
			return brainRequests.CreateVMDefault(client,"vmd-name",true, brain.VMDefaultSpec{
				VMDefault: brain.VMDefault{
					CdromURL:         "",
					Cores:            1,
					Memory:           0,
					Name:             "vm",
					HardwareProfile:  "",
					ZoneName:         "",
				},
				Discs: brain.Discs{
					brain.Disc{
						StorageGrade: "sata",
						Size:         1024,
						BackupSchedules: brain.BackupSchedules{{
							Interval: 604800,
							Capacity: 1,
						}},
					},
				},
				Reimage: &brain.ImageInstall{
					Distribution: 	  "image",
					FirstbootScript:  "script",
					RootPassword: 	  "",
					PublicKeys:       "",
				},
			})
		})
}