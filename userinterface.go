package animaterm

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/mobile/exp/sprite/clock"
)

// UserInterface ...
type UserInterface struct {
	absBorderLeft   int
	absBorderRight  int
	absBorderTop    int
	absBorderBottom int
	pixels          [][]string
	pixelsMutex     sync.RWMutex
	dirtyRegions    [][]bool
	dirtyMutex      sync.RWMutex
	height          int
	width           int
	msPerFrame      int64
	frameMutex      sync.RWMutex
}

// CreateUI creates and initializes a new UserInterface instance.
// It automatically detects terminal dimensions and initializes the pixel buffer.
// If terminal size is below minimum requirements (130x33), it displays a warning
// but continues execution for compatibility.
func CreateUI() IUserInterface {
	minWidth := 130
	minHeight := 33
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if Width() < minWidth || Height() < minHeight {
		fr := &UserInterface{}
		_ = fr.ClearScreen()
		_, _ = pl("You should use the UI in a terminal with a resolution bigger than:")
		_, _ = pf("%v columns X %v rows\n", minWidth, minHeight)
		_, _ = pf("Your current resolution is %v columns X %v rows X\n", Color(strconv.Itoa(Width()), COLORPATTERNLIME), Color(strconv.Itoa(Height()), COLORPATTERNLIME))
		_, _ = pl("In- or decrease your terminal's zoom to fit the canvas onto your screen.")
		_, _ = pl("For optimal content presentation set your terminal into fullscreen mode.")
		// os.Exit(0)
	}
	ui := &UserInterface{
		absBorderLeft:   0,
		absBorderRight:  0,
		absBorderTop:    0,
		absBorderBottom: 0,
		msPerFrame:      320,
	}
	_ = ui.initPixels(Height(), Width())

	return ui
}

func (ui *UserInterface) initPixels(height int, width int) error {
	ui.pixelsMutex.Lock()
	ui.dirtyMutex.Lock()
	defer ui.pixelsMutex.Unlock()
	defer ui.dirtyMutex.Unlock()

	ui.height = height
	ui.width = width
	// init pixels
	ui.pixels = make([][]string, height+1)
	ui.dirtyRegions = make([][]bool, height+1)
	for h := 0; h < height; h++ {
		ui.pixels[h] = make([]string, width+1)
		ui.dirtyRegions[h] = make([]bool, width+1)
		for w := 0; w < width; w++ {
			ui.pixels[h][w] = " "
			ui.dirtyRegions[h][w] = true // Initially mark all as dirty
		}
	}
	return nil
}

// StartDrawLoop initiates the main rendering loop in a separate goroutine.
// percentHeight specifies how much of the terminal height to use (0-100).
// Returns a channel to stop the loop and a WaitGroup to wait for cleanup.
// The loop automatically adjusts frame rate based on dirty regions for optimal performance.
func (ui *UserInterface) StartDrawLoop(percentHeight int) (chan int, *sync.WaitGroup) {
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan int)
	_, _ = p("\033[?25l")

	go ui.drawLoop(Height()*percentHeight/100, Width(), ch, &wg)

	return ch, &wg
}

// Draw ...
func (ui *UserInterface) drawLoop(height int, width int, ch chan int, wg *sync.WaitGroup) {
	_ = time.Now()
	_ = time.Duration(5)
	lineBuffer := ""
	screenBuffer := "" // make([]string, height)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				wg.Done()
				_, _ = p("\033[?25h")
				return
			}
		default:
			_ = time.Now()

			// Check if any regions are dirty before rebuilding buffer
			ui.dirtyMutex.RLock()
			hasDirty := false
			for h := 0; h < height && !hasDirty; h++ {
				for w := 0; w < width && !hasDirty; w++ {
					if ui.dirtyRegions[h][w] {
						hasDirty = true
					}
				}
			}
			ui.dirtyMutex.RUnlock()

			if hasDirty {
				oldbuff := screenBuffer
				screenBuffer = ""
				lastContentRow := height
				ui.pixelsMutex.RLock()
				for h := 0; h < height; h++ {
					for w := 0; w < width; w++ {
						lineBuffer += ui.pixels[h][w]
					}
					if strings.Count(lineBuffer, " ") == width {
						screenBuffer += "\033[1B"
					} else {
						screenBuffer += lineBuffer
						lastContentRow = h + 2
					}
					lineBuffer = ""
				}
				ui.pixelsMutex.RUnlock()

				if oldbuff != screenBuffer {
					_ = ui.moveCursorTo(0, 0)
					ui.frameMutex.Lock()
					ui.msPerFrame = 30
					ui.frameMutex.Unlock()
					_, _ = pl(screenBuffer)
					_ = ui.moveCursorTo(0, lastContentRow)
					ui.clearDirtyRegions()
				}
			} else {
				// No dirty regions, use slower frame rate
				ui.frameMutex.Lock()
				ui.msPerFrame = 320
				ui.frameMutex.Unlock()
			}
			_ = time.Since(time.Now())
			ui.frameMutex.RLock()
			frameRate := ui.msPerFrame
			ui.frameMutex.RUnlock()
			time.Sleep(time.Duration(frameRate) * time.Millisecond)
		}
	}
}

