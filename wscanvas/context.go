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
	buffer *bytes.Buffer
}

type CMDDrawBoxFill struct {
	CMDBase
	Rect  Rect   `json:"rect"`
	Color *Color `json:"color,omitempty"`
}

func (ctx *Context) BufferStart() {
	if ctx.buffer == nil {
		ctx.buffer = bytes.NewBuffer(nil)
	}
}

func (ctx *Context) BufferWrite() error {
	cmd := ctx.buffer.Bytes()
	err := ctx.writeRPCCMD(cmd)
	if err == nil {
		ctx.BufferClear()
	}
	return err
}

func (ctx *Context) BufferClear() {
	ctx.buffer = nil
}

func (ctx *Context) writeRPCCMD(cmd []byte) error {
	return ctx.conn.WriteMessage(websocket.TextMessage, cmd)
}

func (ctx *Context) sendCMD(cmd []byte) error {
	if ctx.buffer != nil {
		ctx.buffer.WriteRune('\n')
		ctx.buffer.Write(cmd)
		return nil
	}
	return ctx.writeRPCCMD(cmd)
}

func (ctx *Context) DrawBoxFill(r Rect, c *Color) error {
	cmd, _ := json.Marshal(CMDDrawBoxFill{
		CMDBase: CMDBase{
			Name: "drawBoxFill",
		},
		Rect:  r,
		Color: c,
	})
	return ctx.sendCMD(cmd)
}

type CMDSetCanvasSize struct {
	CMDBase
	W int `json:"w"`
	H int `json:"h"`
}

func (ctx *Context) SetCanvasSize(size Point) error {
	cmd, _ := json.Marshal(CMDSetCanvasSize{
		CMDBase: CMDBase{
			Name: "setCanvasSize",
		},
		W: size.X,
		H: size.Y,
	})
	return ctx.sendCMD(cmd)
}

type CMDSetButtons struct {
	CMDBase
	Buttons []Button `json:"buttons"`
}

func (ctx *Context) SetButtons(buttons []Button) error {
	cmd, _ := json.Marshal(CMDSetButtons{
		CMDBase: CMDBase{
			Name: "setButtons",
		},
		Buttons: buttons,
	})
	return ctx.sendCMD(cmd)
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
