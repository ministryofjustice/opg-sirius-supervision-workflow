package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBond_GetDeputyName(t *testing.T) {
	assert.Equal(t, "Angela White, Gary Black", Bond{Deputies: []string{"Angela White", "Gary Black"}}.GetDeputiesList())
	assert.Equal(t, "Jo Jane", Bond{Deputies: []string{"", "Jo Jane"}}.GetDeputiesList())
	assert.Equal(t, "Amy Samuel, Zoe William", Bond{Deputies: []string{"Zoe William", "", "Amy Samuel"}}.GetDeputiesList())
}