// SetBorderLeft ...
// Set a global border on left side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorderLeft(percent int) error {
	if percent < 0 || percent > 50 {
		return fmt.Errorf("border percent must be between 0 and 50, got %d", percent)
	}
	ui.absBorderLeft = Width() * percent / 100
	return nil
}

// SetBorderRight ...
// Set a global border on right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorderRight(percent int) error {
	if percent < 0 || percent > 50 {
		return fmt.Errorf("border percent must be between 0 and 50, got %d", percent)
	}
	ui.absBorderRight = Width() * percent / 100
	return nil
}

// SetBorderTop ...
// Set a global border on top of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorderTop(percent int) error {
	if percent < 0 || percent > 50 {
		return fmt.Errorf("border percent must be between 0 and 50, got %d", percent)
	}
	ui.absBorderTop = Height() * percent / 100
	return nil
}

// SetBorderBottom ...
// Set a global border on bottom of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorderBottom(percent int) error {
	if percent < 0 || percent > 50 {
		return fmt.Errorf("border percent must be between 0 and 50, got %d", percent)
	}
	ui.absBorderBottom = Height() * percent / 100
	return nil
}

// SetBorderSides ...
// Set a global border on left and right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorderSides(percent int) error {
	_ = ui.SetBorderLeft(percent)
	_ = ui.SetBorderRight(percent)
	return nil
}

// SetBorderTopBottom ...
// Set a global border on left and right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorderTopBottom(percent int) error {
	_ = ui.SetBorderTop(percent)
	_ = ui.SetBorderBottom(percent)
	return nil
}

// SetBorder ...
// Set a global border on left and right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBorder(percent int) error {
	_ = ui.SetBorderSides(percent)
	_ = ui.SetBorderTopBottom(percent)
	return nil
}

// DrawTable ...
func (ui *UserInterface) DrawTable(pos IRelativePosition, table [][]string, positions []int, colors []int) int {
	y := 0
	for _, s := range table {
		y1 := ui.DrawElementsHorizontal(pos, s, positions, colors)
		if y1 > y {
			y = y1
		}
		pos.IncrementOffset()
	}
	return y
}

// DrawElementsHorizontal ...
func (ui *UserInterface) DrawElementsHorizontal(pos IRelativePosition, texts []string, positions []int, colors []int) int {
	y, y1 := 0, 0
	newPos := CreatePos(pos.GetX(), pos.GetY())
	for k, s := range texts {
		newPos.SetOffset(pos.GetOffset() + int(k/len(positions)))
		y1 = ui.DrawElement(newPos.AddDistance(CreatePos(positions[k%len(positions)], 0)), s, colors[k%len(positions)])
		if y1 > y {
			y = y1
		}
	}
	return y
}

// DrawElement renders text at the specified position with the given color.
// Supports multi-line text and automatically handles line wrapping.
// Returns the Y coordinate of the last rendered line.
func (ui *UserInterface) DrawElement(pos IRelativePosition, text string, color int) int {
	x, y := 0, 0
	for k, line := range getLines(text, color == BLANK) {
		for l, c := range line {
			y = (ui.PercentToAbsoluteHeightInFrame(pos.GetY()) + pos.GetOffset() + ui.absBorderTop) % ui.height
			x = (ui.PercentToAbsoluteWidthInFrame(pos.GetX()) + l + ui.absBorderLeft) % ui.width

			ui.setPixel(x, y, Color(string(c), color+k))
		}
		pos.IncrementOffset()
	}
	return y - 1
}

