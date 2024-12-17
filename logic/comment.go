package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/models"
	"NothingBlog/package/utils"
	"NothingBlog/settings"
	"errors"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type LogicComment struct {
}

var (
	ErrCreateComment       = errors.New("添加新评论失败")
	ErrDelentComment       = errors.New("删除评论失败")
	ErrUpdateCommentStatus = errors.New("更新评论数据失败")
	ErrCommentParamInvalid = errors.New("非法数据")
	ErrQueryComment        = errors.New("查询评论数据失败")
	ErrCommentNotFound     = errors.New("评论记录未找到")
)

var commentDb = mysql.DaoComment{}

// 用于后端管理获取评论列表 - 值返回当前评论，不用返回对应的子评论
func (LogicComment) GetWithPageForAdmin(info *models.CommentWithPageParams) ([]models.CommentWithName, int, int, error) {
	// 查数据库
	// 检查page和size的有效性
	if info.Size > uint(settings.Confg.PageSize) {
		return nil, 0, 0, ErrCommentParamInvalid
	}

	res, total, err := commentDb.GetWithPageForAdmin(info)
	if err != nil {
		zap.L().Debug(ErrQueryComment.Error(), zap.Error(err))
		return nil, 0, 0, ErrQueryComment
	}

	var cur_page = 1
	var total_page = 1
	if info.PageIdx > 0 && info.Size > 0 && total != 0 {
		cur_page = int(info.PageIdx)
		total_page = utils.GetTotalPage(int(info.Size), int(total))
	}

	return res, cur_page, total_page, nil
}

// 用于客户端返回所有评论 - 包括当前评论的子评论
func (LogicComment) GetWithPageForClient(atcId string) ([]models.ResponseCommentListForClient, error) {
	// 查数据库
	res, err := commentDb.GetWithPageForClient(atcId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Debug(ErrCommentNotFound.Error(), zap.Error(err))
			return nil, ErrCommentNotFound
		}
		zap.L().Debug(ErrQueryComment.Error(), zap.Error(err))
		return nil, ErrQueryComment
	}

	return res, nil
}

// CreateComment 创建新评论
func (lc *LogicComment) CreateComment(comment *models.Comment) error {
	// 调用数据库的方法 写入一条新的评论
	if err := commentDb.InsertIntoDb(comment); err != nil {
		zap.L().Debug(ErrCreateComment.Error(), zap.Error(err))
		return ErrCreateComment
	}

	// 读取邮箱设置,如果配置了邮箱服务 则进行下面的邮件通知工作
	// 1. 如果本条评论没有父评论，就向注册用户（博主）发送通知邮件
	// 2. 如果本条评论有父评论，就向父评论的邮箱发送邮件通知
	// to do

	return nil
}

func (lc *LogicComment) DeleteCommentById(id string) error {
	// 这里有个问题：如果删除的评论包含子评论，该如何处理 - to do
	if err := commentDb.DeleteById(id); err != nil {
		zap.L().Debug(ErrDelentComment.Error(), zap.Error(err))
		return ErrDelentComment
	}
	return nil
}

func (lc *LogicComment) DeleteCommentByIds(ids []uint) error {
	// 这里有个问题：如果删除的评论包含子评论，该如何处理 - to do
	if err := commentDb.DeleteByIds(ids); err != nil {
		zap.L().Debug(ErrDelentComment.Error(), zap.Error(err))
		return ErrDelentComment
	}
	return nil
}

func (lc *LogicComment) UpdateCommentStatus(id uint, state models.CommentStatusType) error {
	// 验证字段值是否合法
	if state != models.StatusCheckFailed && state != models.StatusCheckSuccess && state != models.StatusChecking {
		return ErrCommentParamInvalid
	}

	// 更新数据库
	if err := commentDb.UpdateStatus(id, uint8(state)); err != nil {
		zap.L().Debug(ErrDelentComment.Error(), zap.Error(err))
		return ErrUpdateCommentStatus
	}
	return nil
}
