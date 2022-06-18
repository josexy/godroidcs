// Copyright [2021] [josexy]
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package window

import (
	"image"
	"image/color"
	"sync/atomic"

	"github.com/josexy/godroidcli/util"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	delayMs = 16
	scaleX  = 0.35
	scaleY  = 0.35
)

type Rect struct {
	Width  int
	Height int
}

type Window struct {
	Title      string
	DeviceSize Rect
	scaledSize Rect
	window     *sdl.Window
	surface    *sdl.Surface
	closed     int32
	imgBufChan chan image.Image
}

func NewWindow(title string, r Rect) (*Window, error) {
	sdl.SetHint(sdl.HINT_QUIT_ON_LAST_WINDOW_CLOSE, "1")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	scaleW := int32(float32(r.Width) * scaleX)
	scaleH := int32(float32(r.Height) * scaleY)

	util.Info("scaled window size: %dx%d", scaleW, scaleH)

	window, err := sdl.CreateWindow(title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		scaleW,
		scaleH,
		sdl.WINDOW_SHOWN)

	if err != nil {
		return nil, err
	}

	surface, err := window.GetSurface()
	if err != nil {
		return nil, err
	}

	w := &Window{
		Title:      title,
		DeviceSize: r,
		scaledSize: Rect{int(scaleW), int(scaleH)},
		window:     window,
		surface:    surface,
		imgBufChan: make(chan image.Image, 1024),
	}

	return w, nil
}

func (w *Window) Update(img image.Image) {
	if atomic.LoadInt32(&w.closed) == 1 {
		return
	}
	w.imgBufChan <- img
}

func (w *Window) IsReady() bool {
	return atomic.LoadInt32(&w.closed) == 0
}

func (w *Window) Run() {

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch ev := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.WindowEvent:
				if ev.Event == sdl.WINDOWEVENT_CLOSE {
					running = false
				}
			}
		}

		select {
		case img, ok := <-w.imgBufChan:
			if !ok {
				running = false
			} else {
				w.refreshSurface(img)
			}
		default:
			sdl.Delay(delayMs)
		}
	}
}

func (w *Window) refreshSurface(img image.Image) {
	srcb := img.Bounds()

	sx := float64(w.scaledSize.Width) / float64(srcb.Dx())
	sy := float64(w.scaledSize.Height) / float64(srcb.Dy())

	I.Scale(sx, sy).Transform(w.scaledSize.Width, w.scaledSize.Height, img, Bilinear,
		func(x, y int, c color.Color) {
			w.surface.Set(x, y, c)
		},
	)
	w.window.UpdateSurface()
}

func (w *Window) Destroy() {
	if atomic.LoadInt32(&w.closed) == 1 {
		return
	}
	atomic.StoreInt32(&w.closed, 1)
	close(w.imgBufChan)

	if w.surface != nil {
		w.surface.Free()
	}
	if w.window != nil {
		w.window.Destroy()
	}
	sdl.Quit()
	util.Info("Close window")
}
