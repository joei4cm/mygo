package ui

import (
	"image/color"
	"testing"
	"time"

	"github.com/egoist/mygo/internal/media"
)

// The settings of a player of a system's engine set it and the engine,
// each clamped where the settings have a range.
func TestVideoPlayerSettings(t *testing.T) {
	e := &fakeEngine{colors: []color.RGBA{videoRed}, shown: -1, state: media.Ready}
	p := NewEngineVideoPlayer(e)
	visible := true
	tt, _ := videoTester(p, &visible)
	// The engine's states answer Ended and Buffering, read on a frame.
	if p.Buffering() || p.Ended() {
		t.Errorf("ready: Buffering %v, Ended %v", p.Buffering(), p.Ended())
	}
	e.setState(media.Buffering)
	tt.Frame()
	if !p.Buffering() {
		t.Error("Buffering did not follow the engine")
	}
	e.setState(media.Ended)
	tt.Frame()
	if !p.Ended() {
		t.Error("Ended did not follow the engine")
	}
	if p.Buffering() {
		t.Error("Ended reports Buffering")
	}
	e.setState(media.Ready)
	tt.Frame()

	if p.Loop() {
		t.Error("loop without a SetLoop")
	}
	if p.SetLoop(true); !p.Loop() || !e.loop {
		t.Errorf("SetLoop(true): player %v, engine %v", p.Loop(), e.loop)
	}
	p.SetRate(2)
	if p.Rate() != 2 {
		t.Errorf("Rate = %v, want 2", p.Rate())
	}
	// A rate out of range is 1.
	p.SetRate(0)
	if p.Rate() != 1 {
		t.Errorf("SetRate(0): Rate = %v, want 1", p.Rate())
	}
	p.SetVolume(0.5)
	if p.Volume() != 0.5 {
		t.Errorf("Volume = %v, want 0.5", p.Volume())
	}
	// The volume clamps to 0..1.
	p.SetVolume(1.5)
	if p.Volume() != 1 {
		t.Errorf("SetVolume(1.5): Volume = %v, want 1", p.Volume())
	}
	p.SetMuted(true)
	if !p.Muted() || !e.muted {
		t.Errorf("SetMuted(true): player %v, engine %v", p.Muted(), e.muted)
	}
	p.SetMuted(false)
	if p.Muted() || e.muted {
		t.Errorf("SetMuted(false): player %v, engine %v", p.Muted(), e.muted)
	}
}

// A player of a system's engine pauses, seeks and reads its position
// from it, and does none after it closed.
func TestVideoPlayerPausePositionSeek(t *testing.T) {
	e := &fakeEngine{colors: []color.RGBA{videoRed, videoGreen}, shown: -1, state: media.Ready}
	p := NewEngineVideoPlayer(e)
	visible := true
	tt, _ := videoTester(p, &visible)
	// Shown, it plays.
	if !p.Playing() {
		t.Error("not playing with a Video showing it")
	}
	// Pause asks the engine to pause.
	p.Pause()
	if p.Playing() || e.playing {
		t.Errorf("paused: player %v, engine %v", p.Playing(), e.playing)
	}
	p.Play()
	if !p.Playing() {
		t.Error("not playing again")
	}
	// The position is the engine's, and a seek goes to it, clamped at 0.
	e.base = 3 * time.Second
	if p.Position() != 3*time.Second {
		t.Errorf("Position = %v, want 3s", p.Position())
	}
	p.Seek(-time.Second)
	if got := e.seeks[len(e.seeks)-1]; got != 0 {
		t.Errorf("Seek(-1s) = %v, want 0", got)
	}
	p.Seek(time.Second)
	if got := e.seeks[len(e.seeks)-1]; got != time.Second {
		t.Errorf("Seek(1s) = %v", got)
	}
	// A seek from the end goes on playing.
	e.setState(media.Ended)
	tt.Frame()
	p.Seek(2 * time.Second)
	if p.Ended() {
		t.Error("the video stands Ended after a seek")
	}
	// Closed, it does none of them.
	p.Close()
	plays := e.plays
	p.Pause()
	p.Play()
	p.Seek(4 * time.Second)
	if p.Position() != 0 {
		t.Errorf("a closed player reports %v", p.Position())
	}
	if e.plays != plays {
		t.Error("a closed player played")
	}
}
