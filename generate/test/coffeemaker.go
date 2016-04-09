package coffeemaker

type Heater struct {
}

type Pump interface {
}

//go:generate go run ../main.go -type CoffeeMaker
type CoffeeMaker struct {
	Heater *Heater `inject:""`
	Pump   Pump    `inject:""`
}
