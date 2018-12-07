package flags

import (
	"flag"
	"reflect"
)

func cloneFlagValue(flag flag.Value) flag.Value {
	if flag == nil {
		return nil
	}
	// for the rest of this function, variables ending in Value are
	// `reflect.Value`s and those ending in Type are `reflect.Type`s
	flagType := reflect.TypeOf(flag)
	flagValue := reflect.ValueOf(flag)
	isPtr := flagType.Kind() == reflect.Ptr

	// if flag is a pointer to something, get the underlying value & type
	// n.b. it should always be a pointer (reason below) but in case
	// someone comes along and makes a non-pointer that implements
	// PreprocesserFlag I'd prefer that clonePF doesn't cause a panic
	//
	// reason it should always be a pointer: flag.Value.Set and Preprocess are
	// both almost-useless when their receivers aren't pointers.
	if isPtr {
		flagType = flagType.Elem()
		flagValue = reflect.Indirect(flagValue)
	}

	// Make a new value
	newFlagValue := reflect.New(flagType).Elem()
	// and set its fields to the same as the old one
	newFlagValue.Set(flagValue)
	// ensure we have a pointer, if we used to have a pointer
	if isPtr {
		newFlagValue = newFlagValue.Addr()
	}

	// Get out of reflection-land and back to the real world,
	// then cast our interface{} to PreprocessorFlag
	clonedFlag := newFlagValue.Interface().(flag.Value)
	return clonedFlag
}
