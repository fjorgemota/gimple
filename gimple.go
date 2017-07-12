package gimple

import "errors"

func NewGimple() Container {
	return NewGimpleWithValues(make(map[string]interface{}))
}
func NewGimpleWithValues(values map[string]interface{}) Container {
	instances := make(map[string]interface{})
	protected := make(map[string]struct{}, 0)
	factories := make(map[string]struct{}, 0)
	return &Gimple{
		items:     values,
		instances: instances,
		protected: protected,
		factories: factories,
	}
}

func New(values ...map[string]interface{}) Container {
	v := make(map[string]interface{})
	if len(values) > 0 {
		for _, val := range values {
			for key, x := range val {
				v[key] = x
			}
		}
	}
	return NewGimpleWithValues(v)
}

func NewWithValues(values map[string]interface{}) Container {
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
func (container *Gimple) MustGet(key string) interface{} {
	if x, err := container.Get(key); err != nil {
		panic(err)
	} else {
		return x
	}
}
func (container *Gimple) Get(key string) (interface{}, error) {
	item, ok := container.items[key]
	if !ok {
		return nil, notDefined(key)
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
	return obj, nil
}

func (container *Gimple) ExtendFunc(key string, f ExtenderFunc) error {
	return container.Extend(key, f)
}

func (container *Gimple) MustExtendFunc(key string, f ExtenderFunc) {
	container.MustExtend(key, f)
}

func (container *Gimple) MustExtend(key string, e Extender) {
	if err := container.Extend(key, e); err != nil {
		panic(err)
	}
}
func (container *Gimple) Extend(key string, e Extender) error {
	originalItem, exists := container.items[key]
	if !exists {
		return notDefined(key)
	}
	if !isServiceDefinition(originalItem) {
		return newGimpleError("Identifier '%s' does not contain an object definition", key)
	}
	callable := toServiceDefinition(originalItem)
	container.items[key] = func(container Container) interface{} {
		return e.Extend(callable(container), container)
	}
	return nil
}

func (container *Gimple) Factory(fn func(c Container) interface{}) func(c Container) interface{} {
	// We are already receiving a func(c GimpleContainer) interface{}, so just ignore "error" here..
	name := getServiceDefinitionName(fn)
	container.factories[name] = struct{}{}
	return fn
}

func (container *Gimple) Protect(fn func(c Container) interface{}) func(c Container) interface{} {
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

var (
	Undefined = errors.New("Not defined")
)

func (container *Gimple) MustRaw(key string) interface{} {
	if x, err := container.Raw(key); err != nil {
		panic(err)
	} else {
		return x
	}
}
func notDefined(key string) error {
	return newGimpleError("Identifier '%s' is not defined.", key)
}
func (container *Gimple) Raw(key string) (interface{}, error) {
	if item, exists := container.items[key]; exists {
		return item, nil
	}
	return nil, notDefined(key)
}

func (container *Gimple) Register(provider Registerer) {
	provider.Register(container)
}
func (container *Gimple) RegisterFunc(fn RegisterFunc) {
	container.Register(fn)
}

func (container *Gimple) Set(key string, val interface{}) {
	container.items[key] = val
}
