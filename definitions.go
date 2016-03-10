/**
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

type GimpleProvider interface {
	Register(container GimpleContainer)
}
type GimpleContainer interface {
	Get(key string) interface{}
	Set(key string, val interface{})
	Has(key string) bool
	Keys() []string
	Factory(fn func(c GimpleContainer) interface{}) func(c GimpleContainer) interface{}
	Protect(fn func(c GimpleContainer) interface{}) func(c GimpleContainer) interface{}
	Extend(key string, fn GimpleExtender) error
	Register(provider GimpleProvider)
	Raw(key string) interface{}
}

type GimpleError struct {
	err string
}

func (self GimpleError) Error() string {
	return self.err
}

type GimpleExtender func(old interface{}, c GimpleContainer) interface{}

type Gimple struct {
	items     map[string]interface{}
	instances map[string]interface{}
	protected map[string]struct{}
	factories map[string]struct{}
}
