package websocket

//
//import (
//	"github.com/gorilla/websocket"
//	"github.com/labstack/echo/v4"
//	"net/http"
//)
//
//type Handler struct {
//	hub *Hub
//}
//
//func NewHandler(h *Hub) *Handler {
//	return &Handler{
//		hub: h,
//	}
//}
//
//type CreateRoomRequest struct {
//	ID       string `json:"id"`
//	RoomName string `json:"roomName"`
//}
//
//func (h *Handler) CreateRoom(c echo.Context) error {
//	var req CreateRoomRequest
//	if err := c.Bind(&req); err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
//	}
//	h.hub.Rooms[req.ID] = &Room{
//		ID:      req.ID,
//		Name:    req.RoomName,
//		Clients: make(map[string]*Client),
//	}
//	return c.JSON(http.StatusOK, req)
//}
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin: func(r *http.Request) bool {
//		return r.Header.Get("Origin") == "http://localhost:5173"
//	},
//}
//
//func (h *Handler) JoinRoom(c echo.Context) error {
//	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
//	}
//	roomId := c.Param("roomId")
//	clientID := c.QueryParam("userId")
//	username := c.QueryParam("username")
//
//	cl := &Client{
//		Conn:     conn,
//		Message:  make(chan *Message, 10),
//		ID:       clientID,
//		RoomID:   roomId,
//		Username: username,
//	}
//
//	m := &Message{
//		Content:  "New user has joined a Room",
//		RoomID:   roomId,
//		Username: username,
//	}
//
//	h.hub.Register <- cl
//
//	h.hub.Broadcast <- m
//
//	go cl.writeMessage()
//	cl.readMessage(h.hub)
//	return err
//}
//
//type RoomResponse struct {
//	ID       string `json:"id"`
//	RoomName string `json:"roomName"`
//}
//
//func (h *Handler) GetRooms(c echo.Context) error {
//	rooms := make([]*RoomResponse, 0)
//
//	for _, room := range h.hub.Rooms {
//		rooms = append(rooms, &RoomResponse{
//			ID:       room.ID,
//			RoomName: room.Name,
//		})
//	}
//	return c.JSON(http.StatusOK, rooms)
//}
//
//type ClientsResponse struct {
//	ID       string `json:"id"`
//	Username string `json:"username"`
//}
//
//func (h *Handler) GetClients(c echo.Context) error {
//	var clients []*ClientsResponse
//	roomId := c.Param("roomId")
//
//	if _, ok := h.hub.Rooms[roomId]; !ok {
//		clients = make([]*ClientsResponse, 0)
//		return c.JSON(http.StatusOK, clients)
//	}
//
//	for _, c := range h.hub.Rooms[roomId].Clients {
//		clients = append(clients, &ClientsResponse{
//			ID:       c.ID,
//			Username: c.Username,
//		})
//	}
//
//	return c.JSON(http.StatusOK, clients)
//}
