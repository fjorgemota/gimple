package gimple

import (
	"reflect"
	"runtime"
)

func newGimpleError(err string) GimpleError {
	return GimpleError{err}
}

func isServiceDefinition(item interface{}) bool {
	_, ok := item.(func(c GimpleContainer) interface{})
	return ok
}

func toServiceDefinition(item interface{}) func(container GimpleContainer) interface{} {
	return item.(func(container GimpleContainer) interface{})
}

func getServiceDefinitionName(item interface{}) (string, error) {
	if !isServiceDefinition(item) {
		return "", newGimpleError("Invalid service definition")
	}
	pointer := reflect.ValueOf(item).Pointer()
	fn := runtime.FuncForPC(pointer)
	return fn.Name(), nil
}
