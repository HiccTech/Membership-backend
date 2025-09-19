package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hiccpet/service/response"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	// 配置环境变量OSS_ACCESS_KEY_ID。
	accessKeyId = ""
	// 配置环境变量OSS_ACCESS_KEY_SECRET。
	accessKeySecret = ""
	// host的格式为bucketname.endpoint。将${your-bucket}替换为Bucket名称。将${your-endpoint}替换为OSS Endpoint，例如oss-cn-hangzhou.aliyuncs.com。
	host = "https://pet-img.oss-ap-southeast-1.aliyuncs.com"
	// 指定上传到OSS的文件前缀。
	uploadDir = ""
	// 指定过期时间，单位为秒。
	expireTime = int64(3600)
)

func init() {
	godotenv.Load(".oss")
	accessKeyId = os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("OSS_ACCESS_KEY_SECRET")

}

type ConfigStruct struct {
	Expiration string          `json:"expiration"`
	Conditions [][]interface{} `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"ossAccessKeyId"`
	Host        string `json:"host"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func getGMTISO8601(expireEnd int64) string {
	return time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
}

func getPolicyToken() string {
	now := time.Now().Unix()
	expireEnd := now + expireTime
	tokenExpire := getGMTISO8601(expireEnd)

	var config ConfigStruct
	config.Expiration = tokenExpire

	// 添加文件前缀限制
	config.Conditions = append(config.Conditions, []interface{}{"starts-with", "$key", uploadDir})

	// 添加文件大小限制，例如1KB到20MB
	minSize := int64(1024)
	maxSize := int64(20 * 1024 * 1024)
	config.Conditions = append(config.Conditions, []interface{}{"content-length-range", minSize, maxSize})

	result, err := json.Marshal(config)
	if err != nil {
		fmt.Println("callback json err:", err)
		return ""
	}

	encodedResult := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	io.WriteString(h, encodedResult)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	policyToken := PolicyToken{
		AccessKeyId: accessKeyId,
		Host:        host,
		Signature:   signedStr,
		Policy:      encodedResult,
		Directory:   uploadDir,
	}

	response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
		return ""
	}

	return string(response)
}

func GetPostSignatureForOssUpload(c *gin.Context) {

	response.Success(c, getPolicyToken())
}
