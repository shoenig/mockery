package test

import (
	"github.com/shoenig/mockery3/v3/libmockery/fixtures/http"
)

type HasConflictingNestedImports interface {
	RequesterNS
	Z() http.MyStruct
}
