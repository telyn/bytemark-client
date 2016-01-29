package lib

import (
	"github.com/cheekybits/is"
	"testing"
)

func TestBadRequestError(t *testing.T) {
	is := is.New(t)
	type test struct {
		input  string
		output []string
	}
	tests := []test{
		test{`{"discs":[{},{},{},{},{"size":["must be greater than or equal to 50"]}]}`,
			[]string{
				"• Disc 5 (vde) - size must be greater than or equal to 50",
			},
		},
		test{`{"name":["can't be blank"],"memory":["is not included in the list","is not a number"]}`,
			[]string{
				"• Name cannot be blank",
				"• Memory amount was not set",
			},
		},
		test{`{"name":["is invalid","is too short (minimum is 3 characters)"],"memory":["is not included in the list","is not a number"]}`,
			[]string{
				"• Name is too short (minimum is 3 characters)",
				"• Memory amount was not set",
			},
		},
	}

	for i, d := range tests {
		err := newBadRequestError([]byte(input))
		is.Equal(d.output, err.FormatProblems())
	}
}
