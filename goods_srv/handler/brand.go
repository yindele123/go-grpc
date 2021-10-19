package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/goods_srv/model"
	"project/goods_srv/proto"
	"time"
)

type BrandServer struct {
}

func (b *BrandServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var offset int32 = 0
	var limit int32 = 10
	if request.Pages != 0 {
		limit = request.Pages
	}
	if request.PagePerNums != 0 {
		offset = limit * (request.PagePerNums - 1)
	}
	brandList, brandRow, brandErr := model.GetBrandsList("", []interface{}{}, "id,name,logo", int(offset), int(limit))
	if brandErr != nil {
		zap.S().Error("服务器内部出错", brandErr.Error())
		return &proto.BrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, countErr := model.GetBrandsCount("", []interface{}{})
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.BrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	result := make([]*proto.BrandInfoResponse, 0)
	if brandRow != 0 {
		for _, value := range brandList {
			res := &proto.BrandInfoResponse{
				Id:   value.ID,
				Name: value.Name,
				Logo: value.Logo,
			}
			result = append(result, res)
		}
	}
	return &proto.BrandListResponse{Total: total, Data: result}, nil
}

func (b *BrandServer) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	_, brandsRow, brandsErr := model.GetBrandsFirst("name=?", []interface{}{request.Name}, "id")
	if brandsErr != nil {
		zap.S().Error("服务器内部出错", brandsErr.Error())
		return &proto.BrandInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if brandsRow != 0 {
		return &proto.BrandInfoResponse{}, status.Errorf(codes.AlreadyExists, "记录已经存在")
	}

	resCreate, er := model.CreateBrands(model.Brands{
		Name:      request.Name,
		Logo:      request.Logo,
		CreatedAt: uint32(time.Now().Unix()),
		UpdatedAt: uint32(time.Now().Unix()),
	})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.BrandInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.BrandInfoResponse{
		Id:   resCreate.ID,
		Name: resCreate.Name,
		Logo: resCreate.Logo,
	}, nil
}

func (b *BrandServer) DeleteBrand(ctx context.Context, rq *proto.BrandRequest) (*proto.Empty, error) {
	brandFirst, brandRows, brandErr := model.GetBrandsFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if brandErr != nil {
		zap.S().Error("服务器内部出错", brandErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if brandRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	er := model.UpdateBrands(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{brandFirst.ID})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (b *BrandServer) UpdateBrand(ctx context.Context, rq *proto.BrandRequest) (*proto.Empty, error) {
	brandFirst, brandRows, brandErr := model.GetBrandsFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if brandErr != nil {
		zap.S().Error("服务器内部出错", brandErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if brandRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}

	updateGoodsIsErr := model.UpdateBrands(model.Brands{
		Name:      rq.Name,
		Logo:      rq.Logo,
		UpdatedAt: uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{brandFirst.ID})
	if updateGoodsIsErr != nil {
		zap.S().Error("服务器内部出错", updateGoodsIsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}

func (b *BrandServer) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var offset int32 = 0
	var limit int32 = 10
	if request.Pages != 0 {
		limit = request.Pages
	}
	if request.PagePerNums != 0 {
		offset = limit * (request.PagePerNums - 1)
	}
	fmt.Println(offset)
	handler := model.GetSearchModelHandler(&model.Goodscategorybrand{})
	handler.Search()
	//goodsList, rows, goodsErr := model.GetGoodsList("",[]interface{}{}, "", int(offset), int(limit))
	return &proto.CategoryBrandListResponse{},nil
}

func (b *BrandServer) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	panic("implement me")
}

func (b *BrandServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	panic("implement me")
}

func (b *BrandServer) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.Empty, error) {
	panic("implement me")
}

func (b *BrandServer) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.Empty, error) {
	panic("implement me")
}
