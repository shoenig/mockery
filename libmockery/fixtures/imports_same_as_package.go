package test

import (
	"github.com/shoenig/mockery3/libmockery/fixtures/test"
)

type C int

type ImportsSameAsPackage interface {
	A() test.B
	B() KeyManager
	C(C)
}
