package util

import (
	"time"
	"fmt"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

// DateTimeFlag holds datatime in iso8601 format
type DateTimeFlag string

// Set takes user input and attempts to parse the datetime from any format to iso8601
func (dtf *DateTimeFlag) Set(value string) error {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	when, err := w.Parse(value, time.Now())

	*dtf = DateTimeFlag(when.Time.Format("2006-01-02T15:04:05-0700"))

	return err
}

func (dtf *DateTimeFlag) String() string {
	return fmt.Sprintf("%s", *dtf)
}
