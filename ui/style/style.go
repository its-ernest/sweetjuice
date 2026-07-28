package style

type Weight string

const (
	Normal Weight = "normal"
	Bold   Weight = "bold"
	Light  Weight = "light"
)

type Alignment string

const (
	Center Alignment = "center"
	Start  Alignment = "start"
	End    Alignment = "end"
)

// View style properties
type View struct {
	Flex            int       `json:"flex,omitempty"`
	Padding         float64   `json:"padding,omitempty"`
	PaddingVertical float64   `json:"paddingVertical,omitempty"`
	PaddingHorizontal float64 `json:"paddingHorizontal,omitempty"`
	BackgroundColor string    `json:"backgroundColor,omitempty"`
	CornerRadius    float64   `json:"cornerRadius,omitempty"`
	AlignItems      Alignment `json:"alignItems,omitempty"`
	JustifyContent  Alignment `json:"justifyContent,omitempty"`
	Width           float64   `json:"width,omitempty"`
	Height          float64   `json:"height,omitempty"`
}

// Text style properties
type Text struct {
	FontSize float64 `json:"fontSize,omitempty"`
	Weight   Weight  `json:"weight,omitempty"`
	Color    string  `json:"color,omitempty"`
}

// Button style properties
type Button struct {
	BackgroundColor   string  `json:"backgroundColor,omitempty"`
	PaddingVertical   float64 `json:"paddingVertical,omitempty"`
	PaddingHorizontal float64 `json:"paddingHorizontal,omitempty"`
	CornerRadius      float64 `json:"cornerRadius,omitempty"`
	Color             string  `json:"color,omitempty"`
}

type TextButton struct {
	Button
}

type OutlinedButton struct {
	Button
	StrokeWidth float64 `json:"strokeWidth,omitempty"`
	StrokeColor string  `json:"strokeColor,omitempty"`
}

type TonalButton struct {
	Button
}

type ElevatedButton struct {
	Button
}

type IconButton struct {
	Icon    string  `json:"icon,omitempty"`
	Size    float64 `json:"size,omitempty"`
	Tonal   bool    `json:"tonal,omitempty"`
	Outlined bool   `json:"outlined,omitempty"`
}

type FabStyle struct {
	Extended bool   `json:"extended,omitempty"`
	Text     string `json:"text,omitempty"`
	Size     string `json:"size,omitempty"` // small, normal, large
}

type SegmentedButton struct {
	Options          []string `json:"options,omitempty"`
	Selected         string   `json:"selected,omitempty"`
	SingleSelection  bool     `json:"singleSelection,omitempty"`
}

type ButtonGroup struct {
	Orientation string `json:"orientation,omitempty"` // horizontal, vertical
}
