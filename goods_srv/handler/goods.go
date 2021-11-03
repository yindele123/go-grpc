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

type GoodsServer struct {
}

func ConvertGoodsCategory(goodsList []model.Goods, rows int64) (result []*proto.GoodsInfoResponse, err error) {
	categoryIds := make([]uint32, 0)
	BrandIds := make([]uint32, 0)

	if rows != 0 {
		for _, value := range goodsList {
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

	for _, value := range goodsList {
		brandData := model.Brands{}
		categoryData := model.Category{}
		if _, ok := brandsConvert[fmt.Sprint(value.BrandId)]; ok {
			brandData = brandsConvert[fmt.Sprint(value.BrandId)][0].(model.Brands)
		}
		if _, ok := categoryConvert[fmt.Sprint(value.CategoryId)]; ok {
			categoryData = categoryConvert[fmt.Sprint(value.CategoryId)][0].(model.Category)
		}
		res := ConvertGoodsToRsp(value, brandData, categoryData)
		result = append(result, &res)
	}
	return result, nil
}

func (g *GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	var where = make(map[string]interface{}, 0)

	if request.IsHot {
		where["is_hot"] = 1
	}
	if request.IsNew {
		where["is_new"] = 1
	}
	if len(request.KeyWords) != 0 {
		where["name like"] = "%" + request.KeyWords + "%"
	}

	if request.PriceMin != 0 {
		where["shop_price >="] = request.PriceMin
	}
	if request.PriceMax != 0 {
		where["shop_price <="] = request.PriceMin
	}

	if request.Brand != 0 {
		where["brand_id"] = request.Brand

	}
	if request.TopCategory != 0 {
		ids := utils.GetMenuIds(request.TopCategory)
		ids = append(ids, request.TopCategory)
		where["category_id in"] = ids

	}

	whereSql, vals, _ := WhereBuild(where)
	var offset int32 = 0
	var limit int32 = 10
	if request.PagePerNums != 0 {
		limit = request.PagePerNums
	}
	if request.Pages != 0 {
		offset = limit * (request.Pages - 1)
	}
	goodsList, rows, goodsErr := model.GetGoodsList(whereSql, vals, "", int(offset), int(limit))
	if goodsErr != nil {
		zap.S().Error("服务器内部出错", goodsErr.Error())
		return &proto.GoodsListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, countErr := model.GetGoodsCount(whereSql, vals)
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.GoodsListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	result, resErr := ConvertGoodsCategory(goodsList, rows)
	if resErr != nil {
		return &proto.GoodsListResponse{}, status.Errorf(codes.Internal, "服务器内部出错", resErr.Error())
	}
	return &proto.GoodsListResponse{Total: total, Data: result}, nil
}

func (g *GoodsServer) BatchGetGoods(ctx context.Context, request *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	if len(request.Id) == 0 {
		return &proto.GoodsListResponse{}, nil
	}

	goodsList, rows, goodsErr := model.GetGoodsList("id in ?", []interface{}{request.Id}, "", 0, 0)
	if goodsErr != nil {
		zap.S().Error("服务器内部出错", goodsErr.Error())
		return &proto.GoodsListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, countErr := model.GetGoodsCount("id in ?", []interface{}{request.Id})
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.GoodsListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	result, resErr := ConvertGoodsCategory(goodsList, rows)
	if resErr != nil {
		return &proto.GoodsListResponse{}, status.Errorf(codes.Internal, "服务器内部出错", resErr.Error())
	}
	return &proto.GoodsListResponse{Total: total, Data: result}, nil
}

func (g *GoodsServer) CreateGoods(ctx context.Context, request *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {

	categoryFirst, categoryRows, categoryErr := model.GetCategoryFirst("id = ?", []interface{}{request.CategoryId}, "id,name")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if categoryRows == 0 {
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	brandsFirst, brandsRows, brandsErr := model.GetBrandsFirst("id = ?", []interface{}{request.BrandId}, "id,name,logo")
	if brandsErr != nil {
		zap.S().Error("服务器内部出错", brandsErr.Error())
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if brandsRows == 0 {
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.NotFound, "品牌不存在")
	}
	images, _ := json.Marshal(request.Images)
	descImages, _ := json.Marshal(request.DescImages)
	resCreate, er := model.CreateGoods(model.Goods{
		CategoryId:      categoryFirst.ID,
		BrandId:         brandsFirst.ID,
		GoodsSn:         request.GoodsSn,
		Name:            request.Name,
		MarketPrice:     request.MarketPrice,
		ShopPrice:       request.ShopPrice,
		GoodsBrief:      request.GoodsBrief,
		ShipFree:        request.ShipFree,
		Images:          string(images),
		DescImages:      string(descImages),
		GoodsFrontImage: request.GoodsFrontImage,
		OnSale:          request.OnSale,
		IsNew:           request.IsNew,
		IsHot:           request.IsHot,
		CreatedAt:       uint32(time.Now().Unix()),
		UpdatedAt:       uint32(time.Now().Unix()),
	})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	res := ConvertGoodsToRsp(resCreate, brandsFirst, categoryFirst)
	return &res, nil
}

func (g *GoodsServer) DeleteGoods(ctx context.Context, rq *proto.DeleteGoodsInfo) (*proto.Empty, error) {
	goodsFirst, goodsRows, goodsErr := model.GetGoodsFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if goodsErr != nil {
		zap.S().Error("服务器内部出错", goodsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if goodsRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "商品不存在")
	}
	er := model.UpdateGoods(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{goodsFirst.ID})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (g *GoodsServer) UpdateGoods(ctx context.Context, rq *proto.CreateGoodsInfo) (*proto.Empty, error) {
	if rq.CategoryId != 0 {
		_, categoryRows, categoryErr := model.GetCategoryFirst("id = ?", []interface{}{rq.CategoryId}, "id,name")

		if categoryErr != nil {
			zap.S().Error("服务器内部出错", categoryErr.Error())
			return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
		}

		if categoryRows == 0 {
			return &proto.Empty{}, status.Errorf(codes.NotFound, "商品分类不存在")
		}
	}

	if rq.BrandId != 0 {
		_, brandsRows, brandsErr := model.GetBrandsFirst("id = ?", []interface{}{rq.BrandId}, "id,name,logo")

		if brandsErr != nil {
			zap.S().Error("服务器内部出错", brandsErr.Error())
			return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
		}

		if brandsRows == 0 {
			return &proto.Empty{}, status.Errorf(codes.NotFound, "品牌不存在")
		}
	}

	goodsFirst, goodsRows, goodsErr := model.GetGoodsFirst("id = ?", []interface{}{rq.Id}, "id")
	if goodsErr != nil {
		zap.S().Error("服务器内部出错", goodsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if goodsRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "商品不存在")
	}
	images, _ := json.Marshal(rq.Images)
	descImages, _ := json.Marshal(rq.DescImages)
	updateGoodsIsErr := model.UpdateGoods(model.Goods{
		CategoryId:      rq.CategoryId,
		BrandId:         rq.BrandId,
		GoodsSn:         rq.GoodsSn,
		Name:            rq.Name,
		MarketPrice:     rq.MarketPrice,
		ShopPrice:       rq.ShopPrice,
		GoodsBrief:      rq.GoodsBrief,
		ShipFree:        rq.ShipFree,
		Images:          string(images),
		DescImages:      string(descImages),
		GoodsFrontImage: rq.GoodsFrontImage,
		OnSale:          rq.OnSale,
		IsNew:           rq.IsNew,
		IsHot:           rq.IsHot,
		UpdatedAt:       uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{goodsFirst.ID})
	if updateGoodsIsErr != nil {
		zap.S().Error("服务器内部出错", updateGoodsIsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}

func (g *GoodsServer) GetGoodsDetail(ctx context.Context, rq *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	goodsFirst, goodsRows, goodsErr := model.GetGoodsFirst("id = ?", []interface{}{rq.Id}, "")
	if goodsErr != nil {
		zap.S().Error("服务器内部出错", goodsErr.Error())
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if goodsRows == 0 {
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.NotFound, "商品不存在")
	}

	var brands model.Brands
	var category model.Category

	category, _, categoryErr := model.GetCategoryFirst("id = ?", []interface{}{goodsFirst.CategoryId}, "id,name")

	if categoryErr != nil {
		zap.S().Error("服务器内部出错", categoryErr.Error())
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	brands, _, brandsErr := model.GetBrandsFirst("id = ?", []interface{}{goodsFirst.BrandId}, "id,name,logo")

	if brandsErr != nil {
		zap.S().Error("服务器内部出错", brandsErr.Error())
		return &proto.GoodsInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	res := ConvertGoodsToRsp(goodsFirst, brands, category)
	return &res, nil
}
