package layer

import (
	"errors"
)

type Layer string

const (
	L0 Layer = "l0"
	L1 Layer = "l1"
)

func ParseString(str string) (lay Layer, err error) {
	var l Layer

	if str == "l0" {
		l = L0
		return l, nil
	}
	if str == "l1" {
		l = L1
		return l, nil
	}

	return l, errors.New("Unknown layer")
}