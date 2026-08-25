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

func TestParseKeywords(t *testing.T) {
	assert.True(t, pkg.Gesture == pkg.ParseLine(0, "gesture = 3, horizontal, workspace").Type)
	assert.True(t, pkg.Permission == pkg.ParseLine(0, "permission = /usr/(bin|local/bin)/grim, screencopy, allow").Type)
	assert.True(t, pkg.WorkspaceRule == pkg.ParseLine(0, "workspace = w[tv1], gapsout:0, gapsin:0").Type)
	assert.True(t, pkg.LayerRule == pkg.ParseLine(0, "layerrule = blur, waybar").Type)
	assert.True(t, pkg.Source == pkg.ParseLine(0, "source = ~/.config/hypr/myColors.conf").Type)
	assert.True(t, pkg.WindowRule == pkg.ParseLine(0, "windowrule = float, ^(kitty)$").Type)
}

func TestBlocksTakePrecedenceOverKeywords(t *testing.T) {
	assert.True(t, pkg.BlockStart == pkg.ParseLine(0, "windowrule {").Type)
	assert.True(t, pkg.BlockStart == pkg.ParseLine(0, "monitorv2 {").Type)
	assert.True(t, pkg.BlockStart == pkg.ParseLine(0, "gestures {").Type)
}

func TestParseValueContainingEquals(t *testing.T) {
	l := pkg.ParseLine(0, "env = GTK_THEME,Adwaita=dark")
	k, v := l.Parse(true)
	assert.Equal(t, "GTK_THEME", k.Key)
	assert.Equal(t, "Adwaita=dark", v.ToString())
}

func TestParseExampleFile(t *testing.T) {
	cfg, err := pkg.ReadConfig("example_config.conf")
	assert.Equal(t, nil, err)

	cfg.Show()

	grouped := pkg.BuildGroupedConfig(cfg)

	assert.True(t, len(grouped.Variables) > 0)
	assert.True(t, len(grouped.Monitors) > 0)
	assert.True(t, len(grouped.Env) > 0)
	assert.True(t, len(grouped.Bindings) > 0)
	assert.True(t, len(grouped.Gestures) > 0)
	// windowrule {} blocks land under Branches alongside general, decoration, etc.
	assert.True(t, len(grouped.Branches) >= 10)
}

func TestParseLegacyExampleFile(t *testing.T) {
	cfg, err := pkg.ReadConfig("example_config_legacy.conf")
	assert.Equal(t, nil, err)

	grouped := pkg.BuildGroupedConfig(cfg)

	assert.True(t, len(grouped.Variables) > 0)
	assert.True(t, len(grouped.WindowRules) > 0)
}
