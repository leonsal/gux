package window

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type FontStyleType int

const (
	FontRegular FontStyleType = iota
	FontBold
	FontItalic
	FontMedium
	FontMediumItalic
	FontBoldItalic
	FontMono
	FontMonoBold
	FontMonoItalic
	FontMonoBoldItalic
	FontCustom
)

const (
	FontMaxSmaller = 4 // Maximum number of font faces smaller than normal
	FontMaxLarger  = 8 // Maximum number of font faces larger than normal
)

type fontInfo struct {
	fontData []byte       // Font data as TrueType or OpenFont used to build the font faces
	faces    []*FontAtlas // List of FontAtlases for font faces from sizes: -smaller to +larger
}

type FontManager struct {
	runeSets   [][]rune                    // Unicode codepoints range tables for the fonts
	normalSize float64                     // The normal font size in 'points'
	smaller    int                         // Number of font sizes smaller than the normal size
	larger     int                         // Number of font sizes greater than the normal size
	styles     map[FontStyleType]*fontInfo // Maps font style to font info
}

// NewFontManager creates and returns a new empty FontManager.
// The fonts contained by this FontManager will all have the specified unicode range tables
// Each font will have 'smaller' ...
func NewFontManager(normalSize float64, smaller, larger int, runeSets ...[]rune) (*FontManager, error) {

	if smaller < 0 || smaller > FontMaxSmaller {
		return nil, fmt.Errorf("invalid smaller font sizes")
	}
	if larger < 0 || larger > FontMaxLarger {
		return nil, fmt.Errorf("invalid larger font sizes")
	}
	if len(runeSets) == 0 {
		return nil, fmt.Errorf("at least one runeSet must be supplied")
	}

	fm := new(FontManager)
	fm.normalSize = normalSize
	fm.smaller = smaller
	fm.larger = larger
	fm.styles = make(map[FontStyleType]*fontInfo)
	fm.runeSets = append(fm.runeSets, runeSets...)
	return fm, nil
}

// AddStyle adds the specified font style to this FontManager.
// fontData must be a valid TTF/OpenType font description.
func (fm *FontManager) AddStyle(ff FontStyleType, fontData []byte) error {

	_, ok := fm.styles[ff]
	if ok {
		return fmt.Errorf("FontStyle:%d already added to the FontManager", ff)
	}
	fm.styles[ff] = &fontInfo{fontData: fontData}
	return nil
}

// BuildFonts builds the font atlases for each family and each size in this FontManager.
func (fm *FontManager) BuildFonts(w *Window) error {

	// Scale the normal font size from the window.

	for _, fi := range fm.styles {

		// If already built, continue with next family
		if len(fi.faces) > 0 {
			continue
		}

		// Creates font atlas for each relative size
		for relSize := -fm.smaller; relSize <= fm.larger; relSize++ {
			opts := opentype.FaceOptions{
				Size:    fm.normalSize + float64(relSize),
				DPI:     72,
				Hinting: font.HintingNone,
			}
			fa, err := NewFontAtlas(w, fi.fontData, &opts, fm.runeSets...)
			if err != nil {
				return err
			}
			fi.faces = append(fi.faces, fa)
		}
	}
	return nil
}

// DestroyFonts destroys all font atlases created previously for this FontManager.
// It normally should be called before the window is closed.
func (fm *FontManager) DestroyFonts(w *Window) {

	for _, fi := range fm.styles {
		for _, fa := range fi.faces {
			fa.Destroy(w)
		}
		fi.faces = nil
	}
}

// Font return pointer to the Font for the specified font family type and relative size.
// The relative size is 0 for normal, +1, +2, ... for larger and -1, -2, ... for smaller font faces.
func (fm *FontManager) Font(ff FontStyleType, relSize int) *FontAtlas {

	fi, ok := fm.styles[ff]
	if !ok {
		return nil
	}
	var index int
	if relSize > fm.larger {
		index = len(fi.faces) - 1
	} else if relSize < -fm.smaller {
		index = 0
	} else {
		index = relSize + fm.smaller
	}
	return fi.faces[index]
}
