package hsdp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrune(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 5}
	c := prunePorts(a, b)
	assert.Equal(t, []int{2, 3, 4}, c)
	d := prunePorts(c, []int{3})
	assert.Equal(t, []int{2, 4}, d)
}
