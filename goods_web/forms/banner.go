package forms

type BannerForm struct {
	Image      string `form:"image" json:"image" binding:"url"`
	Index      uint32    `form:"index" json:"index" binding:"required"`
	Url      string `form:"url" json:"url" binding:"url"`
}
