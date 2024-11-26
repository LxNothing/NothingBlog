package models

const (
	BlogOrderByTime  = "time"
	BlogOrderByScore = "score"
)

// binding标签是validator库用的，用来验证对应的参数是否合法
// 具体见博客：https://www.liwenzhou.com/posts/Go/validator-usages/
// 这个结构体是用来接收前端传递的参数 - 注册参数
type SignUpParams struct {
	UserName   string            `json:"username" binding:"required"`                     // 用户名
	Password   string            `json:"password" binding:"required"`                     // 密码
	RePassword string            `json:"re_password" binding:"required,eqfield=Password"` // 确认密码
	Email      string            `json:"email"`                                           // 邮箱
	VerifyCode *VerifyCodeParams `json:"verify_code" binding:"required"`                  // 验证码
}

// 用户登录参数
type LoginParams struct {
	UserName   string            `json:"username" binding:"required"`
	Password   string            `json:"password" binding:"required"`
	VerifyCode *VerifyCodeParams `json:"verify_code" binding:"required"` // 验证码
}

// 用户对文章点赞投票的参数
type VoteDateParams struct {
	BlogId int64 `json:"blog_id,string" binding:"required"` // 点赞的文章id
	// -1 反对 0 取消 1 赞成，oneof表示这个字段只能取-1，0，1
	// 删除required 是因为传入0时会被validate这个库当前没有传值
	Direction int8 `json:"direction" binding:"oneof=-1 0 1"`
}

// 用于接收前端传递查询文章的参数- 包含排序字段
type BlogOrderListParams struct {
	Page  int64  `form:"page"`
	Szie  int64  `form:"size"`
	Order string `form:"order"`
}

// 验证码参数
type VerifyCodeParams struct {
	Id   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}
