package zml

import (
	"math"
	"strconv"
	"strings"
)

// Color is a RGB set of ints; for a nice picker
// see https://www.w3schools.com/colors/colors_picker.asp
type Color struct {
	Red, Green, Blue int
}

// colornames maps SVG color names to RGB triples.
var colornames = map[string]Color{
	"aliceblue":            {240, 248, 255},
	"antiquewhite":         {250, 235, 215},
	"aqua":                 {0, 255, 255},
	"aquamarine":           {127, 255, 212},
	"azure":                {240, 255, 255},
	"beige":                {245, 245, 220},
	"bisque":               {255, 228, 196},
	"black":                {0, 0, 0},
	"charlestongreen":      {35, 43, 43},
	"eerieblack":           {27, 27, 27},
	"jetblack":             {52, 52, 52},
	"blanchedalmond":       {255, 235, 205},
	"blue":                 {0, 0, 255},
	"blueviolet":           {138, 43, 226},
	"brown":                {165, 42, 42},
	"burlywood":            {222, 184, 135},
	"cadetblue":            {95, 158, 160},
	"chartreuse":           {127, 255, 0},
	"chocolate":            {210, 105, 30},
	"coral":                {255, 127, 80},
	"cornflowerblue":       {100, 149, 237},
	"cornsilk":             {255, 248, 220},
	"crimson":              {220, 20, 60},
	"cyan":                 {0, 255, 255},
	"darkblue":             {0, 0, 139},
	"darkcyan":             {0, 139, 139},
	"darkgoldenrod":        {184, 134, 11},
	"darkgray":             {169, 169, 169},
	"darkgreen":            {0, 100, 0},
	"darkgrey":             {169, 169, 169},
	"darkkhaki":            {189, 183, 107},
	"darkmagenta":          {139, 0, 139},
	"darkolivegreen":       {85, 107, 47},
	"darkorange":           {255, 140, 0},
	"darkorchid":           {153, 50, 204},
	"darkred":              {139, 0, 0},
	"darksalmon":           {233, 150, 122},
	"darkseagreen":         {143, 188, 143},
	"darkslateblue":        {72, 61, 139},
	"darkslategray":        {47, 79, 79},
	"darkslategrey":        {47, 79, 79},
	"darkturquoise":        {0, 206, 209},
	"darkviolet":           {148, 0, 211},
	"deeppink":             {255, 20, 147},
	"deepskyblue":          {0, 191, 255},
	"dimgray":              {105, 105, 105},
	"dimgrey":              {105, 105, 105},
	"dodgerblue":           {30, 144, 255},
	"firebrick":            {178, 34, 34},
	"floralwhite":          {255, 250, 240},
	"forestgreen":          {34, 139, 34},
	"fuchsia":              {255, 0, 255},
	"gainsboro":            {220, 220, 220},
	"ghostwhite":           {248, 248, 255},
	"gold":                 {255, 215, 0},
	"goldenrod":            {218, 165, 32},
	"gray":                 {128, 128, 128},
	"green":                {0, 128, 0},
	"greenyellow":          {173, 255, 47},
	"grey":                 {128, 128, 128},
	"honeydew":             {240, 255, 240},
	"hotpink":              {255, 105, 180},
	"indianred":            {205, 92, 92},
	"indigo":               {75, 0, 130},
	"ivory":                {255, 255, 240},
	"khaki":                {240, 230, 140},
	"lavender":             {230, 230, 250},
	"lavenderblush":        {255, 240, 245},
	"lawngreen":            {124, 252, 0},
	"lemonchiffon":         {255, 250, 205},
	"lightblue":            {173, 216, 230},
	"lightcoral":           {240, 128, 128},
	"lightcyan":            {224, 255, 255},
	"lightgoldenrodyellow": {250, 250, 210},
	"lightgray":            {211, 211, 211},
	"lightgreen":           {144, 238, 144},
	"lightgrey":            {211, 211, 211},
	"lightpink":            {255, 182, 193},
	"lightsalmon":          {255, 160, 122},
	"lightseagreen":        {32, 178, 170},
	"lightskyblue":         {135, 206, 250},
	"lightslategray":       {119, 136, 153},
	"lightslategrey":       {119, 136, 153},
	"lightsteelblue":       {176, 196, 222},
	"lightyellow":          {255, 255, 224},
	"lime":                 {0, 255, 0},
	"limegreen":            {50, 205, 50},
	"linen":                {250, 240, 230},
	"magenta":              {255, 0, 255},
	"maroon":               {128, 0, 0},
	"mediumaquamarine":     {102, 205, 170},
	"mediumblue":           {0, 0, 205},
	"mediumorchid":         {186, 85, 211},
	"mediumpurple":         {147, 112, 219},
	"mediumseagreen":       {60, 179, 113},
	"mediumslateblue":      {123, 104, 238},
	"mediumspringgreen":    {0, 250, 154},
	"mediumturquoise":      {72, 209, 204},
	"mediumvioletred":      {199, 21, 133},
	"midnightblue":         {25, 25, 112},
	"mintcream":            {245, 255, 250},
	"mistyrose":            {255, 228, 225},
	"moccasin":             {255, 228, 181},
	"navajowhite":          {255, 222, 173},
	"navy":                 {0, 0, 128},
	"oldlace":              {253, 245, 230},
	"olive":                {128, 128, 0},
	"olivedrab":            {107, 142, 35},
	"orange":               {255, 165, 0},
	"orangered":            {255, 69, 0},
	"orchid":               {218, 112, 214},
	"palegoldenrod":        {238, 232, 170},
	"palegreen":            {152, 251, 152},
	"paleturquoise":        {175, 238, 238},
	"palevioletred":        {219, 112, 147},
	"papayawhip":           {255, 239, 213},
	"peachpuff":            {255, 218, 185},
	"peru":                 {205, 133, 63},
	"pink":                 {255, 192, 203},
	"plum":                 {221, 160, 221},
	"powderblue":           {176, 224, 230},
	"purple":               {128, 0, 128},
	"red":                  {255, 0, 0},
	"platered":             {255, 80, 80},
	"rosybrown":            {188, 143, 143},
	"royalblue":            {65, 105, 225},
	"saddlebrown":          {139, 69, 19},
	"salmon":               {250, 128, 114},
	"sandybrown":           {244, 164, 96},
	"seagreen":             {46, 139, 87},
	"seashell":             {255, 245, 238},
	"sienna":               {160, 82, 45},
	"silver":               {192, 192, 192},
	"skyblue":              {135, 206, 235},
	"slateblue":            {106, 90, 205},
	"slategray":            {112, 128, 144},
	"slategrey":            {112, 128, 144},
	"snow":                 {255, 250, 250},
	"springgreen":          {0, 255, 127},
	"steelblue":            {70, 130, 180},
	"tan":                  {210, 180, 140},
	"teal":                 {0, 128, 128},
	"thistle":              {216, 191, 216},
	"tomato":               {255, 99, 71},
	"turquoise":            {64, 224, 208},
	"violet":               {238, 130, 238},
	"wheat":                {245, 222, 179},
	"white":                {255, 255, 255},
	"whitesmoke":           {245, 245, 245},
	"yellow":               {255, 255, 0},
	"yellowgreen":          {154, 205, 50},
}

