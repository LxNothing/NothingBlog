package redis

// 定义redis与文章投票点赞相关的键的常量
const (
	KeyPrefix              = "wevision."
	KeyBlogTimeZset        = "blog.time"   // zset，文章的创建时间
	KeyBlogScoreZset       = "blog.score"  // zset，文章的投票分数
	KeyBlogVotedZsetPrefix = "blog.voted." // zset，记录单个文章的点赞用户，这个key不完整，还要拼接对应的文章id
)

// 获取带前缀的redis key
func getKeyWithPrefix(key string) string {
	return KeyPrefix + key
}
