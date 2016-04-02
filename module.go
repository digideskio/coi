package dagger

// A Module contributes bindings to the object graph.
//
//   package application
//
//   type AppModule struct {}
//
//   // ProvideHeater is invoked whenever a heater is required.
//   func (*AppModule) ProvideHeater() *Heater {
//     return &electricHeather{}
//   }
//
//   // Providers can have dependencies of their own.
//   // ProvidePump is invoked with a thermosiphon whenever a pump is required.
//   func (*AppModule) ProvidePump(pump *Thermosiphon) Pump {
//     return pump
//   }
type Module interface{}
