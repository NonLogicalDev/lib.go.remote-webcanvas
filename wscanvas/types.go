package wscanvas

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	A int `json:"a"`
}

type Style struct {
	FillColor   *Color `json:"fill_color,omitempty"`
	StrokeColor *Color `json:"stroke_color,omitempty"`
	StrokeWidth int    `json:"stroke_width"`
}

type Text struct {
	Text     string `json:"text"`
	Font     string `json:"font"`
	Size     int    `json:"size"`
	Align    string `json:"align"`
	Baseline string `json:"baseline"`
}

type Button struct {
	Name       string `json:"name"`
	CallbackID string `json:"callback_id"`
}
