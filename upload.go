package agin

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiNiu struct {
	Zone          string 	`mapstructure:"zone" json:"zone" yaml:"zone"`
	Bucket        string 	`mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	ImgPath       string 	`mapstructure:"img-path" json:"imgPath" yaml:"img-path"`
	UseHTTPS      bool   	`mapstructure:"use-https" json:"useHttps" yaml:"use-https"`
	AccessKey     string 	`mapstructure:"access-key" json:"accessKey" yaml:"access-key"`
	SecretKey     string 	`mapstructure:"secret-key" json:"secretKey" yaml:"secret-key"`
	UseCdnDomains bool   	`mapstructure:"use-cdn-domains" json:"useCdnDomains" yaml:"use-cdn-domains"`
}

func (q QiNiu) GetUpToken() string {
	bucket:="window-media"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func (q QiNiu) Upload(fileName string) (url string, err error)  {
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	//putExtra := storage.PutExtra{
	//	Params: map[string]string{
	//		"x:name": "github logo",
	//	},
	//}
	//key := strconv.FormatInt(time.Now().Unix(), 8) + "_" + fileName
	key := fileName
	err = formUploader.PutFile(context.Background(), &ret, q.GetUpToken(), key, fileName, nil)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	//fmt.Println(ret.Key,ret.Hash)
	return "http://img.fishfit.fun/" + ret.Key, nil
}