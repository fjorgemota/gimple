package gimple_test

import (
	"math/rand"
	"reflect"
	"runtime"

	. "github.com/alxarch/gimple"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type provider struct {
	called bool
	gimple Container
}

func (prov *provider) Register(app Container) {
	Expect(app).To(Equal(prov.gimple))
	prov.called = true
}

func Symbol(app Container) interface{} {
	return rand.Int()
}

func isFunctionEqual(fn1, fn2 interface{}) bool {
	pointer1 := reflect.ValueOf(fn1).Pointer()
	name1 := runtime.FuncForPC(pointer1).Name()
	pointer2 := reflect.ValueOf(fn2).Pointer()
	name2 := runtime.FuncForPC(pointer2).Name()
	return name1 == name2
}

var _ = Describe("Gimple", func() {
	Describe("#constructor()", func() {
		It("should support passing no parameters", func() {
			gimple := New()
			var g Gimple
			Expect(gimple).To(BeAssignableToTypeOf(&g))
			Expect(gimple.Keys()).To(BeEmpty())
		})
		It("should support passing some parameters", func() {
			values := map[string]interface{}{"name": "xpto", "age": 19}
			gimple := New(values)
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Has("name")).To(BeTrue())
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("eita")).To(BeFalse())
			Expect(gimple.Has("non-existent-key")).To(BeFalse())
			Expect(gimple.MustGet("name")).To(Equal("xpto"))
			Expect(gimple.MustGet("age")).To(Equal(19))
		})
		It("should support passing some services", func() {
			values := map[string]interface{}{
				"n2": func(c Container) interface{} {
					val := c.MustGet("n").(int)
					return val + 1
				},
				"n": func(c Container) interface{} {
					return 19
				}}
			gimple := New(values)
			Expect(gimple.Keys()).To(ContainElement("n"))
			Expect(gimple.Keys()).To(ContainElement("n2"))
			n2 := gimple.MustGet("n2")
			Expect(n2).To(Equal(20))
			n := gimple.MustGet("n")
			Expect(n).To(Equal(19))
			RawN := gimple.MustRaw("n")
			var s func(container Container) interface{}
			Expect(RawN).To(BeAssignableToTypeOf(s))
		})
		It("should support passing some services and parameters", func() {
			values := map[string]interface{}{
				"n2": func(app Container) interface{} {
					val := app.MustGet("n").(int)
					return val + 1
				},
				"n": 19}
			gimple := NewWithValues(values)
			Expect(gimple.Keys()).To(ContainElement("n"))
			Expect(gimple.Keys()).To(ContainElement("n2"))
			n2 := gimple.MustGet("n2")
			n := gimple.MustGet("n")
			RawN := gimple.MustRaw("n")
			Expect(n2).To(Equal(20))
			Expect(n).To(Equal(19))
			Expect(RawN).To(Equal(19))
		})
	})
	Describe("#get()", func() {
		It("should throw an exception when getting non existent key", func() {
			gimple := New()
			Expect(func() {
				gimple.MustGet("non-existent-key")
			}).To(Panic())
			err := make(chan error, 0)
			go func() {
				defer func() {
					err <- recover().(error)
				}()
				gimple.MustGet("non-existent-key")
			}()
			Expect(<-err).To(MatchError(Equal("Identifier 'non-existent-key' is not defined.")))
		})
		Measure("should get parameters fast", func(b Benchmarker) {
			values := map[string]interface{}{"age": 19, "name": "xpto"}
			gimple := NewWithValues(values)
			GetInteger := b.Time("GetInteger", func() {
				Expect(gimple.MustGet("age")).To(Equal(19))
			})
			Expect(GetInteger.Seconds()).To(BeNumerically("<", 0.2), "Get() for integers shouldn't take too long.")
			GetString := b.Time("GetString", func() {
				Expect(gimple.MustGet("name")).To(Equal("xpto"))
			})
			Expect(GetString.Seconds()).To(BeNumerically("<", 0.2), "Get() for strings shouldn't take too long.")
		}, 1000)
		It("should support getting parameters", func() {
			values := map[string]interface{}{"age": 19, "name": "xpto"}
			gimple := NewWithValues(values)
			Expect(gimple.MustGet("age")).To(Equal(19))
			Expect(gimple.MustGet("name")).To(Equal("xpto"))
		})
		Measure("should get services fast", func(b Benchmarker) {
			values := map[string]interface{}{
				"age": func(app Container) interface{} { return 19 }}
			gimple := NewWithValues(values)
			GetService := b.Time("GetService", func() {
				Expect(gimple.MustGet("age")).To(Equal(19))
			})
			Expect(GetService.Seconds()).To(BeNumerically("<", 0.2), "GetService() shouldn't take too long.")
		}, 1000)
		It("should support getting services", func() {
			values := map[string]interface{}{
				"age": func(app Container) interface{} { return 19 }}
			gimple := NewWithValues(values)
			Expect(gimple.MustGet("age")).To(Equal(19))
		})
		It("should cache values of the services", func() {
			values := map[string]interface{}{
				"symbol": Symbol}
			gimple := NewWithValues(values)
			val := gimple.MustGet("symbol")
			val2 := gimple.MustGet("symbol")
			Expect(val).To(Equal(val2))
		})
		It("should not cache values of factories", func() {
			gimple := New()
			gimple.Set("symbol", gimple.Factory(Symbol))
			val := gimple.MustGet("symbol")
			val2 := gimple.MustGet("symbol")
			value := val.(int)
			value2 := val2.(int)
			Expect(value).To(Not(Equal(value2)))
		})
		It("should return raw values of protected closures", func() {
			gimple := New()
			gimple.Set("symbol", gimple.Protect(Symbol))
			val := gimple.MustGet("symbol")
			converted := val.(func(c Container) interface{})
			Expect(isFunctionEqual(converted, Symbol)).To(BeTrue())
		})
	})
	Describe("#set()", func() {
		It("should support saving parameters", func() {
			gimple := New()
			gimple.Set("age", 19)
			gimple.Set("name", "xpto")
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.MustGet("age")).To(Equal(19))
			Expect(gimple.MustGet("name")).To(Equal("xpto"))
		})
		It("should support saving services", func() {
			gimple := New()
			gimple.Set("age", func(app Container) interface{} { return 19 })
			gimple.Set("name", func(app Container) interface{} { return "xpto" })
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeTrue())
			age := gimple.MustGet("age")
			Expect(age).To(Equal(19))
			name := gimple.MustGet("name")
			Expect(name).To(Equal("xpto"))
		})
	})
	Describe("#raw()", func() {
		It("should throw an exception when getting non existent key", func() {
			gimple := New()
			Expect(func() {
				gimple.MustRaw("non-existent-key")
			}).To(Panic())
		})
		It("should return raw parameters", func() {
			gimple := New()
			gimple.Set("age", 19)
			gimple.Set("name", "xpto")
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.MustRaw("age")).To(Equal(19))
			Expect(gimple.MustRaw("name")).To(Equal("xpto"))
		})
		It("should return raw services", func() {
			gimple := New()
			a := func(app Container) interface{} {
				return 19
			}
			gimple.Set("symbol", Symbol)
			gimple.Set("age", a)
			Expect(gimple.Keys()).To(ContainElement("symbol"))
			Expect(gimple.Keys()).To(ContainElement("age"))
			age := gimple.MustGet("age")
			ageRaw := gimple.MustRaw("age")
			ageFunc := ageRaw.(func(c Container) interface{})
			Expect(age).To(Equal(19))
			Expect(isFunctionEqual(ageFunc, a)).To(BeTrue())
			Expect(ageFunc(nil)).To(Equal(19))
			val := gimple.MustGet("symbol")
			val2 := gimple.MustGet("symbol")
			raw := gimple.MustRaw("symbol")
			Expect(val).To(Equal(val2))
			Expect(isFunctionEqual(raw, Symbol)).To(BeTrue())
		})
	})
	Describe("#protect()", func() {
		It("should return raw services", func() {
			gimple := New()
			age := func(app Container) interface{} { return 19 }
			gimple.Set("symbol", gimple.Protect(Symbol))
			gimple.Set("age", gimple.Protect(age))
			Expect(gimple.Keys()).To(ContainElement("symbol"))
			Expect(gimple.Keys()).To(ContainElement("age"))
			ageGetted := gimple.MustGet("age")
			ageFunc := ageGetted.(func(c Container) interface{})
			Expect(isFunctionEqual(ageFunc, age)).To(BeTrue())
			Expect(ageFunc(nil)).To(Equal(19))
			sym := gimple.MustGet("symbol")
			sym2 := gimple.MustGet("symbol")
			Expect(isFunctionEqual(sym2, sym)).To(BeTrue())
			Expect(isFunctionEqual(sym, Symbol)).To(BeTrue())
		})
	})
	Describe("#keys()", func() {
		It("should return keys of parameters", func() {
			gimple := New()
			Expect(gimple.Keys()).To(BeEmpty())
			gimple.Set("age", 19)
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(HaveLen(1))
			gimple.Set("name", "xpto")
			Expect(gimple.Keys()).To(HaveLen(2))
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
		})
		It("should return keys of services", func() {
			gimple := New()
			Expect(gimple.Keys()).To(BeEmpty())
			gimple.Set("age", func() int { return 19 })
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(HaveLen(1))
			gimple.Set("name", func() string { return "xpto" })
			Expect(gimple.Keys()).To(HaveLen(2))
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
		})
		It("should return keys of services and parameters", func() {
			gimple := New()
			Expect(gimple.Keys()).To(BeEmpty())
			gimple.Set("age", 19)
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(HaveLen(1))
			gimple.Set("name", func() string { return "xpto" })
			Expect(gimple.Keys()).To(HaveLen(2))
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
		})
	})
	Describe("#has()", func() {
		It("should recognize parameters", func() {
			gimple := New()
			Expect(gimple.Has("age")).To(BeFalse())
			Expect(gimple.Has("name")).To(BeFalse())
			gimple.Set("age", 19)
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeFalse())
			gimple.Set("name", "xpto")
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeTrue())
		})
		It("should recognize services", func() {
			gimple := New()
			Expect(gimple.Has("age")).To(BeFalse())
			Expect(gimple.Has("name")).To(BeFalse())
			gimple.Set("age", func() int { return 19 })
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeFalse())
			gimple.Set("name", func() string { return "xpto" })
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeTrue())
		})
		It("should return keys of services and parameters", func() {
			gimple := New()
			Expect(gimple.Has("age")).To(BeFalse())
			Expect(gimple.Has("name")).To(BeFalse())
			gimple.Set("age", 19)
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeFalse())
			gimple.Set("name", func() string { return "xpto" })
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeTrue())
		})
	})
	Describe("#register()", func() {
		It("should call register() method on object", func() {
			gimple := New()
			prov := &provider{called: false, gimple: gimple}
			gimple.Register(prov)
			Expect(prov.called).To(BeTrue())
		})
	})
	Describe("#ExtendFunc()", func() {
		It("should throw an error on non-existent key", func() {
			gimple := New()
			err := gimple.Extend("not-found-key", ExtenderFunc(func(val interface{}, container Container) interface{} {
				return nil
			}))
			Expect(err).To(Not(Succeed()))
		})
		It("should throw an error on parameter key", func() {
			gimple := New()
			gimple.Set("age", 19)
			err := gimple.ExtendFunc("age", func(val interface{}, container Container) interface{} {
				return nil
			})
			Expect(err).To(Not(Succeed()))
		})
		It("should overwrite service correctly", func() {
			gimple := New()
			gimple.Set("age", func(c Container) interface{} {
				return 19
			})
			gimple.Set("one", 1)
			age := gimple.MustGet("age")
			Expect(age).To(Equal(19))
			gimple.ExtendFunc("age", func(result interface{}, app Container) interface{} {
				n := result.(int)
				one := app.MustGet("one").(int)
				return n + one
			})
			newAge := gimple.MustGet("age")
			Expect(newAge).To(Equal(20))
		})
	})

})
