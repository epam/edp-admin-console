package webapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalizeAll(t *testing.T) {
	s := "foo"
	gotStr := CapitalizeAll(s)
	expectedStr := "FOO"
	assert.Equal(t, expectedStr, gotStr)
}

func TestCapitalizeFirstLetter(t *testing.T) {
	s := "foo"
	gotStr := CapitalizeFirstLetter(s)
	expectedStr := "Foo"
	assert.Equal(t, expectedStr, gotStr)
}

func TestLowercaseAll(t *testing.T) {
	s := "FOO"
	gotStr := LowercaseAll(s)
	expectedStr := "foo"
	assert.Equal(t, expectedStr, gotStr)
}
