package services

import (
	"akali/app"
	"akali/app/models"
	"akali/app/repositories"
	"akali/app/requests"
	"akali/app/resources"
	"akali/global/errInfos"
	"encoding/json"
	"time"
)

type UserService struct {
	BaseService
	app            *app.App
	userRepository *repositories.UserRepository
	userResource   *resources.UserResource
}

func NewUserService(app *app.App) *UserService {
	return &UserService{
		app:            app,
		userRepository: repositories.NewUserRepository(app),
		userResource:   resources.NewUserResource(app),
	}
}

func (s *UserService) Get(req *requests.UserGetRequest) (datas any, errInfo *errInfos.Err, err error) {
	time.Sleep(10 * time.Second)
	user, err := s.userRepository.Get(models.User{ID: uint(req.ID)})
	if err != nil {
		return nil, errInfos.New(errInfos.SQL_ERROR), err
	}

	if user.ID == 0 {
		return nil, errInfos.New(errInfos.USER_NOT_FOUNT), nil
	}

	response, err := s.userResource.Get(user)
	if err != nil {
		return nil, errInfos.New(errInfos.FORMAT_RESOURCE_ERROR), err
	}

	return response, nil, nil
}

func (s *UserService) Create(req *requests.UserCreateRequest) (datas any, errInfo *errInfos.Err, err error) {
	ips, err := json.Marshal(req.Ips)
	if err != nil {
		return nil, errInfos.New(errInfos.JSON_ENCODE_ERROR), err
	}

	user, err := s.userRepository.Create(models.User{
		Name:       req.Name,
		Ips:        string(ips),
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	})
	if err != nil {
		return nil, errInfos.New(errInfos.SQL_ERROR), err
	}

	response, err := s.userResource.Create(user)
	if err != nil {
		return nil, errInfos.New(errInfos.FORMAT_RESOURCE_ERROR), err
	}

	return response, nil, nil
}

func (s *UserService) Update(req *requests.UserUpdateRequest) (datas any, errInfo *errInfos.Err, err error) {
	user, err := s.userRepository.Get(models.User{ID: uint(req.Id)})
	if err != nil {
		return nil, errInfos.New(errInfos.SQL_ERROR), err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if merge, err := s.app.Tools.MergeJSONStrings(user.Ips, req.Ips); err != nil {
		return nil, errInfos.New(errInfos.JSON_PROCESS_ERROR), err
	} else {
		user.Ips = merge
	}

	user.UpdateTime = s.app.Tools.NowUnix()

	err = s.userRepository.UpdateById(user)
	if err != nil {
		return nil, errInfos.New(errInfos.SQL_ERROR), err
	}

	response, err := s.userResource.Update(user)
	if err != nil {
		return nil, errInfos.New(errInfos.FORMAT_RESOURCE_ERROR), err
	}

	return response, nil, nil
}
