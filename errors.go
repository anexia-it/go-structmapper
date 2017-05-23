package structmapper

import "errors"

var (
	// ErrTagNameEmpty designates that the passed tag name is empty
	ErrTagNameEmpty = errors.New("Tag name is empty")

	// ErrNotAStruct designates that the passed value is not a struct
	ErrNotAStruct = errors.New("Not a struct")

	// ErrInvalidMap designates that the passed value is not a valid map
	ErrInvalidMap = errors.New("Invalid map")

	// ErrFieldIsInterface designates that a field is an interface
	ErrFieldIsInterface = errors.New("Field is interface")

	// ErrMapIsNil designates that the passed map is nil
	ErrMapIsNil = errors.New("Map is nil")

	// ErrNotAStructPointer designates that the passed value is not a pointer to a struct
	ErrNotAStructPointer = errors.New("Not a struct pointer")
)