// Colorlookup returns a RGB triple corresponding to the named color, "rgb(r,g,b)" or "#rrggbb" string.
// On error, return black.
func Colorlookup(s string) (r int, g int, b int) {
	color, ok := colornames[s]
	if ok {
		return color.Red, color.Green, color.Blue
	}
	var red, green, blue int
	ls := len(s)
	// rgb(r, g, b)
	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") && ls > 5 {
		v := colorNumbers(s)
		if len(v) == 3 {
			red, _ = strconv.Atoi(v[0])
			green, _ = strconv.Atoi(v[1])
			blue, _ = strconv.Atoi(v[2])
		}
		return red, green, blue
	}
	// #rrggbb
	if strings.HasPrefix(s, "#") && ls == 7 {
		r, _ := strconv.ParseInt(s[1:3], 16, 32)
		g, _ := strconv.ParseInt(s[3:5], 16, 32)
		b, _ := strconv.ParseInt(s[5:7], 16, 32)
		return int(r), int(g), int(b)
	}
	// hsv(hue, saturation, value)
	if strings.HasPrefix(s, "hsv(") && strings.HasSuffix(s, ")") && ls > 5 {
		v := colorNumbers(s)
		if len(v) == 3 {
			hue, _ := strconv.ParseFloat(v[0], 64)
			sat, _ := strconv.ParseFloat(v[1], 64)
			value, _ := strconv.ParseFloat(v[2], 64)
			red, green, blue = hsv2rgb(hue, sat, value)
		}
		return red, green, blue
	}
	return 0, 0, 0
}

// colorNumbers returns a list of numbers from a comma separated list,
// in the form of xxx(n1, n2, n3), after removing tabs and spaces.
func colorNumbers(s string) []string {
	return strings.Split(strings.NewReplacer(" ", "", "\t", "").Replace(s[4:len(s)-1]), ",")
}

// hsv2rgb converts hsv(h (0-360), s (0-100), v (0-100)) to rgb
// reference: https://en.wikipedia.org/wiki/HSL_and_HSV#HSV_to_RGB
func hsv2rgb(h, s, v float64) (int, int, int) {
	s /= 100
	v /= 100
	if s > 1 || v > 1 {
		return 0, 0, 0
	}
	h = math.Mod(h, 360)
	c := v * s
	section := h / 60
	x := c * (1 - math.Abs(math.Mod(section, 2)-1))

	var r, g, b float64
	switch {
	case section >= 0 && section <= 1:
		r = c
		g = x
		b = 0
	case section > 1 && section <= 2:
		r = x
		g = c
		b = 0
	case section > 2 && section <= 3:
		r = 0
		g = c
		b = x
	case section > 3 && section <= 4:
		r = 0
		g = x
		b = c
	case section > 4 && section <= 5:
		r = x
		g = 0
		b = c
	case section > 5 && section <= 6:
		r = c
		g = 0
		b = x
	default:
		return 0, 0, 0
	}
	m := v - c
	r += m
	g += m
	b += m
	return int(r * 255), int(g * 255), int(b * 255)
}
