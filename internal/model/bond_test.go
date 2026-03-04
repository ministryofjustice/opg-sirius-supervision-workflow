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

func TestBond_GetDeputyName(t *testing.T) {
	assert.Equal(t, "Angela White, Gary Black", Bond{Deputies: map[string]string{"deputy1": "Angela White", "deputy2": "Gary Black"}}.GetDeputiesList())
	assert.Equal(t, "Jo Jane", Bond{Deputies: map[string]string{"deputy1": "", "deputy2": "Jo Jane"}}.GetDeputiesList())
	assert.Equal(t, "Amy Samuel, Zoe William", Bond{Deputies: map[string]string{"deputy1": "Zoe William", "deputy2": "", "deputy3": "Amy Samuel"}}.GetDeputiesList())
}
