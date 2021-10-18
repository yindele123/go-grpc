package handler

import (
	"fmt"
	"project/goods_srv/model"
	"project/goods_srv/proto"
	"reflect"
	"strings"
)


type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)



func WhereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}



func ConvertGoodsToRsp(goods model.Goods,brands model.Brands,category model.Category) (goodsInfoRsp proto.GoodsInfoResponse) {
	brandData := ConvertBrandsToRsp(brands)
	categoryData := ConvertCategoryToRsp(category)
	goodsInfoRsp.Id=goods.ID
	goodsInfoRsp.CategoryId=goods.CategoryId
	goodsInfoRsp.Name=goods.Name
	goodsInfoRsp.GoodsSn=goods.GoodsSn
	goodsInfoRsp.ClickNum=goods.ClickNum
	goodsInfoRsp.SoldNum=goods.SoldNum
	goodsInfoRsp.FavNum=goods.FavNum
	goodsInfoRsp.MarketPrice=goods.MarketPrice
	goodsInfoRsp.ShopPrice=goods.ShopPrice
	goodsInfoRsp.GoodsBrief=goods.GoodsBrief
	goodsInfoRsp.ShipFree=goods.ShipFree
	goodsInfoRsp.Images=goods.Images
	goodsInfoRsp.DescImages=goods.DescImages
	goodsInfoRsp.GoodsFrontImage=goods.GoodsFrontImage
	goodsInfoRsp.IsNew=goods.IsNew
	goodsInfoRsp.IsHot=goods.IsHot
	goodsInfoRsp.OnSale=goods.OnSale
	goodsInfoRsp.Brand=&brandData
	goodsInfoRsp.Category=&categoryData
	return goodsInfoRsp
}

func ConvertBrandsToRsp(brands model.Brands)   (brandInfoRsp proto.BrandInfoResponse) {
	brandInfoRsp.Id=brands.ID
	brandInfoRsp.Name=brands.Name
	brandInfoRsp.Logo=brands.Logo
	return  brandInfoRsp
}


func ConvertCategoryToRsp(category model.Category)   (categoryInfoRsp proto.CategoryBriefInfoResponse) {
	categoryInfoRsp.Id=category.ID
	categoryInfoRsp.Name=category.Name
	return  categoryInfoRsp
}

func StructSliceToMap(source interface{},filedName string) map[string][]interface{}{
	filedIndex := 0
	v := reflect.ValueOf(source) // 判断，interface转为[]interface{}
	resMap := make(map[string][]interface{})
	if v.Kind() != reflect.Slice {
		return resMap
	}
	l := v.Len()
	retList := make([]interface{}, l)
	for i := 0; i < l; i++ {
		retList[i] = v.Index(i).Interface()
	}
	if len(retList) > 0 {
		firstObj := retList[0]
		objT := reflect.TypeOf(firstObj)
		for i := 0; i < objT.NumField(); i++ {
			if objT.Field(i).Name == filedName {
				filedIndex = i
			}
		}
	}
	for _, elem := range retList {
		key := reflect.ValueOf(elem).Field(filedIndex).Interface()
		value :=  make([]interface{}, 0)
		resMap[fmt.Sprint(key)] = value
	}

	for _, elem := range retList {
		key := reflect.ValueOf(elem).Field(filedIndex).Interface()
		resMap[fmt.Sprint(key)] = append(resMap[fmt.Sprint(key)], elem)
	}
	return resMap
}




func RemoveDuplicateElement(originals interface{}) (interface{}) {
	temp := map[string]struct{}{}
	switch slice := originals.(type) {
	case []string:
		result := make([]string, 0, len(originals.([]string)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []int64:
		result := make([]int64, 0, len(originals.([]int64)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	case []uint32:
		result := make([]uint32, 0, len(originals.([]uint32)))
		for _, item := range slice {
			key := fmt.Sprint(item)
			if _, ok := temp[key]; !ok {
				temp[key] = struct{}{}
				result = append(result, item)
			}
		}
		return result
	default:
		return originals
	}
}