package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/goods_srv/model"
	"project/goods_srv/proto"
	"project/goods_srv/utils"
	"time"
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

	resCreate, er := model.CreateCategory(model.Category{
		Name:             request.Name,
		ParentCategoryId: request.ParentCategory,
		Level:            request.Level,
		IsTab:            request.IsTab,
		CreatedAt:        uint32(time.Now().Unix()),
		UpdatedAt:        uint32(time.Now().Unix()),
	})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.CategoryInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.CategoryInfoResponse{
		Id:             resCreate.ID,
		Name:           resCreate.Name,
		ParentCategory: resCreate.ParentCategoryId,
		Level:          resCreate.Level,
		IsTab:          resCreate.IsTab,
	}, nil
}

func (c *CategoryServer) DeleteCategory(ctx context.Context, rq *proto.DeleteCategoryRequest) (*proto.Empty, error) {
	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if categoryRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "分类不存在")
	}
	er := model.UpdateCategory(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{categoryFirst.ID})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (c *CategoryServer) UpdateCategory(ctx context.Context, rq *proto.CategoryInfoRequest) (*proto.Empty, error) {
	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id,name")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if categoryRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	updateGoodsIsErr := model.UpdateCategory(model.Category{
		ID:               rq.Id,
		Name:             rq.Name,
		ParentCategoryId: rq.ParentCategory,
		Level:            rq.Level,
		IsTab:            rq.IsTab,
		UpdatedAt:        uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{categoryFirst.ID})
	if updateGoodsIsErr != nil {
		zap.S().Error("服务器内部出错", updateGoodsIsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}
