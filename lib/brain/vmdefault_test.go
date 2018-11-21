package brain

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"
	"github.com/cheekybits/is"
)

func getFixtureVMD() (vm VirtualMachineDefault) {
	return VirtualMachineDefault{
		Name:   "vmd-name",
		Public: true,
		ServerSettings: VirtualMachineSpec{
			VirtualMachine: VirtualMachine{
				CdromURL:        "test-url",
				Cores:           1,
				Memory:          1,
				HardwareProfile: "test-profile",
				ZoneName:        "test-zone",
				Discs:           nil,
			},
			Discs: Discs{
				Disc{
					StorageGrade: "sata",
					Size:         1,
					BackupSchedules: BackupSchedules{{
						Interval: 604800,
						Capacity: 1,
					}},
				},
			},
			Reimage: &ImageInstall{
				Distribution:    "test-image",
				FirstbootScript: "test-script",
			},
		},
	}
}

func TestStringFields(t *testing.T) {
	is := is.New(t)
	vmd := getFixtureVMD()
	is.Equal("vmd-name", vmd.Name)
	is.Equal("test-profile", vmd.ServerSettings.VirtualMachine.HardwareProfile)
	is.Equal("test-zone", vmd.ServerSettings.VirtualMachine.ZoneName)
}

func TestVMDJSON(t *testing.T) {
	tests := []struct {
		vmd      VirtualMachine
		expected map[string]interface{}
	}{
		{
			vmd: VirtualMachine{CdromURL: "test_url"},
			expected: map[string]interface{}{
				"cdrom_url": "test_url",
			},
		},
	}

	for i, test := range tests {
		js, err := json.Marshal(test.vmd)
		if err != nil {
			t.Fatalf("TestVMDJSON #%d json.Marshal failed: %v\r\n", i, err.Error())
		}
		unmarshalled := make(map[string]interface{})
		err = json.Unmarshal(js, &unmarshalled)
		if err != nil {
			t.Fatalf("TestVMDJSON #%d json.Unmarshal failed: %v\r\n", i, err.Error())
		}
		if !reflect.DeepEqual(test.expected, unmarshalled) {
			t.Fatalf("TestVMDJSON #%d failed.\r\nEXPECTED\r\n%#v\r\nACTUAL\r\n%#v", i, test.expected, unmarshalled)
		}
	}
}

func TestFormatVMD(t *testing.T) {
	is := is.New(t)
	b := new(bytes.Buffer)
	vmd := getFixtureVMD()

	tests := []struct {
		in     VirtualMachineDefault
		detail prettyprint.DetailLevel
		expt   string
	}{
		{
			in:     vmd,
			detail: prettyprint.SingleLine,
			expt:   " ▸ vmd-name with public => true",
		},
	}

	for _, test := range tests {
		b.Truncate(0)
		err := test.in.PrettyPrint(b, test.detail)
		if err != nil {
			t.Error(err)
		}
		is.Equal(test.expt, b.String())
	}
}
