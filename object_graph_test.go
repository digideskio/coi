package dagger_test

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/dagger-go"
)

type TestModule struct{}

func (t *TestModule) ProvideHttpClient() *http.Client {
	return http.DefaultClient
}

func (t *TestModule) ProvideLogger(out io.Writer) *log.Logger {
	return log.New(out, "", log.LstdFlags)
}

func TestGraph(t *testing.T) {
	dagger.NewObjectGraph(&TestModule{})
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
