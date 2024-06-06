package websocket

//
//type Room struct {
//	ID      string             `json:"id"`
//	Name    string             `json:"name"`
//	Clients map[string]*Client `json:"clients"`
//}
//
//type Hub struct {
//	Rooms      map[string]*Room `json:"room"`
//	Register   chan *Client
//	Unregister chan *Client
//	Broadcast  chan *Message
//}
//
//func NewHub() *Hub {
//	return &Hub{
//		Rooms:      make(map[string]*Room),
//		Register:   make(chan *Client),
//		Unregister: make(chan *Client),
//		Broadcast:  make(chan *Message, 5),
//	}
//}
//
//func (hub *Hub) Run() {
//	for {
//		select {
//		case cl := <-hub.Register:
//			if _, ok := hub.Rooms[cl.RoomID]; ok {
//				r := hub.Rooms[cl.RoomID]
//
//				if _, ok := r.Clients[cl.ID]; !ok {
//					r.Clients[r.ID] = cl
//				}
//			}
//
//		case cl := <-hub.Unregister:
//			if _, ok := hub.Rooms[cl.RoomID]; ok {
//				if _, ok := hub.Rooms[cl.RoomID].Clients[cl.ID]; ok {
//					if len(hub.Rooms[cl.RoomID].Clients) != 0 {
//						hub.Broadcast <- &Message{
//							Content:  "user left a chat",
//							RoomID:   cl.RoomID,
//							Username: cl.Username,
//						}
//					}
//					delete(hub.Rooms[cl.RoomID].Clients, cl.ID)
//					close(cl.Message)
//				}
//			}
//
//		case m := <-hub.Broadcast:
//			if _, ok := hub.Rooms[m.RoomID]; ok {
//
//				for _, cl := range hub.Rooms[m.RoomID].Clients {
//					cl.Message <- m
//				}
//			}
//		}
//	}
//}
