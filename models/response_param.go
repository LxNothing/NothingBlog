package models

/*
	与Article文章相关的接口返回数据模型
*/

/*
	与Class种类相关的接口返回数据模型
*/

// 客户端 获取所有的类别时按照这种格式进行返回
// type ResponseClassAllForClient struct {
// 	CurClassName string               // 当前的类别名称
// 	CurTagName   string               // 当前的tag名称
// 	BriefClasses []ResponseClassBrief // 所有的类别信息
// 	CurTagList   []ResponseTagBrief   // 对应class下所包含的tag
// }

// /*
// 	与Tag标签相关的接口返回数据模型
// */
// // 标签Tag的简略信息
// type ResponseTagBrief struct {
// 	TagId    int64  `json:"tag_id,string"` // 标签ID - 由应用层生成
// 	AtcCount int32  `json:"atc_count"`     // 该标签下包含的文章数量
// 	Name     string `json:"name"`          // 类别名称
// }

// 标签Tag的详细信息
// type ResponseTagDetail struct {
// 	ResponseTagBrief
// 	Desc      string    `json:"desc"`       // 标签的描述信息
// 	CreatedAt time.Time `json:"created_at"` // 创建时间
// 	UpdatedAt time.Time `json:"updated_at"` // 更新时间
// }

/*
	与Comment相关的接口返回数据模型
*/
