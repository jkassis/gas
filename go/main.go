// Author: Jeremy Kassis (jkassis@gmail.com).
//
// A frogger game.

package main

import (
	"fmt"
	"frogger/gas"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func CHECK(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// main let's go
func main() {
	rand.Seed(time.Now().UnixNano())

	CHECK(gas.Init())
	defer gas.Destroy()

	view, err := gas.MakeView(800, 600, "Frogger")
	CHECK(err)
	defer view.Destroy()

	s, err := gas.MakeStage(view)
	CHECK(err)
	s.BGColor = sdl.Color{R: 0x01, G: 0xb3, B: 0x35, A: 0xff}

	montserrat96, err := view.FontLoad("fonts/Montserrat-Regular.ttf", 96)
	CHECK(err)
	montserrat48, err := view.FontLoad("fonts/Montserrat-Regular.ttf", 48)
	CHECK(err)

	intro := func() chan struct{} {
		done := make(chan struct{})

		// the order added to the canvas established the z-layer
		bg, _ := s.Root.Spawn("img/bg.png")
		heart1, _ := s.Root.Spawn("img/heart1.png")
		heart2, _ := s.Root.Spawn("img/heart2.png")
		heart2.Exit()
		heart3, _ := s.Root.Spawn("img/heart3.png")
		heart3.Exit()
		frog, _ := s.Root.Spawn("img/frog.png")
		credit, _ := s.Root.Spawn("")
		title, _ := s.Root.Spawn("")

		// bg
		bg.Move(400, 300)

		// title
		title.TxtFillOut("Frogger", gas.SDLC(0x00ff00ff), montserrat96, 4, gas.SDLC(0x333333ff))
		title.Scale = .5
		title.Move(800, 300).
			MoveTo(400, 300, 2*time.Second, gas.EaseInOutSin)

		// frog
		frog.Scale = .05
		frog.
			Move(0, 200).
			MoveTo(120, 300, 2*time.Second, gas.EaseInOutSin).
			Then(func(d *gas.Dob) {
				// move and zoom
				d.MoveTo(300, 120, 2*time.Second, nil)
				d.ZoomTo(4, 2*time.Second, nil).
					Then(func(d *gas.Dob) {
						// move and zoom again. note how these race to Exit
						d.ZoomTo(.25, 3*time.Second, gas.EaseInOutSin).Exit()
						d.MoveTo(330, 280, 3*time.Second, nil)
					})
			})

		// credit
		credit.TxtFillOut("jkassis ©2023", gas.SDLC(0xffff33dd), montserrat48, 2, gas.SDLC(0x003300dd))
		credit.Zoom(.01)
		credit.Move(533, 400)

		// hearts
		heart1.Scale = .1
		heart2.Scale = .1
		heart3.Scale = .2
		heart1.
			Move(0, 200).
			MoveTo(120, 300, 2*time.Second, gas.EaseInOutSin).
			MoveTo(533, 400, 3*time.Second, gas.EaseInOutSin).
			Then(func(d *gas.Dob) {
				credit.
					ZoomTo(1, 3*time.Second, gas.EaseInOutSin).Then(func(d *gas.Dob) {
					// note how we trigger this title anim when the logo anim completes
					title.
						ZoomTo(2, 200*time.Millisecond, nil).
						ZoomTo(1, 400*time.Millisecond, nil).
						Then(func(d *gas.Dob) {
							close(done)
						})
				})
			})

		heart1.Emit(heart1, 20, 500*time.Millisecond, 3*time.Second, nil, gas.EaseInOutSinInv,
			func(d *gas.Dob) {
				if rand.Intn(100) < 25 {
					d.Texture = heart3.Texture
				}
				spinDst := -90 + 180*rand.Float64()
				spinDuration := time.Second + gas.RandDuration(4*time.Second)
				d.SpinTo(spinDst, spinDuration, nil)

				x := float32(view.W) * rand.Float32()
				y := float32(view.H) * rand.Float32()
				moveDur := 2*time.Second + gas.RandDuration(4*time.Second)
				d.MoveTo(x, y, moveDur, nil).Exit()
			})

		return done
	}

	go func() {
		for {
			fmt.Println("looping...")
			<-intro()
			time.Sleep(time.Second)
			s.Root.DobsClear()
			s.Root.AnSetClear()
		}
	}()

	s.Play(30)
}
