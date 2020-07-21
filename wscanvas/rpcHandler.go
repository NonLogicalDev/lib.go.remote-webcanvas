package wscanvas

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type rpcHandler struct {
	onRpcConnection func(ctx context.Context, conn *websocket.Conn)
}

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (rpc *rpcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	rpc.onRpcConnection(r.Context(), conn)

}
