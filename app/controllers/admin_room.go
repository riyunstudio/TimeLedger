package controllers

import (
	"timeLedger/app"
	"timeLedger/app/resources"
	"timeLedger/app/services"

	"github.com/gin-gonic/gin"
)

type AdminRoomController struct {
	BaseController
	app          *app.App
	roomService  *services.RoomService
	roomResource *resources.RoomResource
}

func NewAdminRoomController(app *app.App) *AdminRoomController {
	return &AdminRoomController{
		app:          app,
		roomService:  services.NewRoomService(app),
		roomResource: resources.NewRoomResource(app),
	}
}

// GetRooms 取得教室列表
// @Summary 取得教室列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.RoomResponse}
// @Router /api/v1/admin/rooms [get]
func (ctl *AdminRoomController) GetRooms(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	// 防止瀏覽器快取
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	rooms, errInfo, err := ctl.roomService.GetRooms(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.roomResource.ToRoomResponses(rooms)
	helper.Success(responses)
}

// CreateRoom 新增教室
// @Summary 新增教室
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateRoomRequest true "教室資訊"
// @Success 200 {object} global.ApiResponse{data=resources.RoomResponse}
// @Router /api/v1/admin/rooms [post]
func (ctl *AdminRoomController) CreateRoom(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	var req services.CreateRoomRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	room, errInfo, err := ctl.roomService.CreateRoom(ctx.Request.Context(), centerID, adminID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.roomResource.ToRoomResponse(*room)
	helper.Success(response)
}

// UpdateRoom 更新教室
// @Summary 更新教室
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room_id path int true "Room ID"
// @Param request body services.UpdateRoomRequest true "教室資訊"
// @Success 200 {object} global.ApiResponse{data=resources.RoomResponse}
// @Router /api/v1/admin/rooms/{room_id} [put]
func (ctl *AdminRoomController) UpdateRoom(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	roomID := helper.MustParamUint("room_id")
	if roomID == 0 {
		return
	}

	var req services.UpdateRoomRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	room, errInfo, err := ctl.roomService.UpdateRoom(ctx.Request.Context(), centerID, adminID, roomID, &req)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	response := ctl.roomResource.ToRoomResponse(*room)
	helper.Success(response)
}

// GetActiveRooms 取得已啟用的教室列表
// @Summary 取得已啟用的教室列表
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} global.ApiResponse{data=[]resources.RoomResponse}
// @Router /api/v1/admin/rooms/active [get]
func (ctl *AdminRoomController) GetActiveRooms(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	// 防止瀏覽器快取
	ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "0")

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	rooms, errInfo, err := ctl.roomService.GetActiveRooms(ctx.Request.Context(), centerID)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	responses := ctl.roomResource.ToRoomResponses(rooms)
	helper.Success(responses)
}

// ToggleRoomActive 切換教室啟用狀態
// @Summary 切換教室啟用狀態
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param room_id path int true "Room ID"
// @Param request body services.ToggleActiveRequest true "啟用狀態"
// @Success 200 {object} global.ApiResponse
// @Router /api/v1/admin/rooms/{room_id}/toggle-active [patch]
func (ctl *AdminRoomController) ToggleRoomActive(ctx *gin.Context) {
	helper := NewContextHelper(ctx)

	centerID := helper.MustCenterID()
	if centerID == 0 {
		return
	}

	roomID := helper.MustParamUint("room_id")
	if roomID == 0 {
		return
	}

	var req services.ToggleActiveRequest
	if !helper.MustBindJSON(&req) {
		return
	}

	adminID := helper.MustUserID()
	if adminID == 0 {
		return
	}

	errInfo, err := ctl.roomService.ToggleRoomActive(ctx.Request.Context(), centerID, adminID, roomID, req.IsActive)
	if err != nil {
		helper.ErrorWithInfo(errInfo)
		return
	}

	helper.Success(nil)
}
