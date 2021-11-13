package recv

import (
	"fmt"
)

// returned when a converter doesn't convert input successfully
type ConversionError struct {
	// The arg that the converter belongs to
	Arg *CommandArg
	// The input that failed to convert
	Input string
	// The position of the arg in the command arguments - indexed on 1
	ArgPosition int
}

func (e ConversionError) Error() string {
	return fmt.Sprintf("converter \"%s\" failed to convert input \"%s\"", e.Arg.Converter.Name, e.Input)
}

// returned when a required argument is missing
type MissingRequiredArgumentError struct {
	Arg         *CommandArg
	ArgPosition int
}

func (e MissingRequiredArgumentError) Error() string {
	return fmt.Sprintf("required argument \"%s\" at position %v is missing", e.Arg.Name, e.ArgPosition)
}

// returned when a command check fails 
type CheckError struct {
	// the check that failed
	Check *CommandCheck
}

func (e CheckError) Error() string {
	return fmt.Sprintf("command check \"%s\" failed", e.Check.Name)
}
