package utils

import (
	"fmt"
	"mime/multipart"
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// OSS 配置（实际项目中应从配置文件或环境变量读取）
const (
	Endpoint        = "oss-cn-guangzhou.aliyuncs.com" // 你的地域节点
	AccessKeyID     = "你的AccessKey"
	AccessKeySecret = "你的Secret"
	BucketName      = "moxin-assets"
)

// UploadImage 上传图片到 OSS 并返回 URL
func UploadImage(file *multipart.FileHeader) (string, error) {
	// 1. 初始化客户端
	client, err := oss.New(Endpoint, AccessKeyID, AccessKeySecret)
	if err != nil {
		return "", err
	}

	// 2. 获取存储空间
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		return "", err
	}

	// 3. 生成唯一文件名 (防止冲突)
	ext := path.Ext(file.Filename)
	newFileName := fmt.Sprintf("articles/%s%s", uuid.New().String(), ext)

	// 4. 打开上传文件流
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 5. 执行上传
	err = bucket.PutObject(newFileName, src)
	if err != nil {
		return "", err
	}

	// 6. 返回图片访问地址
	url := fmt.Sprintf("https://%s.%s/%s", BucketName, Endpoint, newFileName)
	return url, nil
}
