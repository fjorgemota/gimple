package gimple

import (
	"fmt"
	"reflect"
	"runtime"
)

func newGimpleError(err string, x ...interface{}) Error {
	return Error{fmt.Sprintf(err, x...)}
}

func isServiceDefinition(item interface{}) bool {
	_, ok := item.(func(c Container) interface{})
	return ok
}

func toServiceDefinition(item interface{}) func(container Container) interface{} {
	return item.(func(container Container) interface{})
}

func getServiceDefinitionName(item func(container Container) interface{}) string {
	pointer := reflect.ValueOf(item).Pointer()
	fn := runtime.FuncForPC(pointer)
	return fn.Name()
}
