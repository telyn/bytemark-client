package mocks

import (
	"time"

	mock "github.com/maraino/go-mock"
	"github.com/urfave/cli"
)

type CliContext struct {
	mock.Mock
}

func (c *CliContext) Args() cli.Args {
	r := c.Called()
	return r.Get(0).(cli.Args)
}
func (c *CliContext) Bool(name string) bool {
	r := c.Called(name)
	return r.Bool(0)
}
func (c *CliContext) BoolT(name string) bool {
	r := c.Called(name)
	return r.Bool(0)
}
func (c *CliContext) Duration(name string) time.Duration {
	r := c.Called(name)
	return r.Get(0).(time.Duration)
}
func (c *CliContext) FlagNames() (names []string) {
	r := c.Called()
	return r.Get(0).([]string)
}
func (c *CliContext) Float64(name string) float64 {
	r := c.Called(name)
	return r.Float64(0)
}
func (c *CliContext) Generic(name string) interface{} {
	r := c.Called(name)
	return r.Get(0)
}
func (c *CliContext) GlobalBool(name string) bool {
	r := c.Called(name)
	return r.Bool(0)
}
func (c *CliContext) GlobalBoolT(name string) bool {
	r := c.Called(name)
	return r.Bool(0)
}
func (c *CliContext) GlobalDuration(name string) time.Duration {
	r := c.Called(name)
	return r.Get(0).(time.Duration)
}
func (c *CliContext) GlobalFlagNames() (names []string) {
	r := c.Called()
	return r.Get(0).([]string)
}
func (c *CliContext) GlobalFloat64(name string) float64 {
	r := c.Called(name)
	return r.Float64(0)
}
func (c *CliContext) GlobalGeneric(name string) interface{} {
	r := c.Called(name)
	return r.Get(0)
}
func (c *CliContext) GlobalInt(name string) int {
	r := c.Called(name)
	return r.Int(0)
}
func (c *CliContext) GlobalInt64(name string) int64 {
	r := c.Called(name)
	return r.Int64(0)
}
func (c *CliContext) GlobalInt64Slice(name string) []int64 {
	r := c.Called(name)
	return r.Get(0).([]int64)
}
func (c *CliContext) GlobalIntSlice(name string) []int {
	r := c.Called(name)
	return r.Get(0).([]int)
}
func (c *CliContext) GlobalIsSet(name string) bool {
	r := c.Called(name)
	return r.Bool(0)
}
func (c *CliContext) GlobalSet(name, value string) error {
	r := c.Called(name, value)
	return r.Error(0)
}
func (c *CliContext) GlobalString(name string) string {
	r := c.Called(name)
	return r.String(0)
}
func (c *CliContext) GlobalStringSlice(name string) []string {
	r := c.Called(name)
	return r.Get(0).([]string)
}
func (c *CliContext) GlobalUint(name string) uint {
	r := c.Called(name)
	return r.Get(0).(uint)
}
func (c *CliContext) GlobalUint64(name string) uint64 {
	r := c.Called(name)
	return r.Get(0).(uint64)
}
func (c *CliContext) Int(name string) int {
	r := c.Called(name)
	return r.Int(0)
}
func (c *CliContext) Int64(name string) int64 {
	r := c.Called(name)
	return r.Int64(0)
}
func (c *CliContext) Int64Slice(name string) []int64 {
	r := c.Called(name)
	return r.Get(0).([]int64)
}
func (c *CliContext) IntSlice(name string) []int {
	r := c.Called(name)
	return r.Get(0).([]int)
}
func (c *CliContext) IsSet(name string) bool {
	r := c.Called(name)
	return r.Bool(0)
}
func (c *CliContext) NArg() int {
	r := c.Called()
	return r.Int(0)
}
func (c *CliContext) NumFlags() int {
	r := c.Called()
	return r.Int(0)
}
func (c *CliContext) Parent() *cli.Context {
	r := c.Called()
	return r.Get(0).(*cli.Context)
}
func (c *CliContext) Set(name, value string) error {
	r := c.Called(name, value)
	return r.Error(0)
}
func (c *CliContext) String(name string) string {
	r := c.Called(name)
	return r.String(0)
}
func (c *CliContext) StringSlice(name string) []string {
	r := c.Called(name)
	return r.Get(0).([]string)
}
func (c *CliContext) Uint(name string) uint {
	r := c.Called(name)
	return r.Get(0).(uint)
}
func (c *CliContext) Uint64(name string) uint64 {
	r := c.Called(name)
	return r.Get(0).(uint64)
}

func (c *CliContext) App() *cli.App {
	r := c.Called()
	return r.Get(0).(*cli.App)
}
func (c *CliContext) Command() cli.Command {
	r := c.Called()
	return r.Get(0).(cli.Command)
}
