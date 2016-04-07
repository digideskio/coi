package dagger

// A Module contributes bindings to the object graph.
//
//   package coffeemaker
//
//   type DripCoffeeModule struct {}
//
//   // ProvideHeater is invoked whenever a heater is required.
//   func (*DripCoffeeModule) ProvideHeater() *Heater {
//     return &ElectricHeater{}
//   }
//
//   // Provider methods can have dependencies of their own.
//   // ProvidePump is invoked with a thermosiphon whenever a pump is required.
//   func (*DripCoffeeModule) ProvidePump(pump *Thermosiphon) Pump {
//     return pump
//   }
type Module interface{}
