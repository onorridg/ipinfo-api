package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIP(t *testing.T) {
	assert.Equal(t, true, IP("192.168.1.1"))
	assert.Equal(t, false, IP("256.256.256.256"))
	assert.Equal(t, true, IP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	assert.Equal(t, false, IP("test:test:test:test:test:test:test:test"))
}
