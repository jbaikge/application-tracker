package main

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"log"
	"time"
)

var (
	appChan = make(chan App)
)

func getActiveApp(X *xgbutil.XUtil) (app App) {
	active, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		log.Fatal(err)
	}

	app.PID, err = ewmh.WmPidGet(X, active)
	if err != nil {
		log.Println("Couldn't get PID of window (0x%x)", active)
	}

	app.Name, err = ewmh.WmNameGet(X, active)
	if err != nil || len(app.Name) == 0 {
		app.Name, err = icccm.WmNameGet(X, active)
		// If we still can't find anything, give up.
		if err != nil || len(app.Name) == 0 {
			app.Name = "N/A"
		}
	}
	return
}

func sendNewNames(X *xgbutil.XUtil) {
	var lastName string
	for {
		app := getActiveApp(X)
		if app.Name != lastName {
			appChan <- app
		}
		lastName = app.Name
		time.Sleep(250 * time.Millisecond)
	}
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	go sendNewNames(X)

	go func() {
		for {
			select {
			case app := <-appChan:
				log.Printf("%5d - %-15s - %s\n", app.PID, app.ProcessName(), app.Name)
			}
		}
	}()

	select {}
}
