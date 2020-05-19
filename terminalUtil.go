package animaterm

import (
	"fmt"
	"math/rand"
	"os"

	CB "golang.org/x/mobile/exp/sprite/clock"

	"golang.org/x/sys/unix"
)

var pl = fmt.Println
var p = fmt.Print
var pf = fmt.Printf

// Animation ...
type Animation struct {
	AnimationType AnimationType
	Duration      int64
	Direction     Direction
	GradientV     bool
	GradientH     bool
}

// Direction ...
type Direction int

//  Up ...
const (
	Right Direction = iota
	Left
	Down
	Up
)

// AnimationType ...
type AnimationType int

// EaseIn ...
const (
	EaseIn AnimationType = iota
	EaseOut
	EaseInOut
	Custom
	Ikea
)

func getAnimation(animationType AnimationType) func(t0, t1, t CB.Time) float32 {
	switch animationType {
	case EaseIn:
		return CB.CubicBezier(0.42, 0, 1, 1)
	case EaseOut:
		return CB.CubicBezier(0, 0, 0.58, 1)
	case EaseInOut:
		return CB.CubicBezier(0.42, 0, 0.58, 1)
	case Ikea:
		return CB.CubicBezier(0.6, 1.2, 0.3, 0.9)
	case Custom:
		return CB.CubicBezier(0, 1, 1, 0)
	default:
		return CB.CubicBezier(0.6, 1.2, 0.3, 0.9)
	}
}

// ControlSequence enum for coloring output
type ControlSequence int

// List of control sequeces colors
const (
	BLACK                       int = 0
	RED                             = 1
	GREEN                           = 2
	YELLOW                          = 3
	BLUE                            = 4
	VIOLET                          = 5
	WHITE                           = 7
	GREY                            = 231
	LIGHTGREY                       = 234
	DARKBLUE                        = 17
	TERMINALGREEN                   = 112
	ORANGE                          = 202
	RED2                            = 160
	PINK                            = 177
	COLORPATTERNPASTEL              = 34
	COLORPATTERNSKYLIGHT            = 39
	COLORPATTERNMEADOWS1            = 40
	COLORPATTERNMEADOWS2            = 46
	COLORPATTERNNEON1               = 64
	COLORPATTERNNEON2               = 70
	COLORPATTERNNEON3               = 76
	COLORPATTERNGREENFOUNDATION     = 77
	COLORPATTERNLIME                = 82
	COLORPATTERNGREY                = 106
	COLORPATTERNSPLITMEADOWS        = 42
	COLORPATTERNBABYSTEPS1          = 74
	COLORPATTERNBABYSTEPS2          = 81
	COLORPATTERNGOINGGREY1          = 100
	COLORPATTERNGOINGGREY2          = 101
	COLORPATTERNGOINGGREY3          = 102
	COLORPATTERNGOINGGREY4          = 103
	COLORPATTERNGOINGGREY5          = 104
	COLORPATTERNGOINGGREY6          = 105
	RESET                           = 500
	RESETLINE                       = 501
	RANDOMGREY                      = 502
	RANDOM                          = 503
	BLANK                           = 504
	ALREADYCOLORED                  = 505
)

func getControlSequence(sequence int) string {
	switch code := sequence; {
	case code >= 0 && code < 256:
		return fmt.Sprintf("\033[38;5;%03dm", code)
	case code == 500:
		// Reset all custom styles
		return "\033[0m"
	case code == 501:
		// Return cursor to start of line and clean it
		return "\r\033[K"
	case code == 502:
		// Return random greyscale color
		return fmt.Sprintf("\033[38;5;%03dm", rand.Intn(22)+231)
	case code == 503:
		// Return random color
		return fmt.Sprintf("\033[38;5;%03dm", rand.Intn(255))
	default:
		return "\033[37m"
	}
}

// Color apply
func Color(str string, color int) string {
	if color == BLANK || color == ALREADYCOLORED {
		return str
	}
	return fmt.Sprintf("%s%s%s", getControlSequence(color), str, getControlSequence(RESET))
}

func getWinsize() (*unix.Winsize, error) {

	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return nil, os.NewSyscallError("GetWinsize", err)
	}

	return ws, nil
}

func SetWinsize(width int, height int) error {
	uws := &unix.Winsize{Row: uint16(height), Col: uint16(width), Xpixel: 0, Ypixel: 0}
	return unix.IoctlSetWinsize(int(os.Stdout.Fd()), unix.TIOCSWINSZ, uws)
}

// Height get full height of terminal window
func Height() int {
	ws, err := getWinsize()
	if err != nil {
		return -1
	}
	return int(ws.Row)
}

// Width get full width of terminal window
func Width() int {
	ws, err := getWinsize()
	if err != nil {
		return -1
	}
	return int(ws.Col)
}
