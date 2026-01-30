package resources

import (
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
)

// RoomResource 教室資源轉換
type RoomResource struct {
	app *app.App
}

// NewRoomResource 建立 RoomResource 實例
func NewRoomResource(appInstance *app.App) *RoomResource {
	return &RoomResource{
		app: appInstance,
	}
}

// RoomResponse 教室響應結構
type RoomResponse struct {
	ID        uint      `json:"id"`
	CenterID  uint      `json:"center_id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToRoomResponse 將教室模型轉換為響應格式
func (r *RoomResource) ToRoomResponse(room models.Room) *RoomResponse {
	return &RoomResponse{
		ID:        room.ID,
		CenterID:  room.CenterID,
		Name:      room.Name,
		Capacity:  room.Capacity,
		IsActive:  room.IsActive,
		CreatedAt: room.CreatedAt,
	}
}

// ToRoomResponses 批量將教室模型轉換為響應格式
func (r *RoomResource) ToRoomResponses(rooms []models.Room) []RoomResponse {
	if rooms == nil {
		return nil
	}

	responses := make([]RoomResponse, len(rooms))
	for i, room := range rooms {
		responses[i] = *r.ToRoomResponse(room)
	}
	return responses
}
