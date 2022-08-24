package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/go-vgo/robotgo"
	"golang.design/x/hotkey"
	"time"
)

var (
	clickingPeriodInput  int32 = 1000
	clickingPeriod       int32 = 1000
	isAutoclickerRunning       = false
	stopClicking               = make(chan bool)
)

func setClickPeriod() {
	clickingPeriod = clickingPeriodInput
	g.Update()
}

func loop() {
	g.SingleWindow().Layout(
		g.Label("Hmm.."),
		g.Row(
			g.InputInt(&clickingPeriodInput).Label("Clicking period: ").Size(60),
			g.Button("Set click time.").OnClick(setClickPeriod),
		),
		g.Row(
			g.Label(fmt.Sprintf("%vms between clicks.", clickingPeriod)),
		),
		g.Checkbox("Press Control+Shift+S to enable or disable clicking!", &isAutoclickerRunning),
	)
}

func registerHotkey() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	if err := hk.Register(); err != nil {
		panic("Hotkey registration failed.")
	}

	for range hk.Keydown() {
		fmt.Println("Hotkey was pressed!")
		isAutoclickerRunning = !isAutoclickerRunning
		if isAutoclickerRunning {
			go clickLeftMouseWithPeriod()
		} else {
			stopClicking <- isAutoclickerRunning
		}
		g.Update()
	}
}

func clickLeftMouseWithPeriod() {
	for {
		select {
		case <-stopClicking:
			fmt.Println("Stopped!")
			return
		default:
			go robotgo.Click()
			time.Sleep(time.Duration(clickingPeriod) * time.Millisecond)
		}
	}

}

func main() {
	wnd := g.NewMasterWindow("Go Gadget Autoclicker", 400, 200, 0)

	go registerHotkey()

	wnd.Run(loop)
}
