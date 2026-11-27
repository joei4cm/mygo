package bridge

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/egoist/mygo/internal/platform"
)

func TestScript(t *testing.T) {
	src := Script(Config{
		Platform: "darwin",
		WindowID: 7,
		Version:  "1.2.3",
		Secret:   "s3cret",
		Locale:   "zh-Hans",
		Dir:      "rtl",
		Fit:      &FitContent{MinWidth: 100, MaxHeight: 50},
	})
	if !strings.HasPrefix(src, "(function(__MYGO_CONFIG__){\n") || !strings.HasSuffix(src, ");") {
		t.Fatalf("Script is not the bridge wrapped in a function: %.40s…", src)
	}
	start := strings.LastIndex(src, "})(") + 3
	var cfg Config
	if err := json.Unmarshal([]byte(src[start:len(src)-2]), &cfg); err != nil {
		t.Fatalf("the config is not JSON: %v", err)
	}
	if cfg.Platform != "darwin" || cfg.WindowID != 7 || cfg.Version != "1.2.3" || cfg.Secret != "s3cret" || cfg.Locale != "zh-Hans" || cfg.Dir != "rtl" {
		t.Errorf("the config came back as %+v", cfg)
	}
	if cfg.Fit == nil || cfg.Fit.MinWidth != 100 || cfg.Fit.MaxHeight != 50 {
		t.Errorf("the config's fit came back as %+v", cfg.Fit)
	}
	if cfg := Script(Config{Platform: "web"}); strings.Contains(cfg, `"fit"`) {
		t.Error("a window without fit limits got a fit key")
	}
	// The bridge source itself is there: the IIFE ends with it, and it
	// registers the secret it posts with.
	if !strings.Contains(src, "__mygo") {
		t.Error("the bridge source is not in the script")
	}
}

func TestNewTitleBar(t *testing.T) {
	tb := platform.TitleBar{Height: 56, Left: 16, Right: 80, Bottom: 8}
	for _, c := range []struct {
		zoom float64
		want TitleBar
	}{
		{2, TitleBar{Height: 28, Left: 8, Right: 40, Bottom: 4}},
		{1, TitleBar{Height: 56, Left: 16, Right: 80, Bottom: 8}},
		{0, TitleBar{Height: 56, Left: 16, Right: 80, Bottom: 8}},
		{-3, TitleBar{Height: 56, Left: 16, Right: 80, Bottom: 8}},
	} {
		if got := NewTitleBar(tb, c.zoom); got != c.want {
			t.Errorf("NewTitleBar at zoom %v = %+v, want %+v", c.zoom, got, c.want)
		}
	}
}

func TestNewFitContent(t *testing.T) {
	min, max := platform.Size{Width: 200, Height: 100}, platform.Size{Width: 800, Height: 600}
	for _, c := range []struct {
		zoom float64
		want FitContent
	}{
		{2, FitContent{MinWidth: 100, MinHeight: 50, MaxWidth: 400, MaxHeight: 300}},
		{0, FitContent{MinWidth: 200, MinHeight: 100, MaxWidth: 800, MaxHeight: 600}},
	} {
		if got := NewFitContent(min, max, c.zoom); got != c.want {
			t.Errorf("NewFitContent at zoom %v = %+v, want %+v", c.zoom, got, c.want)
		}
	}
}

func TestScripts(t *testing.T) {
	tb := platform.TitleBar{Height: 56, Left: 16, Right: 80, Bottom: 8}
	s := TitleBarScript(tb, 2)
	if !strings.Contains(s, `"n":"`+TitleBarEvent+`"`) || !strings.Contains(s, `"t":"event"`) {
		t.Errorf("TitleBarScript does not deliver the event: %s", s)
	}
	if !strings.Contains(s, `"height":28`) || !strings.Contains(s, `"left":8`) {
		t.Errorf("TitleBarScript does not deliver CSS pixels: %s", s)
	}
	if !strings.Contains(s, `window.__mygo&&window.__mygo.receive(`) {
		t.Errorf("TitleBarScript does not go through the bridge: %s", s)
	}
	p := PreferencesScript([]byte(`{"accent":"#204060"}`))
	if !strings.Contains(p, `"n":"`+PreferencesEvent+`"`) || !strings.Contains(p, `"p":{"accent":"#204060"}`) {
		t.Errorf("PreferencesScript does not deliver the payload: %s", p)
	}
}
