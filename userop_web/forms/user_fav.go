package forms

type UserFavForm struct {
	GoodsId uint64 `form:"goods" json:"goods" binding:"required"`
}
