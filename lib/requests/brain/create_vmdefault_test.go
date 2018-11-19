package brain_test

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	brainMethods "github.com/BytemarkHosting/bytemark-client/lib/requests/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib/testutil/assert"
)

func TestCreateVMDefault(t *testing.T) {
	tests := []struct {
		name           string
		public         bool
		serverSettings brain.VMDefaultSpec
		expected       map[string]interface{}
		shouldErr      bool
	}{
		{
			name:   "vmd-test",
			public: true,
			serverSettings: brain.VMDefaultSpec{
				VMDefault: brain.VMDefault{
					CdromURL:        "test-url",
					Cores:           1,
					Memory:          1024,
					Name:            "test-vm",
					HardwareProfile: "test-hwp",
					ZoneName:        "test-zone",
				},
				Discs: brain.Discs{
					brain.Disc{
						StorageGrade: "test-grade",
						Size:         1,
						BackupSchedules: brain.BackupSchedules{{
							StartDate: "",
							Interval:  7 * 86400,
							Capacity:  1,
						}},
					},
				},
				Reimage: &brain.ImageInstall{
					Distribution:    "test-image",
					FirstbootScript: "test-script",
					PublicKeys:      "",
					RootPassword:    "",
				},
			},
			expected: map[string]interface{}{
				"name":   "vmd-test",
				"public": true,
				"server_settings": map[string]interface{}{
					"vm_default": map[string]interface{}{
						"cdrom_url":        "test-url",
						"cores":            1,
						"memory":           1024,
						"name":             "test-vm",
						"hardware_profile": "test-hwp",
						"zone_name":        "test-zone"},
					"disc": []interface{}{map[string]interface{}{
						"storage_grade": "test-grade",
						"size":          1,
						"backup_schedules": map[string]interface{}{
							"start_at":         "",
							"interval_seconds": 604800,
							"capacity":         1}}},
					"reimage": map[string]interface{}{
						"distribution":     "test-image",
						"firstboot_script": "test-script",
						"root_password":    "",
						"ssh_public_key":   ""},
				},
			},
		},
	}
	for i, test := range tests {
		testName := testutil.Name(i)
		rts := testutil.RequestTestSpec{
			Method:   "POST",
			Endpoint: lib.BrainEndpoint,
			URL:      fmt.Sprintf("/vm_defaults"),
			// TODO(tom): Implement this properly. Objects contain same keys, but not ordered bc map
			//AssertRequest: assert.BodyUnmarshalEqual(test.expected),
		}
		rts.Run(t, testName, true, func(client lib.Client) {
			err := brainMethods.CreateVMDefault(client, test.name, test.public, test.serverSettings)
			if test.shouldErr {
				assert.NotEqual(t, testName, nil, err)
			} else {
				assert.Equal(t, testName, nil, err)
			}
		})
	}
}
