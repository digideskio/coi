package coi_test

import (
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/coi"
)

type TestModule struct{}

func (t *TestModule) ProvideRoundTripper() http.RoundTripper {
	return http.DefaultTransport
}

func (t *TestModule) ProvideTimeout() time.Duration {
	return 10 * time.Minute
}

func (t *TestModule) ProvideHttpClient(transport http.RoundTripper, timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

type A struct {
	Client  *http.Client  `inject:""`
	Timeout time.Duration `inject:""`
}

func load(o coi.ObjectGraph, v interface{}) {
	valueType := reflect.TypeOf(v)
	if valueType.Kind() != reflect.Ptr {
		panic("can only load into a pointer")
	}
	rawValue := reflect.ValueOf(v).Elem()
	value := o.Get(rawValue.Type())
	rawValue.Set(reflect.ValueOf(value))
}

func TestObjectGraphGet(t *testing.T) {
	graph := coi.NewObjectGraph(&TestModule{})
	var client *http.Client
	load(graph, &client)
	assert.Equal(t, 10*time.Minute, client.Timeout)
}

type B struct {
	ADependency A `inject:""`
}

func TestObjectGraphGetNonModuleDependencies(t *testing.T) {
	graph := coi.NewObjectGraph(&TestModule{})
	var b *B
	load(graph, &b)
	assert.Equal(t, 10*time.Minute, b.ADependency.Client.Timeout)
}

func TestObjectGraphInject(t *testing.T) {
	graph := coi.NewObjectGraph(&TestModule{})
	a := &A{}
	graph.Inject(a)
	assert.Equal(t, 10*time.Minute, a.Client.Timeout)
	assert.Equal(t, 10*time.Minute, a.Timeout)
}

type ProviderMethodReturnsMultiple struct{}

func (*ProviderMethodReturnsMultiple) ProvideMultiple() (*log.Logger, *http.Client) {
	return nil, nil
}

func TestProviderMethodReturningMultiplePanics(t *testing.T) {
	defer func() {
		err := recover().(error)
		assert.Equal(t, "Provider methods can return only one argument. func(*coi_test.ProviderMethodReturnsMultiple) (*log.Logger, *http.Client) returns 2 values.", err.Error())
	}()

	coi.NewObjectGraph(&ProviderMethodReturnsMultiple{})
}
