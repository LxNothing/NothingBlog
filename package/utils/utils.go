package utils

import (
	"crypto/md5"
	"encoding/hex"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

const secret = "xl_ngblog" // md5加密时使用的密钥

// GetTotalPage - 根据传入的每页显示条数，总条数，计算总的页数
func GetTotalPage(pagesize, total int) int {
	total_page := total / pagesize
	if total%pagesize != 0 {
		total_page++
	}
	return total_page
}

// 使用md5算法对内容加密
func EncryptContent(str string) string {
	ecy := md5.New()
	ecy.Write([]byte(secret))
	return hex.EncodeToString(ecy.Sum([]byte(str)))
}

// 通过文件获取文件的类型和后缀
func GetFileTypeAndExtention(file multipart.File) (string, string, error) {
	buffer := make([]byte, 512)
	//设置下一次 Read 或 Write 的偏移量为 offset，
	// whence： 0 表示相对于文件的起始处，1 表示相对于当前的偏移，而 2 表示相对于其结尾处。
	// Seek 返回新的偏移量和一个错误，如果有的话。
	// 这里将文件的读取指针挪到文件开头，避免后续读的时候读到的数据不完整
	defer file.Seek(0, 0)
	if _, err := file.Read(buffer); err != nil {
		zap.L().Debug("读取文件失败", zap.Error(err))
		return "", "", err
	}

	// 判断类型
	fileType := http.DetectContentType(buffer)
	// 获取后缀名称
	extension, _ := mime.ExtensionsByType(fileType)
	if extension == nil {
		return fileType, "", nil
	}
	// 获取类型
	fileType = strings.Split(fileType, "/")[0]

	return fileType, extension[0], nil
}

// func GetVaildPageAndSize(page *int, size *int, totalPage int) {
// 	if *page < 1 {
// 		*page = 1
// 	}

// 	if totalPage != -1 && *page > totalPage {
// 		*page = totalPage
// 	}

// 	if *size < 1 {
// 		*size = int(settings.Confg.SystemConfig.PageSize)
// 	}
// }

// func CheckPageAndSize(page, size uint) bool {
// 	// 客户端设置的页尺寸不能大于配置中的页尺寸
// 	if size > uint(settings.Confg.PageSize) {
// 		return false
// 	}
// 	return true
// }
