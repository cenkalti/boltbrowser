// +build !windows

package boltbrowser

import (
	"os"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/nsf/termbox-go"
)

func Run(db *bolt.DB, readOnly bool) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)

	memBolt := NewModel(db, readOnly)
	style := DefaultStyle()

	screens := defaultScreensForData(memBolt)
	displayScreen := screens[BrowserScreenIndex]
	layoutAndDrawScreen(displayScreen, style)
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyCtrlZ {
				process, _ := os.FindProcess(os.Getpid())
				termbox.Close()
				process.Signal(syscall.SIGSTOP)
				termbox.Init()
			}
			newScreenIndex := displayScreen.handleKeyEvent(event)
			if newScreenIndex < len(screens) {
				displayScreen = screens[newScreenIndex]
				layoutAndDrawScreen(displayScreen, style)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize {
			layoutAndDrawScreen(displayScreen, style)
		}
	}
}
