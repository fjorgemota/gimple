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

func getServiceDefinitionName(item func(container GimpleContainer) interface{}) string {
	pointer := reflect.ValueOf(item).Pointer()
	fn := runtime.FuncForPC(pointer)
	return fn.Name()
}
