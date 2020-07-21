package wscanvas

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

func NewHandler(address string, onConnection func (ctx *Context)) http.Handler {
	c := chi.NewRouter()
	c.Handle("/ws", &rpcHandler{
		onRpcConnection: func(ctx context.Context, conn *websocket.Conn) {
			onConnection(&Context{Context: ctx, conn: conn})
		},
	})
	c.Handle("/*", &staticHandler{
		webSocketAddress: address + "/ws",
	})
	return c
}
