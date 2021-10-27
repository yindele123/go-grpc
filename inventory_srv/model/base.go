package model

import (
	"bytes"
	"fmt"
	"project/inventory_srv/global"
	"project/inventory_srv/utils"
	"strconv"
)

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
	sql=utils.TrimLastChar(buff.String())
	where := fmt.Sprintf(` WHERE id IN (%v)`, ids)
	if err := global.MysqlDb.Debug().Exec(fmt.Sprintf("%s %s",sql,where)).Error; err != nil {
		return err
	}
	return nil

}
