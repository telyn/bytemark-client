package lib

import (
	"fmt"
)

type BigVError struct {
	ThingType string
	Thing     string
	User      string
	Action    string
}

type NotFoundError struct {
	BigVError
}

type NotAuthorizedError struct {
	BigVError
}

func (e BigVError) Error() string {
	return fmt.Sprintf("BigVError of type %T - Should have its own error message. Anyway, some details:\r\nThing: %s\r\nThingType: %s\r\nUser: %s\r\nAction: %s\r\n", e, e.Thing, e.ThingType, e.User, e.Action)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Couldn't find %s %s as user %s", e.ThingType, e.Thing, e.User)
}

func (e NotAuthorizedError) Error() string {
	return fmt.Sprintf("User %s is unauthorised to %s %s on %s %s", e.User, e.Action, e.ThingType, e.Thing)

}

func (bigv *BigVClient) PopulateError(err error, thing, thingType, action string) error {
	switch err.(type) {
	case NotFoundError:
		betterError := err.(NotFoundError)
		betterError.Thing = thing
		betterError.ThingType = thingType
		betterError.User = bigv.authSession.Username
		return betterError

	case NotAuthorizedError:
		betterError := err.(NotAuthorizedError)
		betterError.Thing = thing
		betterError.ThingType = thingType
		betterError.User = bigv.authSession.Username
		return betterError

	}
	return err
}
