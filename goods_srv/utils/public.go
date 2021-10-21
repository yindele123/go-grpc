package utils

import (
	"project/goods_srv/model"
)

// 子级菜单
type UmsMenuNode struct {
	Id           uint32        `json:"id,omitempty"`
	ParentId     uint32        `json:"parent_id"` // 父级ID
	Name         string        `json:"name"`      // 菜单名称
	Level        uint8         `json:"level"`
	IsTab        bool          `json:"is_tab"`
	SubCategorys []UmsMenuNode `json:"subCategorys"`
}

//递归获取树形菜单
/*func GetMenu(parentId uint32) []*proto.CategoryInfoResponse {
	//获取parentId的所有子菜单
	categoryList, rows, _ := model.GetCategoryList("parent_category_id=?", []interface{}{parentId}, "id,parent_category_id,name,level,is_tab", 0, 0)
	tree := make([]*proto.CategoryInfoResponse, 0)
	if rows != 0 {
		for _, item := range categoryList {
			child := GetMenu(item.ID) //获取parentId每一个子菜单的子菜单
			node := &proto.CategoryInfoResponse{
				Id:       item.ID,
				ParentCategory: item.ParentCategoryId,
				Name:     item.Name,
				Level:    item.Level,
				IsTab:    item.IsTab,
			}
			node.SubCategorys = child
			tree = append(tree, node)
		}
	}

	return tree
}*/

func GetMenuIds(id uint32) []uint32 {
	//获取parentId的所有子菜单
	categoryList, rows, _ := model.GetCategoryList("parent_category_id=?", []interface{}{id}, "id,parent_category_id", 0, 0)
	var treeIds []uint32
	if rows != 0 {
		for _, item := range categoryList {
			if item.ParentCategoryId == id {
				treeIds = append(treeIds, item.ID)
				ids := GetMenuIds(item.ID)
				treeIds = append(treeIds, ids...)
			}
		}
	}

	return treeIds
}
