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
      case "drawBox":
        drawBox(ctx, cmd)
        return
      case "drawText":
        drawText(ctx, cmd)
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

  function drawText(ctx, cmd) {
    let p = cmd.point;
    let s = cmd.style;
    let t = cmd.text;

    if (t.font === "") {
      t.font = "sans-serif";
    }
    if (t.size === 0) {
      t.size = 10
    }
    if (t.align === "") {
      t.align = "center"
    }
    if (t.baseline === "") {
      t.baseline = "middle"
    }

    let font = `${t.size}px ${t.font}`;

    ctx.font = font;
    ctx.textAlign = t.align;
    ctx.textBaseline = t.baseline;

    console.log(font, cmd, ctx);

    if (s.fill_color) {
      setFillStyle(ctx, s.fill_color);
      ctx.fillText(t.text, p.x, p.y);
    }
    if (s.stroke_color) {
      setStrokeStyle(ctx, s.stroke_color, s.stroke_width);
      ctx.strokeText(t.text, p.x, p.y);
    }

    window.ctx = ctx;
  }

  function drawBox(ctx, cmd) {
    let r = cmd.rect;
    let s = cmd.style;

    if (s.fill_color) {
      setFillStyle(ctx, s.fill_color);
      ctx.fillRect(r.x, r.y, r.w, r.h);
    }
    if (s.stroke_color) {
      setStrokeStyle(ctx, s.stroke_color, s.stroke_width);
      ctx.strokeRect(r.x, r.y, r.w, r.h);
    }
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

  let setStrokeStyleCache = "";
  function setStrokeStyle(ctx, color, width) {
    let style = colorToStyle(color);
    if (setStrokeStyleCache !== style) {
      ctx.lineWidth = width;
      ctx.strokeStyle = style;
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