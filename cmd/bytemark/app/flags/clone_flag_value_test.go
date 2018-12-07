package flags

import (
	"fmt"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/app"
)

// exampleByVal is a type that implements PreprocessorFlag without using pointer
// receivers.
type exampleByVal struct {
	data string
}

func (ebv exampleByVal) Set(val string) error {
	return nil
}

func (ebv exampleByVal) Preprocess(ctx *app.Context) error {
	return nil
}

func (ebv exampleByVal) String() string {
	return fmt.Sprintf("exampleByVal %s", ebv.data)
}

// exampleByRef is much simpler since it implements PreprocessorFlag on pointer
// receivers.
type exampleByRef struct {
	data string
}

func (ebr *exampleByRef) Set(val string) error {
	ebr.data = val
	return nil
}

func (ebr *exampleByRef) Preprocess(ctx *app.Context) error {
	ebr.data = "PREPROCESSED"
	return nil
}

// String is readonly so doesn't need to be on a pointer receiver.
func (ebr exampleByRef) String() string {
	return fmt.Sprintf("exampleByRef %s", ebr.data)
}

func TestClonePF(t *testing.T) {
	t.Run("by value", func(t *testing.T) {
		ebv := exampleByVal{
			data: "i am become test",
		}
		ebvClone, ok := clonePF(ebv).(exampleByVal)
		if !ok {
			t.Fatal("ebvClone wasn't an exampleByVal")
		}

		if ebvClone.data != "i am become test" {
			t.Error("original's data was not copied to clone")
		}

		// change the clone's data and check it doesn't affect the original's
		ebvClone.data = "destroyer of worlds"
		if ebv.data == "destroyer of worlds" {
			t.Error("changing the clone's data changed the original's")
		}

	})

	t.Run("by reference", func(t *testing.T) {
		ebr := exampleByRef{
			data: "i am become test",
		}
		ebrClone, ok := clonePF(&ebr).(*exampleByRef)
		if !ok {
			t.Fatal("ebrClone wasn't an exampleByRef")
		}
		// change the clone's data and check it doesn't affect the original's
		ebrClone.data = "destroyer of worlds"
		if ebr.data == "destroyer of worlds" {
			t.Error("changing the clone's data changed the original's")
		}

	})

}
