package category

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"project/goods_web/api"
	"project/goods_web/forms"
	"project/goods_web/global"
	"project/goods_web/proto"
	"strconv"
)

// 子级菜单
type UmsMenuNode struct {
	Id           uint32        `json:"id,omitempty"`
	ParentId     uint32        `json:"parent_id"` // 父级ID
	Name         string        `json:"name"`      // 菜单名称
	Level        int32         `json:"level"`
	IsTab        bool          `json:"is_tab"`
	SubCategorys []UmsMenuNode `json:"subCategorys"`
}

//递归获取树形菜单
func ConvertCategoryMenu(categoryList []*proto.CategoryInfoResponse, parentId uint32) []UmsMenuNode {
	treeList := []UmsMenuNode{}
	for _, item := range categoryList {
		if item.ParentCategory == parentId {
			child := ConvertCategoryMenu(categoryList, item.Id)
			node := UmsMenuNode{
				Id:       item.Id,
				ParentId: item.ParentCategory,
				Name:     item.Name,
				Level:    item.Level,
				IsTab:    item.IsTab,
			}
			node.SubCategorys = child
			treeList = append(treeList, node)
		}
	}

	return treeList
}

func List(ctx *gin.Context) {
	r, err := global.CategorySrvClient.GetAllCategorysList(context.Background(), &proto.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	res := ConvertCategoryMenu(r.Data, 0)

	ctx.JSON(http.StatusOK, res)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	//1. 先查询出该分类写的所有子分类
	//2. 将所有的分类全部逻辑删除
	//3. 将该分类下的所有的商品逻辑删除
	_, err = global.CategorySrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	r, err := global.CategorySrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: uint32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	subCategorys := ConvertCategoryMenu(r.SubCategorys, r.Info.Id)
	reMap := make(map[string]interface{})
	reMap["id"] = r.Info.Id
	reMap["name"] = r.Info.Name
	reMap["level"] = r.Info.Level
	reMap["parent_category"] = r.Info.ParentCategory
	reMap["is_tab"] = r.Info.IsTab
	reMap["sub_categorys"] = subCategorys

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.CategorySrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	ctx.JSON(http.StatusOK, request)
}

func Update(ctx *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   uint32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.CategorySrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
