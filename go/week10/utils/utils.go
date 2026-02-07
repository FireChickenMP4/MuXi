package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	uploadDir = "uploads"
)

func SaveImages(c *gin.Context, file *multipart.FileHeader, prefix string) (string, error) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
	//其实也可以直接用限定创建
	ext := filepath.Ext(file.Filename)
	//ext可以提取扩展名
	timestamp := time.Now().Unix()
	randomStr := uuid.New().String()[:8]
	newFileName := fmt.Sprintf("%s_%d_%s%s", prefix, timestamp, randomStr, ext)

	dst := filepath.Join(uploadDir, newFileName)
	//连接路径

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	//没想到gin本身就有保存文件啊，好好gin
	return filepath.ToSlash(dst), nil
	//统一路径分隔符,会根据系统选择\还是/
}

func RemoveFile(filePath string) error {
	if filePath != "" {
		err := os.Remove(filePath)
		return err
	}
	return nil
	//没有也不用管，反正本身就没
}
func GetEnv(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
