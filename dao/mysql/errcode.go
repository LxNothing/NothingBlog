package mysql

import "errors"

var (
	ErrInvalidCommunityId = errors.New("社区id无效")
	ErrBlogIdNotExisted   = errors.New("文章id不存在")
	ErrQueryFailed        = errors.New("数据库查询出错")
	ErrUserExisted        = errors.New("用户已经存在")
	ErrUserNotExisted     = errors.New("用户不存在")
)
