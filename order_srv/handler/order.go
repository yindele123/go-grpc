package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/order_srv/model"
	"project/order_srv/proto"
	"time"
)

type OrderServer struct {
}

func (o *OrderServer) CartItemList(ctx context.Context, info *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var where = make(map[string]interface{}, 0)
	if info.Id != 0 {
		where["user"] = info.Id
	}
	where["is_deleted"] = 0
	whereSql, vals, _ := WhereBuild(where)
	list, rows, err := model.GetShoppingcartList(whereSql, vals, "", 0, 0)
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.CartItemListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	total, countErr := model.GetShoppingcartCount(whereSql, vals)
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.CartItemListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	var result []*proto.ShopCartInfoResponse
	if rows != 0 {
		for _, value := range list {
			res := &proto.ShopCartInfoResponse{
				Id:      value.ID,
				UserId:  value.User,
				GoodsId: value.Goods,
				Nums:    value.Nums,
				Checked: value.Checked,
			}
			result = append(result, res)
		}
	}
	return &proto.CartItemListResponse{Data: result, Total: int32(total)}, nil
}

func (o *OrderServer) CreateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//todo
	/*conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【goods服务失败】")
	}*/

	shoppingcartFirst, shoppingcartRows, shoppingcartErr := model.GetShoppingcartFirst("goods = ? and user=? and is_deleted=?", []interface{}{request.GoodsId, request.UserId, 0}, "id,nums")
	if shoppingcartErr != nil {
		zap.S().Error("服务器内部出错", shoppingcartErr.Error())
		return &proto.ShopCartInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var sqlErr error
	if shoppingcartRows == 0 {
		shoppingcartFirst, sqlErr = model.CreateShoppingcart(model.Shoppingcart{
			Goods:   request.GoodsId,
			User:    request.UserId,
			Nums:    request.Nums,
			Checked: request.Checked,
		})
	} else {
		sqlErr = model.UpdateShoppingcart(model.Shoppingcart{
			Nums:      shoppingcartFirst.Nums + request.Nums,
			UpdatedAt: uint32(time.Now().Unix()),
			Checked:   request.Checked,
		}, "id=?", []interface{}{shoppingcartFirst.ID})
	}
	if sqlErr != nil {
		zap.S().Error("服务器内部出错", sqlErr.Error())
		return &proto.ShopCartInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.ShopCartInfoResponse{Id: shoppingcartFirst.ID}, nil
}

func (o *OrderServer) UpdateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.Empty, error) {
	//todo
	shoppingcartFirst, shoppingcartRows, shoppingcartErr := model.GetShoppingcartFirst("goods = ? and user=? and is_deleted=?", []interface{}{request.GoodsId, request.UserId, 0}, "id,nums")
	if shoppingcartErr != nil {
		zap.S().Error("服务器内部出错", shoppingcartErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if shoppingcartRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	var nums = shoppingcartFirst.Nums
	if request.Nums != 0 {
		nums = request.Nums
	}
	err := model.UpdateShoppingcart(model.Shoppingcart{
		Nums:      nums,
		Checked:   request.Checked,
		UpdatedAt: uint32(time.Now().Unix()),
	}, "id=?", []interface{}{shoppingcartFirst.ID})
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (o *OrderServer) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.Empty, error) {
	//todo
	shoppingcartFirst, shoppingcartRows, shoppingcartErr := model.GetShoppingcartFirst("goods = ? and user=? and is_deleted=?", []interface{}{request.GoodsId, request.UserId, 0}, "id,nums")
	if shoppingcartErr != nil {
		zap.S().Error("服务器内部出错", shoppingcartErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if shoppingcartRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	err := model.UpdateShoppingcart(model.Shoppingcart{
		IsDeleted: true,
		DeletedAt: uint32(time.Now().Unix()),
	}, "id=?", []interface{}{shoppingcartFirst.ID})
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}

func (o *OrderServer) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	panic("implement me")
}

func (o *OrderServer) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var where = make(map[string]interface{}, 0)
	where["is_deleted"] = 0
	if request.UserId != 0 {
		//todo
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
	orderinfoList, orderinfoRow, orderinfoErr := model.GetOrderinfoList(whereSql, vals, "id,user,order_sn,pay_type,status,trade_no,order_mount,pay_time,address,signer_name,singer_mobile,post,created_at", int(offset), int(limit))
	if orderinfoErr != nil {
		zap.S().Error("服务器内部出错", orderinfoErr.Error())
		return &proto.OrderListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	total, countErr := model.GetOrderinfoCount(whereSql, vals)
	if countErr != nil {
		zap.S().Error("服务器内部出错", countErr.Error())
		return &proto.OrderListResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var res []*proto.OrderInfoResponse
	if orderinfoRow != 0 {
		for _, val := range orderinfoList {
			data := &proto.OrderInfoResponse{
				Id:      val.ID,
				UserId:  val.User,
				OrderSn: val.OrderSn,
				PayType: val.PayType,
				Status:  val.Status,
				Post:    val.Post,
				Total:   val.OrderMount,
				Address: val.Address,
				Name:    val.SignerName,
				Mobile:  val.SingerMobile,
				AddTime: val.CreatedAt,
			}
			res = append(res, data)
		}
	}
	return &proto.OrderListResponse{Data: res, Total: int32(total)}, nil
}

func (o *OrderServer) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	orderInfo, orderInfoRow, orderInfoErr := model.GetOrderinfoFirst("is_deleted=? and id=?", []interface{}{0, request.Id}, "id,user,order_sn,pay_type,status,post,order_mount,address,signer_name,singer_mobile")
	if orderInfoErr != nil {
		zap.S().Error("服务器内部出错", orderInfoErr.Error())
		return &proto.OrderInfoDetailResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if orderInfoRow == 0 {
		return &proto.OrderInfoDetailResponse{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	var orderInfoResponse = proto.OrderInfoResponse{
		Id:      orderInfo.ID,
		UserId:  orderInfo.User,
		OrderSn: orderInfo.OrderSn,
		PayType: orderInfo.PayType,
		Status:  orderInfo.Status,
		Post:    orderInfo.Post,
		Total:   orderInfo.OrderMount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
	}
	ordergoodsList, ordergoodsRow, ordergoodsErr := model.GetOrdergoodsList("`order`=? and is_deleted=?", []interface{}{orderInfo.ID, 0}, "id,goods,goods_name,goods_image,goods_price,nums", 0, 0)
	if ordergoodsErr != nil {
		zap.S().Error("服务器内部出错", ordergoodsErr.Error())
		return &proto.OrderInfoDetailResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var resData []*proto.OrderItemResponse
	if ordergoodsRow != 0 {
		for _, val := range ordergoodsList {
			data := &proto.OrderItemResponse{
				GoodsId:    val.Goods,
				GoodsName:  val.GoodsName,
				GoodsImage: val.GoodsImage,
				GoodsPrice: val.GoodsPrice,
				Nums:       val.Nums,
			}
			resData = append(resData, data)
		}
	}
	return &proto.OrderInfoDetailResponse{OrderInfo: &orderInfoResponse, Data: resData}, nil
}

func (o *OrderServer) UpdateOrderStatus(ctx context.Context, request *proto.OrderStatus) (*proto.Empty, error) {
	orderinfo, orderinfoRow, orderinfoErr := model.GetOrderinfoFirst("order_sn=? and is_deleted=?", []interface{}{request.OrderSn, 0}, "id")
	if orderinfoErr != nil {
		zap.S().Error("服务器内部出错", orderinfoErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if orderinfoRow == 0 {
		return &proto.Empty{}, status.Errorf(codes.NotFound, "记录不存在")
	}
	sqlErr := model.UpdateOrderinfo(model.Orderinfo{
		Status: request.Status,
	}, "id=?", []interface{}{orderinfo.ID})
	if sqlErr != nil {
		zap.S().Error("服务器内部出错", sqlErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}
