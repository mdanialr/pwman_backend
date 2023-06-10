package helper_test

import (
	"testing"

	"github.com/magiconair/properties/assert"
	help "github.com/mdanialr/pwman_backend/pkg/helper"
)

var sample = []string{
	"hello",
	"world",
	"and",
	"universe",
}

func TestPad(t *testing.T) {
	const expect = "hello world and universe"
	actual := help.Pad(sample...)
	assert.Equal(t, actual, expect)
}

func BenchmarkPad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		help.Pad(sample...)
	}
}
