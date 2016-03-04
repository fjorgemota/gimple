package gimple

import "fmt"

var _ = fmt.Println

func NewGimple() GimpleContainer {
	return NewGimpleWithValues(make(map[string]interface{}))
}
func NewGimpleWithValues(values map[string]interface{}) GimpleContainer {
	instances := make(map[string]interface{})
	protected := make(map[string]struct{}, 0)
	factories := make(map[string]struct{}, 0)
	return &Gimple{
		items:     values,
		instances: instances,
		protected: protected,
		factories: factories}
}

func (self *Gimple) isProtected(name string) bool {
	_, ok := self.protected[name]
	return ok
}

func (self *Gimple) isFactory(name string) bool {
	_, ok := self.factories[name]
	return ok
}
func (self *Gimple) Get(key string) interface{} {
	item, ok := self.items[key]
	if !ok {
		// We will panic here because, normally, without a key the user cannot proceed in a DI
		panic(newGimpleError("Identifier '" + key + "' is not defined."))
	}
	var obj interface{}
	if isServiceDefinition(item) {
		// We already checked if the item is a service definition, so ignore it here..
		itemFn := toServiceDefinition(item)
		name, _ := getServiceDefinitionName(itemFn)
		protected := self.isProtected(name)
		if protected {
			obj = item
		} else if instance, exists := self.instances[name]; exists {
			obj = instance
		} else {
			obj = itemFn(self)
			if !self.isFactory(name) {
				self.instances[name] = obj
			}
		}
	} else {
		obj = item
	}
	return obj
}

func (self *Gimple) Extend(key string, fn GimpleExtender) error {
	originalItem, exists := self.items[key]
	if !exists {
		return newGimpleError("Identifier '" + key + "' is not defined.")
	}
	if !isServiceDefinition(originalItem) {
		return newGimpleError("Identifier '" + key + "' does not contain an object definition")
	}
	callable := toServiceDefinition(originalItem)
	self.items[key] = func(container GimpleContainer) interface{} {
		return fn(callable(container), container)
	}
	return nil
}

func (self *Gimple) Factory(fn func(c GimpleContainer) interface{}) func(c GimpleContainer) interface{} {
	// We are already receiving a func(c GimpleContainer) interface{}, so just ignore "error" here..
	name, _ := getServiceDefinitionName(fn)
	self.factories[name] = struct{}{}
	return fn
}

func (self *Gimple) Protect(fn func(c GimpleContainer) interface{}) func(c GimpleContainer) interface{} {
	// We are already receiving a func(c GimpleContainer) interface{}, so just ignore "error" here..
	name, _ := getServiceDefinitionName(fn)
	self.protected[name] = struct{}{}
	return fn
}

func (self *Gimple) Has(key string) bool {
	_, ok := self.items[key]
	return ok
}

func (self *Gimple) Keys() []string {
	keys := make([]string, len(self.items))
	i := 0
	for key := range self.items {
		keys[i] = key
		i++
	}
	return keys
}

func (self *Gimple) Raw(key string) interface{} {
	item, exists := self.items[key]
	if !exists {
		panic(newGimpleError("Identifier '" + key + "' is not defined."))
	}
	return item
}

func (self *Gimple) Register(provider GimpleProvider) {
	provider.Register(self)
}

func (self *Gimple) Set(key string, val interface{}) {
	self.items[key] = val
}
