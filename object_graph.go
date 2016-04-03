package dagger

import (
	"fmt"
	"reflect"
)

// ObjectGraph is a graph of objects linked by their dependencies.
type ObjectGraph interface {
	// Inject the fields of `target`.
	Inject(target interface{})
	// Get retrieves an object from the stores the result in the value pointed to by v.
	Get(v interface{})
}

// NewObjectGraph returns a new dependency graph using the given modules.
func NewObjectGraph(modules ...Module) ObjectGraph {
	var graph map[reflect.Type]provider
	for _, module := range modules {
		moduleType := reflect.TypeOf(module)
		moduleValue := reflect.ValueOf(module)
		for i := 0; i < moduleType.NumMethod(); i++ {
			method := moduleType.Method(i)
			if !isMethodExported(method) {
				// todo: log
				continue
			}
			if !isProviderMethod(method) {
				// todo: log
				continue
			}

			if method.Type.NumOut() != 1 {
				panic(fmt.Errorf("Provider methods can return only one argument. %s returns %d values.", method.Type, method.Type.NumOut()))
			}
			binding := method.Type.Out(0)

			var dependencies []reflect.Type
			for j := 1; j < method.Type.NumIn(); j++ {
				dependencies = append(dependencies, method.Type.In(j))
			}

			var p provider
			makeProvider(&p, moduleValue, method, nil)
			graph[binding] = p

			fmt.Println(dependencies, binding)
		}
	}
	return &objectGraph{graph}
}

func makeProvider(fptr interface{}, moduleValue reflect.Value, providerMethod reflect.Method, dependencyProviders ...provider) {
	fn := reflect.ValueOf(fptr).Elem()
	provide := func(in []reflect.Value) []reflect.Value {
		if len(dependencyProviders) == 0 {
			return providerMethod.Func.Call([]reflect.Value{moduleValue})
		}

		var dependencies []reflect.Value
		dependencies = append(dependencies, moduleValue)
		for _, p := range dependencyProviders {
			v := p()
			dependencies = append(dependencies, reflect.ValueOf(v))
		}
		return providerMethod.Func.Call(dependencies)
	}
	v := reflect.MakeFunc(fn.Type(), provide)
	fn.Set(v)
}

type provider func() interface{}

type objectGraph struct {
	graph map[reflect.Type]provider
}

func (o *objectGraph) Inject(target interface{}) {
}

func (o *objectGraph) Get(v interface{}) {
}
