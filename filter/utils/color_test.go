package utils

import (
	"image/color"
	"testing"
)

func TestParseHexColorSupportedForms(t *testing.T) {
	tests := []struct {
		input string
		want  color.RGBA
	}{
		{"#123", color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xff}},
		{"1234", color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0x44}},
		{"#102030", color.RGBA{R: 0x10, G: 0x20, B: 0x30, A: 0xff}},
		{"10203040", color.RGBA{R: 0x10, G: 0x20, B: 0x30, A: 0x40}},
		{"#aBcDeF", color.RGBA{R: 0xab, G: 0xcd, B: 0xef, A: 0xff}},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := ParseHexColor(test.input)
			if err != nil {
				t.Fatalf("ParseHexColor() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("ParseHexColor() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestParseHexColorRejectsMalformedValues(t *testing.T) {
	for _, input := range []string{"", "#12", "#12345", "#gggggg", "##123456", " 123456 "} {
		if _, err := ParseHexColor(input); err == nil {
			t.Errorf("ParseHexColor(%q) unexpectedly succeeded", input)
		}
	}
}

func TestLerpColorSupportsIncreasingAndDecreasingChannels(t *testing.T) {
	from := color.RGBA{R: 200, G: 10, B: 240, A: 255}
	to := color.RGBA{R: 100, G: 210, B: 40, A: 55}
	if got, want := LerpColor(from, to, 0.5), (color.RGBA{R: 150, G: 110, B: 140, A: 155}); got != want {
		t.Fatalf("LerpColor() = %#v, want %#v", got, want)
	}
	if got := LerpColor(from, to, -1); got != from {
		t.Fatalf("LerpColor below range = %#v, want from", got)
	}
	if got := LerpColor(from, to, 2); got != to {
		t.Fatalf("LerpColor above range = %#v, want to", got)
	}
}

func TestColorAdjustmentsClampAndPreserveExpectedAlpha(t *testing.T) {
	c := color.RGBA{R: 100, G: 150, B: 200, A: 77}
	if got := DarkenColor(c, 1); got != (color.RGBA{A: 77}) {
		t.Fatalf("DarkenColor full = %#v", got)
	}
	if got := LightenColor(c, 1); got != (color.RGBA{R: 255, G: 255, B: 255, A: 77}) {
		t.Fatalf("LightenColor full = %#v", got)
	}
	if got := AdjustOpacity(c, 0.5); got != (color.RGBA{R: 100, G: 150, B: 200, A: 127}) {
		t.Fatalf("AdjustOpacity half = %#v", got)
	}
	if got := ToHexString(c); got != "#6496C84D" {
		t.Fatalf("ToHexString() = %q", got)
	}
}
