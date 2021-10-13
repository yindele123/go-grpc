package handler

import (
	"context"
	"project/goods_srv/proto"
)

type GoodsServer struct {
}

func (g *GoodsServer) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	return &proto.GoodsListResponse{},nil
}
func (g *GoodsServer) BatchGetGoods(ctx context.Context, request *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error){
	return &proto.GoodsListResponse{},nil
}
func (g *GoodsServer) CreateGoods(ctx context.Context, request *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error){
	return &proto.GoodsInfoResponse{},nil
}
func (g *GoodsServer) DeleteGoods(ctx context.Context, request *proto.DeleteGoodsInfo) (*proto.Empty, error){
	return &proto.Empty{},nil
}
func (g *GoodsServer) UpdateGoods(ctx context.Context, request *proto.CreateGoodsInfo) (*proto.Empty, error){
	return &proto.Empty{},nil
}
func (g *GoodsServer) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error){
	return &proto.GoodsInfoResponse{},nil
}