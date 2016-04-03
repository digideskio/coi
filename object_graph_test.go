package dagger_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/dagger-go"
	"github.com/f2prateek/go-pointers"
)

type TestModule struct{}

func (t *TestModule) ProvideHttpClient() *http.Client {
	return http.DefaultClient
}

func (t *TestModule) ProvideLogger(out io.Writer) *log.Logger {
	return log.New(out, "", log.LstdFlags)
}

func (t *TestModule) ProvideInt() *int {
	return pointers.Int(4)
}

func (t *TestModule) ProvideString(a *int) *string {
	return pointers.String(fmt.Sprintf("%d", *a))
}

func TestGraph(t *testing.T) {
	graph := dagger.NewObjectGraph(&TestModule{})
	s := graph.Get(reflect.TypeOf((*string)(nil))).(*string)
	fmt.Println("got: ", *s)
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
