package requests

import (
	"akali/app"
	"akali/app/repositories"
	"akali/global/errInfos"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	BaseRequest
	app            *app.App
	userRepository *repositories.UserRepository
}

func NewUserRequest(app *app.App) *UserRequest {
	return &UserRequest{
		app:            app,
		userRepository: repositories.NewUserRepository(app),
	}
}

type UserGetRequest struct {
	ID int `json:"id"`
}

func (r *UserRequest) Get(ctx *gin.Context) (req *UserGetRequest, errInfo *errInfos.Res, err error) {
	req, err = Validate[UserGetRequest](ctx)
	if err != nil {
		return nil, r.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}

	return
}

type UserCreateRequest struct {
	Name  string   `json:"name" binding:"required"`
	Ips   []string `json:"ips"`
	Games []string `json:"game"`
}

func (r *UserRequest) Create(ctx *gin.Context) (req *UserCreateRequest, errInfo *errInfos.Res, err error) {
	req, err = Validate[UserCreateRequest](ctx)
	if err != nil {
		return nil, r.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}

	// 不合法的ip直接過濾
	if len(req.Ips) != 0 {
		var filteredIps []string
		for _, ip := range req.Ips {
			if r.app.Tools.IsIPv4(ip) {
				filteredIps = append(filteredIps, ip)
			}
		}
		req.Ips = filteredIps
	}

	return
}

type UserUpdateRequest struct {
	Id    int      `json:"id" binding:"required"`
	Name  string   `json:"name"`
	Ips   []string `json:"ips"`
	Games []string `json:"game"`
}

func (r *UserRequest) Update(ctx *gin.Context) (req *UserUpdateRequest, errInfo *errInfos.Res, err error) {
	req, err = Validate[UserUpdateRequest](ctx)
	if err != nil {
		return nil, r.app.Err.New(errInfos.PARAMS_VALIDATE_ERROR), err
	}

	// 不合法的ip直接過濾
	if len(req.Ips) != 0 {
		var filteredIps []string
		for _, ip := range req.Ips {
			if r.app.Tools.IsIPv4(ip) {
				filteredIps = append(filteredIps, ip)
			}
		}
		req.Ips = filteredIps
	}

	return
}
