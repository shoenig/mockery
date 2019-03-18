package test

import (
	"net/http"

	my_http "github.com/shoenig/mockery3/v3/libmockery/fixtures/http"
)

// Example is an example
type Example interface {
	A() http.Flusher
	B(fixtureshttp string) my_http.MyStruct
}
