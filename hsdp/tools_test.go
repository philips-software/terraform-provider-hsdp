package hsdp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrune(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 5}
	c := prunePorts(a, b)
	assert.Equal(t, []int{2, 3, 4}, c)
	d := prunePorts(c, []int{3})
	assert.Equal(t, []int{2, 4}, d)
}

func TestSlidingExpiresOn(t *testing.T) {
	now := time.Date(1975, 10, 28, 0, 0, 0, 0, time.UTC)
	expected := time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC)

	sliding := slidingExpiresOn(now)
	assert.Equal(t, expected, sliding)
}
