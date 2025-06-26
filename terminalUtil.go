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

// Animation defines the parameters for animated operations.
// Duration is in milliseconds, Direction controls expansion/movement direction,
// and gradient flags enable color transitions during animation.
type Animation struct {
	AnimationType AnimationType
	Duration      int64
	Direction     Direction
	GradientV     bool
	GradientH     bool
}

// Direction ...
type Direction int

// Up ...
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
	RED                         int = 1
	GREEN                       int = 2
	YELLOW                      int = 3
	BLUE                        int = 4
	VIOLET                      int = 5
	WHITE                       int = 7
	GREY                        int = 231
	LIGHTGREY                   int = 234
	DARKBLUE                    int = 17
	TERMINALGREEN               int = 112
	ORANGE                      int = 202
	RED2                        int = 160
	PINK                        int = 177
	COLORPATTERNPASTEL          int = 34
	COLORPATTERNSKYLIGHT        int = 39
	COLORPATTERNMEADOWS1        int = 40
	COLORPATTERNMEADOWS2        int = 46
	COLORPATTERNNEON1           int = 64
	COLORPATTERNNEON2           int = 70
	COLORPATTERNNEON3           int = 76
	COLORPATTERNGREENFOUNDATION int = 77
	COLORPATTERNLIME            int = 82
	COLORPATTERNGREY            int = 106
	COLORPATTERNSPLITMEADOWS    int = 42
	COLORPATTERNBABYSTEPS1      int = 74
	COLORPATTERNBABYSTEPS2      int = 81
	COLORPATTERNGOINGGREY1      int = 100
	COLORPATTERNGOINGGREY2      int = 101
	COLORPATTERNGOINGGREY3      int = 102
	COLORPATTERNGOINGGREY4      int = 103
	COLORPATTERNGOINGGREY5      int = 104
	COLORPATTERNGOINGGREY6      int = 105
	RESET                       int = 500
	RESETLINE                   int = 501
	RANDOMGREY                  int = 502
	RANDOM                      int = 503
	BLANK                       int = 504
	ALREADYCOLORED              int = 505
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

// Color applies ANSI color codes to text for terminal display.
// Supports 256-color mode, special effects, and handles blank/pre-colored text.
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

// Height returns the terminal height in rows.
// Falls back to 24 rows if terminal size detection fails.
func Height() int {
	ws, err := getWinsize()
	if err != nil {
		// Fallback to reasonable default if terminal size unavailable
		return 24
	}
	if ws.Row == 0 {
		return 24
	}
	return int(ws.Row)
}

// Width returns the terminal width in columns.
// Falls back to 80 columns if terminal size detection fails.
func Width() int {
	ws, err := getWinsize()
	if err != nil {
		// Fallback to reasonable default if terminal size unavailable
		return 80
	}
	if ws.Col == 0 {
		return 80
	}
	return int(ws.Col)
}
