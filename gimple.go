package gimple

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

func New() GimpleContainer {
	return NewGimple()
}

func NewWithValues(values map[string]interface{}) GimpleContainer {
	return NewGimpleWithValues(values)
}

func (container *Gimple) isProtected(name string) bool {
	_, ok := container.protected[name]
	return ok
}

func (container *Gimple) isFactory(name string) bool {
	_, ok := container.factories[name]
	return ok
}
func (container *Gimple) Get(key string) interface{} {
	item, ok := container.items[key]
	if !ok {
		// We will panic here because, normally, without a key the user cannot proceed in a DI
		panic(newGimpleError("Identifier '" + key + "' is not defined."))
	}
	var obj interface{}
	if isServiceDefinition(item) {
		// We already checked if the item is a service definition, so ignore it here..
		itemFn := toServiceDefinition(item)
		name := getServiceDefinitionName(itemFn)
		protected := container.isProtected(name)
		if protected {
			obj = item
		} else if instance, exists := container.instances[name]; exists {
			obj = instance
		} else {
			obj = itemFn(container)
			if !container.isFactory(name) {
				container.instances[name] = obj
			}
		}
	} else {
		obj = item
	}
	return obj
}

func (container *Gimple) Extend(key string, fn GimpleExtender) error {
	originalItem, exists := container.items[key]
	if !exists {
		return newGimpleError("Identifier '" + key + "' is not defined.")
	}
	if !isServiceDefinition(originalItem) {
		return newGimpleError("Identifier '" + key + "' does not contain an object definition")
	}
	callable := toServiceDefinition(originalItem)
	container.items[key] = func(container GimpleContainer) interface{} {
		return fn(callable(container), container)
	}
	return nil
}

func (container *Gimple) Factory(fn func(c GimpleContainer) interface{}) func(c GimpleContainer) interface{} {
	// We are already receiving a func(c GimpleContainer) interface{}, so just ignore "error" here..
	name := getServiceDefinitionName(fn)
	container.factories[name] = struct{}{}
	return fn
}

func (container *Gimple) Protect(fn func(c GimpleContainer) interface{}) func(c GimpleContainer) interface{} {
	// We are already receiving a func(c GimpleContainer) interface{}, so just ignore "error" here..
	name := getServiceDefinitionName(fn)
	container.protected[name] = struct{}{}
	return fn
}

func (container *Gimple) Has(key string) bool {
	_, ok := container.items[key]
	return ok
}

func (container *Gimple) Keys() []string {
	keys := make([]string, len(container.items))
	i := 0
	for key := range container.items {
		keys[i] = key
		i++
	}
	return keys
}

func (container *Gimple) Raw(key string) interface{} {
	item, exists := container.items[key]
	if !exists {
		panic(newGimpleError("Identifier '" + key + "' is not defined."))
	}
	return item
}

func (container *Gimple) Register(provider GimpleProvider) {
	provider.Register(container)
}

func (container *Gimple) Set(key string, val interface{}) {
	container.items[key] = val
}
