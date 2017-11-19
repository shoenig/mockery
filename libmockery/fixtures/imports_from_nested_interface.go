package test

import (
	"github.com/shoenig/mockery/libmockery/fixtures/http"
)

type HasConflictingNestedImports interface {
	RequesterNS
	Z() http.MyStruct
}
