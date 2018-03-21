package assert

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Offset within the call stack where the caller is
const callerOffset = 2

// AssertionPanic represents the details of a panic which is triggered by a failed runtime assertion
type AssertionPanic struct {
	Message string
}

func (assertionPanic AssertionPanic) String() string {
	return assertionPanic.Message
}

// AssertNotNil checks that the argument is not nil, and if it is nil, will panic with an AssertionPanic
// argName is required to properly describe the argument being checked in the assertion
func AssertNotNil(arg interface{}, argName string) {
	// NOTE: In order for the caller stack level to remain the same (and thus use the same callerOffset value), this
	// assertion must contain the same logic as the other assertion function. Previously this function called the
	// other assertion function, but that would have a different call stack, and thus result in incorrect reporting
	// on where the assertion failure occurred in the caller.
	if arg == nil {
		panic(AssertionPanic{
			Message: getFormattedMessage("%s may not be nil", argName),
		})
	}
}

// Assert checks a condition evaluates to true, otherwise will panic with an AssertionPanic
// A message describing the required condition must be provided, along with any values to be formatted into the message
func Assert(assertion bool, msg string, args ...interface{}) {
	if !assertion {
		panic(AssertionPanic{
			Message: getFormattedMessage(msg, args...),
		})
	}
}

func getFormattedMessage(msg string, args ...interface{}) string {
	message := msg

	if args != nil {
		message = fmt.Sprintf(msg, args)
	}

	_, filename, line, ok := runtime.Caller(callerOffset)

	if ok {
		dir, file := filepath.Split(filename)
		parent := filepath.Base(dir)
		return fmt.Sprintf("%s. %v/%v:%v", message, parent, file, line)
	}
	return message
}
