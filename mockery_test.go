package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func configFromCommandLine(str string) flags {
	return parseFlags(strings.Split(str, " "))
}

func TestParseConfigDefaults(t *testing.T) {
	config := configFromCommandLine("mockery")
	assert.Equal(t, false, config.version)
	assert.Equal(t, "", config.iface)
	assert.Equal(t, false, config.stdout)
	assert.Equal(t, "", config.pkgname)
	assert.Equal(t, "", config.comment)
}

func TestParseConfigFlippingValues(t *testing.T) {
	config := configFromCommandLine("mockery -interface=MyInterface -stdout=true -package=mypackage -comment=blahblah -version=true")
	assert.Equal(t, true, config.version)
	assert.Equal(t, "MyInterface", config.iface)
	assert.Equal(t, true, config.stdout)
	assert.Equal(t, "mypackage", config.pkgname)
	assert.Equal(t, "blahblah", config.comment)
}
