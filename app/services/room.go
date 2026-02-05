package services

import (
	"context"
	"time"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/global/errInfos"
)

// RoomService 教室管理相關業務邏輯
type RoomService struct {
	BaseService
	app            *app.App
	roomRepository *repositories.RoomRepository
	auditLogRepo   *repositories.AuditLogRepository
	cacheService   *CacheService
}

// NewRoomService 建立教室服務
func NewRoomService(app *app.App) *RoomService {
	return &RoomService{
		app:            app,
		roomRepository: repositories.NewRoomRepository(app),
		auditLogRepo:   repositories.NewAuditLogRepository(app),
		cacheService:   NewCacheService(app),
	}
}

// CreateRoomRequest 建立教室請求
type CreateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required"`
}

// UpdateRoomRequest 更新教室請求
type UpdateRoomRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int    `json:"capacity" binding:"required"`
}

// ToggleActiveRequest 啟用/停用請求
type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

// GetRooms 取得教室列表（使用快取）
func (s *RoomService) GetRooms(ctx context.Context, centerID uint, query string, page, limit int) ([]models.Room, int64, *errInfos.Res, error) {
	// 如果有查詢參數，直接跳過快取查詢資料庫
	if query != "" {
		s.Logger.Debug("room search query, skipping cache", "center_id", centerID, "query", query)
		rooms, total, err := s.roomRepository.SearchByNamePaginated(ctx, centerID, query, page, limit)
		if err != nil {
			return nil, 0, s.app.Err.New(errInfos.SQL_ERROR), err
		}
		return rooms, total, nil, nil
	}

	// 先從快取取得
	cached, err := s.cacheService.GetRoomList(ctx, centerID)
	if err == nil && cached != nil {
		// 將快取項目轉換為 models.Room
		rooms := make([]models.Room, 0, len(cached))
		for _, item := range cached {
			rooms = append(rooms, models.Room{
				ID:       item.ID,
				CenterID: centerID,
				Name:     item.Name,
				Capacity: item.Capacity,
				IsActive: item.IsActive,
			})
		}
		return rooms, int64(len(rooms)), nil, nil
	}

	// 快取未命中，從資料庫取得
	rooms, err := s.roomRepository.ListByCenterID(ctx, centerID)
	if err != nil {
		return nil, 0, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 存入快取
	cacheItems := make([]RoomCacheItem, 0, len(rooms))
	for _, r := range rooms {
		cacheItems = append(cacheItems, RoomCacheItem{
			ID:       r.ID,
			Name:     r.Name,
			Capacity: r.Capacity,
			IsActive: r.IsActive,
		})
	}
	_ = s.cacheService.SetRoomList(ctx, centerID, cacheItems)

	return rooms, int64(len(rooms)), nil, nil
}

// GetActiveRooms 取得已啟用的教室列表
func (s *RoomService) GetActiveRooms(ctx context.Context, centerID uint) ([]models.Room, *errInfos.Res, error) {
	rooms, err := s.roomRepository.ListActiveByCenterID(ctx, centerID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}
	return rooms, nil, nil
}

// CreateRoom 新增教室
func (s *RoomService) CreateRoom(ctx context.Context, centerID, adminID uint, req *CreateRoomRequest) (*models.Room, *errInfos.Res, error) {
	room := models.Room{
		CenterID:  centerID,
		Name:      req.Name,
		Capacity:  req.Capacity,
		CreatedAt: time.Now(),
	}

	createdRoom, err := s.roomRepository.Create(ctx, room)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除教室列表快取
	_ = s.cacheService.InvalidateRoomList(ctx, centerID)

	// 記錄審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "CREATE_ROOM",
		TargetType: "Room",
		TargetID:   createdRoom.ID,
		Payload: models.AuditPayload{
			After: room,
		},
	})

	return &createdRoom, nil, nil
}

// UpdateRoom 更新教室
func (s *RoomService) UpdateRoom(ctx context.Context, centerID, adminID, roomID uint, req *UpdateRoomRequest) (*models.Room, *errInfos.Res, error) {
	// 查詢現有教室
	room, err := s.roomRepository.GetByID(ctx, roomID)
	if err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 驗證權限
	if room.CenterID != centerID {
		return nil, s.app.Err.New(errInfos.FORBIDDEN), nil
	}

	// 更新
	room.Name = req.Name
	room.Capacity = req.Capacity

	if err := s.roomRepository.Update(ctx, room); err != nil {
		return nil, s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除教室列表快取
	_ = s.cacheService.InvalidateRoomList(ctx, centerID)

	// 記錄審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "UPDATE_ROOM",
		TargetType: "Room",
		TargetID:   room.ID,
		Payload: models.AuditPayload{
			Before: models.Room{Name: "", Capacity: 0},
			After:  room,
		},
	})

	return &room, nil, nil
}

// ToggleRoomActive 切換教室啟用狀態
func (s *RoomService) ToggleRoomActive(ctx context.Context, centerID, adminID, roomID uint, isActive bool) (*errInfos.Res, error) {
	if err := s.roomRepository.ToggleActive(ctx, roomID, centerID, isActive); err != nil {
		return s.app.Err.New(errInfos.SQL_ERROR), err
	}

	// 清除教室列表快取
	_ = s.cacheService.InvalidateRoomList(ctx, centerID)

	// 記錄審核日誌
	s.auditLogRepo.Create(ctx, models.AuditLog{
		CenterID:   centerID,
		ActorType:  "ADMIN",
		ActorID:    adminID,
		Action:     "TOGGLE_ROOM_ACTIVE",
		TargetType: "Room",
		TargetID:   roomID,
		Payload: models.AuditPayload{
			After: map[string]interface{}{
				"is_active": isActive,
			},
		},
	})

	return nil, nil
}
