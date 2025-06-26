package animaterm

import (
	"testing"
)

func TestCreatePos(t *testing.T) {
	tests := []struct {
		name                 string
		x, y                 int
		expectedX, expectedY int
	}{
		{"Valid coordinates", 50, 75, 50, 75},
		{"Boundary values", 0, 100, 0, 100},
		{"Over max values", 150, 200, 100, 100},
		{"Under min values", -150, -200, -100, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := CreatePos(tt.x, tt.y)
			if pos.GetX() != tt.expectedX {
				t.Errorf("CreatePos(%d, %d).GetX() = %d, want %d", tt.x, tt.y, pos.GetX(), tt.expectedX)
			}
			if pos.GetY() != tt.expectedY {
				t.Errorf("CreatePos(%d, %d).GetY() = %d, want %d", tt.x, tt.y, pos.GetY(), tt.expectedY)
			}
		})
	}
}

func TestPositionSetX(t *testing.T) {
	pos := CreatePos(0, 0)

	tests := []struct {
		name     string
		x        int
		expected int
	}{
		{"Valid positive", 50, 50},
		{"Valid zero", 0, 0},
		{"Valid negative", -50, -50},
		{"Over max", 150, 100},
		{"Under min", -150, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pos.SetX(tt.x)
			if err != nil {
				t.Errorf("SetX(%d) returned error: %v", tt.x, err)
			}
			if pos.GetX() != tt.expected {
				t.Errorf("SetX(%d) = %d, want %d", tt.x, pos.GetX(), tt.expected)
			}
		})
	}
}

func TestPositionDistance(t *testing.T) {
	pos1 := CreatePos(10, 20)
	pos2 := CreatePos(30, 50)

	distance := pos1.DistanceTo(pos2)

	if distance.GetX() != 20 {
		t.Errorf("Distance X = %d, want 20", distance.GetX())
	}
	if distance.GetY() != 30 {
		t.Errorf("Distance Y = %d, want 30", distance.GetY())
	}
}

func TestPositionAddDistance(t *testing.T) {
	pos1 := CreatePos(10, 20)
	pos2 := CreatePos(5, 15)

	result := pos1.AddDistance(pos2)

	if result.GetX() != 15 {
		t.Errorf("AddDistance X = %d, want 15", result.GetX())
	}
	if result.GetY() != 35 {
		t.Errorf("AddDistance Y = %d, want 35", result.GetY())
	}
}

func TestPositionMultiplyWith(t *testing.T) {
	pos := CreatePos(10, 20)

	result := pos.MultiplyWith(2.5)

	if result.GetX() != 25 {
		t.Errorf("MultiplyWith(2.5) X = %d, want 25", result.GetX())
	}
	if result.GetY() != 50 {
		t.Errorf("MultiplyWith(2.5) Y = %d, want 50", result.GetY())
	}
}

func TestPositionOffset(t *testing.T) {
	pos := CreatePos(10, 20)

	// Test initial offset
	if pos.GetOffset() != 0 {
		t.Errorf("Initial offset = %d, want 0", pos.GetOffset())
	}

	// Test increment
	pos.IncrementOffset()
	if pos.GetOffset() != 1 {
		t.Errorf("After increment offset = %d, want 1", pos.GetOffset())
	}

	// Test set offset
	pos.SetOffset(5)
	if pos.GetOffset() != 5 {
		t.Errorf("After SetOffset(5) offset = %d, want 5", pos.GetOffset())
	}

	// Test reset
	pos.ResetOffset()
	if pos.GetOffset() != 0 {
		t.Errorf("After reset offset = %d, want 0", pos.GetOffset())
	}
}
