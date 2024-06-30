package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidCep(t *testing.T) {
	validCep := "11111111"
	invalidCep := "asdf123"
	shortCep := "123"

	assert.True(t, isValidCep(validCep))
	assert.False(t, isValidCep(invalidCep))
	assert.False(t, isValidCep(shortCep))
}
