package animaterm

import (
	"sync"
)

// IUserInterface provides the main interface for terminal-based animations and UI rendering.
// It handles drawing elements, animations, and managing the terminal display buffer.
// All coordinates use a percentage-based system (0-100) for both X and Y axes.
type IUserInterface interface {
	SetBorderLeft(percent int) error
	SetBorderRight(percent int) error
	SetBorderTop(percent int) error
	SetBorderBottom(percent int) error
	SetBorderSides(percent int) error
	SetBorderTopBottom(percent int) error
	SetBorder(percent int) error
	// 	ClearScreen rest screen
	ClearScreen() error
	// StartDrawLoop initiates a loop, that renders the animations in a fixed loop
	// the returning channel must be closed (i.e close(ch)) in order to stop the
	// rendering; to allow the goroutine to stop all work the wg.Wait() command shopuld be used after
	// sending the stop signal
	StartDrawLoop(percentHeight int) (chan int, *sync.WaitGroup)
	DrawElement(pos IRelativePosition, text string, color int) int
	DrawElementsHorizontal(pos IRelativePosition, texts []string, positions []int, colors []int) int
	DrawTable(pos IRelativePosition, table [][]string, positions []int, colors []int) int
	DrawPattern(startPos IRelativePosition, expansion int, text string, color int, animation Animation) int
	MoveElement(startPos IRelativePosition, endPos IRelativePosition, text string, color int, animation Animation) error

	// PercentToAbsoluteWidth returns the absolute width of percentage in frame (disregarding the absolute position)
	PercentToAbsoluteWidth(percent int) int
	// PercentToAbsoluteWidth returns the absolute height of percentage in frame (disregarding the absolute position)
	PercentToAbsoluteHeight(percent int) int
	// PercentToAbsoluteXPostion returns the absolute x coordinate of percentage in frame
	PercentToAbsoluteXPostion(percent int) int
	// PercentToAbsoluteYPostion returns the absolute y coordinate of percentage in frame
	PercentToAbsoluteYPostion(percent int) int
}

// IRelativePosition represents a 2D position using percentage coordinates (0-100).
// It supports coordinate transformations, distance calculations, and offset management
// for multi-line text rendering.
type IRelativePosition interface {
	// GetX returns x coodinate in percent
	GetX() int
	// GetY returns y coodinate in percent
	GetY() int
	// GetXandY returns x and y coodinates in percent
	GetXandY() (int, int)
	// SetX sets x coordiante in percent
	SetX(percentx int) error
	// SetY see IRelativePosition
	SetY(percenty int) error
	// SetXandY see IRelativePosition
	SetXandY(percentx int, percenty int) error
	// IncrementOffset see IRelativePosition
	IncrementOffset() IRelativePosition
	// DistanceTo calculate distance between two points
	// and returns a new position representing the distance
	DistanceTo(p2 IRelativePosition) IRelativePosition
	// AddDistance calculates the resulting distance when adding a position
	// and returns a new position
	AddDistance(p2 IRelativePosition) IRelativePosition
	// MultiplyWith multiplies coordinates of point with factor
	// and returns a new position with (rounded) x and y coordinates
	MultiplyWith(factor float32) IRelativePosition
	ResetOffset() IRelativePosition
	GetOffset() int
	SetOffset(offset int) IRelativePosition
}
