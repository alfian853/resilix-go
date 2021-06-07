package util

import "fmt"

type UnhandledError struct {
	error
	Message interface{}
}

func (e *UnhandledError) Error() string {

	preFormat := "ResilixExecutor encountered an unhandled error"

	switch e.Message.(type) {
	case error:
		return fmt.Sprintf(preFormat+": %s\n", e.Message.(error).Error())
	case string:
		return fmt.Sprintf(preFormat+": %s\n", e.Message.(string))
	}

	var args []interface{}
	canBeString := false
	if _, ok := e.Message.(fmt.Stringer); ok {
		canBeString = true
		preFormat += ", String(): %s"
		args = append(args, e.Message.(fmt.Stringer).String())
	}

	if _, ok := e.Message.(fmt.GoStringer); ok {
		canBeString = true
		preFormat += ", GoString(): %s"
		args = append(args, e.Message.(fmt.Stringer).String())
	}

	if canBeString {
		return fmt.Sprintf(preFormat+"\n", args)
	}

	return fmt.Sprintf(preFormat+", %%#v: %#v\n", e.Message)
}
