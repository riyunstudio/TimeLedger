package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"timeLedger/app"
	"timeLedger/app/models"
)

type UserResource struct {
	BaseResource
	app *app.App
}

func NewUserResource(app *app.App) *UserResource {
	return &UserResource{
		app: app,
	}
}

type UserGetResource struct {
	ID         uint     `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Ips        []string `json:"ips,omitempty"`
	CreateTime int64    `json:"create_time,omitempty"`
	UpdateTime int64    `json:"update_time,omitempty"`
}

func (rs *UserResource) Get(ctx context.Context, datas models.User) (*UserGetResource, error) {
	ips := []string{}
	if err := json.Unmarshal([]byte(datas.Ips), &ips); err != nil {
		return nil, fmt.Errorf("JsonDecode ips err, Err: %v", err)
	}

	return &UserGetResource{
		ID:         datas.ID,
		Name:       datas.Name,
		Ips:        ips,
		CreateTime: datas.CreateTime,
		UpdateTime: datas.UpdateTime,
	}, nil
}

type UserCreateResource struct {
	ID        uint     `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Ips       []string `json:"ips,omitempty"`
	CreatedAt int64    `json:"created_at,omitempty"`
}

func (rs *UserResource) Create(ctx context.Context, datas models.User) (*UserCreateResource, error) {
	ips := []string{}
	if err := json.Unmarshal([]byte(datas.Ips), &ips); err != nil {
		return nil, fmt.Errorf("JsonDecode ips err, Err: %v", err)
	}

	return &UserCreateResource{
		ID:        datas.ID,
		Name:      datas.Name,
		Ips:       ips,
		CreatedAt: datas.CreateTime,
	}, nil
}

type UserUpdateResource struct {
	Name      string   `json:"name,omitempty"`
	Ips       []string `json:"ips,omitempty"`
	UpdatedAt int64    `json:"updated_at,omitempty"`
}

func (rs *UserResource) Update(ctx context.Context, datas models.User) (*UserUpdateResource, error) {
	ips := []string{}
	if err := json.Unmarshal([]byte(datas.Ips), &ips); err != nil {
		return nil, fmt.Errorf("JsonDecode ips err, Err: %v", err)
	}

	return &UserUpdateResource{
		Name:      datas.Name,
		Ips:       ips,
		UpdatedAt: datas.UpdateTime,
	}, nil
}
