package handler

import (
	"context"
	"encoding/json"
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

func ConvertCategoryInfoResponse(categoryList []model.Category, rows int64) (result []*proto.CategoryInfoResponse, err error) {
	if rows != 0 {
		for _, value := range categoryList {
			res := proto.CategoryInfoResponse{
				Id:             value.ID,
				Name:           value.Name,
				Level:          value.Level,
				IsTab:          value.IsTab,
				ParentCategory: value.ParentCategoryId,
			}

			result = append(result, &res)
		}
	}
	return result, nil
}

func (c *CategoryServer) GetAllCategorysList(ctx context.Context, request *proto.Empty) (*proto.CategoryListResponse, error) {
	categoryList, categoryRows, categoryErr := model.GetCategoryList("", []interface{}{}, "id,name,parent_category_id,level,is_tab", 0, 0)
	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.CategoryListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	result, resErr := ConvertCategoryInfoResponse(categoryList, categoryRows)

	if resErr != nil {
		return &proto.CategoryListResponse{}, status.Errorf(codes.Internal, "服务器内部出错", resErr.Error())
	}

	jsonStr, _ := json.Marshal(result)
	return &proto.CategoryListResponse{
		Data:     result,
		JsonData: string(jsonStr),
	}, nil
}

func (c *CategoryServer) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id=? and is_deleted=?", []interface{}{request.Id, 0}, "id,name,parent_category_id,level,is_tab")
	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.SubCategoryListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if categoryRows == 0 {
		return &proto.SubCategoryListResponse{}, status.Errorf(codes.NotFound, "分类不存在")
	}

	var where = make(map[string]interface{}, 0)
	ids := utils.GetMenuIds(request.Id)
	where["id in"] = ids
	whereSql, vals, _ := WhereBuild(where)
	categoryList, categoryRows, categoryErr := model.GetCategoryList(whereSql, vals, "id,name,parent_category_id,level,is_tab", 0, 0)

	resultSub, resErr := ConvertCategoryInfoResponse(categoryList, categoryRows)

	if resErr != nil {
		return &proto.SubCategoryListResponse{}, status.Errorf(codes.Internal, "服务器内部出错", resErr.Error())
	}
	resInfo := ConvertCategoryToRsp(categoryFirst)

	return &proto.SubCategoryListResponse{SubCategorys: resultSub, Info: &resInfo}, nil
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
