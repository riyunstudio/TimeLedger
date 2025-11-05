package services

import (
	"akali/app"
	"akali/app/models"
	"akali/app/repositories"
	"akali/grpc/proto/user"
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

type User struct {
	user.UnimplementedUserServiceServer
	BaseService
	App            *app.App
	UserRepository *repositories.UserRepository
}

func (s *User) TestTimeout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return RunWithTimeout(ctx, 5*time.Second, func(ctx context.Context, do func(func() error) error) (*emptypb.Empty, error) {
		for i := 1; i <= 15; i++ {
			if err := do(func() error {
				time.Sleep(1 * time.Second)
				fmt.Println(i)
				return nil
			}); err != nil {
				return nil, err
			}
		}
		return &emptypb.Empty{}, nil
	})
}

func (s *User) Get(ctx context.Context, req *user.GetRequest) (*user.GetResponse, error) {
	return RunWithTimeout(ctx, 5*time.Second, func(ctx context.Context, do func(func() error) error) (*user.GetResponse, error) {
		data, err := s.UserRepository.Get(models.User{ID: uint(req.GetID())})
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
	})
}
