package client

import (
	"context"
	"fmt"

	protobuffer "github.com/mxbikes/protobuf/user"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	Client protobuffer.UserServiceClient
}

func InitUserServiceClient(url string) UserServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := UserServiceClient{
		Client: protobuffer.NewUserServiceClient(cc),
	}

	return c
}

func (c *UserServiceClient) GetUserByID(id string) (*protobuffer.GetUserByIDResponse, error) {
	req := &protobuffer.GetUserByIDRequest{
		ID: id,
	}

	return c.Client.GetUserByID(context.Background(), req)
}
