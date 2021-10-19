package model

import (
	"encoding/json"
	"fmt"
	"project/goods_srv/global"
	"reflect"
)

// SearchModel 搜索接口
type SearchModel interface {
	TableName() string
}

// SearchModelHandler 存储一些查询过程中的必要信息
type SearchModelHandler struct {
	Model SearchModel
}

// GetSearchModelHandler 获取处理器
func GetSearchModelHandler(model SearchModel) *SearchModelHandler {
	return &SearchModelHandler{
		Model: model,
	}
}

// Search 查找
func (s *SearchModelHandler) Search() string {
	query := global.MysqlDb.Model(s.Model)
	itemPtrType := reflect.TypeOf(s.Model)
	if itemPtrType.Kind() != reflect.Ptr {
		itemPtrType = reflect.PtrTo(itemPtrType)
	}
	itemSlice := reflect.SliceOf(itemPtrType)
	res := reflect.New(itemSlice)
	err := query.Debug().Find(res.Interface()).Error
	if err != nil {
		// 这里不要学我
		panic("error")
	}

	ret, _ := json.Marshal(res)
	fmt.Println(string(ret))
	return string(ret)
}