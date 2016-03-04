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
