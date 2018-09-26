package internal

import "testing"

func TestToSnake(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{{
		in:  "hello_friend",
		out: "hello_friend",
	}, {
		in:  "iHaveNeverBeenHappier",
		out: "i_have_never_been_happier",
	}, {
		in:  "TheWorldIsBRILLIANT",
		out: "the_world_is_brilliant",
	}, {
		in:  "JOYIsForAll",
		out: "joy_is_for_all",
	}, {
		in:  "ThisISTheBestTIME",
		out: "this_is_the_best_time",
	}}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			out := ToSnake(test.in)
			if out != test.out {
				t.Errorf("expected %q, got %q", test.out, out)
			}
		})
	}
}
