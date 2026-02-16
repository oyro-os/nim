package cmd

import "testing"

func TestResolveInputOutput(t *testing.T) {
	t.Run("rejects too many args", func(t *testing.T) {
		_, _, err := resolveInputOutput([]string{"a", "b", "c"}, "", "")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("uses positional input and output", func(t *testing.T) {
		in, out, err := resolveInputOutput([]string{"in.jpg", "out.png"}, "", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in != "in.jpg" || out != "out.png" {
			t.Fatalf("unexpected files: %s -> %s", in, out)
		}
	})

	t.Run("single positional arg sets output", func(t *testing.T) {
		in, out, err := resolveInputOutput([]string{"out.png"}, "in.jpg", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if in != "in.jpg" || out != "out.png" {
			t.Fatalf("unexpected files: %s -> %s", in, out)
		}
	})

	t.Run("requires input", func(t *testing.T) {
		_, _, err := resolveInputOutput(nil, "", "out.png")
		if err == nil {
			t.Fatal("expected input validation error")
		}
	})
}

func TestParseSizeArg(t *testing.T) {
	width, height, err := parseSizeArg("300X200", 800, 512)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if width != 300 || height != 200 {
		t.Fatalf("unexpected size: %dx%d", width, height)
	}

	_, _, err = parseSizeArg("300", 800, 512)
	if err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestParseResizeMode(t *testing.T) {
	mode, err := parseResizeMode("fit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mode != "fit" {
		t.Fatalf("unexpected mode: %s", mode)
	}

	_, err = parseResizeMode("weird")
	if err == nil {
		t.Fatal("expected invalid mode error")
	}
}

func TestParsePadColor(t *testing.T) {
	color, err := parsePadColor("#0aFF10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if color != [3]uint8{0x0a, 0xFF, 0x10} {
		t.Fatalf("unexpected color: %#v", color)
	}

	defaultColor, err := parsePadColor("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if defaultColor != [3]uint8{255, 255, 255} {
		t.Fatalf("unexpected default color: %#v", defaultColor)
	}

	_, err = parsePadColor("#fff")
	if err == nil {
		t.Fatal("expected invalid color format error")
	}
}

func TestValidateDimensionAndQuality(t *testing.T) {
	if err := validateDimension("width", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateDimension("width", 0); err == nil {
		t.Fatal("expected width validation error")
	}

	if err := validateQuality(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateQuality(101); err == nil {
		t.Fatal("expected quality validation error")
	}
}
