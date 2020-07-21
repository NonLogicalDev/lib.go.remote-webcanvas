(function (window, document) {
  let HOST_ADDRESS = "ws://%%WEB_SOCKET_ADDRESS%%"

  function createUI() {
    let body = document.querySelector("body")

    let ui = {
      canvas: document.createElement("canvas"),
      buttons: document.createElement("div"),
    }

    body.appendChild(ui.canvas);
    body.appendChild(ui.buttons);

    return ui;
  }

  function connectRPC(url, onCmd) {
    let ws = new WebSocket(url);
    ws.onopen = function (evt) {
      ws.send(JSON.stringify({
        type: "event.client.lifecycle",
        ready: true,
      }))
    }
    ws.onmessage = function (evt) {
      let rawCmds = evt.data.split("\n");
      for (const rawCmd of rawCmds) {
        if (rawCmd.trim().length > 0) {
          onCmd(ws, JSON.parse(rawCmd));
        }
      }
    }
    ws.onerror = function (evt) {
      console.error("WS: ERROR: " + evt.data);
    }
  }

  function rpcCMD(ws, ui, cmd) {
    let ctx = ui.canvas.getContext('2d');

    switch (cmd.name) {
      case "setCanvasSize":
        setCanvasSize(ctx, cmd)
        return
      case "drawBoxFill":
        drawBoxFill(ctx, cmd)
        return
      case "setButtons":
        setButtons(ws, ui, cmd)
        return
    }
  }

  function setButtons(ws, ui, cmd) {
    // Clear out old buttons:
    ui.buttons.querySelectorAll('*').forEach(n => n.remove());

    for(const button of cmd.buttons) {
      let b = document.createElement("button");
      b.textContent = button.name;
      b.onclick = function (e) {
        ws.send(JSON.stringify({
          type: "event.button",
          source: button.callback_id,
        }))
      }
      ui.buttons.appendChild(b);
    }
  }

  function setCanvasSize(ctx, cmd) {
    ctx.canvas.width = cmd.w;
    ctx.canvas.height = cmd.h;
  }

  function drawBoxFill(ctx, cmd) {
    let r = cmd.rect;
    let c = cmd.color;
    setFillStyle(ctx, c);
    ctx.fillRect(r.x, r.y, r.w, r.h);
  }

  // Utiltities:
  function colorToStyle(c) {
    return `rgba(${c.r}, ${c.g}, ${c.b}, ${c.a / 255})`;
  }

  let setFillStyleCache = "";
  function setFillStyle(ctx, color) {
    let style = colorToStyle(color);
    if (setFillStyleCache !== style) {
      ctx.fillStyle = style;
    }
  }


  // Main
  function main() {
    let ui = createUI();

    // EXPORT FOR DEBUGGING:
    window.ui = ui;

    connectRPC(HOST_ADDRESS,
      (ws, cmd) => {
        rpcCMD(ws, ui, cmd)
      },
    );
  }

  main();
})(window, document)