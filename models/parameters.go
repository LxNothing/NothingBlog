package models

const (
	BlogOrderByTime  = "time"
	BlogOrderByScore = "score"
)

// binding标签是validator库用的，用来验证对应的参数是否合法
// 具体见博客：https://www.liwenzhou.com/posts/Go/validator-usages/
// 这个结构体是用来接收前端传递的参数 - 注册参数
type SignUpParams struct {
	UserName   string            `json:"username" binding:"required"`                     // 用户名 - 用户名作为唯一标识，不允许重复
	Password   string            `json:"password" binding:"required"`                     // 密码
	RePassword string            `json:"re_password" binding:"required,eqfield=Password"` // 确认密码
	Email      string            `json:"email"`                                           // 邮箱 - 仅在重置密码时，接收验证码使用
	VerifyCode *VerifyCodeParams `json:"verify_code" binding:"required"`                  // 验证码
}

// 用户登录参数
type LoginParams struct {
	UserName   string            `json:"username" binding:"required"`
	Password   string            `json:"password" binding:"required"`
	VerifyCode *VerifyCodeParams `json:"verify_code" binding:"required"` // 验证码
}

// 验证码参数
type VerifyCodeParams struct {
	Id   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// 重置密码的参数
type ResetPasswordParams struct {
	UserName string `json:"username" binding:"required"` // 用户名
	Email    string `json:"email" binding:"email"`       // 邮箱
}

// 修改密码的参数
type ModifyPasswordParams struct {
	UserName    string            `json:"username" binding:"required"`     // 用户名
	OldPassword string            `json:"old_password" binding:"required"` // 旧密码
	NewPassword string            `json:"new_password" binding:"required"` // 新密码
	VerifyCode  *VerifyCodeParams `json:"verify_code" binding:"required"`  // 验证码
}

/*
==========================================

	与文章相关的请求参数结构定义
	 数据格式：
	 {
	    "class_id":"123",
	    "top_flag":true,
	    "en_comment":false,
	    "status":0,
	    "privilege":1,
	    "title":"第一篇文章-golang",
	    "summary":"hello world",
	    "image":"no image",
	    "Content":"golang is a program language, yes",
	    "tag_id_list":[
	        {"id":"123456"},
	        {"id":"3456"},
	        {"id":"2222"}]
	    },
	}

==========================================
*/
// admin管理端 - 分页获取文章列表
type ArticleWithPageParams struct {
	Keyword   string `form:"n"`  // 关键字
	Privilege string `form:"pr"` // 权限
	Class     string `form:"c"`  // 文章类别
	Status    string `form:"st"` // 文章的状态
	Page      uint   `form:"p"`  // 第几页
	Size      uint   `form:"s"`  // 每页的大小
	Tag       string `form:"t"`  // tag
}

// 新建文章的前端数据
type NewArticleFormsParams struct {
	TagIdList []TagFormsParams `json:"tag_id_list"`                        // 文章标签-标签可以为空
	ClassId   int64            `json:"class_id,string" binding:"required"` // 文章所属的分类 - 比如教程，分享等
	TopFlag   bool             `json:"top_flag"`                           // 是否置顶标志 false-没有置顶 true-置顶
	EnComment bool             `json:"en_comment"`                         // 是否允许评论 false-不允许 true-允许
	Status    uint8            `json:"status"`                             // 文章状态 0-草稿 1-发布 2-删除
	Privilege uint8            `json:"privilege"`                          // 文章权限 0-公开 1-私有
	Title     string           `json:"title" binding:"required"`           // 文章标题
	Summary   string           `json:"summary"`                            // 文章的摘要信息
	Image     string           `json:"image"`                              // 文章缩略图
	Content   string           `json:"content" binding:"required"`         // 文章内容
}

// 更新文章的前端数据
type UpdateArticleFormsParams struct {
	ArticleId int64 `json:"article_id,string" binding:"required"`
	NewArticleFormsParams
}

type DeleteMultiArticleParams struct {
	Ids []int64 `json:"ids" binding:"required"` // 待删除的文章ID列表
}

// 软删除的文章参数
type SoftDeleteArticleParams struct {
	Ids     []int64 `json:"ids" binding:"required"`      // 待删除的文章ID列表
	DelFalg bool    `json:"del_flag" binding:"required"` // true - 删除  false - 恢复
}

/*
==========================================

	与标签相关的请求参数结构定义

==========================================
*/
type TagFormsParams struct {
	Id int64 `json:"id,string" binding:"required"` // tag id
}

// 创建tag时的参数
type TagCreateFormParams struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

type DeleteMultiTagParams struct {
	Ids []int64 `json:"ids" binding:"required"` // 待删除的TagID列表
}

// 更新Tag的参数
type UpdateTagParams struct {
	TagId int64  `json:"tag_id,string" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Desc  string `json:"desc"`
}

/*
==========================================

	与分类相关的请求参数结构定义

==========================================
*/
type ClassFormsParams struct {
	Id int64 `json:"id,string" binding:"required"` // tag id
}

type ClassCreateFormParams struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc"`
}

// 更新Tag的参数
type UpdateClassParams struct {
	ClassId int64  `json:"class_id,string" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Desc    string `json:"desc"`
}

/*
==========================================

	与评论相关的请求参数结构定义

==========================================
*/
// 删除多个评论的参数定义
type CommentDeleteFormParams struct {
	Ids []uint `json:"ids"` // id列表
}

// 更新评论的状态的参数
type CommentUpdateStateParams struct {
	Id    uint  `json:"id"`
	Value uint8 `json:"value"`
}

// 评论投票的参数
type CommentVoteParams struct {
}

// 分页获取评论的参数 - 后台管理端
type CommentWithPageParams struct {
	Keyword     string `form:"keyword"`       // 关键字
	Type        uint8  `form:"type"`          // 类型
	AtcId       int64  `form:"atc_id,string"` // 文章ID
	IsParentCmt uint8  `form:"is_parent_cmt"` // 是否为父评论
	Status      uint8  `form:"status"`        // 评论的状态
	PageIdx     uint   `form:"page_idx"`      // 第几页
	Size        uint   `form:"size"`          // 每页的大小
}

// 创建评论的表单数据
type CommentCreateFormParams struct {
	ArticleId int64 `json:"article_id,string"` // 被评论的文章ID
	//RootCommentId   uint   `json:"root_comment_id,string"`   // 根评论的ID
	ParentCommentId uint   `json:"parent_comment_id,string"` // 父评论的ID
	UserName        string `json:"user_name" `               // 评论人名字 - 强制要求
	Email           string `json:"email"`                    // 评论人的邮箱 - 不强制要求
	Content         string `json:"content"`                  // 评论内容
	Status          uint8  `json:"status"`                   // 评论状态：0-审核中 1-审核通过 2-审核未通过
	Type            uint8  `json:"type"`                     // 评论种类 1- 表示文章评论
}

// ParamToDbModel 将参数格式转换为数据库模型的模式
func (c *CommentCreateFormParams) ParamToDbModel() *Comment {
	return &Comment{
		ArticleId: c.ArticleId,
		//RootCommentId:   c.RootCommentId,
		ParentCommentId: c.ParentCommentId,
		UserName:        c.UserName,
		Email:           c.Email,
		Content:         c.Content,
		Type:            c.Type,
		Status:          CommentStatusType(c.Status),
	}
}
