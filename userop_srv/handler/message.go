package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/userop_srv/model"
	"project/userop_srv/proto"
	"time"
)

type MessageServer struct {
}

func (m *MessageServer) MessageList(ctx context.Context, request *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var where = make(map[string]interface{}, 0)
	where["is_deleted"] = 0
	if request.UserId != 0 {
		where["user"] = request.UserId
	}
	var offset int32 = 0
	var limit int32 = 10
	if request.PagePerNums != 0 {
		limit = request.PagePerNums
	}
	if request.Pages != 0 {
		offset = limit * (request.Pages - 1)
	}
	whereSql, vals, _ := WhereBuild(where)
	messagesList, messagesRow, err := model.GetMessagesList(whereSql, vals, "id,user,message_type,subject,message,file,created_at", int(offset), int(limit))
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.MessageListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, err := model.GetMessagesCount(whereSql, vals)
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.MessageListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var resData []*proto.MessageResponse
	if messagesRow != 0 {
		for _, val := range messagesList {
			res := &proto.MessageResponse{
				Id:          val.ID,
				UserId:      val.User,
				MessageType: val.MessageType,
				Subject:     val.Subject,
				Message:     val.Message,
				File:        val.File,
				AddTime:     val.CreatedAt,
			}
			resData = append(resData, res)
		}
	}
	return &proto.MessageListResponse{Total: int32(total), Data: resData}, nil
}

func (m *MessageServer) CreateMessage(ctx context.Context, request *proto.MessageRequest) (*proto.MessageResponse, error) {
	messagesFind, err := model.CreateMessages(model.Leavingmessages{
		User:        request.UserId,
		MessageType: request.MessageType,
		Subject:     request.Subject,
		Message:     request.Message,
		File:        request.File,
		CreatedAt:   uint32(time.Now().Unix()),
	})
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.MessageResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.MessageResponse{
		Id:          messagesFind.ID,
		MessageType: messagesFind.MessageType,
		Subject:     messagesFind.Subject,
		Message:     messagesFind.Message,
		File:        messagesFind.File,
		AddTime:     messagesFind.CreatedAt,
	}, nil
}
