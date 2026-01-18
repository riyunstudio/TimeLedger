package services

import (
	"context"
	"timeLedger/app"
	"timeLedger/app/models"
	"timeLedger/app/repositories"
	"timeLedger/grpc/proto/user"
)

type User struct {
	user.UnimplementedUserServiceServer
	BaseService
	App            *app.App
	UserRepository *repositories.UserRepository
}

func (s *User) Get(ctx context.Context, req *user.GetRequest) (*user.GetResponse, error) {
	data, err := s.UserRepository.Get(ctx, models.User{ID: uint(req.GetID())})
	if err != nil {
		return &user.GetResponse{Code: 100, Msg: err.Error()}, err
	}

	return &user.GetResponse{
		Msg: "OK",
		Datas: &user.GetResponseDatas{
			ID:   int64(data.ID),
			Name: data.Name,
			Ips:  data.Ips,
		},
	}, nil
}