// MoveElement animates text movement from startPos to endPos over the specified duration.
// Supports gradient effects and various animation curves (EaseIn, EaseOut, etc.).
// The animation blocks until completion.
func (ui *UserInterface) MoveElement(startPos IRelativePosition, endPos IRelativePosition, text string, color int, animation Animation) error {
	if startPos == nil || endPos == nil {
		return fmt.Errorf("start and end positions cannot be nil")
	}
	if text == "" {
		return fmt.Errorf("text cannot be empty")
	}
	if animation.Duration < 0 {
		return fmt.Errorf("animation duration cannot be negative")
	}

	ui.frameMutex.RLock()
	frameRate := ui.msPerFrame
	ui.frameMutex.RUnlock()
	frames := int(animation.Duration / frameRate)
	// currentPos := CreatePos(startPos.GetX(), startPos.GetY())
	// currentPos.ResetOffset()

	distance := startPos.DistanceTo(endPos)

	var factor float32 = 0
	for i := 0; i <= frames; i++ {
		// delete previous frame
		ui.DrawElement(startPos.AddDistance(distance.MultiplyWith(factor)), text, BLANK)
		factor = getAnimation(animation.AnimationType)(clock.Time(0), clock.Time(frames), clock.Time(i))

		if animation.GradientV || animation.GradientH {
			ui.DrawElement(startPos.AddDistance(distance.MultiplyWith(factor)), text, color+36*int(float32(5)*factor))
		} else {
			ui.DrawElement(startPos.AddDistance(distance.MultiplyWith(factor)), text, color)
		}

		ui.frameMutex.RLock()
		frameRate := ui.msPerFrame
		ui.frameMutex.RUnlock()
		time.Sleep(time.Duration(frameRate) * time.Millisecond)
	}
	return nil
}

// DrawPattern creates expanding patterns with animation support.
// expansion controls how far the pattern extends (0-200 percent).
// Supports directional expansion (Right, Left, Down, Up) and gradient effects.
// Returns the Y coordinate of the last drawn element.
func (ui *UserInterface) DrawPattern(startPos IRelativePosition, expansion int, text string, color int, animation Animation) int {
	if startPos == nil {
		return -1
	}
	if expansion < 0 || expansion > 200 {
		return -1
	}
	if text == "" {
		return -1
	}

	y := 0
	startAbsWidth := ui.PercentToAbsoluteWidthInFrame(startPos.GetX()) + ui.absBorderLeft
	startAbsHeight := ui.PercentToAbsoluteHeightInFrame(startPos.GetY()) + ui.absBorderTop + startPos.GetOffset()

	drawPixel := func(h int, w int, factorColor float32, expander []int) {
		basecolor := color
		if animation.GradientV {
			basecolor = color + int(float32(5)*factorColor)
		}
		for k, line := range getLines(text, false) {
			expH := expander[0] * k
			expW := expander[1] * k
			line = string([]rune(line)[0])
			yPos := (h + expH) % ui.height
			xPos := (w + expW) % ui.width
			var coloredLine string
			if animation.GradientH {
				coloredLine = Color(line, (basecolor+(k*36))%255)
			} else {
				coloredLine = Color(line, basecolor)
			}
			ui.setPixel(xPos, yPos, coloredLine)
			if h+expH > y {
				y = h + expH
			}
		}
	}

	switchDir := func(expInPercent int) {
		expAbs := 0
		var expDir1, expDir2 []int

		switch animation.Direction {
		case 0:
			// right
			expAbs = ui.PercentToAbsoluteWidthInFrame(expInPercent)
			expDir1 = []int{1, 0}
			expDir2 = []int{0, 1}
		case 1:
			// left()
			expAbs = ui.PercentToAbsoluteWidthInFrame(expInPercent)
			expDir1 = []int{1, 0}
			expDir2 = []int{0, -1}
		case 2:
			// down
			expAbs = ui.PercentToAbsoluteHeightInFrame(expInPercent)
			expDir1 = []int{0, 1}
			expDir2 = []int{1, 0}
		case 3:
			// up
			expAbs = ui.PercentToAbsoluteHeightInFrame(expInPercent)
			expDir1 = []int{0, 1}
			expDir2 = []int{-1, 0}
		default:
			return
		}

		for counter := 0; counter <= expAbs; counter++ {
			factorColor := getAnimation(animation.AnimationType)(clock.Time(0), clock.Time(expAbs), clock.Time(counter))
			drawPixel(startAbsHeight+expDir2[0]*counter, startAbsWidth+expDir2[1]*counter, factorColor, expDir1)
		}

	}

	if animation.Duration > 0 {
		ui.frameMutex.RLock()
		frameRate := ui.msPerFrame
		ui.frameMutex.RUnlock()
		frames := int(animation.Duration / frameRate)
		// animation
		for i := 0; i <= frames; i++ {
			factor := getAnimation(animation.AnimationType)(clock.Time(0), clock.Time(frames), clock.Time(i))
			currentExpansionPercent := int(float32(expansion)*factor + 0.5)
			switchDir(currentExpansionPercent)
			ui.frameMutex.RLock()
			frameRate := ui.msPerFrame
			ui.frameMutex.RUnlock()
			time.Sleep(time.Duration(frameRate) * time.Millisecond)
		}
	} else {
		switchDir(expansion)
		ui.frameMutex.RLock()
		frameRate := ui.msPerFrame
		ui.frameMutex.RUnlock()
		time.Sleep(time.Duration(frameRate) * time.Millisecond)
	}

	// time.Sleep(time.Duration(msPerFrame) * time.Millisecond)
	return y - 1
}

