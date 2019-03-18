package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func configFromCommandLine(str string) flags {
	return parseFlags(strings.Split(str, " "))
}

func TestParseConfigDefaults(t *testing.T) {
	config := configFromCommandLine("mockery")
	require.Equal(t, false, config.version)
	require.Equal(t, "", config.iface)
	require.Equal(t, false, config.stdout)
	require.Equal(t, "", config.pkgname)
	require.Equal(t, "", config.comment)
}

func TestParseConfigFlippingValues(t *testing.T) {
	config := configFromCommandLine("mockery -interface=MyInterface -stdout=true -package=mypackage -comment=blahblah -version=true")
	require.Equal(t, true, config.version)
	require.Equal(t, "MyInterface", config.iface)
	require.Equal(t, true, config.stdout)
	require.Equal(t, "mypackage", config.pkgname)
	require.Equal(t, "blahblah", config.comment)
}
