package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"project/order_srv/global"
	"project/order_srv/model"
	"project/order_srv/proto"
	"project/order_srv/utils/register"
	"strings"
	"time"
)

type OrderServer struct {
}

func GenerateOrderDn(userid uint32) string {
	t := time.Now()
	return fmt.Sprintf("%d%d%d%d%d", t.Year(), t.Month(), t.Day(), userid, rand.Intn(100))
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
				Checked: *value.Checked,
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
			Checked: &request.Checked,
		})
	} else {
		sqlErr = model.UpdateShoppingcart(model.Shoppingcart{
			Nums:      shoppingcartFirst.Nums + request.Nums,
			UpdatedAt: uint32(time.Now().Unix()),
			Checked:   &request.Checked,
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
		Checked:   &request.Checked,
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
	shoppingcartList, shoppingcartRow, shoppingcartErr := model.GetShoppingcartList("user = ? and checked =? and is_deleted=?", []interface{}{request.UserId, 1, 0}, "", 0, 0)
	if shoppingcartErr != nil {
		zap.S().Error("服务器内部出错", shoppingcartErr.Error())
		return &proto.OrderInfoResponse{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if shoppingcartRow == 0 {
		return &proto.OrderInfoResponse{}, status.Errorf(codes.NotFound, "请提交商品")
	}
	var goodids []uint64
	var orderAmount float32
	var goodsNums = make(map[uint64]uint32)
	var orderGoodsList []model.Ordergoods
	var goodsSellInfo []*proto.GoodsInvInfo
	for _, val := range shoppingcartList {
		goodids = append(goodids, val.Goods)
		goodsNums[val.Goods] = val.Nums
	}
	var consulRegister register.Register = register.ConsulRegister{
		Host: global.ServerConfig.GoodsSrv.Consul.Host,
		Port: global.ServerConfig.GoodsSrv.Consul.Port,
	}
	goodsClientConn, goodsClientErr := consulRegister.GetServuce(global.ServerConfig.GoodsSrv.Name)
	if goodsClientErr != nil {
		zap.S().Error("商品服务不可用", goodsClientErr.Error())
		return &proto.OrderInfoResponse{}, status.Errorf(codes.Unavailable, "商品服务不可用")
	}
	goodsSrvClient := proto.NewGoodsClient(goodsClientConn)

	goodsList, goodsListErr := goodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodids,
	})
	if goodsListErr != nil {
		zap.S().Error("商品服务不可用", goodsListErr.Error())
		return &proto.OrderInfoResponse{}, status.Errorf(codes.Unavailable, "商品服务不可用")
	}
	for _, goodsInfo := range goodsList.Data {
		orderAmount += goodsInfo.ShopPrice * float32(goodsNums[goodsInfo.Id])
		orderGoods := model.Ordergoods{
			Goods:      goodsInfo.Id,
			GoodsName:  goodsInfo.Name,
			GoodsImage: goodsInfo.GoodsFrontImage,
			GoodsPrice: goodsInfo.ShopPrice,
			Nums:       goodsNums[goodsInfo.Id],
		}
		orderGoodsList = append(orderGoodsList, orderGoods)

		goodsSellInfoD := &proto.GoodsInvInfo{
			GoodsId: goodsInfo.Id,
			Num:     goodsNums[goodsInfo.Id],
		}
		goodsSellInfo = append(goodsSellInfo, goodsSellInfoD)
	}
	orderSn := GenerateOrderDn(request.UserId)
	//库存服务
	consulRegister = register.ConsulRegister{
		Host: global.ServerConfig.InvSrv.Consul.Host,
		Port: global.ServerConfig.InvSrv.Consul.Port,
	}
	invClientConn, invClientErr := consulRegister.GetServuce(global.ServerConfig.InvSrv.Name)
	if invClientErr != nil {
		zap.S().Error("库存服务不可用", invClientErr.Error())
		return &proto.OrderInfoResponse{}, status.Errorf(codes.Unavailable, "库存服务不可用")
	}
	invSrvClient := proto.NewInventoryClient(invClientConn)
	_, invSellErr := invSrvClient.Sell(context.Background(), &proto.SellInfo{GoodsInfo: goodsSellInfo, OrderSn: orderSn})
	if invSellErr != nil {
		e, _ := status.FromError(invSellErr)
		zap.S().Error("扣减库存失败", invSellErr.Error())
		if e.Code() == codes.ResourceExhausted {
			return &proto.OrderInfoResponse{}, status.Errorf(codes.ResourceExhausted, e.Message())
		} else {
			return &proto.OrderInfoResponse{}, status.Errorf(codes.Unavailable, "扣减库存失败")
		}

	}
	var timeData = uint32(time.Now().Unix())
	var orderinfoFind = model.Orderinfo{
		OrderSn:      orderSn,
		OrderMount:   orderAmount,
		Address:      request.Address,
		SignerName:   request.Name,
		SingerMobile: request.Mobile,
		Post:         request.Post,
		User:         request.UserId,
		CreatedAt:    timeData,
	}
	tx := global.MysqlDb.Begin()
	orderinfoFindErr := tx.Create(&orderinfoFind).Error
	if orderinfoFindErr != nil {
		tx.Rollback()
		zap.S().Error("服务器内部出错", orderinfoFindErr.Error())
		return &proto.OrderInfoResponse{}, status.Errorf(codes.Internal, "订单创建失败")
	}
	var orderGoodsSql string
	if len(orderGoodsList) != 0 {
		for _, val := range orderGoodsList {
			orderGoodsSql += fmt.Sprintf("('%d','%d','%s','%s','%v','%d','%d'),", orderinfoFind.ID, val.Goods, val.GoodsName, val.GoodsImage, val.GoodsPrice, val.Nums, timeData)
		}
		fields := "`order`,`goods`,`goods_name`,`goods_image`,`goods_price`,`nums`,`created_at`"
		if len(orderGoodsSql) != 0 {
			orderGoodsSql = strings.Trim(orderGoodsSql, ",")
			insertOk := model.BatchSave(tx, "ordergoods", fields, orderGoodsSql)
			//insertOk = false
			if !insertOk {
				tx.Rollback()
				return &proto.OrderInfoResponse{}, status.Errorf(codes.Internal, "订单创建失败")
			}
		}
	}

	shoppingcartE := tx.Model(&model.Shoppingcart{}).Where("user=? and checked=?", request.UserId, true).Updates(model.Shoppingcart{IsDeleted: true, DeletedAt: timeData}).Debug().Error
	if shoppingcartE != nil {
		fmt.Println(shoppingcartE.Error())
		tx.Rollback()
		return &proto.OrderInfoResponse{}, status.Errorf(codes.Internal, "订单创建失败")
	}
	tx.Commit()
	return &proto.OrderInfoResponse{Id: orderinfoFind.ID, OrderSn: orderSn, Total: orderAmount}, nil
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
	var where = make(map[string]interface{}, 0)
	where["is_deleted"] = 0
	where["id"] = request.Id
	if request.UserId != 0 {
		where["user"] = request.UserId
	}

	whereSql, vals, _ := WhereBuild(where)
	orderInfo, orderInfoRow, orderInfoErr := model.GetOrderinfoFirst(whereSql, vals, "id,user,order_sn,pay_type,status,post,order_mount,address,signer_name,singer_mobile")
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
