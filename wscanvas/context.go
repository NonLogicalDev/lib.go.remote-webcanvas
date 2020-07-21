package wscanvas

import (
	"bytes"
	"context"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type CMDBase struct {
	Name string `json:"name"`
}

type Context struct {
	context.Context
	conn   *websocket.Conn
	buffer []interface{}
}

/*
	Commands:
*/

type CMDDrawBox struct {
	CMDBase
	Rect  Rect   `json:"rect"`
	Style Style `json:"style"`
}

func (ctx *Context) DrawBox(r Rect, s Style) error {
	return ctx.sendCMD(CMDDrawBox{
		CMDBase: CMDBase{
			Name: "drawBox",
		},
		Rect:  r,
		Style: s,
	})
}

type CMDDrawText struct {
	CMDBase
	Text  Text  `json:"text"`
	Point Point `json:"point"`
	Style Style `json:"style"`
}

func (ctx *Context) DrawText(text Text, point Point, s Style) error {
	return ctx.sendCMD(CMDDrawText{
		CMDBase: CMDBase{
			Name: "drawText",
		},
		Text:    text,
		Point:   point,
		Style:   s,
	})
}

type CMDSetCanvasSize struct {
	CMDBase
	W int `json:"w"`
	H int `json:"h"`
}

func (ctx *Context) SetCanvasSize(size Point) error {
	return ctx.sendCMD(CMDSetCanvasSize{
		CMDBase: CMDBase{
			Name: "setCanvasSize",
		},
		W: size.X,
		H: size.Y,
	})
}

type CMDSetButtons struct {
	CMDBase
	Buttons []Button `json:"buttons"`
}

func (ctx *Context) SetButtons(buttons []Button) error {
	return ctx.sendCMD(CMDSetButtons{
		CMDBase: CMDBase{
			Name: "setButtons",
		},
		Buttons: buttons,
	})
}

type RPCEvent struct {
	Type string `json:"type"`
}

func (ctx *Context) OnEvent(f func(eventType string, eventData []byte) error) error {
	_, data, err := ctx.conn.ReadMessage()
	if err != nil {
		return err
	}
	event := new(RPCEvent)
	err = json.Unmarshal(data, event)
	if err != nil {
		return err
	}
	return f(event.Type, data)
}

/*
	Buffer Implementation:
*/

func (ctx *Context) BufferStart() {
	if ctx.buffer == nil {
		ctx.buffer = make([]interface{}, 0)
	}
}

func (ctx *Context) BufferWrite() error {
	err := ctx.writeRPCCMD(ctx.buffer...)
	if err == nil {
		ctx.BufferClear()
	}
	return err
}

func (ctx *Context) BufferClear() {
	ctx.buffer = nil
}

func (ctx *Context) sendCMD(cmd interface{}) error {
	if ctx.buffer != nil {
		ctx.buffer = append(ctx.buffer, cmd)
		return nil
	}
	return ctx.writeRPCCMD(cmd)
}

/*
	Wire Implementation:
*/

func (ctx *Context) writeRPCCMD(cmds... interface{}) error {
	b := bytes.NewBuffer(nil)

	// JSON Implementation:
	e := json.NewEncoder(b)
	for _, cmd := range cmds {
		e.Encode(cmd)
		b.WriteRune('\n')
	}

	return ctx.conn.WriteMessage(websocket.TextMessage, b.Bytes())
}

