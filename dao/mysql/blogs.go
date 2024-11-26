package mysql

import (
	"NothingBlog/models"
)

func InsertNewBlog(b *models.BlogsArch) error {
	// sqlStr := "insert into blog(blog_id, title, content, author_id, community_id) values(?,?,?,?,?)"
	// _, err := db.Exec(sqlStr, b.Id, b.Title, b.Content, b.AuthorId, b.CommmunityId)
	// return err
	return nil
}

// 根据文章id查询文章详情
func QueryBlogDetailById(id int64) (blog *models.BlogsArch, err error) {
	// sqlStr := `select
	// 		blog_id, title, content, author_id, community_id, create_time, update_time
	// 		from blog
	// 		where blog_id = ?`
	// blog = new(models.BlogsArch)
	// err = db.Get(blog, sqlStr, id)
	// if err == sql.ErrNoRows {
	// 	err = ErrBlogIdNotExisted
	// }
	return nil, nil
}

func QueryBlogList(page int64, size int64) (blogs []*models.BlogsArch, err error) {
	// blogs = make([]*models.BlogsArch, 0, size)
	// sqlStr := `select
	// 		blog_id, title, content, author_id, community_id, create_time, update_time
	// 		from blog
	// 		order by create_time
	// 		desc
	// 		limit ?,?` // 按照时间由新到旧进行返回
	// err = db.Select(&blogs, sqlStr, (page-1)*size, size)
	return nil, nil
}

// 根据传递的Ids进行查询数据
func GetBlogWithOrder(ids []string) ([]*models.BlogsArch, error) {
	// sqlStr := `select
	// 		blog_id, title, content, author_id, community_id, create_time, update_time
	// 		from blog
	// 		where blog_id in (?)
	// 		order by FIND_IN_SET(blog_id, ?)`
	// q, a, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	// if err != nil {
	// 	return nil, err
	// }
	// q = sqlx.Rebind(sqlx.QUESTION, q)
	// res := make([]*models.BlogsArch, 0, len(ids))
	// if err := db.Select(&res, q, a...); err != nil {
	// 	return nil, err
	// }
	//return res, err
	return nil, nil
}
