//go:build !darwin

package texture

import (
	"runtime"
	"testing"

	"github.com/egoist/mygo/internal/scene"
)

func TestRegistry(t *testing.T) {
	tt := New()
	if tt.ID() == 0 {
		t.Fatal("a texture has no ID")
	}
	if got := Lookup(tt.ID()); got != tt {
		t.Fatalf("Lookup(%d) = %v, want the texture", tt.ID(), got)
	}
	if got := Lookup(1 << 40); got != nil {
		t.Fatalf("Lookup of a made-up ID = %v, want nil", got)
	}
	id := tt.ID()
	tt = nil
	runtime.GC()
	if got := Lookup(id); got != nil {
		t.Fatalf("Lookup after the texture is gone = %v, want nil", got)
	}
}

func TestFrameBeforeAPresent(t *testing.T) {
	tt := New()
	if img := tt.Frame("w", func() {}); img != nil {
		t.Fatalf("Frame before any present = %v, want nil", img)
	}
	if w, h := tt.Size(); w != 0 || h != 0 {
		t.Fatalf("Size before any present = %vx%v, want 0x0", w, h)
	}
}

func TestPresentWakesWatchers(t *testing.T) {
	tt := New()
	calls := 0
	if img := tt.Frame("w", func() { calls++ }); img != nil {
		t.Fatalf("Frame before any present = %v, want nil", img)
	}
	tt.Present(0) // a nil surface shows nothing
	if calls != 0 {
		t.Fatalf("Present of a nil surface woke the watcher %d times", calls)
	}
	tt.Present(1)
	if calls != 1 {
		t.Fatalf("a present woke the watcher %d times, want 1", calls)
	}
	// The window draws: the frame shows the surface presented last, and
	// the watcher now counts as having drawn that present.
	img := tt.Frame("w", func() { calls++ })
	if img == nil {
		t.Fatal("Frame after a present = nil")
	}
	if calls != 1 {
		t.Fatalf("drawing woke the watcher itself %d times", calls)
	}
	tt.Present(2)
	if calls != 2 {
		t.Fatalf("a second present woke the watcher %d times, want 2", calls)
	}
	// Drawing again with a new callback hands the watch over to it.
	second := 0
	if img = tt.Frame("w", func() { second++ }); img == nil {
		t.Fatal("Frame of a showing texture = nil")
	}
	tt.Present(3)
	if calls != 2 || second != 1 {
		t.Fatalf("presents went to the old watcher (%d) or none (%d)", calls, second)
	}
	if v := img.Version(); v == 0 {
		t.Error("a frame that changed has no version")
	}
	// The CPU renderer reads the surface's pixels through the frame's
	// restore, which asks the texture.
	if img.Pixels() != nil {
		t.Error("a frame of a stub surface has pixels")
	}
}

func TestStaleWatchers(t *testing.T) {
	tt := New()
	calls := 0
	tt.Frame("w", func() { calls++ })
	// A window that stops drawing for staleAfter presents stops hearing
	// of them.
	for i := 0; i < staleAfter; i++ {
		tt.Present(1)
	}
	if calls != staleAfter {
		t.Fatalf("the watcher heard of %d of %d presents", calls, staleAfter)
	}
	for i := 0; i < 4*staleAfter; i++ {
		tt.Present(1)
	}
	if calls != staleAfter {
		t.Fatalf("a stale watcher still heard of presents (%d calls)", calls)
	}
	// Drawing again watches again.
	tt.Frame("w", func() { calls++ })
	tt.Present(1)
	if calls != staleAfter+1 {
		t.Fatalf("a redrawn texture woke its watcher %d times, want %d", calls, staleAfter+1)
	}
}

func TestGPUTexture(t *testing.T) {
	tt := New()
	if got := tt.s.GPUTexture(scene.Metal, 0); got != 0 {
		t.Fatalf("GPUTexture before any frame = %x, want 0", got)
	}
	tt.Present(1)
	img := tt.Frame("w", func() {})
	if got := tt.s.GPUTexture(scene.OpenGL, 0); got != 0 {
		t.Fatalf("GPUTexture for OpenGL = %x, want 0", got)
	}
	if got := tt.s.GPUTexture(scene.Metal, 0); got != 0 {
		t.Fatalf("GPUTexture for Metal of a stub surface = %x, want 0", got)
	}
	if img == nil {
		t.Fatal("Frame after a present = nil")
	}
}
