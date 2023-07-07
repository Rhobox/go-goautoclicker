package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"github.com/go-vgo/robotgo"
	"golang.design/x/hotkey"
	"log"
	"time"
)

var (
	version                    = "Version 1.0.5"
	clickingPeriodInput  int32 = 1000
	clickingPeriod       int32 = 1000
	isAutoclickerRunning       = false
	isLeftClicking             = true
	isRightClicking            = false
	isSpacePressing            = false
	isCtrlPressing             = false
	isFPressing                = false
	isEPressing                = false
	isCPressing                = false
	isAltTabPressing           = false
	stopClicking               = make(chan bool)
	inputTask                  = make(chan inputParameters, 10)
	isInputTaskRunning         = false
)

type inputParameters struct {
	isLeftClicking   bool
	isRightClicking  bool
	isSpacePressing  bool
	isCtrlPressing   bool
	isFPressing      bool
	isEPressing      bool
	isCPressing      bool
	isAltTabPressing bool
}

func setClickPeriod() {
	clickingPeriod = clickingPeriodInput
	g.Update()
}

func runningCheckbox() {
	isAutoclickerRunning = !isAutoclickerRunning
	return
}

func loop() {
	g.SingleWindow().Layout(
		g.Label(version),
		g.Row(
			g.InputInt(&clickingPeriodInput).Label("Clicking period: ").Size(60),
			g.Button("Set click time.").OnClick(setClickPeriod),
		),
		g.Row(
			g.Label(fmt.Sprintf("%vms between clicks.", clickingPeriod)),
		),
		g.Checkbox("Press Control+Shift+S to enable or disable clicking!", &isInputTaskRunning).OnChange(runningCheckbox),
		g.Checkbox("Left click!", &isLeftClicking),
		g.Checkbox("Right click, too!!", &isRightClicking),
		g.Checkbox("Click spacebar!", &isSpacePressing),
		g.Checkbox("Click F!", &isFPressing),
		g.Checkbox("Click E!", &isEPressing),
		g.Checkbox("Click C!", &isCPressing),
		g.Checkbox("Click Ctrl!", &isCtrlPressing),
		g.Label("If you're clicking slower than 250:"),
		g.Checkbox("Alt Tab!", &isAltTabPressing),
	)
}

func registerHotkey() {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyS)
	if err := hk.Register(); err != nil {
		log.Fatalf("Hotkey: Failed to register hotkey: %v", err)
	}

	for range hk.Keydown() {
		fmt.Println("Hotkey was pressed!")
		isAutoclickerRunning = !isAutoclickerRunning
		if isAutoclickerRunning {
			go clickStuffWithPeriod()
		} else {
			isInputTaskRunning = false
			stopClicking <- isAutoclickerRunning
		}
		g.Update()
	}
}

func clickStuffWithPeriod() {
	if !isInputTaskRunning {
		isInputTaskRunning = true
		go inputLoop()
	}
	ticker := time.NewTicker(time.Duration(clickingPeriod) * time.Millisecond)
	for {
		if !isAutoclickerRunning {
			return
		}
		select {
		case <-stopClicking:
			fmt.Println("Stopped!")
			return
		case <-ticker.C:
			fmt.Println("Adding to queue...")
			task := inputParameters{
				isLeftClicking:   isLeftClicking,
				isRightClicking:  isRightClicking,
				isSpacePressing:  isSpacePressing,
				isCtrlPressing:   isCtrlPressing,
				isFPressing:      isFPressing,
				isEPressing:      isEPressing,
				isCPressing:      isCPressing,
				isAltTabPressing: isAltTabPressing,
			}
			inputTask <- task
		}
	}

}

func inputLoop() {
	isInputTaskRunning = true
	for {
		if !isInputTaskRunning {
			fmt.Println("Input task not running.")
			return
		}
		select {
		case <-stopClicking:
			fmt.Println("Stopped the input loop.")
			isInputTaskRunning = false
			return
		case inputs := <-inputTask:
			go performInputFunctions(inputs)
		}
	}
}

func performInputFunctions(input inputParameters) {
	if input.isLeftClicking {
		robotgo.Click()
	}
	if isRightClicking {
		robotgo.Click("right")
	}
	if isFPressing {
		robotgo.KeyTap("f")
	}
	if isEPressing {
		robotgo.KeyTap("e")
	}
	if isCPressing {
		robotgo.KeyTap("c")
	}
	if isSpacePressing {
		robotgo.KeyTap("space")
	}
	if isCtrlPressing {
		robotgo.KeyDown("ctrl")
		time.Sleep(25 * time.Millisecond)
		robotgo.KeyUp("ctrl")
	}
	if isAltTabPressing && clickingPeriod > 250 {
		time.Sleep(1 * time.Millisecond)
		robotgo.KeyTap("tab", "alt")
	}
}

func populateKeys() {
	var keyStrings []string
	for key, _ := range robotgo.Keycode {
		keyStrings = append(keyStrings, key)
	}
}

func main() {
	wnd := g.NewMasterWindow("Go Gadget Autoclicker", 400, 300, 0)

	go registerHotkey()

	wnd.Run(loop)
}
