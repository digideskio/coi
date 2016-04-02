package dagger

// ObjectGraph is a graph of objects linked by their dependencies.
type ObjectGraph interface {
	// Inject the fields of `target`.
	Inject(target interface{})
	// Get retrieves an object from the stores the result in the value pointed to by v.
	Get(v interface{})
}

// NewObjectGraph returns a new dependency graph using the given modules.
func NewObjectGraph(modules ...Module) ObjectGraph {
	return &objectGraph{}
}

type objectGraph struct {
}

func (o *objectGraph) Inject(target interface{}) {
}

func (o *objectGraph) Get(v interface{}) {
}
