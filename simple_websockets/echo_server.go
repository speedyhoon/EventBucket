func EchoServer(ws *websocket.Conn) {
  input := make([]byte, 100);
  for {
    sz, err := ws.Reader(input)
    ws.Write([]byte(fmt.Sprintf("You sent: %i", sz)))
  }
}