package animaterm

import (
	"testing"
)

func TestCreateUI(t *testing.T) {
	ui := CreateUI()
	if ui == nil {
		t.Fatal("CreateUI() returned nil")
	}
	
	// Test that UI is properly initialized
	uiImpl := ui.(*UserInterface)
	if uiImpl.msPerFrame != 320 {
		t.Errorf("Initial msPerFrame = %d, want 320", uiImpl.msPerFrame)
	}
	
	if uiImpl.pixels == nil {
		t.Error("Pixels array not initialized")
	}
	
	if uiImpl.dirtyRegions == nil {
		t.Error("Dirty regions array not initialized")
	}
}

func TestSetBorderMethods(t *testing.T) {
	ui := CreateUI().(*UserInterface)
	
	tests := []struct {
		name string
		percent int
		shouldError bool
	}{
		{"Valid border", 10, false},
		{"Zero border", 0, false},
		{"Max border", 50, false},
		{"Over max border", 60, true},
		{"Negative border", -10, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ui.SetBorderLeft(tt.percent)
			if (err != nil) != tt.shouldError {
				t.Errorf("SetBorderLeft(%d) error = %v, wantErr %v", tt.percent, err, tt.shouldError)
			}
			
			err = ui.SetBorderRight(tt.percent)
			if (err != nil) != tt.shouldError {
				t.Errorf("SetBorderRight(%d) error = %v, wantErr %v", tt.percent, err, tt.shouldError)
			}
			
			err = ui.SetBorderTop(tt.percent)
			if (err != nil) != tt.shouldError {
				t.Errorf("SetBorderTop(%d) error = %v, wantErr %v", tt.percent, err, tt.shouldError)
			}
			
			err = ui.SetBorderBottom(tt.percent)
			if (err != nil) != tt.shouldError {
				t.Errorf("SetBorderBottom(%d) error = %v, wantErr %v", tt.percent, err, tt.shouldError)
			}
		})
	}
}

func TestPercentToAbsoluteConversion(t *testing.T) {
	ui := CreateUI().(*UserInterface)
	
	// Test width conversion
	width50 := ui.PercentToAbsoluteWidth(50)
	expectedWidth := Width() * 50 / 100
	if width50 != expectedWidth {
		t.Errorf("PercentToAbsoluteWidth(50) = %d, want %d", width50, expectedWidth)
	}
	
	// Test height conversion
	height50 := ui.PercentToAbsoluteHeight(50)
	expectedHeight := Height() * 50 / 100
	if height50 != expectedHeight {
		t.Errorf("PercentToAbsoluteHeight(50) = %d, want %d", height50, expectedHeight)
	}
}

func TestMoveElementValidation(t *testing.T) {
	ui := CreateUI().(*UserInterface)
	
	pos1 := CreatePos(10, 10)
	pos2 := CreatePos(50, 50)
	animation := Animation{Duration: 1000}
	
	// Test nil positions
	err := ui.MoveElement(nil, pos2, "test", RED, animation)
	if err == nil {
		t.Error("MoveElement with nil startPos should return error")
	}
	
	err = ui.MoveElement(pos1, nil, "test", RED, animation)
	if err == nil {
		t.Error("MoveElement with nil endPos should return error")
	}
	
	// Test empty text
	err = ui.MoveElement(pos1, pos2, "", RED, animation)
	if err == nil {
		t.Error("MoveElement with empty text should return error")
	}
	
	// Test negative duration
	badAnimation := Animation{Duration: -100}
	err = ui.MoveElement(pos1, pos2, "test", RED, badAnimation)
	if err == nil {
		t.Error("MoveElement with negative duration should return error")
	}
	
	// Test valid parameters
	err = ui.MoveElement(pos1, pos2, "test", RED, animation)
	if err != nil {
		t.Errorf("MoveElement with valid parameters should not return error: %v", err)
	}
}

func TestDrawPatternValidation(t *testing.T) {
	ui := CreateUI().(*UserInterface)
	
	pos := CreatePos(10, 10)
	animation := Animation{Duration: 1000}
	
	// Test nil position
	result := ui.DrawPattern(nil, 50, "█", RED, animation)
	if result != -1 {
		t.Error("DrawPattern with nil position should return -1")
	}
	
	// Test invalid expansion
	result = ui.DrawPattern(pos, -10, "█", RED, animation)
	if result != -1 {
		t.Error("DrawPattern with negative expansion should return -1")
	}
	
	result = ui.DrawPattern(pos, 300, "█", RED, animation)
	if result != -1 {
		t.Error("DrawPattern with excessive expansion should return -1")
	}
	
	// Test empty text
	result = ui.DrawPattern(pos, 50, "", RED, animation)
	if result != -1 {
		t.Error("DrawPattern with empty text should return -1")
	}
}

func TestSetPixelSafety(t *testing.T) {
	ui := CreateUI().(*UserInterface)
	
	// Test out of bounds access
	ui.setPixel(-1, 0, "X")
	ui.setPixel(0, -1, "X")
	ui.setPixel(ui.width+10, ui.height+10, "X")
	
	// Test valid access
	ui.setPixel(1, 1, "X")
	
	// Verify pixel was set
	ui.pixelsMutex.RLock()
	if ui.pixels[1][1] != "X" {
		t.Error("Valid pixel was not set correctly")
	}
	ui.pixelsMutex.RUnlock()
}