package core

import (
	"image/color"
	"testing"
)

func TestDrawTextOptionSizingMethodsResetOppositeMode(t *testing.T) {
	adapted := (DrawTextOption{FontSize: 42}).SetAdaptedSize(120, 30)
	if adapted.FontSize != 0 || adapted.MaxWidth != 120 || adapted.MaxHeight != 30 {
		t.Fatalf("SetAdaptedSize produced %#v", adapted)
	}

	static := adapted.SetStaticSize(18)
	if static.FontSize != 18 || static.MaxWidth != 0 || static.MaxHeight != 0 {
		t.Fatalf("SetStaticSize produced %#v", static)
	}
}

func TestDrawTextOptionEffectExpansion(t *testing.T) {
	base := DrawTextOption{
		FontColor: color.White,
		Text:      "effect",
		XOffset:   10,
		YOffset:   20,
	}.SetStaticSize(18)

	shadow := base.AddShadow(color.Black, 3).ToEffectGroup()
	if len(shadow) != 1 {
		t.Fatalf("shadow expansion count = %d, want 1", len(shadow))
	}
	if shadow[0].XOffset != 13 || shadow[0].YOffset != 23 || shadow[0].FontColor != color.Black {
		t.Fatalf("shadow expansion = %#v", shadow[0])
	}

	outline := base.AddOutline(color.Black, 1).ToEffectGroup()
	if len(outline) != 5 {
		t.Fatalf("outline radius 1 expansion count = %d, want 5", len(outline))
	}
	for _, option := range outline {
		if option.FontColor != color.Black {
			t.Fatalf("outline color = %#v, want black", option.FontColor)
		}
	}
}
