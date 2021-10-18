package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"project/goods_srv/proto"
	"project/goods_srv/utils"
)

type CategoryServer struct {
}

func (c *CategoryServer) GetAllCategorysList(ctx context.Context, request *proto.CategoryListRequest) (*proto.CategoryListResponse, error) {
	categoryList := utils.GetMenu(request.Id)
	fmt.Println(categoryList)
	jsonStr, _ := json.Marshal(categoryList)
	return &proto.CategoryListResponse{
		Data:     categoryList,
		JsonData: string(jsonStr),
	}, nil
}

func (c *CategoryServer) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	panic("implement me")
}

func (c *CategoryServer) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*proto.Empty, error) {
	panic("implement me")
}

func (c *CategoryServer) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.Empty, error) {
	panic("implement me")
}
