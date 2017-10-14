package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvert(t *testing.T) {
	for _, tc := range [][2]string{
		{"", ""},
		{"ooɟ", "foo"},
		{"oןɐzuob", "gonzalo"},
		{"dıןɟǝןqɐʇ", "tableflip"},
	} {
		assert.Equal(t, tc[0], invert(tc[1]))
	}
}
