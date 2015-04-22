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
