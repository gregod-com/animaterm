package animaterm

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	CB "golang.org/x/mobile/exp/sprite/clock"
)

// UserInterface ...
type UserInterface struct {
	absBorderLeft   int
	absBorderRight  int
	absBorderTop    int
	absBorderBottom int
}

// CreateUI return an instance of UserInterface
// There has to be a delay of 100ms at startup, since the docker environemt needs a few moments
// to get the terminal width and hight from its host.
// Todo: remove the sleep as soon as other boot-sequences take longer that 100ms anyways
func CreateUI() IUserInterface {
	minWidth := 130
	minHeight := 33
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if Width() < minWidth || Height() < minHeight {
		fr := &UserInterface{}
		fr.ClearScreen()
		pl("You should use the UI in a terminal with a resolution bigger than:")
		pf("%v columns X %v rows\n", minWidth, minHeight)
		pf("Your current resolution is %v columns X %v rows X\n", Color(strconv.Itoa(Width()), COLORPATTERNLIME), Color(strconv.Itoa(Height()), COLORPATTERNLIME))
		pl("In- or decrease your terminal's zoom to fit the canvas onto your screen.")
		pl("For optimal content presentation set your terminal into fullscreen mode.")
		// os.Exit(0)
	}
	ui := &UserInterface{
		absBorderLeft:   0,
		absBorderRight:  0,
		absBorderTop:    0,
		absBorderBottom: 0,
	}
	ui.initPixels(Height(), Width())

	return ui
}

var msPerFrame int64 = 320
var pixels = make([][]string, 0)

func (ui *UserInterface) initPixels(height int, width int) error {
	// init pixels
	pixels = make([][]string, height+1)
	for h := 0; h < height; h++ {
		pixels[h] = make([]string, width+1)
		for w := 0; w < width; w++ {
			pixels[h][w] = " "
		}
	}
	return nil
}

// StartDrawLoop ...
func (ui *UserInterface) StartDrawLoop(percentHeight int) (chan int, *sync.WaitGroup) {
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan int)
	p("\033[?25l")

	go ui.drawLoop(Height()*percentHeight/100, Width(), ch, &wg)

	return ch, &wg
}

// Draw ...
func (ui *UserInterface) drawLoop(height int, width int, ch chan int, wg *sync.WaitGroup) {
	start := time.Now()
	took := time.Duration(5)
	lineBuffer := ""
	screenBuffer := "" // make([]string, height)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				wg.Done()
				p("\033[?25h")
				return
			}
		default:
			start = time.Now()

			// lineBuffer = ""
			oldbuff := screenBuffer
			screenBuffer = ""
			lastContentRow := height
			for h := 0; h < height; h++ {
				for w := 0; w < width; w++ {
					lineBuffer += pixels[h][w]
				}
				if strings.Count(lineBuffer, " ") == width {
					screenBuffer += "\033[1B"
				} else {
					screenBuffer += lineBuffer
					lastContentRow = h + 2
				}
				lineBuffer = ""

			}
			if oldbuff != screenBuffer {
				ui.moveCursorTo(0, 0)
				msPerFrame = 30
				pl(screenBuffer)
				ui.moveCursorTo(0, lastContentRow)
			}
			took = time.Since(start)
			time.Sleep(time.Duration(msPerFrame)*time.Millisecond - took)
		}
	}
}

// SetBoarderLeft ...
// Set a global boarder on left side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarderLeft(percent int) error {
	ui.absBorderLeft = Width() * percent / 100
	return nil
}

// SetBoarderRight ...
// Set a global boarder on right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarderRight(percent int) error {
	ui.absBorderRight = Width() * percent / 100
	return nil
}

// SetBoarderTop ...
// Set a global boarder on top of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarderTop(percent int) error {
	ui.absBorderTop = Height() * percent / 100
	return nil
}

// SetBoarderBottom ...
// Set a global boarder on bottom of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarderBottom(percent int) error {
	ui.absBorderBottom = Height() * percent / 100
	return nil
}

// SetBoarderSides ...
// Set a global boarder on left and right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarderSides(percent int) error {
	ui.SetBoarderLeft(percent)
	ui.SetBoarderRight(percent)
	return nil
}

