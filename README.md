# Go_Laravel_Session
Go Laravel Session
Read Laravel's SessionID support for golang.

# Examples

```
func (this *WebSocketController) Join() {

    userId, err := ls.GetUserId(this.Ctx.GetCookie("laravel_session"))
    if err != nil {
        beego.Debug("Error read session: ", err)
        this.StopRun()
    }
    beego.Debug("Connecting User ID: ", userId)
	// Upgrade from http request to WebSocket.
	ws, err = websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
    clients[userId] = append(clients[userId], ws)


    if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
            beego.Debug("Disconnecting User ID: ", userId)
			this.StopRun()
		}
        beego.Debug(string(p))
	}
}
```
