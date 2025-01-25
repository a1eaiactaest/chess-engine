package engine_test

import (
	"chess-engine/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDSForcedCheckmate(t *testing.T) {
	a, err := engine.NewGame("r1b2k1r/pp1pQppp/3P4/5P2/8/5N2/4KP1P/qN3B1R b - - 3 19")
	assert.NoError(t, err)

	b, err := engine.NewGame("r1b2k1r/pp1pQppp/3P4/5P2/8/5N2/4KP1P/qN3B1R b - - 3 19")
	assert.NoError(t, err)

	ma := a.IDS(5, false)
	mb := b.IDS(5, false)
	expected := "f8g8"
	assert.Equal(t, ma, expected)
	assert.Equal(t, mb, expected)

}
