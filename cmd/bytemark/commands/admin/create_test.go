package admin_test

import (
	"strings"
	"testing"

	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/commands/admin"
	"github.com/BytemarkHosting/bytemark-client/cmd/bytemark/testutil"
	"github.com/BytemarkHosting/bytemark-client/lib"
	"github.com/BytemarkHosting/bytemark-client/lib/brain"
	"github.com/BytemarkHosting/bytemark-client/lib/util"
	"github.com/BytemarkHosting/bytemark-client/mocks"
)

func TestCreateMigration(t *testing.T) {
	tests := []struct {
		name      string
		args      string
		exp       brain.MigrationJob
		shouldErr bool
	}{
		{
			name: "OneDisc",
			args: "--disc 1",
			exp: brain.MigrationJob{
				Args: brain.MigrationJobSpec{
					Sources: brain.MigrationJobLocations{
						Discs: []util.NumberOrString{"1"},
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
					Destinations: brain.MigrationJobLocations{
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
				},
				Queue: brain.MigrationJobQueue{
					Discs: []int{1},
				},
			},
		},
		{
			name: "OnePool",
			args: "--pool 1",
			exp: brain.MigrationJob{
				Args: brain.MigrationJobSpec{
					Sources: brain.MigrationJobLocations{
						Discs: []util.NumberOrString{},
						Pools: []util.NumberOrString{"1"},
						Tails: []util.NumberOrString{},
					},
					Destinations: brain.MigrationJobLocations{
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
				},
				Queue: brain.MigrationJobQueue{
					Discs: []int{1},
				},
			},
		},
		{
			name: "OneTail",
			args: "--tail 1",
			exp: brain.MigrationJob{
				Args: brain.MigrationJobSpec{
					Sources: brain.MigrationJobLocations{
						Discs: []util.NumberOrString{},
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{"1"},
					},
					Destinations: brain.MigrationJobLocations{
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
				},
				Queue: brain.MigrationJobQueue{
					Discs: []int{1},
				},
			},
		},
		{
			name: "TwoDiscsToOnePool",
			args: "--disc 1 --disc 2 --to-pool 1",
			exp: brain.MigrationJob{
				Args: brain.MigrationJobSpec{
					Sources: brain.MigrationJobLocations{
						Discs: []util.NumberOrString{"1", "2"},
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{},
					},
					Destinations: brain.MigrationJobLocations{
						Pools: []util.NumberOrString{"1"},
						Tails: []util.NumberOrString{},
					},
				},
				Queue: brain.MigrationJobQueue{
					Discs: []int{1, 2},
				},
			},
		},
		{
			name: "OneDiscAndTwoPoolsToOneTailHighPriority",
			args: "--priority 100 --disc 1 --pool 1 --pool 2 --to-tail 1",
			exp: brain.MigrationJob{
				Args: brain.MigrationJobSpec{
					Sources: brain.MigrationJobLocations{
						Discs: []util.NumberOrString{"1"},
						Pools: []util.NumberOrString{"1", "2"},
						Tails: []util.NumberOrString{},
					},
					Destinations: brain.MigrationJobLocations{
						Pools: []util.NumberOrString{},
						Tails: []util.NumberOrString{"1"},
					},
					Options: brain.MigrationJobOptions{
						Priority: 100,
					},
				},
				Queue: brain.MigrationJobQueue{
					Discs: []int{1, 2, 3},
				},
				Priority: 100,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(y *testing.T) {
			_, client, app := testutil.BaseTestAuthSetup(t, true, admin.Commands)
			postReq := &mocks.Request{
				T:              t,
				StatusCode:     201,
				ResponseObject: test.exp,
			}
			client.When("BuildRequest", "POST", lib.BrainEndpoint, "/admin/migration_jobs%s", []string{""}).Return(postReq).Times(1)
			args := strings.Split("bytemark --admin create migration "+test.args, " ")
			err := app.Run(args)
			if !test.shouldErr && err != nil {
				t.Errorf("shouldn't err, but did: %T{%s}", err, err.Error())
			} else if test.shouldErr && err == nil {
				t.Fatal("should err, but didn't")
				return
			}
			if ok, err := client.Verify(); !ok {
				t.Fatal(err)
			}
			postReq.AssertRequestObjectEqual(test.exp.Args)

		})
	}
}
