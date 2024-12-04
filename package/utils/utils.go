package utils

func GetTotalPage(pagesize, total int64) int64 {
	total_page := total / int64(pagesize)
	if total%int64(pagesize) != 0 {
		total_page++
	}
	return total_page
}
