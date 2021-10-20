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

func ConvertCategoryBrandResponse(GoodscategorybrandList []model.Goodscategorybrand, goodscategorybrandRows uint32) (result []*proto.CategoryBrandResponse, err error) {
	categoryIds := make([]uint32, 0)
	BrandIds := make([]uint32, 0)
	if goodscategorybrandRows != 0 {
		for _, value := range GoodscategorybrandList {
			categoryIds = append(categoryIds, value.CategoryId)
			BrandIds = append(BrandIds, value.BrandId)
		}
	}
	brandsConvert, brandsErr := ConvertBrands(RemoveDuplicateElement(BrandIds), "id,name,logo", "Id")
	if brandsErr != nil {
		zap.S().Error(brandsErr.Error())
		return result, brandsErr
	}

	categoryConvert, categoryErr := ConvertCategory(RemoveDuplicateElement(categoryIds), "id,name,parent_category_id,level,is_tab", "Id")
	if categoryErr != nil {
		zap.S().Error(categoryErr.Error())
		return result, categoryErr
	}

	resdata := make([]*proto.CategoryBrandResponse, 0)
	for _, value := range GoodscategorybrandList {
		brandData := model.Brands{}
		categoryData := model.Category{}
		if _, ok := brandsConvert[fmt.Sprint(value.BrandId)]; ok {
			brandData = brandsConvert[fmt.Sprint(value.BrandId)][0].(model.Brands)
		}
		if _, ok := categoryConvert[fmt.Sprint(value.CategoryId)]; ok {
			categoryData = categoryConvert[fmt.Sprint(value.CategoryId)][0].(model.Category)
		}
		brandsFind := ConvertBrandsToRsp(brandData)
		categoryFind := ConvertCategoryToRsp(categoryData)
		res := &proto.CategoryBrandResponse{
			Id:       value.ID,
			Brand:    &brandsFind,
			Category: &categoryFind,
		}
		resdata = append(resdata, res)
	}
	return resdata, nil
}

func ConvertBrandInfoResponse(GoodscategorybrandList []model.Goodscategorybrand, goodscategorybrandRows uint32) (result []*proto.BrandInfoResponse, err error) {
	BrandIds := make([]uint32, 0)
	if goodscategorybrandRows != 0 {
		for _, value := range GoodscategorybrandList {
			BrandIds = append(BrandIds, value.BrandId)
		}
	}
	brandsConvert, brandsErr := ConvertBrands(RemoveDuplicateElement(BrandIds), "id,name,logo", "Id")
	if brandsErr != nil {
		zap.S().Error(brandsErr.Error())
		return result, brandsErr
	}
	resdata := make([]*proto.BrandInfoResponse, 0)
	for _, value := range GoodscategorybrandList {
		brandData := model.Brands{}
		if _, ok := brandsConvert[fmt.Sprint(value.BrandId)]; ok {
			brandData = brandsConvert[fmt.Sprint(value.BrandId)][0].(model.Brands)
		}
		brandsFind := ConvertBrandsToRsp(brandData)
		resdata = append(resdata, &brandsFind)
	}
	return resdata, nil
}

