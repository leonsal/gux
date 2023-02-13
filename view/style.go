package view

import (
	"github.com/leonsal/gux/color"
	"github.com/leonsal/gux/window"
)

// StyleType is the type for all style configurations except for colors
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

// StyleColorType is the type for all style color configuration
type StyleColorType int

const (
	StyleColorText StyleColorType = iota
	StyleColorTextDisabled
	StyleColorUser
)

// StyleMap maps a style configuration to its current value
type StyleMap map[StyleType]any

// StyleColorMap maps a style color configuration to its current color value
type StyleColorMap map[StyleColorType]color.Color

// windowStyleMap maps Windows to its current styles map
var windowStyleMap = map[*window.Window]StyleMap{}

// windowStyleColorMap maps Windows to its current style color map
var windowStyleColorMap = map[*window.Window]StyleColorMap{}

// SetWindowStyle sets the the style map for the window
func SetWindowStyle(w *window.Window, styles StyleMap) {

	windowStyleMap[w] = styles
}

// SetWindowStyleColor sets the style color map for the window
func SetWindowStyleColor(w *window.Window, scolors StyleColorMap) {

	windowStyleColorMap[w] = scolors
}

// Alpha returns the current alpha value from the style map
func (sm StyleMap) Alpha() float32 {

	return sm[StyleAlpha].(float32)
}

func (sm StyleMap) SetAlpha(alpha float32) {

	sm[StyleAlpha] = alpha
}

func (v *View) StyleAlpha(w *window.Window) float32 {

	return getStyle(w, v, StyleAlpha).(float32)
}

func (v *View) SetStyleAlpha(alpha float32) {

	v.setStyle(StyleAlpha, alpha)
}

func (v *View) DelStyleAlpha() {

	v.deleteStyle(StyleAlpha)
}

// StyleColor returns the current style color for the view, window and color configuration
func (v *View) StyleColor(w *window.Window, scolor StyleColorType) color.Color {

	if v.styleColor != nil {
		s, ok := v.styleColor[scolor]
		if ok {
			return s
		}
	}
	return windowStyleColorMap[w][scolor]
}

// SetStyleColor sets specific style color for the view
func (v *View) SetStyleColor(scolor StyleColorType, c color.Color) {

	if v.styleColor == nil {
		v.styleColor = StyleColorMap{}
	}
	v.styleColor[scolor] = c
}

// DelStyleColor deletes specific style color configuration for the view.
func (v *View) DelStyleColor(scolor StyleColorType) {

	if v.styleColor != nil {
		delete(v.styleColor, scolor)
	}
}

func getStyle(w *window.Window, v *View, stype StyleType) any {

	if v.style != nil {
		s, ok := v.style[stype]
		if ok {
			return s
		}
	}
	return windowStyleMap[w][stype]
}

func (v *View) setStyle(style StyleType, value any) {

	if v.style == nil {
		v.style = StyleMap{}
	}
	v.style[style] = value
}

func (v *View) deleteStyle(style StyleType) {

	if v.style == nil {
		delete(v.style, style)
	}
}

// type Style struct {
// 	Alpha            float32
// 	DisabledAlpha    float32
// 	WindowPadding    gb.Vec2 // Padding within a window.
// 	WindowRounding   float32 // Radius of window corners rounding. Set to 0.0f to have rectangular windows. Large values tend to lead to variety of artifacts and are not recommended.
// 	WindowBorderSize float32 // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
// 	WindowMinSize    float32 // Minimum window size. This is a global setting. If you want to constrain individual windows, use SetNextWindowSizeConstraints().
// }

// enum ImGuiCol_
// {
//     ImGuiCol_Text,
//     ImGuiCol_TextDisabled,
//     ImGuiCol_WindowBg,              // Background of normal windows
//     ImGuiCol_ChildBg,               // Background of child windows
//     ImGuiCol_PopupBg,               // Background of popups, menus, tooltips windows
//     ImGuiCol_Border,
//     ImGuiCol_BorderShadow,
//     ImGuiCol_FrameBg,               // Background of checkbox, radio button, plot, slider, text input
//     ImGuiCol_FrameBgHovered,
//     ImGuiCol_FrameBgActive,
//     ImGuiCol_TitleBg,
//     ImGuiCol_TitleBgActive,
//     ImGuiCol_TitleBgCollapsed,
//     ImGuiCol_MenuBarBg,
//     ImGuiCol_ScrollbarBg,
//     ImGuiCol_ScrollbarGrab,
//     ImGuiCol_ScrollbarGrabHovered,
//     ImGuiCol_ScrollbarGrabActive,
//     ImGuiCol_CheckMark,
//     ImGuiCol_SliderGrab,
//     ImGuiCol_SliderGrabActive,
//     ImGuiCol_Button,
//     ImGuiCol_ButtonHovered,
//     ImGuiCol_ButtonActive,
//     ImGuiCol_Header,                // Header* colors are used for CollapsingHeader, TreeNode, Selectable, MenuItem
//     ImGuiCol_HeaderHovered,
//     ImGuiCol_HeaderActive,
//     ImGuiCol_Separator,
//     ImGuiCol_SeparatorHovered,
//     ImGuiCol_SeparatorActive,
//     ImGuiCol_ResizeGrip,            // Resize grip in lower-right and lower-left corners of windows.
//     ImGuiCol_ResizeGripHovered,
//     ImGuiCol_ResizeGripActive,
//     ImGuiCol_Tab,                   // TabItem in a TabBar
//     ImGuiCol_TabHovered,
//     ImGuiCol_TabActive,
//     ImGuiCol_TabUnfocused,
//     ImGuiCol_TabUnfocusedActive,
//     ImGuiCol_DockingPreview,        // Preview overlay color when about to docking something
//     ImGuiCol_DockingEmptyBg,        // Background color for empty node (e.g. CentralNode with no window docked into it)
//     ImGuiCol_PlotLines,
//     ImGuiCol_PlotLinesHovered,
//     ImGuiCol_PlotHistogram,
//     ImGuiCol_PlotHistogramHovered,
//     ImGuiCol_TableHeaderBg,         // Table header background
//     ImGuiCol_TableBorderStrong,     // Table outer and header borders (prefer using Alpha=1.0 here)
//     ImGuiCol_TableBorderLight,      // Table inner borders (prefer using Alpha=1.0 here)
//     ImGuiCol_TableRowBg,            // Table row background (even rows)
//     ImGuiCol_TableRowBgAlt,         // Table row background (odd rows)
//     ImGuiCol_TextSelectedBg,
//     ImGuiCol_DragDropTarget,        // Rectangle highlighting a drop target
//     ImGuiCol_NavHighlight,          // Gamepad/keyboard: current highlighted item
//     ImGuiCol_NavWindowingHighlight, // Highlight window when using CTRL+TAB
//     ImGuiCol_NavWindowingDimBg,     // Darken/colorize entire screen behind the CTRL+TAB window list, when active
//     ImGuiCol_ModalWindowDimBg,      // Darken/colorize entire screen behind a modal window, when one is active
//     ImGuiCol_COUNT
// };
//
