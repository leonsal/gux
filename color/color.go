package color

import "github.com/leonsal/gux/gb"

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

func (c Color) RGBA() gb.RGBA {

	return gb.MakeColor(byte(c.R*255.0), byte(c.G*255.0), byte(c.B*255.0), byte(c.A*255.0))
}

func ColorFromRGBA(rgba gb.RGBA) Color {

	return Color{
		R: float32(((rgba >> gb.RGBAShiftR) & 0xFF) * 1.0 / 255.0),
		G: float32(((rgba >> gb.RGBAShiftG) & 0xFF) * 1.0 / 255.0),
		B: float32(((rgba >> gb.RGBAShiftB) & 0xFF) * 1.0 / 255.0),
		A: float32(((rgba >> gb.RGBAShiftA) & 0xFF) * 1.0 / 255.0),
	}
}

var (
	Aliceblue            = Color{0.941, 0.973, 1.000, 1.0}
	Antiquewhite         = Color{0.980, 0.922, 0.843, 1.0}
	Aqua                 = Color{0.000, 1.000, 1.000, 1.0}
	Aquamarine           = Color{0.498, 1.000, 0.831, 1.0}
	Azure                = Color{0.941, 1.000, 1.000, 1.0}
	Beige                = Color{0.961, 0.961, 0.863, 1.0}
	Bisque               = Color{1.000, 0.894, 0.769, 1.0}
	Black                = Color{0.000, 0.000, 0.000, 1.0}
	Blanchedalmond       = Color{1.000, 0.922, 0.804, 1.0}
	Blue                 = Color{0.000, 0.000, 1.000, 1.0}
	Blueviolet           = Color{0.541, 0.169, 0.886, 1.0}
	Brown                = Color{0.647, 0.165, 0.165, 1.0}
	Burlywood            = Color{0.871, 0.722, 0.529, 1.0}
	Cadetblue            = Color{0.373, 0.620, 0.627, 1.0}
	Chartreuse           = Color{0.498, 1.000, 0.000, 1.0}
	Chocolate            = Color{0.824, 0.412, 0.118, 1.0}
	Coral                = Color{1.000, 0.498, 0.314, 1.0}
	Cornflowerblue       = Color{0.392, 0.584, 0.929, 1.0}
	Cornsilk             = Color{1.000, 0.973, 0.863, 1.0}
	Crimson              = Color{0.863, 0.078, 0.235, 1.0}
	Cyan                 = Color{0.000, 1.000, 1.000, 1.0}
	Darkblue             = Color{0.000, 0.000, 0.545, 1.0}
	Darkcyan             = Color{0.000, 0.545, 0.545, 1.0}
	Darkgoldenrod        = Color{0.722, 0.525, 0.043, 1.0}
	Darkgray             = Color{0.663, 0.663, 0.663, 1.0}
	Darkgreen            = Color{0.000, 0.392, 0.000, 1.0}
	Darkgrey             = Color{0.663, 0.663, 0.663, 1.0}
	Darkkhaki            = Color{0.741, 0.718, 0.420, 1.0}
	Darkmagenta          = Color{0.545, 0.000, 0.545, 1.0}
	Darkolivegreen       = Color{0.333, 0.420, 0.184, 1.0}
	Darkorange           = Color{1.000, 0.549, 0.000, 1.0}
	Darkorchid           = Color{0.600, 0.196, 0.800, 1.0}
	Darkred              = Color{0.545, 0.000, 0.000, 1.0}
	Darksalmon           = Color{0.914, 0.588, 0.478, 1.0}
	Darkseagreen         = Color{0.561, 0.737, 0.561, 1.0}
	Darkslateblue        = Color{0.282, 0.239, 0.545, 1.0}
	Darkslategray        = Color{0.184, 0.310, 0.310, 1.0}
	Darkslategrey        = Color{0.184, 0.310, 0.310, 1.0}
	Darkturquoise        = Color{0.000, 0.808, 0.820, 1.0}
	Darkviolet           = Color{0.580, 0.000, 0.827, 1.0}
	Deeppink             = Color{1.000, 0.078, 0.576, 1.0}
	Deepskyblue          = Color{0.000, 0.749, 1.000, 1.0}
	Dimgray              = Color{0.412, 0.412, 0.412, 1.0}
	Dimgrey              = Color{0.412, 0.412, 0.412, 1.0}
	Dodgerblue           = Color{0.118, 0.565, 1.000, 1.0}
	Firebrick            = Color{0.698, 0.133, 0.133, 1.0}
	Floralwhite          = Color{1.000, 0.980, 0.941, 1.0}
	Forestgreen          = Color{0.133, 0.545, 0.133, 1.0}
	Fuchsia              = Color{1.000, 0.000, 1.000, 1.0}
	Gainsboro            = Color{0.863, 0.863, 0.863, 1.0}
	Ghostwhite           = Color{0.973, 0.973, 1.000, 1.0}
	Gold                 = Color{1.000, 0.843, 0.000, 1.0}
	Goldenrod            = Color{0.855, 0.647, 0.125, 1.0}
	Gray                 = Color{0.502, 0.502, 0.502, 1.0}
	Green                = Color{0.000, 0.502, 0.000, 1.0}
	Greenyellow          = Color{0.678, 1.000, 0.184, 1.0}
	Grey                 = Color{0.502, 0.502, 0.502, 1.0}
	Honeydew             = Color{0.941, 1.000, 0.941, 1.0}
	Hotpink              = Color{1.000, 0.412, 0.706, 1.0}
	Indianred            = Color{0.804, 0.361, 0.361, 1.0}
	Indigo               = Color{0.294, 0.000, 0.510, 1.0}
	Ivory                = Color{1.000, 1.000, 0.941, 1.0}
	Khaki                = Color{0.941, 0.902, 0.549, 1.0}
	Lavender             = Color{0.902, 0.902, 0.980, 1.0}
	Lavenderblush        = Color{1.000, 0.941, 0.961, 1.0}
	Lawngreen            = Color{0.486, 0.988, 0.000, 1.0}
	Lemonchiffon         = Color{1.000, 0.980, 0.804, 1.0}
	Lightblue            = Color{0.678, 0.847, 0.902, 1.0}
	Lightcoral           = Color{0.941, 0.502, 0.502, 1.0}
	Lightcyan            = Color{0.878, 1.000, 1.000, 1.0}
	Lightgoldenrodyellow = Color{0.980, 0.980, 0.824, 1.0}
	Lightgray            = Color{0.827, 0.827, 0.827, 1.0}
	Lightgreen           = Color{0.565, 0.933, 0.565, 1.0}
	Lightgrey            = Color{0.827, 0.827, 0.827, 1.0}
	Lightpink            = Color{1.000, 0.714, 0.757, 1.0}
	Lightsalmon          = Color{1.000, 0.627, 0.478, 1.0}
	Lightseagreen        = Color{0.125, 0.698, 0.667, 1.0}
	Lightskyblue         = Color{0.529, 0.808, 0.980, 1.0}
	Lightslategray       = Color{0.467, 0.533, 0.600, 1.0}
	Lightslategrey       = Color{0.467, 0.533, 0.600, 1.0}
	Lightsteelblue       = Color{0.690, 0.769, 0.871, 1.0}
	Lightyellow          = Color{1.000, 1.000, 0.878, 1.0}
	Lime                 = Color{0.000, 1.000, 0.000, 1.0}
	Limegreen            = Color{0.196, 0.804, 0.196, 1.0}
	Linen                = Color{0.980, 0.941, 0.902, 1.0}
	Magenta              = Color{1.000, 0.000, 1.000, 1.0}
	Maroon               = Color{0.502, 0.000, 0.000, 1.0}
	Mediumaquamarine     = Color{0.400, 0.804, 0.667, 1.0}
	Mediumblue           = Color{0.000, 0.000, 0.804, 1.0}
	Mediumorchid         = Color{0.729, 0.333, 0.827, 1.0}
	Mediumpurple         = Color{0.576, 0.439, 0.859, 1.0}
	Mediumseagreen       = Color{0.235, 0.702, 0.443, 1.0}
	Mediumslateblue      = Color{0.482, 0.408, 0.933, 1.0}
	Mediumspringgreen    = Color{0.000, 0.980, 0.604, 1.0}
	Mediumturquoise      = Color{0.282, 0.820, 0.800, 1.0}
	Mediumvioletred      = Color{0.780, 0.082, 0.522, 1.0}
	Midnightblue         = Color{0.098, 0.098, 0.439, 1.0}
	Mintcream            = Color{0.961, 1.000, 0.980, 1.0}
	Mistyrose            = Color{1.000, 0.894, 0.882, 1.0}
	Moccasin             = Color{1.000, 0.894, 0.710, 1.0}
	Navajowhite          = Color{1.000, 0.871, 0.678, 1.0}
	Navy                 = Color{0.000, 0.000, 0.502, 1.0}
	Oldlace              = Color{0.992, 0.961, 0.902, 1.0}
	Olive                = Color{0.502, 0.502, 0.000, 1.0}
	Olivedrab            = Color{0.420, 0.557, 0.137, 1.0}
	Orange               = Color{1.000, 0.647, 0.000, 1.0}
	Orangered            = Color{1.000, 0.271, 0.000, 1.0}
	Orchid               = Color{0.855, 0.439, 0.839, 1.0}
	Palegoldenrod        = Color{0.933, 0.910, 0.667, 1.0}
	Palegreen            = Color{0.596, 0.984, 0.596, 1.0}
	Paleturquoise        = Color{0.686, 0.933, 0.933, 1.0}
	Palevioletred        = Color{0.859, 0.439, 0.576, 1.0}
	Papayawhip           = Color{1.000, 0.937, 0.835, 1.0}
	Peachpuff            = Color{1.000, 0.855, 0.725, 1.0}
	Peru                 = Color{0.804, 0.522, 0.247, 1.0}
	Pink                 = Color{1.000, 0.753, 0.796, 1.0}
	Plum                 = Color{0.867, 0.627, 0.867, 1.0}
	Powderblue           = Color{0.690, 0.878, 0.902, 1.0}
	Purple               = Color{0.502, 0.000, 0.502, 1.0}
	Red                  = Color{1.000, 0.000, 0.000, 1.0}
	Rosybrown            = Color{0.737, 0.561, 0.561, 1.0}
	Royalblue            = Color{0.255, 0.412, 0.882, 1.0}
	Saddlebrown          = Color{0.545, 0.271, 0.075, 1.0}
	Salmon               = Color{0.980, 0.502, 0.447, 1.0}
	Sandybrown           = Color{0.957, 0.643, 0.376, 1.0}
	Seagreen             = Color{0.180, 0.545, 0.341, 1.0}
	Seashell             = Color{1.000, 0.961, 0.933, 1.0}
	Sienna               = Color{0.627, 0.322, 0.176, 1.0}
	Silver               = Color{0.753, 0.753, 0.753, 1.0}
	Skyblue              = Color{0.529, 0.808, 0.922, 1.0}
	Slateblue            = Color{0.416, 0.353, 0.804, 1.0}
	Slategray            = Color{0.439, 0.502, 0.565, 1.0}
	Slategrey            = Color{0.439, 0.502, 0.565, 1.0}
	Snow                 = Color{1.000, 0.980, 0.980, 1.0}
	Springgreen          = Color{0.000, 1.000, 0.498, 1.0}
	Steelblue            = Color{0.275, 0.510, 0.706, 1.0}
	Tan                  = Color{0.824, 0.706, 0.549, 1.0}
	Teal                 = Color{0.000, 0.502, 0.502, 1.0}
	Thistle              = Color{0.847, 0.749, 0.847, 1.0}
	Tomato               = Color{1.000, 0.388, 0.278, 1.0}
	Turquoise            = Color{0.251, 0.878, 0.816, 1.0}
	Violet               = Color{0.933, 0.510, 0.933, 1.0}
	Wheat                = Color{0.961, 0.871, 0.702, 1.0}
	White                = Color{1.000, 1.000, 1.000, 1.0}
	Whitesmoke           = Color{0.961, 0.961, 0.961, 1.0}
	Yellow               = Color{1.000, 1.000, 0.000, 1.0}
	Yellowgreen          = Color{0.604, 0.804, 0.196, 1.0}
)
