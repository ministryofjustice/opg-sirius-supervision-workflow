package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBond_GetBondAmount(t *testing.T) {
	assert.Equal(t, "£0.00", Bond{BondAmount: 0}.GetBondAmount())
	assert.Equal(t, "£1.01", Bond{BondAmount: 101}.GetBondAmount())
	assert.Equal(t, "£1.00", Bond{BondAmount: 100}.GetBondAmount())
	assert.Equal(t, "£1.70", Bond{BondAmount: 170}.GetBondAmount())
}
