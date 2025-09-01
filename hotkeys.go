package main

import (
	"fmt"

	hook "github.com/robotn/gohook"
	rt "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) listenHotkeys() {
	fmt.Println("Listening for hotkeys...")

	hook.Register(hook.KeyDown, []string{"ctrl", "o"}, func(e hook.Event) {
		rt.WindowShow(a.ctx)
		rt.WindowSetAlwaysOnTop(a.ctx, true)
		rt.WindowSetAlwaysOnTop(a.ctx, false)
		rt.EventsEmit(a.ctx, "focusCapture")
	})

	hook.Register(hook.KeyDown, []string{"esc"}, func(e hook.Event) {
		if err := a.SendFile(); err != nil {
			fmt.Printf("failed to send file: %v\n", err)
		}
		rt.WindowHide(a.ctx)
	})

	hook.Register(hook.KeyDown, []string{"ctrl", "space"}, func(e hook.Event) {
		rt.WindowShow(a.ctx)
		rt.WindowSetAlwaysOnTop(a.ctx, true)
		rt.WindowSetAlwaysOnTop(a.ctx, false)
		rt.EventsEmit(a.ctx, "focusSearch")
	})

	s := hook.Start()
	select {
	case <-hook.Process(s):
	case <-a.ctx.Done():
		hook.End()
	}
}
