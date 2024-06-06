package websocket

//
//import (
//	"github.com/gorilla/websocket"
//	"github.com/labstack/gommon/log"
//)
//
//type Client struct {
//	Conn     *websocket.Conn
//	Message  chan *Message
//	ID       string `json:"id"`
//	RoomID   string `json:"room_id"`
//	Username string `json:"username"`
//}
//
//type Message struct {
//	Content  string `json:"content"`
//	RoomID   string `json:"roomID"`
//	Username string `json:"username"`
//}
//
//func (c *Client) writeMessage() {
//	defer func() {
//		err := c.Conn.Close()
//		if err != nil {
//			return
//		}
//	}()
//
//	for {
//		message, ok := <-c.Message
//		if !ok {
//			return
//		}
//
//		err := c.Conn.WriteJSON(message)
//		if err != nil {
//			return
//		}
//	}
//}
//
//func (c *Client) readMessage(hub *Hub) {
//	defer func() {
//		hub.Unregister <- c
//		err := c.Conn.Close()
//		if err != nil {
//			return
//		}
//	}()
//	for {
//		_, m, err := c.Conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("error: %v", err)
//			}
//			break
//		}
//		msg := &Message{
//			Content:  string(m),
//			RoomID:   c.RoomID,
//			Username: c.Username,
//		}
//		hub.Broadcast <- msg
//	}
//}
