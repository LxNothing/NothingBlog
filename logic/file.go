package logic

import (
	"NothingBlog/package/utils"
	"NothingBlog/settings"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
)

// const (
// 	timeFormat = "20060102150405" // 时间的格式
// )

var (
	ErrWriteFileToLocal = errors.New("写入文件到本地失败")
)

func WriteFileToLocal(fileHeader *multipart.FileHeader) (string, error) {
	extension := path.Ext(fileHeader.Filename)                    // 获取后缀名
	orgName := strings.TrimSuffix(fileHeader.Filename, extension) // 获取文件名的原始名称

	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Debug("打开文件失败", zap.Error(err))
		return "", ErrWriteFileToLocal
	}
	defer file.Close()

	// 判断文件的类型
	fileType, orgExtention, err := utils.GetFileTypeAndExtention(file)
	if err != nil {
		zap.L().Debug("读取文件失败", zap.Error(err))
		return "", ErrWriteFileToLocal
	}

	// 创建保存对应文件的文件夹
	storePath := settings.Confg.SystemConfig.UploadPath + "/" + fileType
	if err := os.MkdirAll(storePath, os.ModePerm); err != nil {
		zap.L().Debug("创建文件夹失败", zap.Error(err), zap.String("path", settings.Confg.SystemConfig.UploadPath+"/"+fileType))
		return "", ErrWriteFileToLocal
	}

	if extension == "" {
		extension = orgExtention // 如果原文件不包含后缀，则使用默认的后缀
	}

	// 生成新的文件名称
	//newName := utils.EncryptContent(orgName) + "_" + time.Now().Format(timeFormat) + extension
	newName := utils.EncryptContent(orgName) + extension
	zap.L().Debug("写入文件到本地", zap.String("OldName", fileHeader.Filename), zap.String("NewName", newName))

	fileStorePath := storePath + "/" + newName
	visitPath := settings.Confg.SystemConfig.VisitPath + "/" + fileType + "/" + newName
	// 创建要保存的新文件
	outFile, err := os.Create(fileStorePath)
	if err != nil {
		zap.L().Debug("创建文件失败", zap.Error(err), zap.String("file", fileStorePath))
		return "", ErrWriteFileToLocal
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		zap.L().Debug("拷贝文件失败", zap.Error(err))
		return "", ErrWriteFileToLocal
	}

	return visitPath, nil
}
