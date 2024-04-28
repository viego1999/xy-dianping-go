package v1

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"xy-dianping-go/internal/common"
	"xy-dianping-go/internal/constants"
	"xy-dianping-go/pkg/utils"
)

type UploadController struct {
}

func NewUploadController() *UploadController {
	return &UploadController{}
}

func (c *UploadController) UploadImage(w http.ResponseWriter, r *http.Request) {
	// 限制上传大小为 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		panic(fmt.Sprintf("文件上传失败：%+v", err))
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		panic(fmt.Sprintf("无效文件：%+v", err))
	}
	defer func(file multipart.File) {
		if err = file.Close(); err != nil {
			panic(fmt.Sprintf("file close error: %+v", err))
		}
	}(file)

	newFileName := c.generateNewFileName(handler.Filename)
	newFilePath := filepath.Join(constants.IMAGE_UPLOAD_DIR, newFileName)

	tempFile, err := os.CreateTemp(constants.IMAGE_UPLOAD_DIR, "upload-*.png")
	if err != nil {
		panic(fmt.Sprintf("文件上传失败：%+v", err))
	}
	defer func(tempFile *os.File) {
		if err := tempFile.Close(); err != nil {
			panic(fmt.Sprintf("临时文件关闭失败：%+v", err))
		}
	}(tempFile)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("文件上传失败：%+v", err))
	}
	if _, err = tempFile.Write(fileBytes); err != nil {
		panic(fmt.Sprintf("文件写入失败：%+v", err))
	}
	if err = os.Rename(tempFile.Name(), newFilePath); err != nil {
		panic(fmt.Sprintf("文件上传失败：%+v", err))
	}

	common.SendResponse(w, common.OkWithData(newFileName))
}

func (c *UploadController) DeleteBlogImg(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("name")
	filePath := filepath.Join(constants.IMAGE_UPLOAD_DIR, filename)

	if fileInfo, err := os.Stat(filePath); os.IsNotExist(err) || fileInfo.IsDir() {
		panic(fmt.Sprintf("错误的文件名称：%+v", err))
	}

	if err := os.Remove(filePath); err != nil {
		panic(fmt.Sprintf("文件删除失败：%+v", err))
	}
	common.SendResponse(w, common.Ok())
}

func (c *UploadController) generateNewFileName(originalFileName string) string {
	ext, uuid := filepath.Ext(originalFileName), utils.GenerateUUID()
	hash := utils.HashCode(uuid)
	d1, d2 := strconv.Itoa(hash&0xF), strconv.Itoa((hash>>4)&0xF)
	dirPath := filepath.Join(constants.IMAGE_UPLOAD_DIR, "blogs", d1, d2)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			panic(fmt.Sprintf("os.MkdirAll error: %+v", err))
		}
	}
	return filepath.Join("blogs", d1, d2, fmt.Sprintf("%s.%s", uuid, ext))
}
