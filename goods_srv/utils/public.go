package utils

import (
	"project/goods_srv/model"
)

// 子级菜单
type UmsMenuNode struct {
	Id         uint32     `json:"id,omitempty"`
	ParentId   uint32     `json:"parent_id"`   // 父级ID
	Name       string    `json:"name"`        // 菜单名称
	Children []UmsMenuNode `json:"children"`
}


//递归获取树形菜单
func GetMenu(parentId uint32)([]UmsMenuNode,error){
	//获取parentId的所有子菜单
	categoryList, rows, _:=model.GetCategoryList("parent_category_id=?",[]interface{}{parentId},"id,parent_category_id,name",0,0)
	tree := make([]UmsMenuNode,0)
	if rows!=0 {
		for _,item := range categoryList {
			child ,_:= GetMenu(item.ID) //获取parentId每一个子菜单的子菜单
			node := UmsMenuNode{
				Id: item.ID,
				ParentId: item.ParentCategoryId,
				Name: item.Name,
			}
			node.Children = child
			tree = append(tree,node)
		}
	}

	return tree,nil
}


func GetMenuIds(id uint32)([]uint32){
	//获取parentId的所有子菜单
	categoryList, rows, _:=model.GetCategoryList("parent_category_id=?",[]interface{}{id},"id,parent_category_id",0,0)
	var treeIds []uint32
	if rows!=0 {
		for _,item := range categoryList {
			if item.ParentCategoryId==id {
				treeIds = append(treeIds,item.ID)
				ids := GetMenuIds(item.ID)
				treeIds = append(treeIds,ids...)
			}
		}
	}

	return treeIds
}