func getLines(multilineText string, replaceWithBlanks bool) []string {
	var lines = []string{}
	for _, line := range strings.Split(strings.TrimSuffix(multilineText, "\n"), "\n") {
		if replaceWithBlanks {
			line = strings.Repeat(" ", len(fmt.Sprintf("%v", line)))
		}
		lines = append(lines, line)
	}
	return lines
}

// PercentToAbsoluteWidth ...
func (ui *UserInterface) PercentToAbsoluteWidth(percent int) int {
	return Width() * percent / 100
}

// PercentToAbsoluteHeight ...
func (ui *UserInterface) PercentToAbsoluteHeight(percent int) int {
	return (Height() * percent / 100)
}

// GetAbsFrameWidth ...
func (ui *UserInterface) GetAbsFrameWidth() int {
	return Width() - ui.absBorderLeft - ui.absBorderRight
}

// GetAbsFrameHeight ...
func (ui *UserInterface) GetAbsFrameHeight() int {
	return Height() - ui.absBorderTop - ui.absBorderBottom
}

// PercentToAbsoluteWidthInFrame ...
func (ui *UserInterface) PercentToAbsoluteWidthInFrame(percent int) int {
	return (ui.GetAbsFrameWidth() * percent / 100)
}

// PercentToAbsoluteHeightInFrame ...
func (ui *UserInterface) PercentToAbsoluteHeightInFrame(percent int) int {
	return (ui.GetAbsFrameHeight() * percent / 100)
}

// PercentToAbsoluteWidthInFrame ...
func (ui *UserInterface) PercentToAbsoluteXPostion(percent int) int {
	return (ui.GetAbsFrameWidth() * percent / 100) + ui.absBorderLeft
}

// PercentToAbsoluteHeightInFrame ...
func (ui *UserInterface) PercentToAbsoluteYPostion(percent int) int {
	return (ui.GetAbsFrameHeight() * percent / 100) + ui.absBorderTop
}

// ClearScreen ...
func (ui *UserInterface) ClearScreen() error {
	_ = ui.initPixels(Height(), Width())
	cmd := exec.Command("clear", "cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	return nil
}

// setPixel safely sets a pixel and marks the region as dirty
func (ui *UserInterface) setPixel(x, y int, value string) {
	ui.pixelsMutex.Lock()
	ui.dirtyMutex.Lock()
	defer ui.pixelsMutex.Unlock()
	defer ui.dirtyMutex.Unlock()

	if y >= 0 && y < len(ui.pixels) && x >= 0 && x < len(ui.pixels[y]) {
		if ui.pixels[y][x] != value {
			ui.pixels[y][x] = value
			ui.dirtyRegions[y][x] = true
		}
	}
}

// clearDirtyRegions resets all dirty flags
func (ui *UserInterface) clearDirtyRegions() {
	ui.dirtyMutex.Lock()
	defer ui.dirtyMutex.Unlock()

	for h := 0; h < ui.height; h++ {
		for w := 0; w < ui.width; w++ {
			ui.dirtyRegions[h][w] = false
		}
	}
}

func (ui *UserInterface) moveCursorTo(absX int, absY int) error {
	_, _ = p("\033[" + strconv.Itoa(absY) + ";" + strconv.Itoa(absX) + "H")
	return nil
}
