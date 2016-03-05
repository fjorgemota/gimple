package gimple_test

import (
	"math/rand"
	"reflect"
	"runtime"

	. "github.com/fjorgemota/gimple"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type provider struct {
	called bool
	gimple GimpleContainer
}

func (self *provider) Register(app GimpleContainer) {
	Expect(app).To(Equal(self.gimple))
	self.called = true
}

func Symbol(app GimpleContainer) interface{} {
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
			gimple := NewGimple()
			var g Gimple
			// var gc GimpleContainer
			Expect(gimple).To(BeAssignableToTypeOf(&g))
			// Expect(gimple).To(BeAssignableToTypeOf(&gc))
			Expect(gimple.Keys()).To(BeEmpty())
		})
		It("should support passing some parameters", func() {
			values := map[string]interface{}{"name": "xpto", "age": 19}
			gimple := NewGimpleWithValues(values)
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Get("name")).To(Equal("xpto"))
			Expect(gimple.Get("age")).To(Equal(19))
		})
		It("should support passing some services", func() {
			values := map[string]interface{}{
				"n2": func(c GimpleContainer) interface{} {
					val := c.Get("n").(int)
					return val + 1
				},
				"n": func(c GimpleContainer) interface{} {
					return 19
				}}
			gimple := NewGimpleWithValues(values)
			Expect(gimple.Keys()).To(ContainElement("n"))
			Expect(gimple.Keys()).To(ContainElement("n2"))
			n2 := gimple.Get("n2")
			Expect(n2).To(Equal(20))
			n := gimple.Get("n")
			Expect(n).To(Equal(19))
			RawN := gimple.Raw("n")
			var s func(container GimpleContainer) interface{}
			Expect(RawN).To(BeAssignableToTypeOf(s))
		})
		It("should support passing some services and parameters", func() {
			values := map[string]interface{}{
				"n2": func(app GimpleContainer) interface{} {
					val := app.Get("n").(int)
					return val + 1
				},
				"n": 19}
			gimple := NewGimpleWithValues(values)
			Expect(gimple.Keys()).To(ContainElement("n"))
			Expect(gimple.Keys()).To(ContainElement("n2"))
			n2 := gimple.Get("n2")
			n := gimple.Get("n")
			RawN := gimple.Raw("n")
			Expect(n2).To(Equal(20))
			Expect(n).To(Equal(19))
			Expect(RawN).To(Equal(19))
		})
	})
	Describe("#get()", func() {
		It("should throw an exception when getting non existent key", func() {
			gimple := NewGimple()
			Expect(func() {
				gimple.Get("non-existent-key")
			}).To(Panic())
		})
		Measure("should get parameters fast", func(b Benchmarker) {
			values := map[string]interface{}{"age": 19, "name": "xpto"}
			gimple := NewGimpleWithValues(values)
			GetInteger := b.Time("GetInteger", func() {
				Expect(gimple.Get("age")).To(Equal(19))
			})
			Expect(GetInteger.Seconds()).To(BeNumerically("<", 0.2), "Get() for integers shouldn't take too long.")
			GetString := b.Time("GetString", func() {
				Expect(gimple.Get("name")).To(Equal("xpto"))
			})
			Expect(GetString.Seconds()).To(BeNumerically("<", 0.2), "Get() for strings shouldn't take too long.")
		}, 1000)
		It("should support getting parameters", func() {
			values := map[string]interface{}{"age": 19, "name": "xpto"}
			gimple := NewGimpleWithValues(values)
			Expect(gimple.Get("age")).To(Equal(19))
			Expect(gimple.Get("name")).To(Equal("xpto"))
		})
		Measure("should get services fast", func(b Benchmarker) {
			values := map[string]interface{}{
				"age": func(app GimpleContainer) interface{} { return 19 }}
			gimple := NewGimpleWithValues(values)
			GetService := b.Time("GetService", func() {
				Expect(gimple.Get("age")).To(Equal(19))
			})
			Expect(GetService.Seconds()).To(BeNumerically("<", 0.2), "GetService() shouldn't take too long.")
		}, 1000)
		It("should support getting services", func() {
			values := map[string]interface{}{
				"age": func(app GimpleContainer) interface{} { return 19 }}
			gimple := NewGimpleWithValues(values)
			Expect(gimple.Get("age")).To(Equal(19))
		})
		It("should cache values of the services", func() {
			values := map[string]interface{}{
				"symbol": Symbol}
			gimple := NewGimpleWithValues(values)
			val := gimple.Get("symbol")
			val2 := gimple.Get("symbol")
			Expect(val).To(Equal(val2))
		})
		It("should not cache values of factories", func() {
			gimple := NewGimple()
			gimple.Set("symbol", gimple.Factory(Symbol))
			val := gimple.Get("symbol")
			val2 := gimple.Get("symbol")
			value := val.(int)
			value2 := val2.(int)
			Expect(value).To(Not(Equal(value2)))
		})
		It("should return raw values of protected closures", func() {
			gimple := NewGimple()
			gimple.Set("symbol", gimple.Protect(Symbol))
			val := gimple.Get("symbol")
			converted := val.(func(c GimpleContainer) interface{})
			Expect(isFunctionEqual(converted, Symbol)).To(BeTrue())
		})
	})
	Describe("#set()", func() {
		It("should support saving parameters", func() {
			gimple := NewGimple()
			gimple.Set("age", 19)
			gimple.Set("name", "xpto")
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.Get("age")).To(Equal(19))
			Expect(gimple.Get("name")).To(Equal("xpto"))
		})
		It("should support saving services", func() {
			gimple := NewGimple()
			gimple.Set("age", func(app GimpleContainer) interface{} { return 19 })
			gimple.Set("name", func(app GimpleContainer) interface{} { return "xpto" })
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.Has("age")).To(BeTrue())
			Expect(gimple.Has("name")).To(BeTrue())
			age := gimple.Get("age")
			Expect(age).To(Equal(19))
			name := gimple.Get("name")
			Expect(name).To(Equal("xpto"))
		})
	})
	Describe("#raw()", func() {
		It("should throw an exception when getting non existent key", func() {
			gimple := NewGimple()
			Expect(func() {
				gimple.Raw("non-existent-key")
			}).To(Panic())
		})
		It("should return raw parameters", func() {
			gimple := NewGimple()
			gimple.Set("age", 19)
			gimple.Set("name", "xpto")
			Expect(gimple.Keys()).To(ContainElement("age"))
			Expect(gimple.Keys()).To(ContainElement("name"))
			Expect(gimple.Raw("age")).To(Equal(19))
			Expect(gimple.Raw("name")).To(Equal("xpto"))
		})
		It("should return raw services", func() {
			gimple := NewGimple()
			a := func(app GimpleContainer) interface{} {
				return 19
			}
			gimple.Set("symbol", Symbol)
			gimple.Set("age", a)
			Expect(gimple.Keys()).To(ContainElement("symbol"))
			Expect(gimple.Keys()).To(ContainElement("age"))
			age := gimple.Get("age")
			ageRaw := gimple.Raw("age")
			ageFunc := ageRaw.(func(c GimpleContainer) interface{})
			Expect(age).To(Equal(19))
			Expect(isFunctionEqual(ageFunc, a)).To(BeTrue())
			Expect(ageFunc(nil)).To(Equal(19))
			val := gimple.Get("symbol")
			val2 := gimple.Get("symbol")
			raw := gimple.Raw("symbol")
			Expect(val).To(Equal(val2))
			Expect(isFunctionEqual(raw, Symbol)).To(BeTrue())
		})
	})
	Describe("#protect()", func() {
		It("should return raw services", func() {
			gimple := NewGimple()
			age := func(app GimpleContainer) interface{} { return 19 }
			gimple.Set("symbol", gimple.Protect(Symbol))
			gimple.Set("age", gimple.Protect(age))
			Expect(gimple.Keys()).To(ContainElement("symbol"))
			Expect(gimple.Keys()).To(ContainElement("age"))
			ageGetted := gimple.Get("age")
			ageFunc := ageGetted.(func(c GimpleContainer) interface{})
			Expect(isFunctionEqual(ageFunc, age)).To(BeTrue())
			Expect(ageFunc(nil)).To(Equal(19))
			sym := gimple.Get("symbol")
			sym2 := gimple.Get("symbol")
			Expect(isFunctionEqual(sym2, sym)).To(BeTrue())
			Expect(isFunctionEqual(sym, Symbol)).To(BeTrue())
		})
	})
	Describe("#keys()", func() {
		It("should return keys of parameters", func() {
			gimple := NewGimple()
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
			gimple := NewGimple()
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
			gimple := NewGimple()
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
			gimple := NewGimple()
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
			gimple := NewGimple()
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
			gimple := NewGimple()
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
			gimple := NewGimple()
			prov := &provider{called: false, gimple: gimple}
			gimple.Register(prov)
			Expect(prov.called).To(BeTrue())
		})
	})
	Describe("#extend()", func() {
		It("should throw an error on non-existent key", func() {
			gimple := NewGimple()
			err := gimple.Extend("not-found-key", func(val interface{}, container GimpleContainer) interface{} {
				return nil
			})
			Expect(err).To(Not(Succeed()))
		})
		It("should throw an error on parameter key", func() {
			gimple := NewGimple()
			gimple.Set("age", 19)
			err := gimple.Extend("age", func(val interface{}, container GimpleContainer) interface{} {
				return nil
			})
			Expect(err).To(Not(Succeed()))
		})
		It("should overwrite service correctly", func() {
			gimple := NewGimple()
			gimple.Set("age", func(c GimpleContainer) interface{} {
				return 19
			})
			gimple.Set("one", 1)
			age := gimple.Get("age")
			Expect(age).To(Equal(19))
			gimple.Extend("age", func(result interface{}, app GimpleContainer) interface{} {
				n := result.(int)
				one := app.Get("one").(int)
				return n + one
			})
			newAge := gimple.Get("age")
			Expect(newAge).To(Equal(20))
		})
	})

})
