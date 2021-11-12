package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/userop_srv/global"
	"project/userop_srv/model"
	"project/userop_srv/proto"
	"time"
)

type UserFavServer struct {
}

func (u *UserFavServer) GetFavList(ctx context.Context, request *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var userFav = make([]model.Userfav,0)
	var where = make(map[string]interface{}, 0)
	if request.UserId != 0 {
		where["user"] = request.UserId
	}
	if request.GoodsId != 0 {
		where["goods"] = request.GoodsId
	}
	whereSql, vals, _ := WhereBuild(where)
	result := global.MysqlDb.Where(whereSql, vals...).Select("id,user,goods,created_at").Find(&userFav)
	if result.Error != nil {
		zap.S().Error("服务器内部出错", result.Error.Error())
		return &proto.UserFavListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var total int64
	err := global.MysqlDb.Model(&model.Userfav{}).Where(whereSql, vals...).Count(&total).Error
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.UserFavListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var resData = make([]*proto.UserFavResponse, 0)
	if result.RowsAffected != 0 {
		for _, val := range userFav {
			res := &proto.UserFavResponse{
				UserId:  val.User,
				GoodsId: val.Goods,
			}
			resData = append(resData, res)
		}
	}
	return &proto.UserFavListResponse{Total: int32(total), Data: resData}, nil
}

func (u *UserFavServer) AddUserFav(ctx context.Context, request *proto.UserFavRequest) (*proto.Empty, error) {
	//todo
	//这里可以访问用户跟商品服务来确定信息准确度(我这就不处理了，交给web端处理)
	var userFav = model.Userfav{
		User:      request.UserId,
		Goods:     request.GoodsId,
		CreatedAt: uint32(time.Now().Unix()),
	}
	result := global.MysqlDb.Where("user=? and goods=?", request.UserId, request.GoodsId).Limit(1).Find(&model.Userfav{})
	if result.Error != nil {
		zap.S().Error("服务器内部出错", result.Error.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if result.RowsAffected != 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录已存在")
	}
	err := global.MysqlDb.Create(&userFav).Error
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (u *UserFavServer) DeleteUserFav(ctx context.Context, request *proto.UserFavRequest) (*proto.Empty, error) {
	var userFav = model.Userfav{}
	result := global.MysqlDb.Where("user=? and goods=?", request.UserId, request.GoodsId).Limit(1).Find(&userFav)
	if result.Error != nil {
		zap.S().Error("服务器内部出错", result.Error.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if result.RowsAffected == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	err := global.MysqlDb.Where("id=?", userFav.ID).Delete(&model.Userfav{}).Error
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (u *UserFavServer) GetUserFavDetail(ctx context.Context, request *proto.UserFavRequest) (*proto.Empty, error) {
	var userFav = model.Userfav{}
	result := global.MysqlDb.Where("user=? and goods=?", request.UserId, request.GoodsId).Limit(1).Find(&userFav)
	if result.Error != nil {
		zap.S().Error("服务器内部出错", result.Error.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if result.RowsAffected == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &proto.Empty{}, nil
}
