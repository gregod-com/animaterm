package animaterm

import (
	"testing"
)

func TestColor(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		color    int
		expected string
	}{
		{
			name:     "Red color",
			text:     "Hello",
			color:    RED,
			expected: "\033[38;5;001mHello\033[0m",
		},
		{
			name:     "Blank color",
			text:     "Hello",
			color:    BLANK,
			expected: "Hello",
		},
		{
			name:     "Already colored",
			text:     "Hello",
			color:    ALREADYCOLORED,
			expected: "Hello",
		},
		{
			name:     "Valid 256 color",
			text:     "Test",
			color:    100,
			expected: "\033[38;5;100mTest\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Color(tt.text, tt.color)
			if result != tt.expected {
				t.Errorf("Color(%q, %d) = %q, want %q", tt.text, tt.color, result, tt.expected)
			}
		})
	}
}

func TestGetControlSequence(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{
			name:     "Valid 256 color",
			code:     42,
			expected: "\033[38;5;042m",
		},
		{
			name:     "Reset",
			code:     RESET,
			expected: "\033[0m",
		},
		{
			name:     "Reset line",
			code:     RESETLINE,
			expected: "\r\033[K",
		},
		{
			name:     "Default fallback",
			code:     999,
			expected: "\033[37m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getControlSequence(tt.code)
			if result != tt.expected {
				t.Errorf("getControlSequence(%d) = %q, want %q", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRandomColors(t *testing.T) {
	// Test that random colors return valid sequences
	randomGreySeq := getControlSequence(RANDOMGREY)
	if len(randomGreySeq) == 0 {
		t.Error("RANDOMGREY should return non-empty sequence")
	}
	
	randomSeq := getControlSequence(RANDOM)
	if len(randomSeq) == 0 {
		t.Error("RANDOM should return non-empty sequence")
	}
}

func TestHeightWidth(t *testing.T) {
	height := Height()
	width := Width()
	
	// Should return reasonable defaults even if terminal detection fails
	if height < 1 {
		t.Errorf("Height() = %d, should be positive", height)
	}
	
	if width < 1 {
		t.Errorf("Width() = %d, should be positive", width)
	}
	
	// Should be reasonable terminal sizes
	if height > 1000 {
		t.Errorf("Height() = %d, seems unreasonably large", height)
	}
	
	if width > 1000 {
		t.Errorf("Width() = %d, seems unreasonably large", width)
	}
}

func TestAnimationTypes(t *testing.T) {
	// Test that all animation types return valid functions
	animationTypes := []AnimationType{EaseIn, EaseOut, EaseInOut, Custom, Ikea}
	
	for _, animType := range animationTypes {
		t.Run(animType.String(), func(t *testing.T) {
			animFunc := getAnimation(animType)
			if animFunc == nil {
				t.Errorf("getAnimation(%v) returned nil", animType)
			}
		})
	}
}

// Add String method for AnimationType to support test names
func (a AnimationType) String() string {
	switch a {
	case EaseIn:
		return "EaseIn"
	case EaseOut:
		return "EaseOut"
	case EaseInOut:
		return "EaseInOut"
	case Custom:
		return "Custom"
	case Ikea:
		return "Ikea"
	default:
		return "Unknown"
	}
}