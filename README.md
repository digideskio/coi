# coi

CoI is a dependency injector for Go.

Typically components wire up their own dependencies.

```go
func NewCoffeeApp() *CoffeeApp {
  return &CoffeeApp{
    Heater: &ElectricHeater{},
    Pump: &Thermosiphon{},
  }
}
```

Dependency injection inverses this flow, and makes it easy to create reusable, interchangeable modules.

```go
func NewCoffeeApp(heater *Heater, pump Pump) *CoffeeApp {
  return &CoffeeApp{
    Heater: heater,
    Pump: pump,
  }
}
```

Dependency injection facilitates creating different environments, and makes it easy to swap out your `ProdLogger` for an in-memory `TestLogger`.  Or swap it for `DevLogger` during development.

Wiring up this graph throughout your app is tedious and boring. CoI eliminates the clumsy duct that wires your dependencies. Simply declare your dependencies, specify how to satisfy them, and profit!

## Providing Dependencies

Provider methods satisfy dependencies. A provider method is an exported method beginning with `Provide`.

```go
 // ProvideHeater is invoked whenever a heater is required.
 func (*module) ProvideHeater() *Heater {
   return &ElectricHeater{}
 }
 ```

 Provider methods can have dependencies of their own.

 ```go
 // ProvidePump is invoked with a thermosiphon whenever a pump is required.
 func (*DripCoffeeModule) ProvidePump(pump *Thermosiphon) Pump {
   return pump
 }
 ```

## Modules

Modules are reusable components that logically group provider methods.

```go
type DripCoffeeModule struct {}

func (*DripCoffeeModule) ProvideHeater() *Heater {
  return &ElectricHeater{}
}

func (*DripCoffeeModule) ProvidePump(pump *Thermosiphon) Pump {
  return pump
}
```

## ObjectGraph

The `ObjectGraph` forms a graph of dependencies. It composes modules and their provider methods to satisfy dependencies when requested.

```go
objectGraph := coi.NewObjectGraph(&DripCoffeeModule{})

coffeeApp := &CoffeeApp{}
objectGraph.Inject(coffeeApp)
```
