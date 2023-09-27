package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var secretID string
var secretKey string
var bucketURL string

// 腾讯云对象存储相关操作

func GetTencentPreSignedUrl(user_id string, file_name string) (string, name string, public_url string) {

	secretID = viper.GetString("credentials.secretID")
	secretKey = viper.GetString("credentials.secretKey")
	bucketURL = viper.GetString("bucketURL")

	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: secretID, // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: secretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})

	// 通过环境变量获取密钥
	// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
	ak := secretID // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
	// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
	sk := secretKey // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140

	name = fmt.Sprintf("/%s/%s", user_id, file_name) //应该用不同人的 usrid 去拼凑出不同的 key
	ctx := context.Background()

	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodPut, name, ak, sk, time.Hour, nil)
	if err != nil {
		panic(err)
	}

	public_url = fmt.Sprintf("%s%s", bucketURL, name)

	return presignedURL.String(), name, public_url

}
