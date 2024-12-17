package mysql

import (
	"NothingBlog/models"
	"errors"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type DaoComment struct {
}

func (dm *DaoComment) GetCount() {

}

// 用于后端管理获取评论列表 - 只返回当前评论，不用返回对应的子评论
func (DaoComment) GetWithPageForAdmin(info *models.CommentWithPageParams) ([]models.CommentWithName, int64, error) {
	var commentList []models.CommentWithName

	// 首先查出所有的评论和评论对应的文章名称,以及父评论的用户姓名，并且按照创建时间降序排列（最新的排在前面）
	tempQuery := Db.Table("comments").Select("comments.*, articles.title article_name, c.user_name parent_comment_name").
		Joins("left join articles on articles.article_id = comments.article_id").
		Joins("left join comments c on c.id=comments.parent_comment_id").
		Order("created_at desc").Find(&commentList)
	if tempQuery.Error != nil {
		return nil, 0, tempQuery.Error
	}

	// 对 tempQuery 查出的数据进一步过滤
	if info.AtcId > 0 { // 筛选对应文章的评论
		tempQuery = tempQuery.Where("comments.article_id = ?", info.AtcId)
	}

	if info.Keyword != "" { // 筛选对应关键字的评论
		tempQuery = tempQuery.Where("content like concat('%',?,'%')", info.Keyword)
	}
	// Type > 0; 表示所有的类型的评论，0表示所有，1表示文章的评论 其他未实现
	if info.Type > 0 { // 筛选对应类型的评论
		tempQuery = tempQuery.Where("type = ?", info.Type)
	}
	// Status > 0; 表示所有状态的评论 1-审核中 2-审核通过 3-审核未通过
	if info.Status > 0 {
		tempQuery = tempQuery.Where("status = ?", info.Status)
	}

	var total int64
	if err := tempQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := int(info.PageIdx)
	size := int(info.Size)
	if page > 0 && size > 0 { // 只有在传递的page和size都大于0时才按页查找
		tempQuery = tempQuery.Limit(size).Offset(size * (page - 1))
	}

	if err := tempQuery.Find(&commentList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Debug("评论记录未找到")
			return make([]models.CommentWithName, 0), 0, nil
		}
		return nil, 0, err
	}

	return commentList, total, nil
}

// 返回以给定pid为parent_comment_id的评论
func findAllChildComment(pid uint, comments []models.ResponseCommentListForClient) []models.ResponseCommentListForClient {
	res := make([]models.ResponseCommentListForClient, 0)
	// 这里的优化思路：每次递归的循环次数最大都是一定的，应该可以找到某一条评论的位置之后就可以删掉该评论以减少
	// 后续的递归中循环的次数
	for _, v := range comments {
		// 找到了一条子评论
		if v.ParentCommentId == pid {
			// 递归查找子评论的子评论
			v.SubComment = findAllChildComment(v.Id, comments)
			if v.SubComment == nil { // 加这个是为了避免为空时json序列化成null，前端需要[]
				v.SubComment = make([]models.ResponseCommentListForClient, 0)
			}
			res = append(res, v)
		}
	}
	if len(res) == 0 {
		return nil
	}

	return res
}

// 用于客户端获取评论列表 - 返回所有的评论内容，包括子评论
func (DaoComment) GetWithPageForClient(atcId string) ([]models.ResponseCommentListForClient, error) {
	var tempList []models.ResponseCommentListForClient
	commentList := make([]models.ResponseCommentListForClient, 0) // 最终的返回结果

	// 首先根据文章ID筛选出所有的评论
	//query := Db.Table("comments").Select("comments.*").Where("article_id = ? and status = ? and type = ?", atcId, 1, 1).Find(&tempList)

	// 查询指定的字段 - 客户端不需要所有字段
	// 这里的问题：parent_comment_id可能为null，会导致后面的值无法正常接收，这里相当于判断如果是null的话就赋值0
	rows, err := Db.Raw(`SELECT id, created_at, user_name, ifnull(parent_comment_id,0), content, icon, agree, disagree 
				FROM comments
				WHERE article_id = ? and status = ? and type = ?`, atcId, 1, 1).Rows()
	if err != nil {
		zap.L().Debug("查询出错或者无数据", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	// 循环读取所有的数据
	for rows.Next() {
		var tmp models.ResponseCommentListForClient
		rows.Scan(&tmp.Id, &tmp.CreatedAt, &tmp.UserName, &tmp.ParentCommentId, &tmp.Content, &tmp.Icon, &tmp.Agree, &tmp.Disagree)
		tempList = append(tempList, tmp)
	}

	for _, cmt := range tempList {
		// 如果是顶级评论 parent_comment_id == 0
		if cmt.ParentCommentId == 0 {
			cmt.SubComment = findAllChildComment(cmt.Id, tempList)
			if cmt.SubComment == nil { // 加这个是为了避免为空时json序列化成null，前端需要[]
				cmt.SubComment = make([]models.ResponseCommentListForClient, 0)
			}
			commentList = append(commentList, cmt)
		}
	}

	return commentList, nil
}

// 获取评论列表 - 服务于客户端展示，需要展示对应的子评论
// 客户端的评论只能根据对应的主体ID进行查询，比如文章ID
func GetWithPageForClient(articleId int64) error {
	// 查询出所有的根评论

	// 查询一级根评论的子评论

	return nil
}

func (dm DaoComment) InsertIntoDb(comment *models.Comment) error {
	return Db.Create(comment).Error
}

func (dm DaoComment) DeleteById(id string) error {
	return Db.Where("id = ?", id).Unscoped().Delete(&models.Comment{}).Error
}

func (dm DaoComment) DeleteByIds(ids []uint) error {
	return Db.Where("id in (?)", ids).Unscoped().Delete(&models.Comment{}).Error
}

func (dm DaoComment) UpdateStatus(id uint, state uint8) error {
	return Db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("status", state).Error
}
