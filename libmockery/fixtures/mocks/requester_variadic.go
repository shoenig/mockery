package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"

// RequesterVariadic is an autogenerated mock type for the RequesterVariadic type
type RequesterVariadic struct {
	mock.Mock
}

// Get provides a mock function with given fields: values
func (mockerySelf *RequesterVariadic) Get(values ...string) bool {
	mockeryVariadicArg := make([]interface{}, len(values))
	for mockeryI := range values {
		mockeryVariadicArg[mockeryI] = values[mockeryI]
	}
	var mockeryCalledArg []interface{}
	mockeryCalledArg = append(mockeryCalledArg, mockeryVariadicArg...)
	ret := mockerySelf.Called(mockeryCalledArg...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(...string) bool); ok {
		r0 = rf(values...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MultiWriteToFile provides a mock function with given fields: filename, w
func (mockerySelf *RequesterVariadic) MultiWriteToFile(filename string, w ...io.Writer) string {
	mockeryVariadicArg := make([]interface{}, len(w))
	for mockeryI := range w {
		mockeryVariadicArg[mockeryI] = w[mockeryI]
	}
	var mockeryCalledArg []interface{}
	mockeryCalledArg = append(mockeryCalledArg, filename)
	mockeryCalledArg = append(mockeryCalledArg, mockeryVariadicArg...)
	ret := mockerySelf.Called(mockeryCalledArg...)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, ...io.Writer) string); ok {
		r0 = rf(filename, w...)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// OneInterface provides a mock function with given fields: a
func (mockerySelf *RequesterVariadic) OneInterface(a ...interface{}) bool {
	var mockeryCalledArg []interface{}
	mockeryCalledArg = append(mockeryCalledArg, a...)
	ret := mockerySelf.Called(mockeryCalledArg...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(...interface{}) bool); ok {
		r0 = rf(a...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Sprintf provides a mock function with given fields: format, a
func (mockerySelf *RequesterVariadic) Sprintf(format string, a ...interface{}) string {
	var mockeryCalledArg []interface{}
	mockeryCalledArg = append(mockeryCalledArg, format)
	mockeryCalledArg = append(mockeryCalledArg, a...)
	ret := mockerySelf.Called(mockeryCalledArg...)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, ...interface{}) string); ok {
		r0 = rf(format, a...)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}