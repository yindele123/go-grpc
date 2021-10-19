package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/goods_srv/model"
	"project/goods_srv/proto"
	"time"
)

type BannerServer struct {
}

func (b *BannerServer) BannerList(ctx context.Context, empty *proto.Empty) (*proto.BannerListResponse, error) {
	bannerList, bannerRow, bannerErr := model.GetBannersList("", []interface{}{}, "id,image,url,`index`", 0, 0, "`index` desc")
	if bannerErr != nil {
		zap.S().Error("服务器内部出错", bannerErr.Error())
		return &proto.BannerListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, countErr := model.GetBannersCount("", []interface{}{})
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.BannerListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	result := make([]*proto.BannerResponse, 0)
	if bannerRow != 0 {
		for _, value := range bannerList {
			res := &proto.BannerResponse{
				Id:    value.ID,
				Index: value.Index,
				Image: value.Image,
				Url:   value.Url,
			}
			result = append(result, res)
		}
	}
	return &proto.BannerListResponse{Total: total, Data: result}, nil
}

func (b *BannerServer) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	resCreate, er := model.CreateBanners(model.Banners{
		Image:     request.Image,
		Index:     request.Index,
		Url:       request.Url,
		CreatedAt: uint32(time.Now().Unix()),
		UpdatedAt: uint32(time.Now().Unix()),
	})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.BannerResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.BannerResponse{
		Id:    resCreate.ID,
		Image: resCreate.Image,
		Url:   resCreate.Url,
		Index: resCreate.Index,
	}, nil
}

func (b *BannerServer) DeleteBanner(ctx context.Context, rq *proto.BannerRequest) (*proto.Empty, error) {
	bannerFirst, bannerRows, bannerErr := model.GetBannersFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if bannerErr != nil {
		zap.S().Error("服务器内部出错", bannerErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if bannerRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	er := model.UpdateBanners(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{bannerFirst.ID})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (b *BannerServer) UpdateBanner(ctx context.Context, rq *proto.BannerRequest) (*proto.Empty, error) {
	bannerFirst, bannerRows, bannerErr := model.GetBannersFirst("id = ? and is_deleted=?", []interface{}{rq.Id, 0}, "id")

	if bannerErr != nil {
		zap.S().Error("服务器内部出错", bannerErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	if bannerRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}

	updateGoodsIsErr := model.UpdateBanners(model.Banners{
		Image:     rq.Image,
		Index:     rq.Index,
		Url:       rq.Url,
		UpdatedAt: uint32(time.Now().Unix()),
	}, "id = ?", []interface{}{bannerFirst.ID})
	if updateGoodsIsErr != nil {
		zap.S().Error("服务器内部出错", updateGoodsIsErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}