// SetBoarderTopBottom ...
// Set a global boarder on left and right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarderTopBottom(percent int) error {
	ui.SetBoarderTop(percent)
	ui.SetBoarderBottom(percent)
	return nil
}

// SetBoarder ...
// Set a global boarder on left and right side of screen that forces
// all elements to be printed outside it's boundaries
func (ui *UserInterface) SetBoarder(percent int) error {
	ui.SetBoarderSides(percent)
	ui.SetBoarderTopBottom(percent)
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

// DrawElement ...
// Render a text box at given coordinates
func (ui *UserInterface) DrawElement(pos IRelativePosition, text string, color int) int {
	x, y := 0, 0
	for k, line := range getLines(text, color == BLANK) {
		for l, c := range line {
			y = (ui.PercentToAbsoluteHeightInFrame(pos.GetY()) + pos.GetOffset() + ui.absBorderTop) % Height()
			x = (ui.PercentToAbsoluteWidthInFrame(pos.GetX()) + l + ui.absBorderLeft) % Width()

			pixels[y][x] = Color(string(c), color+k)
		}
		pos.IncrementOffset()
	}
	return y - 1
}

// MoveElement ...
func (ui *UserInterface) MoveElement(startPos IRelativePosition, endPos IRelativePosition, text string, color int, animation Animation) error {

	frames := int(animation.Duration / msPerFrame)
	// currentPos := CreatePos(startPos.GetX(), startPos.GetY())
	// currentPos.ResetOffset()

	distance := startPos.DistanceTo(endPos)

	var factor float32 = 0
	for i := 0; i <= frames; i++ {
		// delete previous frame
		ui.DrawElement(startPos.AddDistance(distance.MultiplyWith(factor)), text, BLANK)
		factor = getAnimation(animation.AnimationType)(CB.Time(0), CB.Time(frames), CB.Time(i))

		if animation.GradientV || animation.GradientH {
			ui.DrawElement(startPos.AddDistance(distance.MultiplyWith(factor)), text, color+36*int(float32(5)*factor))
		} else {
			ui.DrawElement(startPos.AddDistance(distance.MultiplyWith(factor)), text, color)
		}

		time.Sleep(time.Duration(msPerFrame) * time.Millisecond)
	}
	return nil
}

// DrawPattern ...
func (ui *UserInterface) DrawPattern(startPos IRelativePosition, expansion int, text string, color int, animation Animation) int {

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
			if animation.GradientH {
				pixels[(h+expH)%Height()][w+expW%Height()] = Color(line, (basecolor+(k*36))%255)
			} else {
				pixels[h+expH][w+expW] = Color(line, basecolor)
			}
			if h+expH > y {
				y = h + expH
			}
		}
	}

	switchDir := func(expInPercent int) {
		expAbs := 0
		expDir1 := []int{0, 0}
		expDir2 := []int{0, 0}

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
			factorColor := getAnimation(animation.AnimationType)(CB.Time(0), CB.Time(expAbs), CB.Time(counter))
			drawPixel(startAbsHeight+expDir2[0]*counter, startAbsWidth+expDir2[1]*counter, factorColor, expDir1)
		}

	}

	if animation.Duration > 0 {
		frames := int(animation.Duration / msPerFrame)
		// animation
		for i := 0; i <= frames; i++ {
			factor := getAnimation(animation.AnimationType)(CB.Time(0), CB.Time(frames), CB.Time(i))
			currentExpansionPercent := int(float32(expansion)*factor + 0.5)
			switchDir(currentExpansionPercent)
			time.Sleep(time.Duration(msPerFrame) * time.Millisecond)
		}
	} else {
		switchDir(expansion)
		time.Sleep(time.Duration(msPerFrame) * time.Millisecond)
	}

	// time.Sleep(time.Duration(msPerFrame) * time.Millisecond)
	return y - 1
}

func replaceCharAt(str string, replace string, index int) string {
	return str[:index] + replace + str[index+1:]
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
	ui.initPixels(Height(), Width())
	cmd := exec.Command("clear", "cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}

func (ui *UserInterface) moveCursorTo(absX int, absY int) error {
	p("\033[" + strconv.Itoa(absY) + ";" + strconv.Itoa(absX) + "H")
	return nil
}
