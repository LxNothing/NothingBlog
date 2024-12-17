package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type CommentStatusType uint8

const (
	StatusChecking     CommentStatusType = 1 + iota // 评论审核中
	StatusCheckSuccess                              // 评论审核通过，可以对外展示
	StatusCheckFailed                               // 评论审核未通过，不对外展示
)

// 评论的数据库存储模型
type Comment struct {
	gorm.Model
	ArticleId int64 `json:"article_id"` // 被评论的文章ID
	//RootCommentId   uint              `json:"root_comment_id"`                                  // 根评论的ID
	ParentCommentId uint              `json:"parent_comment_id" gorm:"default:null"`            // 父评论的ID
	ChildComment    []Comment         `json:"child_comment" gorm:"-"`                           // 子评论 gorm:"-"表示gorm忽略这个字段
	UserName        string            `json:"user_name" gorm:"type:varchar(64);not null;"`      // 评论人名字 - 强制要求
	Email           string            `json:"email" gorm:"type:varchar(64)"`                    // 评论人的邮箱 - 不强制要求
	Content         string            `json:"content" gorm:"type:Text;not null"`                // 评论内容
	Icon            string            `json:"icon" gorm:"type:varchar(256)"`                    // 评论人的图标。类似头像
	Status          CommentStatusType `json:"status" gorm:"type:tinyint(1);unsigned;default:1"` // 评论状态：1-审核中 2-审核通过 3-审核未通过
	Type            uint8             `json:"type" gorm:"type:tinyint(1);unsigned;default:1"`   // 评论种类 1- 表示文章评论
	Agree           uint16            `json:"agree" gorm:"unsigned; default:0"`                 // 对该条评论的点赞数
	Disagree        uint16            `json:"disagree" gorm:"unsigned; default:0"`              // 对该条评论的不赞同数量
	// 外键约束
	//RootComment   *Comment `json:"root_comment" gorm:"ForeignKey:RootCommentId"`      // 根评论
	ParentComment *Comment `json:"parient_comment" gorm:"ForeignKey:ParentCommentId"` // 父评论
	Article       *Article `json:"article" gorm:"references:ArticleId"`               // 文章
}

// 包含对应文章名称、父评论名称的评论结构
type CommentWithName struct {
	Comment
	ArticleName       string `json:"article_name"`
	ParentCommentName string `json:"parent_comment_name"`
}

// 后台管理端的评论列表
type ResponseCommentListForAdmin struct {
	Id                uint              `json:"id"`                  // 评论的ID
	CreatedAt         time.Time         `json:"create_at"`           // 评论的创建时间
	UserName          string            `json:"user_name"`           // 评论人名字
	Email             string            `json:"email"`               // 评论人的邮箱
	ArticleName       string            `json:"article_name"`        // 被评论的主体名称，比如是文章，或者某个页面，目前支持文章
	ParentCommentName string            `json:"parent_comment_name"` // 该条评论的父评论用户名
	Content           string            `json:"content"`             // 评论内容
	Agree             uint16            `json:"agree"`               // 对该条评论的点赞数
	Disagree          uint16            `json:"disagree"`            // 对该条评论的不赞同数量
	Type              uint8             `json:"type"`                // 评论种类 1- 表示文章评论
	Status            CommentStatusType `json:"status"`              // 评论状态：1-审核中 2-审核通过 3-审核未通过
}

// 客户端评论列表
type ResponseCommentListForClient struct {
	Id              uint                           `json:"id"`                // 评论的ID
	CreatedAt       time.Time                      `json:"create_at"`         // 评论的创建时间
	UserName        string                         `json:"user_name"`         // 评论人名字
	Icon            string                         `json:"icon"`              // 头像
	ParentCommentId uint                           `json:"parent_comment_id"` // 该条评论的父评论ID
	Content         string                         `json:"content"`           // 评论内容
	Agree           uint16                         `json:"agree"`             // 对该条评论的点赞数
	Disagree        uint16                         `json:"disagree"`          // 对该条评论的不赞同数量
	SubComment      []ResponseCommentListForClient `json:"sub_comment"`       // 子评论
}

// 将Comment 类型转为 ResponseCommentListForAdmin这种类型
func (c *CommentWithName) BindToResponseForAdmin() *ResponseCommentListForAdmin {
	return &ResponseCommentListForAdmin{
		Id:                c.ID,
		CreatedAt:         c.CreatedAt,
		UserName:          c.UserName,
		Email:             c.Email,
		Content:           c.Content,
		Agree:             c.Agree,
		Disagree:          c.Disagree,
		Type:              c.Type,
		Status:            c.Status,
		ArticleName:       c.ArticleName,
		ParentCommentName: c.ParentCommentName,
	}
}

func (c *Comment) BindToResponseForClient() *ResponseCommentListForClient {
	return &ResponseCommentListForClient{
		Id:              c.ID,
		CreatedAt:       c.CreatedAt,
		UserName:        c.UserName,
		ParentCommentId: c.ParentCommentId,
		Content:         c.Content,
		Agree:           c.Agree,
		Disagree:        c.Disagree,
	}
}
