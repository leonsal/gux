package gux

import (
	"fmt"

	"github.com/leonsal/gux/gb"
)

//const TexLinesWidthMax = 63

const TexLinesWidthMax = 9

type Window struct {
	gbw        *gb.Window                    // Graphics backend native window reference
	dl         gb.DrawList                   // Draw list to render
	texUvLines [TexLinesWidthMax + 1]gb.Vec4 // UV coordinates for textured lines
}

func NewWindow(title string, width, height int) (*Window, error) {

	// Creates graphics backend native window
	w := new(Window)
	var err error
	w.gbw, err = gb.CreateWindow(title, width, height)
	if err != nil {
		return nil, err
	}

	// Create line texture and transfer to backend
	w.buildRenderLinesTexData()

	return w, nil
}

func (w *Window) StartFrame(timeout float64) bool {

	w.dl.Clear()
	return w.gbw.StartFrame(timeout)
}

func (w *Window) RenderFrame(view IView) {

	view.Render(w)
	w.gbw.RenderFrame(&w.dl)
}

// Adds specified draw list to this Window's draw list
func (w *Window) AddList(src gb.DrawList) {

	w.dl.AddList(src)
}

func (w *Window) Destroy() {

	w.gbw.Destroy()
}

// buildRenderLinesTexData generates a texture with a triangular shape with various line widths
// stacked on top of each other to allow interpolation between them.
func (w *Window) buildRenderLinesTexData() {

	/*
		Example for TexLinesWidthMax = 9
		T - transparent white
		O - opaque white

		Line Width		Texels
		0				TTTTTTTTTTT
		1               TTTTTOTTTTT
		2               TTTTOOTTTTT
		3               TTTTOOOTTTT
		4               TTTOOOOOTTT
		5               TTTOOOOOTTT
		6               TTOOOOOOTTT
		7               TTOOOOOOOTT
		8               TOOOOOOOOTT
		9               TOOOOOOOOOT
		...
	*/

	width := TexLinesWidthMax + 2
	height := TexLinesWidthMax + 1
	rect := make([]gb.Color, width*height)
	uvScale := gb.Vec2{1 / float32(width), 1 / float32(height)}
	for n := 0; n < height; n++ {

		// Each line consists of at least one transparent pixel at each side, with a line of solid pixels in the middle
		lineWidth := n
		padLeft := (width - lineWidth) / 2
		padRight := width - (padLeft + lineWidth)
		fmt.Println("lineWidth", lineWidth, "padLeft", padLeft, "padRight", padRight)
		pos := n * width

		for i := 0; i < padLeft; i++ {
			rect[pos+i] = gb.MakeColor(255, 255, 255, 0)
		}
		for i := 0; i < lineWidth; i++ {
			rect[pos+padLeft+i] = gb.MakeColor(255, 255, 255, 255)
		}
		for i := 0; i < padRight; i++ {
			rect[pos+padLeft+lineWidth+i] = gb.MakeColor(255, 255, 255, 0)
		}

		// Calculate UVs for this line
		uv0 := gb.Vec2Mult(gb.Vec2{float32(padLeft - 1), float32(n)}, uvScale)
		uv1 := gb.Vec2Mult(gb.Vec2{float32(padLeft + lineWidth + 1), float32(n + 1)}, uvScale)
		//halfV := (uv0.Y + uv1.Y) * 0.5 // Calculate a constant V in the middle of the row to avoid sampling artifacts
		//w.texUvLines[n] = gb.Vec4{uv0.X, halfV, uv1.X, halfV}
		w.texUvLines[n] = gb.Vec4{uv0.X, uv0.Y, uv1.X, uv1.Y}
	}

	// Print image data
	for n := 0; n < height; n++ {
		pos := n * width
		for c := 0; c < width; c++ {
			fmt.Printf("%d ", rect[pos+c])
		}
		fmt.Println()
	}

	// Print UVs
	for n := 0; n < height; n++ {
		fmt.Println(w.texUvLines[n])
	}
}

//static void ImFontAtlasBuildRenderLinesTexData(ImFontAtlas* atlas)
//{
//    if (atlas->Flags & ImFontAtlasFlags_NoBakedLines)
//        return;
//
//    // This generates a triangular shape in the texture, with the various line widths stacked on top of each other to allow interpolation between them
//    ImFontAtlasCustomRect* r = atlas->GetCustomRectByIndex(atlas->PackIdLines);
//    IM_ASSERT(r->IsPacked());
//    for (unsigned int n = 0; n < IM_DRAWLIST_TEX_LINES_WIDTH_MAX + 1; n++) // +1 because of the zero-width row
//    {
//        // Each line consists of at least two empty pixels at the ends, with a line of solid pixels in the middle
//        unsigned int y = n;
//        unsigned int line_width = n;
//        unsigned int pad_left = (r->Width - line_width) / 2;
//        unsigned int pad_right = r->Width - (pad_left + line_width);
//
//        // Write each slice
//        IM_ASSERT(pad_left + line_width + pad_right == r->Width && y < r->Height); // Make sure we're inside the texture bounds before we start writing pixels
//        if (atlas->TexPixelsAlpha8 != NULL)
//        {
//            unsigned char* write_ptr = &atlas->TexPixelsAlpha8[r->X + ((r->Y + y) * atlas->TexWidth)];
//            for (unsigned int i = 0; i < pad_left; i++)
//                *(write_ptr + i) = 0x00;
//
//            for (unsigned int i = 0; i < line_width; i++)
//                *(write_ptr + pad_left + i) = 0xFF;
//
//            for (unsigned int i = 0; i < pad_right; i++)
//                *(write_ptr + pad_left + line_width + i) = 0x00;
//        }
//        else
//        {
//            unsigned int* write_ptr = &atlas->TexPixelsRGBA32[r->X + ((r->Y + y) * atlas->TexWidth)];
//            for (unsigned int i = 0; i < pad_left; i++)
//                *(write_ptr + i) = IM_COL32(255, 255, 255, 0);
//
//            for (unsigned int i = 0; i < line_width; i++)
//                *(write_ptr + pad_left + i) = IM_COL32_WHITE;
//
//            for (unsigned int i = 0; i < pad_right; i++)
//                *(write_ptr + pad_left + line_width + i) = IM_COL32(255, 255, 255, 0);
//        }
//
//        // Calculate UVs for this line
//        ImVec2 uv0 = ImVec2((float)(r->X + pad_left - 1), (float)(r->Y + y)) * atlas->TexUvScale;
//        ImVec2 uv1 = ImVec2((float)(r->X + pad_left + line_width + 1), (float)(r->Y + y + 1)) * atlas->TexUvScale;
//        float half_v = (uv0.y + uv1.y) * 0.5f; // Calculate a constant V in the middle of the row to avoid sampling artifacts
//        atlas->TexUvLines[n] = ImVec4(uv0.x, half_v, uv1.x, half_v);
//    }
//}
