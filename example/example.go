package main

import (
	"os"
	"strconv"
	"time"

	"github.com/common-nighthawk/go-figure"
	ui "github.com/gregorpirolt/animaterm"
)

func main() {

	myUI := ui.CreateUI()
	myUI.ClearScreen()
	myUI.SetBoarder(0)

	var duration int64 = 800

	ch, wg := myUI.StartDrawLoop(70)
	iamASCII := figure.NewFigure("Staging", "standard", true)

	myUI.MoveElement(
		ui.CreatePos(70, 15),
		ui.CreatePos(35, 15),
		iamASCII.String(),
		ui.COLORPATTERNMEADOWS1,
		ui.Animation{
			AnimationType: ui.Ikea,
			Duration:      2400,
		})

	close(ch)
	wg.Wait()
	time.Sleep(time.Duration(500) * time.Millisecond)
	os.Exit(0)
	// for i := 10; i <= 100; i += 10 {
	// 	go myUI.DrawPattern(ui.CreatePos(0, i-10), i, "█\n█\n", ui.COLORPATTERNNEON1,
	// 		ui.Animation{
	// 			Duration:      int64(i) * 10,
	// 			AnimationType: ui.Ikea,
	// 			GradientV:     true,
	// 			GradientH:     true,
	// 			Direction:     ui.Right,
	// 		})
	// 	myUI.DrawPattern(ui.CreatePos(100, i-8), i, "█\n█\n", ui.COLORPATTERNNEON1,
	// 		ui.Animation{
	// 			Duration:      int64(i) * 10,
	// 			AnimationType: ui.Ikea,
	// 			GradientV:     true,
	// 			GradientH:     true,
	// 			Direction:     ui.Left,
	// 		})
	// 	myUI.DrawElement(ui.CreatePos(50, i-6), strconv.Itoa(i), ui.COLORPATTERNNEON1)
	// 	time.Sleep(time.Duration(100) * time.Millisecond)
	// }
	// close(ch)
	// wg.Wait()
	// os.Exit(0)

	// myUI.ClearScreen()

	for i := 10; i <= 100; i += 10 {
		go myUI.DrawPattern(ui.CreatePos(i-10, 0), 50, "█", ui.COLORPATTERNNEON1,
			ui.Animation{
				Duration:      int64(i) * 10,
				AnimationType: ui.Ikea,
				GradientV:     true,
				GradientH:     true,
				Direction:     ui.Down,
			})
		go myUI.DrawPattern(ui.CreatePos(i-8, 100), 50, "█", ui.COLORPATTERNNEON1,
			ui.Animation{
				Duration:      int64(i) * 10,
				AnimationType: ui.Ikea,
				GradientV:     true,
				GradientH:     true,
				Direction:     ui.Up,
			})
		myUI.DrawElement(ui.CreatePos(i-12, 0), strconv.Itoa(i), ui.COLORPATTERNNEON1)
		myUI.DrawElement(ui.CreatePos(i-12, 50), strconv.Itoa(i), ui.COLORPATTERNNEON1)
		myUI.DrawElement(ui.CreatePos(i-12, 100), strconv.Itoa(i), ui.COLORPATTERNNEON1)
	}
	time.Sleep(time.Duration(1000) * time.Millisecond)

	// // left vertical column
	// go myUI.DrawPattern(ui.CreatePos(10, 10), 80, "|\n|\n|\n", ui.COLORPATTERNNEON1,
	// 	ui.Animation{
	// 		Duration:      0,
	// 		AnimationType: ui.Ikea,
	// 		GradientV:     true,
	// 		GradientH:     true,
	// 		Direction:     ui.Down,
	// 	})

	// // right vertical column
	// myUI.DrawPattern(ui.CreatePos(80, 90), 80, "|\n|\n|\n", ui.COLORPATTERNNEON1,
	// 	ui.Animation{
	// 		Duration:      1800,
	// 		AnimationType: ui.Ikea,
	// 		GradientV:     true,
	// 		GradientH:     true,
	// 		Direction:     ui.Up,
	// 	})

	// top bar going left to right
	go myUI.DrawPattern(ui.CreatePos(0, 10), 100, "█\n█\n", ui.COLORPATTERNNEON1,
		ui.Animation{
			Duration:      10,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	// 2nd to top bar going left to right
	myUI.DrawPattern(ui.CreatePos(0, 30), 100, "█\n█\n", ui.COLORPATTERNNEON1,
		ui.Animation{
			Duration:      1900,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	// centoer bar going left to right
	myUI.DrawPattern(ui.CreatePos(100, 50), 100, "█\n█\n█\n█\n█\n", ui.COLORPATTERNMEADOWS1,
		ui.Animation{
			Duration:      1600,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Left,
		})

	// top bar overwriting previous bars (left to right)
	myUI.DrawPattern(
		ui.CreatePos(0, 5),
		100,
		"█\n█\n█\n█\n█\n",
		ui.COLORPATTERNLIME,
		ui.Animation{
			Duration:      600,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	// ascii font going from center to left
	myCLI := figure.NewFigure("animaterm", "standard", true)
	go myUI.MoveElement(ui.CreatePos(40, 5), ui.CreatePos(0, 5), myCLI.String(), ui.COLORPATTERNLIME,
		ui.Animation{
			Duration:      duration,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	myUI.DrawPattern(ui.CreatePos(50, 5), 50, "█\n█\n█\n█\n█\n", ui.COLORPATTERNPASTEL,
		ui.Animation{
			Duration:      1600,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	animation := figure.NewFigure("a go animation", "standard", true)
	go myUI.MoveElement(ui.CreatePos(27, 65), ui.CreatePos(14, 65), animation.String(), ui.COLORPATTERNLIME,
		ui.Animation{
			Duration:      800,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	framework := figure.NewFigure("framework", "standard", true)
	myUI.MoveElement(ui.CreatePos(15, 80), ui.CreatePos(40, 80), framework.String(), ui.COLORPATTERNGREY,
		ui.Animation{
			Duration:      1600,
			AnimationType: ui.Ikea,
			GradientV:     true,
			GradientH:     true,
			Direction:     ui.Right,
		})

	close(ch)
	wg.Wait()

}
