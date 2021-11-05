package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"project/inventory_srv/utils"
	"project/order_srv/global"
	"reflect"
	"strconv"
)

//下面是反射查询(我现在不使用这块)

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

// 获取新的struct切片，返回值 *[]*struct{}
func (s *SearchModelHandler) GetNewModelSlice() interface{} {
	t := reflect.TypeOf(s.Model)
	// return reflect.Indirect(reflect.New(reflect.SliceOf(t))).Addr().Interface()
	list := reflect.New(reflect.SliceOf(t)).Elem()
	list.Set(reflect.MakeSlice(list.Type(), 0, 0))
	return reflect.Indirect(list).Addr().Interface()
}

// Search 查找
func (s *SearchModelHandler) GetList(whereSql string, vals []interface{}, fields string, Offset int, limit int, order string) (resStr string, rows int64, err error) {
	query := global.MysqlDb.Model(s.Model).Limit(limit).Offset(Offset)
	itemPtrType := reflect.TypeOf(s.Model)
	if itemPtrType.Kind() != reflect.Ptr {
		itemPtrType = reflect.PtrTo(itemPtrType)
	}
	itemSlice := reflect.SliceOf(itemPtrType)
	res := reflect.New(itemSlice)

	if len(fields) != 0 {
		query.Select(fields)
	}
	if len(order) != 0 {
		query.Order(order)
	}
	if len(whereSql) != 0 && len(vals) != 0 {
		query.Where(whereSql, vals...)
	}
	result := query.Find(res.Interface())
	ret, _ := json.Marshal(res.Interface())
	return string(ret), result.RowsAffected, result.Error
}

func BatchUpdateData(table string, data map[uint64]map[string]interface{}, lastKey uint64) (err error) {
	var ids string
	buff := bytes.Buffer{}
	sql := fmt.Sprintf(`UPDATE %s%s SET `, global.ServerConfig.MysqlInfo.TablePrefix, table)
	buff.WriteString(sql)
	for k, _ := range data {
		if ids == "" {
			ids = strconv.Itoa(int(k))
		} else {
			ids = ids + "," + strconv.Itoa(int(k))
		}
	}
	for k, _ := range data[lastKey] {
		buff.WriteString(fmt.Sprintf(` %s= CASE id  `, k))
		for key, val := range data {
			buff.WriteString(fmt.Sprintf(`WHEN %d THEN '%s' `, key, fmt.Sprint(val[k])))
		}
		buff.WriteString("END,")
	}
	fmt.Println(ids)
	sql = utils.TrimLastChar(buff.String())
	where := fmt.Sprintf(` WHERE id IN (%v)`, ids)
	if err := global.MysqlDb.Debug().Exec(fmt.Sprintf("%s %s", sql, where)).Error; err != nil {
		return err
	}
	return nil

}

func BatchSave(db *gorm.DB, table, fields, value string) bool {
	buff := bytes.Buffer{}
	sql := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s", global.ServerConfig.MysqlInfo.TablePrefix+table, fields, value)
	buff.WriteString(sql)
	if err := db.Debug().Exec(fmt.Sprintf("%s", sql)).Error; err != nil {
		return false
	}
	return true
}
