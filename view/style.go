package view

import "github.com/leonsal/gux/window"

type StyleType int

const (
	StyleAlpha StyleType = iota
	StyleDisabledAlpha
	StyleWindowPadding
	StyleWindowRounding
	StyleWindowBorderSize
	StyleWindowMinSize
	StyleUser
)

type StyleMap map[StyleType]any

// windowStyles maps Windows pointer to its current styles map
var windowStyles = map[*window.Window]StyleMap{}

func SetWindowStyles(w *window.Window, styles StyleMap) {

	windowStyles[w] = styles
}

func getStyle(w *window.Window, v IView, stype StyleType) any {

	vstyles := v.Styles()
	if vstyles != nil {
		s, ok := vstyles[stype]
		if ok {
			return s
		}
	}
	return windowStyles[w][stype]
}

func Alpha(w *window.Window, v IView) float32 {

	return getStyle(w, v, StyleAlpha).(float32)
}

func (s *StyleMap) Alpha() float32 {

	return (*s)[StyleAlpha].(float32)
}

func (s *StyleMap) SetAlpha(alpha float32) {

	(*s)[StyleAlpha] = alpha
}

func (s *StyleMap) DisabledAlpha() float32 {

	return (*s)[StyleDisabledAlpha].(float32)
}

func (s *StyleMap) SetDisabledAlpha(alpha float32) {

	(*s)[StyleDisabledAlpha] = alpha
}

// type Style struct {
// 	Alpha            float32
// 	DisabledAlpha    float32
// 	WindowPadding    gb.Vec2 // Padding within a window.
// 	WindowRounding   float32 // Radius of window corners rounding. Set to 0.0f to have rectangular windows. Large values tend to lead to variety of artifacts and are not recommended.
// 	WindowBorderSize float32 // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
// 	WindowMinSize    float32 // Minimum window size. This is a global setting. If you want to constrain individual windows, use SetNextWindowSizeConstraints().
// }
