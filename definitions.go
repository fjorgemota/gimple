/*
Package gimple is a Dependency Injection Container developed in Golang with features highly inspired on Pimple,
a micro dependency injection container for PHP.

Some of it's features is:

- Define services;
- Define factories;
- Define parameters easily;
- Allow services to depend directly on interfaces, and not on concrete struct;
- Defining services/parameters/factories from another files - because you should be able to split your configuration easily;
- Simple API;
- Allows extending services easily;
- Allow to get the raw service creator easily;
- Pure Go, no C code envolved;
- Fully tested on each commit;
- I already said that it have a really Simple API? :)
*/
package gimple

type Registerer interface {
	Register(container Container)
}
type RegisterFunc func(c Container)

func (fn RegisterFunc) Register(c Container) {
	fn(c)
}

type Container interface {
	Get(key string) (interface{}, error)
	MustGet(key string) interface{}
	Set(key string, val interface{})
	Has(key string) bool
	Keys() []string
	Factory(fn func(c Container) interface{}) func(c Container) interface{}
	Protect(fn func(c Container) interface{}) func(c Container) interface{}
	Extend(key string, e Extender) error
	MustExtend(key string, e Extender)
	ExtendFunc(key string, f ExtenderFunc) error
	MustExtendFunc(key string, f ExtenderFunc)
	Register(provider Registerer)
	RegisterFunc(fn RegisterFunc)
	Raw(key string) (interface{}, error)
	MustRaw(key string) interface{}
}

type Error struct {
	err string
}

func (self Error) Error() string {
	return self.err
}

type ExtenderFunc func(old interface{}, c Container) interface{}

func (f ExtenderFunc) Extend(old interface{}, c Container) interface{} {
	return f(old, c)
}

type Extender interface {
	Extend(old interface{}, c Container) interface{}
}

type Gimple struct {
	items     map[string]interface{}
	instances map[string]interface{}
	protected map[string]struct{}
	factories map[string]struct{}
}