func (b *BrandServer) BrandList(ctx context.Context, request *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var offset int32 = 0
	var limit int32 = 10

	if request.PagePerNums != 0 {
		limit = request.PagePerNums
	}
	if request.Pages != 0 {
		offset = limit * (request.Pages - 1)
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
	fmt.Println(brandList)
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
	GoodscategorybrandList, goodscategorybrandRows, goodscategorybrandErr := model.GetGoodscategorybrandList("", []interface{}{}, "", int(offset), int(limit))
	if goodscategorybrandErr != nil {
		zap.S().Error("服务器内部出错", goodscategorybrandErr.Error())
		return &proto.CategoryBrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	total, countErr := model.GetGoodscategorybrandCount("", []interface{}{})
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.CategoryBrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	resdata, resErr := ConvertCategoryBrandResponse(GoodscategorybrandList, goodscategorybrandRows)
	if resErr != nil {
		zap.S().Error("服务器内部出错", resErr.Error())
		return &proto.CategoryBrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错", resErr.Error())
	}

	return &proto.CategoryBrandListResponse{Data: resdata, Total: total}, nil
}

func (b *BrandServer) GetCategoryBrandList(ctx context.Context, rq *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.BrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if categoryRows == 0 {
		return &proto.BrandListResponse{}, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	goodscategorybrandList, goodscategorybrandRows, goodscategorybrandErr := model.GetGoodscategorybrandList("category_id=? and is_deleted=?", []interface{}{categoryFirst.ID, 0}, "id,brand_id,category_id", 0, 0)
	if goodscategorybrandErr != nil {
		zap.S().Error("服务器内部出错", goodscategorybrandErr.Error())
		return &proto.BrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, countErr := model.GetGoodscategorybrandCount("category_id=? and is_deleted=?", []interface{}{categoryFirst.ID, 0})
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.BrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	brandList, BrandListErr := ConvertBrandInfoResponse(goodscategorybrandList, goodscategorybrandRows)
	if BrandListErr != nil {
		zap.S().Error("服务器内部出错", BrandListErr.Error())
		return &proto.BrandListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.BrandListResponse{Total: total, Data: brandList}, nil
}

func (b *BrandServer) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id = ?", []interface{}{request.CategoryId}, "id,name,parent_category_id,level,is_tab")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.CategoryBrandResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if categoryRows == 0 {
		return &proto.CategoryBrandResponse{}, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	brandsFirst, brandsRows, brandsErr := model.GetBrandsFirst("id = ?", []interface{}{request.BrandId}, "id,name,logo")
	if brandsErr != nil {
		zap.S().Error("服务器内部出错", brandsErr.Error())
		return &proto.CategoryBrandResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if brandsRows == 0 {
		return &proto.CategoryBrandResponse{}, status.Errorf(codes.NotFound, "品牌不存在")
	}

	resCreate, er := model.CreateGoodscategorybrand(model.Goodscategorybrand{
		CategoryId: categoryFirst.ID,
		BrandId:    brandsFirst.ID,
		CreatedAt:  uint32(time.Now().Unix()),
		UpdatedAt:  uint32(time.Now().Unix()),
	})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.CategoryBrandResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	resBrand := ConvertBrandsToRsp(brandsFirst)
	resCategory := ConvertCategoryToRsp(categoryFirst)

	return &proto.CategoryBrandResponse{Id: resCreate.ID, Brand: &resBrand, Category: &resCategory}, nil
}

func (b *BrandServer) DeleteCategoryBrand(ctx context.Context, rq *proto.CategoryBrandRequest) (*proto.Empty, error) {
	goodscategorybrandFirst, goodscategorybrandRows, goodscategorybrandErr := model.GetGoodscategorybrandFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if goodscategorybrandErr != nil {
		zap.S().Error("服务器内部出错", goodscategorybrandErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if goodscategorybrandRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	er := model.UpdateGoodscategorybrand(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{goodscategorybrandFirst.ID})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (b *BrandServer) UpdateCategoryBrand(ctx context.Context, rq *proto.CategoryBrandRequest) (*proto.Empty, error) {

	_, firstRows, firstErr := model.GetGoodscategorybrandFirst("category_id = ? and brand_id= ?  and is_deleted=?", []interface{}{rq.CategoryId, rq.BrandId, 0}, "id")
	if firstErr != nil {
		zap.S().Error("服务器内部出错", firstErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if firstRows != 0 {
		return &proto.Empty{}, status.Errorf(codes.InvalidArgument, "改信息已经存在")
	}

	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id = ?", []interface{}{rq.CategoryId}, "id,name")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if categoryRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	brandsFirst, brandsRows, brandsErr := model.GetBrandsFirst("id = ?", []interface{}{rq.BrandId}, "id,name,logo")

	if brandsErr != nil {
		zap.S().Error("服务器内部出错", brandsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if brandsRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "品牌不存在")
	}

	goodscategorybrandFirst, goodscategorybrandRows, goodscategorybrandErr := model.GetGoodscategorybrandFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")
	if goodscategorybrandErr != nil {
		zap.S().Error("服务器内部出错", goodscategorybrandErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if goodscategorybrandRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	updateGoodsIsErr := model.UpdateGoodscategorybrand(model.Goodscategorybrand{
		CategoryId: categoryFirst.ID,
		BrandId:    brandsFirst.ID,
		UpdatedAt:  uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{goodscategorybrandFirst.ID})
	if updateGoodsIsErr != nil {
		zap.S().Error("服务器内部出错", updateGoodsIsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}
