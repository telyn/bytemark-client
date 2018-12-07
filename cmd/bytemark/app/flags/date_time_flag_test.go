package flags_test

import (
	"testing"
	"time"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app/flags"
	"github.com/cheekybits/is"
)

func TestDateTimeFlag(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		input  string
		layout string
	}{
		{
			input:  "1/2/2018 15:30:00 +01:00",
			layout: "2006-01-02T15:04:05-07:00",
		},
		{
			input:  "15:30:00",
			layout: "T15:04:05",
		},
		{
			input:  "1/2/2018 15:30",
			layout: "2006-01-02T15:04",
		},
		{
			input:  "1-2-2018",
			layout: "2006-01-02",
		},
		{
			input:  "1st June 2018",
			layout: "2006-01-02",
		},
		{
			input:  "June 1 2018",
			layout: "2006-01-02",
		},
	}

	for _, test := range tests {
		t.Run("datetime flag", func(t *testing.T) {
			var flag flags.DateTimeFlag

			err := flag.Set(test.input)

			is.Nil(err)
			is.OK(flag.String())

			// check format
			formatted, err := time.Parse(test.layout, flag.String())
			is.Nil(err)
			is.OK(formatted)
		})
	}
}
