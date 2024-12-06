package utils

// GetTotalPage - 根据传入的每页显示条数，总条数，计算总的页数
func GetTotalPage(pagesize, total int64) int64 {
	total_page := total / int64(pagesize)
	if total%int64(pagesize) != 0 {
		total_page++
	}
	return total_page
}
