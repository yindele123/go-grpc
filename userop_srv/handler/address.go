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

type AddressServer struct {
}

func (a *AddressServer) GetAddressList(ctx context.Context, request *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var where = make(map[string]interface{}, 0)
	where["is_deleted"] = 0
	if request.UserId != 0 {
		where["user"] = request.UserId
	}
	whereSql, vals, _ := WhereBuild(where)
	addressList, addressRow, err := model.GetAddressList(whereSql, vals, "id,user,province,city,district,address,signer_name,signer_mobile,created_at", 0, 0)
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.AddressListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, err := model.GetAddressCount(whereSql, vals)
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.AddressListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var resData []*proto.AddressResponse
	if addressRow != 0 {
		for _, val := range addressList {
			res := &proto.AddressResponse{
				Id:     val.ID,
				UserId: val.User,

				Province:     val.Province,
				City:         val.City,
				District:     val.District,
				Address:      val.Address,
				SignerName:   val.SignerName,
				SignerMobile: val.SignerMobile,
				AddTime:      val.CreatedAt,
			}
			resData = append(resData, res)
		}
	}
	return &proto.AddressListResponse{Total: int32(total), Data: resData}, nil
}

func (a *AddressServer) CreateAddress(ctx context.Context, request *proto.AddressRequest) (*proto.AddressResponse, error) {
	addressFind, err := model.CreateAddress(model.Address{
		User:         request.UserId,
		Province:     request.Province,
		City:         request.City,
		District:     request.District,
		Address:      request.Address,
		SignerName:   request.SignerName,
		SignerMobile: request.SignerMobile,
		CreatedAt:    uint32(time.Now().Unix()),
	})
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.AddressResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.AddressResponse{Id: addressFind.ID}, nil
}

func (a *AddressServer) DeleteAddress(ctx context.Context, request *proto.AddressRequest) (*proto.Empty, error) {
	addressFirst, addressRow, err := model.GetAddressFirst("id=? and user=?", []interface{}{request.Id, request.UserId}, "id")
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if addressRow == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	err = global.MysqlDb.Where("id=?", addressFirst.ID).Delete(&model.Address{}).Error
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil

}

func (a *AddressServer) UpdateAddress(ctx context.Context, request *proto.AddressRequest) (*proto.Empty, error) {
	addressFirst, addressRow, err := model.GetAddressFirst("id=? and user=?", []interface{}{request.Id, request.UserId}, "id")
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if addressRow == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	var updateData = make(map[string]interface{}, 0)
	if len(request.Province) != 0 {
		updateData["province"] = request.Province
	}
	if len(request.City) != 0 {
		updateData["city"] = request.City
	}
	if len(request.District) != 0 {
		updateData["district"] = request.District
	}
	if len(request.Address) != 0 {
		updateData["address"] = request.Address
	}
	if len(request.SignerName) != 0 {
		updateData["signer_name"] = request.SignerName
	}
	if len(request.SignerMobile) != 0 {
		updateData["signer_mobile"] = request.SignerMobile
	}
	updateData["created_at"] = uint32(time.Now().Unix())

	err = global.MysqlDb.Model(&model.Address{}).Where("id=?", addressFirst.ID).Updates(updateData).Error
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}
