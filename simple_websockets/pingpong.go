func pingpong(ws *websocket.Conn) {
        defer func() {
                log.Println("ws connection closed")
                ws.Close()
        }()
        log.Println("ws connection established")
        buf := make([]byte, 1024)
        msg := "hi client"
        for {
                _, err := ws.Write([]byte(msg))
                log.Println(">", msg)
                if err != nil {
                        log.Println(err)
                        break
                }
                n, err := ws.Read(buf)
                if err != nil {
                        log.Println(err)
                        break
                }
                log.Println("<", string(buf[0:n]))
        }
}