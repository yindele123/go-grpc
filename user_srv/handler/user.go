package handler

import (
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/user_srv/model"
	"project/user_srv/proto"
)

type UserServer struct {
}

func ConvertUserToRsp(user model.User) (user_info_rsp proto.UserInfoResponse) {
	user_info_rsp.Id = user.ID
	user_info_rsp.PassWord = user.Password
	user_info_rsp.Mobile = user.Mobile
	user_info_rsp.Role = user.Role
	user_info_rsp.NickName = user.NickName
	user_info_rsp.Gender = user.Gender
	user_info_rsp.BirthDay = user.Birthday
	return
}

func (u *UserServer) GetUserList(ctx context.Context, request *proto.PageInfo) (*proto.UserListResonse, error) {
	var offset uint32 = 0
	var limit uint32 = 10
	if request.PSize != 0 {
		limit = request.PSize
	}
	if request.Pn != 0 {
		offset = limit * (request.Pn - 1)
	}
	//where map[string]interface{}{
	//		"mobile": "18602058150 or 1=1",
	//	}
	userList, rows, err := model.GetUserList("", "", int(offset), int(limit))
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.UserListResonse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, er := model.GetUserCount("")
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.UserListResonse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	result := make([]*proto.UserInfoResponse, 0)
	if rows != 0 {
		for _, value := range userList {
			res := ConvertUserToRsp(value)
			result = append(result, &res)
		}
	}
	return &proto.UserListResonse{Total: uint32(total), Data: result}, nil
}

func (u *UserServer) GetUserByMobile(ctx context.Context, request *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var res proto.UserInfoResponse
	mobile := request.Mobile
	if len(mobile) == 0 {
		return &proto.UserInfoResponse{}, status.Errorf(codes.InvalidArgument, "Mobile信息无效")
	}
	UserFirst, rows, err := model.GetUserFirst(map[string]interface{}{
		"mobile": mobile,
	}, "")

	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if rows != 0 {
		res = ConvertUserToRsp(UserFirst)
	} else {
		return &res, status.Errorf(codes.NotFound, "用户不存在")
	}
	return &res, nil
}

func (u *UserServer) GetUserById(ctx context.Context, request *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var res proto.UserInfoResponse
	id := request.Id
	if id == 0 {
		return &proto.UserInfoResponse{}, status.Errorf(codes.InvalidArgument, "Id信息无效")
	}
	UserFirst, rows, err := model.GetUserFirst(map[string]interface{}{
		"id": id,
	}, "")

	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if rows != 0 {
		res = ConvertUserToRsp(UserFirst)
	} else {
		return &res, status.Errorf(codes.NotFound, "用户不存在")
	}
	return &res, nil
}

func (u *UserServer) CreateUser(ctx context.Context, request *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	_, rows, err := model.GetUserFirst(map[string]interface{}{"mobile": request.Mobile}, "")
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if rows != 0 {
		return &proto.UserInfoResponse{}, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	if len(request.Mobile) == 0 || len(request.PassWord) == 0 {
		return &proto.UserInfoResponse{}, status.Errorf(codes.InvalidArgument, "Mobile或PassWord信息无效")
	}

	passwordbyte, err := bcrypt.GenerateFromPassword([]byte(request.PassWord), bcrypt.DefaultCost)
	if err != nil {
		zap.S().Error("加密出错了", err.Error())
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	UserFirst, er := model.CreateUser(model.User{
		NickName: request.NickName,
		Mobile:   request.Mobile,
		Password: string(passwordbyte),
	})
	res := ConvertUserToRsp(UserFirst)
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.UserInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &res, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, request *proto.UpdateUserInfo) (*proto.Empty, error) {
	UserFirst, rows, err := model.GetUserFirst(map[string]interface{}{
		"id": request.Id,
	}, "")
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if rows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "用户不存在")
	}
	er := model.UpdateUser(map[string]interface{}{
		"nick_name": request.NickName,
		"gender":    request.Gender,
		"birthday":  request.BirthDay,
	}, map[string]interface{}{
		"id": UserFirst.ID,
	})
	if er != nil {
		zap.S().Error("服务器内部出错", er.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (u *UserServer) CheckPassWord(ctx context.Context, request *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(request.EncryptedPassword), []byte(request.Password)); err != nil {
		return &proto.CheckResponse{Success: false}, status.Errorf(codes.InvalidArgument, "密码比对错误")
	}
	return &proto.CheckResponse{Success: true}, nil

}
