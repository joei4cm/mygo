package ui

import (
	"errors"
	"testing"
	"time"

	"github.com/egoist/mygo/i18n"
)

// What the public Context tells the view: the language and the region of
// the user, the settings of the desktop, the room the system covers, the
// room the window's controls take, the window's zoom, and the keys held.
// What it asks of the system: the clipboard, and the URLs the user follows.
func TestContextEnvironment(t *testing.T) {
	var (
		locale string
		region i18n.Region
		prefs  Preferences
		top    float32
		right  float32
		bottom float32
		left   float32
		vibrancy,
		hidden bool
		zoom      float32
		bar       TitleBar
		menus     MenuPresentation
		mods      Modifiers
		keyA      bool
		keys      []Key
		textW     float32
		textH     float32
		router    *Router
		clip      string
		openErr   error
		menusSet  bool
		clipAsk   bool
		urlAsk    bool
		urlOpened bool
		zoomAsk   bool
		zoomSet   bool
	)
	r := NewRouter("/")
	tt := coreNewTester(func(c *context) {
		r.coreView(c, func(rt *Route) {
			p := makeContext(c)
			locale = p.Locale()
			region = p.Region()
			prefs = p.Preferences()
			top, right, bottom, left = p.SafeArea()
			vibrancy = p.Vibrancy()
			hidden = p.WindowHidden()
			zoom = p.Zoom()
			bar = p.TitleBar()
			menus = p.MenuPresentation()
			mods = p.Modifiers()
			keyA = p.KeyDown(KeyA)
			keys = p.KeysDown()
			textW, textH = p.MeasureText(100, Span{Text: "Hello"})
			router = p.Router()
			if !menusSet {
				// The view asks for the menus of its own look, and the
				// frames after it show them.
				menusSet = true
				p.SetMenuPresentation(PopoverMenus)
			}
			if clipAsk {
				clip = p.ReadClipboard()
				p.WriteClipboard("written")
			}
			if zoomAsk && !zoomSet {
				// The view zooms the window's content itself.
				zoomSet = true
				p.SetZoom(1.25)
			}
			if urlAsk && !urlOpened {
				urlOpened = true
				p.OpenURL("https://example.com")
				p.OpenURLThen("https://example.com/2", func(err error) {
					openErr = err
				})
			}
		})
	}, 200, 200)
	tt.SetLocale("fr")
	tt.SetRegion(i18n.Region{Locale: "fr-FR", FirstWeekday: time.Monday, HourCycle: 24})
	tt.SetPreferences(Preferences{TextScale: 1.5, ReduceMotion: true})
	tt.SetInsets(1, 2, 3, 4)
	tt.SetVibrancy(true)
	tt.SetTitleBar(TitleBar{Height: 28, Left: 70, Right: 70, Bottom: 30})
	tt.KeyDown(KeyA)
	tt.HoldModifiers(Shift)
	tt.Frame()

	if locale != "fr" {
		t.Errorf("the language is %q", locale)
	}
	if region.Locale != "fr-FR" || region.FirstWeekday != time.Monday || region.HourCycle != 24 {
		t.Errorf("the region is %+v", region)
	}
	if prefs.TextScale != 1.5 || !prefs.ReduceMotion {
		t.Errorf("the settings are %+v", prefs)
	}
	if top != 1 || right != 2 || bottom != 3 || left != 4 {
		t.Errorf("the room the system covers is %v, %v, %v, %v", top, right, bottom, left)
	}
	if !vibrancy {
		t.Error("the window shows no material")
	}
	if hidden {
		t.Error("the window is hidden")
	}
	if bar.Height != 28 || bar.Left != 70 || bar.Right != 70 || bar.Bottom != 30 {
		t.Errorf("the room the controls take is %+v", bar)
	}
	if menus != PopoverMenus {
		t.Errorf("the menus are %v", menus)
	}
	if mods&Shift == 0 {
		t.Errorf("the modifiers held are %v", mods)
	}
	if !keyA {
		t.Error("the key held is not reported")
	}
	if len(keys) != 1 || keys[0] != KeyA {
		t.Errorf("the keys held are %v", keys)
	}
	if router != r {
		t.Errorf("the Context tells the view the router %v", router)
	}
	if textW <= 0 || textH <= 0 {
		t.Errorf("the text measures %v by %v", textW, textH)
	}

	// The window zooms its content, and the room the system covers is in
	// the content's DIPs, half of what the window's are.
	tt.SetZoom(2)
	tt.Frame()
	if zoom != 2 {
		t.Errorf("the zoom is %v", zoom)
	}
	if top != 0.5 || right != 1 || bottom != 1.5 || left != 2 {
		t.Errorf("the room the system covers is %v, %v, %v, %v", top, right, bottom, left)
	}

	// The view reads what the clipboard holds, and puts text there.
	tt.SetClipboard("copied")
	clipAsk = true
	tt.Frame()
	if clip != "copied" {
		t.Errorf("the clipboard holds %q", clip)
	}
	if tt.Clipboard() != "written" {
		t.Errorf("the clipboard holds %q", tt.Clipboard())
	}

	// The view opens what the user follows, and hears what came of it.
	tt.FailOpenURL(errors.New("no app"))
	urlAsk = true
	tt.Frame()
	urls := tt.OpenedURLs()
	if len(urls) != 2 || urls[0] != "https://example.com" || urls[1] != "https://example.com/2" {
		t.Errorf("the URLs opened are %v", urls)
	}
	if openErr == nil {
		t.Error("the view heard nothing of the URL it opened")
	}

	// The view zooms the window's content itself.
	zoomAsk = true
	tt.Frame()
	if zoom != 1.25 {
		t.Errorf("the zoom is %v", zoom)
	}

	// Outside a frame the Context tells nothing: it is the view's, and
	// the view is not building.
	p := tt.rt.public
	if p.Locale() != "" || p.Modifiers() != 0 || p.KeyDown(KeyA) || p.KeysDown() != nil {
		t.Error("the Context answers outside a frame")
	}
	if p.Region() != (i18n.Region{}) || p.Preferences() != (Preferences{}) {
		t.Error("the Context answers outside a frame")
	}
	if top, right, bottom, left = p.SafeArea(); top != 0 || right != 0 || bottom != 0 || left != 0 {
		t.Errorf("the room the system covers is %v, %v, %v, %v", top, right, bottom, left)
	}
	if p.Vibrancy() || p.WindowHidden() || p.Zoom() != 0 || p.TitleBar() != (TitleBar{}) {
		t.Error("the Context answers outside a frame")
	}
	if p.MenuPresentation() != SystemMenus || p.Router() != nil {
		t.Error("the Context answers outside a frame")
	}
	if w, h := p.MeasureText(100, Span{Text: "Hello"}); w != 0 || h != 0 {
		t.Errorf("the text measures %v by %v outside a frame", w, h)
	}
	// Outside a frame the view asks for nothing, and nothing changes.
	p.SetMenuPresentation(SystemMenus)
	p.SetZoom(1.5)
	p.WriteClipboard("again")
	if p.ReadClipboard() != "" {
		t.Error("the Context answers outside a frame")
	}
	if tt.Clipboard() != "written" {
		t.Errorf("the clipboard holds %q", tt.Clipboard())
	}
	p.OpenURL("https://example.com/3")
	p.OpenURLThen("https://example.com/4", func(err error) { openErr = nil })
	if urls = tt.OpenedURLs(); len(urls) != 2 {
		t.Errorf("the URLs opened are %v", urls)
	}
	tt.Frame()
	if urls = tt.OpenedURLs(); len(urls) != 2 {
		t.Errorf("the URLs opened are %v", urls)
	}
	if menus != PopoverMenus {
		t.Errorf("the menus are %v", menus)
	}
	if zoom != 1.25 {
		t.Errorf("the zoom is %v", zoom)
	}

	// A window that shows nothing tells the view it is hidden, and shows
	// again as it does.
	tt.SetWindowHidden(true)
	if !hidden {
		t.Error("the window that shows nothing is not reported hidden")
	}
	tt.SetWindowHidden(false)
	if hidden {
		t.Error("the window shown again is still hidden")
	}
}
