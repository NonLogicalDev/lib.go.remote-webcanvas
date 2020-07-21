package wscanvas

type Point struct {
	X int
	Y int
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

type Button struct {
	Name       string `json:"name"`
	CallbackID string `json:"callback_id"`
}
