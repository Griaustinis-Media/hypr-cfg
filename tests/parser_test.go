package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"griaustinismedia/hypr-cfg/pkg"
)

func TestParseComments(t *testing.T) {
	l1 := pkg.ParseLine(0, "# Hello world")
	assert.True(t, pkg.Comment == l1.Type)

	l2 := pkg.ParseLine(0, "   # Hello world")
	assert.True(t, pkg.Comment == l2.Type)
}

func TestParseEmpty(t *testing.T) {
	l1 := pkg.ParseLine(0, "")
	assert.True(t, pkg.Empty == l1.Type)

	l2 := pkg.ParseLine(0, "   ")
	assert.True(t, pkg.Empty == l2.Type)
}

func TestParseVariableDef(t *testing.T) {
	l := pkg.ParseLine(0, "$terminal = kitty")
	assert.True(t, pkg.Variable == l.Type)
}

func TestParseExampleFile(t *testing.T) {
	cfg, err := pkg.ReadConfig("example_config.conf")
	assert.Equal(t, nil, err)

	cfg.Show()
}
