package brain

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/lib/output/prettyprint"

	"github.com/cheekybits/is"
)

func getFixtureVMD() (vm VMDefault) {
	disc := getFixtureDisc()

	return VMDefault{
		CdromURL:        "",
		Cores:           1,
		Memory:          1,
		Name:            "vmd-name",
		HardwareProfile: "virtio_test",
		ZoneName:        "york",
		Discs: []Disc{
			disc,
		},
		ID: 0,
	}
}

func TestStringFields(t *testing.T) {
	is := is.New(t)
	vmd := getFixtureVMD()
	is.Equal("vmd-name", vmd.Name)
	is.Equal("virtio_test", vmd.HardwareProfile)
	is.Equal("york", vmd.ZoneName)
}

// TODO(tom): remove disc test? already tested in VirtualMachine
func TestVMDDiscs(t *testing.T) {
	is := is.New(t)
	discs := getFixtureDiscSet()
	for _, d := range discs {
		d2, err := d.Validate()
		is.Nil(err)

		is.Equal(d.Size, d2.Size)
		switch d.ID {
		case 1, 3:
			is.Equal("sata", d2.StorageGrade)
		case 2:
			is.Equal("archive", d2.StorageGrade)
		}
	}
}

func TestVMDJSON(t *testing.T) {
	tests := []struct {
		vmd      VMDefault
		expected map[string]interface{}
	}{
		{
			vmd: VMDefault{CdromURL: "test_url", Name: "vmd_name"},
			expected: map[string]interface{}{
				"cdrom_url": "test_url",
				"name":      "vmd_name",
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
		in     VMDefault
		detail prettyprint.DetailLevel
		expt   string
	}{
		{
			in:     vmd,
			detail: prettyprint.SingleLine,
			expt:   " ▸ vmd-name in York",
		},
		{
			in:     vmd,
			detail: prettyprint.Medium,
			expt: ` ▸ vmd-name in York
   - 1 core, 1MiB, 25GiB on 1 disc`,
		},
		{
			in:     vmd,
			detail: prettyprint.Full,
			expt: ` ▸ vmd-name in York
   - 1 core, 1MiB, 25GiB on 1 disc

    discs:
      •  - 25GiB, sata grade

`,
		},
		{
			in:     VMDefault{},
			detail: "_discs",
			expt:   "",
		},
		{
			in:     VMDefault{},
			detail: "_spec",
			expt:   "   - 0 cores, 0MiB, no discs",
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
