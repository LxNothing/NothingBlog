package logic

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/dao/redis"
	"NothingBlog/models"
	"NothingBlog/package/snowflake"
	"errors"
	"time"

	"go.uber.org/zap"
)

var (
	ErrBlogIdNotExisted = errors.New("文章id不存在")
	ErrInvalidCommId    = errors.New("社区ID无效")
	ErrInvalidUserId    = errors.New("用户ID无效")
	//ErrNoDataSelected   = errors.New("没有文章需要查询")
)

func CreateNewBlog(b *models.BlogsArch) error {
	// 生成文章id - 雪花算法
	b.Id = snowflake.GetNextId().Int64()

	// 获取创建时间 和 更新时间
	b.CreateTime = time.Now()
	b.UpdateTime = b.CreateTime
	b.VoteScore = b.CreateTime.Unix() // 默认初始分数是unix时间 - 对于点赞直接从0开始

	// 保存到mysql数据库
	if err := mysql.InsertNewBlog(b); err != nil {
		return err
	}

	// 创建对应的redis zset，为了投票点赞功能实现
	return redis.CreateNewBlogZset(b.Id, b.VoteScore)
}

func GetBlogDetailById(id int64) (allinfo *models.ApiBlogDetail, err error) {
	allinfo = new(models.ApiBlogDetail)

	// 查询文章信息
	allinfo.Blog, err = mysql.QueryBlogDetailById(id)
	if errors.Is(err, mysql.ErrBlogIdNotExisted) {
		err = ErrBlogIdNotExisted
		return
	}
	// 查询所属的社区详情
	allinfo.Community, err = mysql.QueryCommunityById(allinfo.Blog.CommmunityId)

	if errors.Is(err, mysql.ErrInvalidCommunityId) {
		err = ErrInvalidCommId
		return
	}
	// 查询对应的用户姓名
	allinfo.AuthorName, err = mysql.QueryUsernameById(allinfo.Blog.AuthorId)
	// if errors.Is(err, mysql.ErrorUserNotExist) {
	// 	err = ErrInvalidUserId
	// 	return
	// }
	return
}

func GetBlogList(page int64, size int64) (blogs []*models.ApiBlogDetail, err error) {
	var blogList []*models.BlogsArch
	blogList, err = mysql.QueryBlogList(page, size)

	blogs = make([]*models.ApiBlogDetail, 0, len(blogList))

	for _, blog := range blogList {
		// 查询所属的社区详情
		community, suberr := mysql.QueryCommunityById(blog.CommmunityId)

		if errors.Is(suberr, mysql.ErrInvalidCommunityId) {
			return nil, ErrInvalidCommId
		}
		// 查询对应的用户姓名
		author, suberr := mysql.QueryUsernameById(blog.AuthorId)
		// if errors.Is(suberr, mysql.ErrorUserNotExist) {
		// 	return nil, ErrInvalidUserId
		// }

		// 组合信息
		blogdetail := &models.ApiBlogDetail{
			AuthorName: author,
			Community:  community,
			Blog:       blog,
		}
		blogs = append(blogs, blogdetail)
	}

	return
}

func GetBlogOrderList(page int64, size int64, od string) (blogs []*models.ApiBlogDetail, err error) {
	// 根据排序方式从redis中获取文章ID列表
	var ids []string
	ids, err = redis.GetBlogIdWithOrder(page, size, od)
	if len(ids) == 0 || err != nil {
		zap.L().Error("redis.GetBlogIdWithOrder 出错", zap.Error(err))
		return nil, nil
	}

	// 根据文章ID从mysql中查找数据
	var blogList []*models.BlogsArch
	blogList, err = mysql.GetBlogWithOrder(ids)
	if err != nil {
		zap.L().Error("mysql.GetBlogWithOrder 出错", zap.Error(err))
		return nil, nil
	}

	// 合并社区信息，创建用户的信息
	blogs = make([]*models.ApiBlogDetail, 0, len(blogList))

	for _, blog := range blogList {
		// 查询所属的社区详情
		community, suberr := mysql.QueryCommunityById(blog.CommmunityId)

		if errors.Is(suberr, mysql.ErrInvalidCommunityId) {
			return nil, ErrInvalidCommId
		}
		// 查询对应的用户姓名
		author, suberr := mysql.QueryUsernameById(blog.AuthorId)
		// if errors.Is(suberr, mysql.ErrorUserNotExist) {
		// 	return nil, ErrInvalidUserId
		// }

		// 组合信息
		blogdetail := &models.ApiBlogDetail{
			AuthorName: author,
			Community:  community,
			Blog:       blog,
		}
		blogs = append(blogs, blogdetail)
	}

	return
}
