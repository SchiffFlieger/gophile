package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringLength(t *testing.T) {
	for k := 0; k < 100; k++ {
		str1 := Rnd.RandomString(k)
		assert.Equal(t, k, len(str1))
	}
}

func TestRandomStringRandomness(t *testing.T) {
	for k := 50; k < 100; k++ {
		str1 := Rnd.RandomString(k)
		str2 := Rnd.RandomString(k)
		assert.NotEqual(t, str1, str2)
	}
}
