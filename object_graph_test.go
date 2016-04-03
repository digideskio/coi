package dagger_test

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/dagger-go"
)

type TestModule struct{}

func (t *TestModule) ProvideRoundTripper() http.RoundTripper {
	return http.DefaultTransport
}

func (t *TestModule) ProvideHttpClient(transport http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Minute,
	}
}

func load(o dagger.ObjectGraph, v interface{}) {
	valueType := reflect.TypeOf(v)
	if valueType.Kind() != reflect.Ptr {
		panic("can only load into a pointer")
	}
	rawValue := reflect.ValueOf(v).Elem()
	value := o.Get(rawValue.Type())
	rawValue.Set(reflect.ValueOf(value))
}

func TestGraph(t *testing.T) {
	graph := dagger.NewObjectGraph(&TestModule{})
	var client *http.Client
	load(graph, &client)
	fmt.Println("http client timeout: ", client.Timeout)
}

type ProviderMethodReturnsMultiple struct{}

func (*ProviderMethodReturnsMultiple) ProvideMultiple() (*log.Logger, *http.Client) {
	return nil, nil
}

func TestProviderMethodReturningMultiplePanics(t *testing.T) {
	defer func() {
		err := recover().(error)
		assert.Equal(t, "Provider methods can return only one argument. func(*dagger_test.ProviderMethodReturnsMultiple) (*log.Logger, *http.Client) returns 2 values.", err.Error())
	}()

	dagger.NewObjectGraph(&ProviderMethodReturnsMultiple{})
}
