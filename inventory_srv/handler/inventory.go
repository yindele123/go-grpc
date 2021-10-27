package handler

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"project/inventory_srv/global"
	"project/inventory_srv/model"
	"project/inventory_srv/proto"
	"sync"
	"time"
)

type InventoryServer struct {
}

func (i InventoryServer) SetInv(ctx context.Context, info *proto.GoodsInvInfo) (*proto.Empty, error) {
	inventoryFirst, inventoryRows, inventoryErr := model.GetInventoryFirst("goods=?", []interface{}{info.GoodsId}, "id")

	if inventoryErr != nil {
		zap.S().Error("服务器内部出错", inventoryErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	var resErr error
	if inventoryRows == 0 {
		_, resErr = model.CreateInventory(model.Inventory{
			Goods:     info.GoodsId,
			Stocks:    info.Num,
			CreatedAt: uint32(time.Now().Unix()),
			UpdatedAt: uint32(time.Now().Unix()),
		})
	} else {
		resErr = model.UpdateInventory(model.Inventory{
			Stocks:    info.Num,
			UpdatedAt: uint32(time.Now().Unix()),
		}, "id=?", []interface{}{inventoryFirst.ID})
	}
	if resErr != nil {
		zap.S().Error("服务器内部出错", resErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}

	return &proto.Empty{}, nil
}

func (i InventoryServer) InvDetail(ctx context.Context, info *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	inventoryFirst, inventoryRows, inventoryErr := model.GetInventoryFirst("goods=?", []interface{}{info.GoodsId}, "id,stocks,goods")
	if inventoryErr != nil {
		zap.S().Error("服务器内部出错", inventoryErr.Error())
		return &proto.GoodsInvInfo{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if inventoryRows == 0 {
		return &proto.GoodsInvInfo{}, status.Errorf(codes.Internal, "记录不存在")
	}
	return &proto.GoodsInvInfo{GoodsId: inventoryFirst.Goods, Num: inventoryFirst.Stocks}, nil
}

func (i InventoryServer) Sell(ctx context.Context, info *proto.SellInfo) (*proto.Empty, error) {
	var wg sync.WaitGroup
	done := make(chan bool, 0)
	locker := redislock.New(global.Rdb)

	ctxLock := context.Background()
	var lock *redislock.Lock
	var err error

	for {
		lock, err = locker.Obtain(ctxLock, "inventory_sell", 30*time.Second, nil)
		if redislock.ErrNotObtained != err {
			break
		}
	}
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	wg.Add(1)
	go func(conduit chan bool) {
		for {
			select {
			case <-conduit:
				wg.Done()
				return
			default:
				ttl, _ := lock.TTL(ctxLock)
				if ttl < 10*time.Second {
					if err := lock.Refresh(ctxLock, 30*time.Second, nil); err != nil {
						zap.S().Error("redis服务内部出错", err.Error())
						return
					}
				}
			}
		}
	}(done)
	defer wg.Wait()
	defer lock.Release(ctxLock)
	defer func() {
		done <- true
		close(done)
	}()
	var goodsIdS []uint64
	for _, value := range info.GoodsInfo {
		goodsIdS = append(goodsIdS, value.GoodsId)
	}
	inventoryList, inventoryRows, inventoryErr := model.GetInventoryList("goods in ?", []interface{}{goodsIdS}, "id,goods,stocks", 0, 0)
	if inventoryErr != nil {
		zap.S().Error("服务器内部出错", inventoryErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if inventoryRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.InvalidArgument, "参数错误(商品ID记录不存在)")
	}
	inventoryListMap := StructSliceToMap(inventoryList, "goods")
	var lastKey uint64
	upateDataAll := make(map[uint64]map[string]interface{}, 0)
	for _, value := range info.GoodsInfo {
		if _, ok := inventoryListMap[fmt.Sprint(value.GoodsId)]; !ok {
			return &proto.Empty{}, status.Errorf(codes.InvalidArgument, fmt.Sprintf("参数错误(商品ID%d记录不存在)", value.GoodsId))
			break
		}
		inventoryFind := inventoryListMap[fmt.Sprint(value.GoodsId)][0].(model.Inventory)
		if inventoryFind.Stocks < value.Num {
			return &proto.Empty{}, status.Errorf(codes.ResourceExhausted, fmt.Sprintf("商品ID%d库存不足", value.GoodsId))
			break
		}
		upateDataAll[inventoryFind.ID] = map[string]interface{}{"stocks": inventoryFind.Stocks - value.Num}
		lastKey = inventoryFind.ID
	}
	batchUpdateErr := model.BatchUpdateData("inventory", upateDataAll, lastKey)
	if batchUpdateErr != nil {
		zap.S().Error("服务器内部出错", batchUpdateErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}

func (i InventoryServer) Reback(ctx context.Context, info *proto.SellInfo) (*proto.Empty, error) {
	var wg sync.WaitGroup
	done := make(chan bool, 0)
	locker := redislock.New(global.Rdb)

	ctxLock := context.Background()
	var lock *redislock.Lock
	var err error

	for {
		lock, err = locker.Obtain(ctxLock, "inventory_sell", 30*time.Second, nil)
		if redislock.ErrNotObtained != err {
			break
		}
	}
	if err != nil {
		zap.S().Error("服务器内部出错", err.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	wg.Add(1)
	go func(conduit chan bool) {
		for {
			select {
			case <-conduit:
				wg.Done()
				return
			default:
				ttl, _ := lock.TTL(ctxLock)
				if ttl < 10*time.Second {
					if err := lock.Refresh(ctxLock, 30*time.Second, nil); err != nil {
						zap.S().Error("redis服务内部出错", err.Error())
						return
					}
				}
			}
		}
	}(done)
	defer wg.Wait()
	defer lock.Release(ctxLock)
	defer func() {
		done <- true
		close(done)
	}()
	var goodsIdS []uint64
	for _, value := range info.GoodsInfo {
		goodsIdS = append(goodsIdS, value.GoodsId)
	}
	inventoryList, inventoryRows, inventoryErr := model.GetInventoryList("goods in ?", []interface{}{goodsIdS}, "id,goods,stocks", 0, 0)
	if inventoryErr != nil {
		zap.S().Error("服务器内部出错", inventoryErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	if inventoryRows == 0 {
		return &proto.Empty{}, status.Errorf(codes.InvalidArgument, "参数错误(商品ID记录不存在)")
	}
	inventoryListMap := StructSliceToMap(inventoryList, "goods")
	var lastKey uint64
	upateDataAll := make(map[uint64]map[string]interface{}, 0)
	for _, value := range info.GoodsInfo {
		if _, ok := inventoryListMap[fmt.Sprint(value.GoodsId)]; !ok {
			return &proto.Empty{}, status.Errorf(codes.InvalidArgument, fmt.Sprintf("参数错误(商品ID%d记录不存在)", value.GoodsId))
			break
		}
		inventoryFind := inventoryListMap[fmt.Sprint(value.GoodsId)][0].(model.Inventory)
		upateDataAll[inventoryFind.ID] = map[string]interface{}{"stocks": inventoryFind.Stocks + value.Num}
		lastKey = inventoryFind.ID
	}
	batchUpdateErr := model.BatchUpdateData("inventory", upateDataAll, lastKey)
	if batchUpdateErr != nil {
		zap.S().Error("服务器内部出错", batchUpdateErr.Error())
		return &proto.Empty{}, status.Errorf(codes.Internal, "服务器内部出错")
	}
	return &proto.Empty{}, nil
}
