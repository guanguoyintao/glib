package uoss

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type S3OSSUrlInfo struct {
	Bucket string
	Region string
}

var cosOSS2CDN = map[string]string{
	"umu-pd-1303248253.cos.ap-beijing.myqcloud.com":    "statics-umu-pd.umucdn.cn",
	"umu-hh-1303248253.cos.ap-beijing.myqcloud.com":    "statics-umu-hh.umucdn.cn",
	"umu-cn-pg-1303248253.cos.ap-beijing.myqcloud.com": "statics-umu-cn-pg.umucdn.cn",
	"umu-cn-1303248253.cos.ap-beijing.myqcloud.com":    "umu-cn.umucdn.cn",
	"umu-test-1303248253.cos.ap-beijing.myqcloud.com":  "umu-test.umucdn.cn",
}
var cosCDN2OSS = make(map[string]string, len(cosOSS2CDN))

var s3Bucket2CDN = map[string]string{
	"umu.co":            "co.umustatic.com",
	"umu.com":           "com.umustatic.com",
	"umu.tw":            "tw.umustatic.com",
	"umu.io":            "resource.umu.io",
	"umu-daihatsu":      "daihatsu.umustatic.com",
	"cdn.umustatic.com": "cdn-hk-bos.s3.ap-northeast-1.amazonaws.com",
}
var s3Bucket2Region = map[string]string{
	"umu.co":            "ap-northeast-1",
	"umu.com":           "us-west-2",
	"umu.tw":            "ap-northeast-1",
	"umu.io":            "eu-west-1",
	"cdn.umustatic.com": "ap-northeast-1",
	// TODO: umu-daihatsu
}
var s3CDN2OSS = make(map[string]S3OSSUrlInfo, len(s3Bucket2CDN))

func init() {
	// cos map init
	for oss, cdn := range cosOSS2CDN {
		cosCDN2OSS[cdn] = oss
	}
	// s3 map init
	for bucket, cdn := range s3Bucket2CDN {
		region, ok := s3Bucket2Region[bucket]
		if !ok {
			continue
		}
		s3CDN2OSS[cdn] = S3OSSUrlInfo{
			Bucket: bucket,
			Region: region,
		}

	}
	s3CDN2OSS["d1dw9odj9f4xrp.cloudfront.net"] = S3OSSUrlInfo{
		Bucket: "umu.co",
		Region: s3Bucket2Region["umu.co"],
	}
	s3CDN2OSS["dz4z9fk4z9ace.cloudfront.net"] = S3OSSUrlInfo{
		Bucket: "umu.com",
		Region: s3Bucket2Region["umu.com"],
	}
	s3CDN2OSS["d1iwrer5wugjum.cloudfront.net"] = S3OSSUrlInfo{
		Bucket: "umu.tw",
		Region: s3Bucket2Region["umu.tw"],
	}
	s3CDN2OSS["d36hh8vvjpeb6i.cloudfront.net"] = S3OSSUrlInfo{
		Bucket: "umu-daihatsu",
		Region: s3Bucket2Region["umu-daihatsu"],
	}
	s3CDN2OSS["resource.umu.io"] = S3OSSUrlInfo{
		Bucket: "umu.io",
		Region: s3Bucket2Region["umu.io"],
	}
}

func CosCovertUrlOSS2CDN(ossURL string) (string, error) {
	u, err := url.Parse(ossURL)
	if err != nil {
		return "", err
	}
	host, ok := cosOSS2CDN[u.Host]
	if ok {
		u.Host = host
	}

	return u.String(), nil
}

func CosCovertUrlCDN2OSS(cdnURL string) (string, error) {
	u, err := url.Parse(cdnURL)
	if err != nil {
		return "", err
	}
	host, ok := cosCDN2OSS[u.Host]
	if ok {
		u.Host = host
	}

	return u.String(), nil
}

func S3CovertUrlOSS2CDN(ossURL string) (string, error) {
	// 从url中提取bucket
	u, err := url.Parse(ossURL)
	if err != nil {
		return "", err
	}
	p := strings.Split(u.Path, "/")
	if len(p) == 0 {
		return ossURL, nil
	}
	s3Bucket := p[1]
	cdn, ok := s3Bucket2CDN[s3Bucket]
	if ok {
		u.Host = cdn
		u.Path = strings.TrimPrefix(u.Path, fmt.Sprintf("/%s", s3Bucket))
		return u.String(), nil
	}

	return ossURL, nil
}

func S3CovertUrlCDN2OSS(cdnURL string) (string, error) {
	// 从url中提取bucket
	u, err := url.Parse(cdnURL)
	if err != nil {
		return "", err
	}
	oss, ok := s3CDN2OSS[u.Host]
	if ok {
		u.Host = fmt.Sprintf("s3.%s.amazonaws.com", oss.Region)
		u.Path = path.Join(oss.Bucket, u.Path)
		return u.String(), nil
	}

	return cdnURL, nil
}

func CovertUrl2Outer(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	cosPattern := `cos\.[a-z0-9\-]+\.myqcloud\.com`
	s3Pattern := `s3\.[a-z0-9\-]+\.amazonaws\.com`
	var cdn string
	if match, _ := regexp.MatchString(cosPattern, u.Host); match {
		// cos
		cdn, err = CosCovertUrlOSS2CDN(urlString)
		if err != nil {
			return "", err
		}
	} else if match, _ := regexp.MatchString(s3Pattern, u.Host); match {
		// s3
		cdn, err = S3CovertUrlOSS2CDN(urlString)
		if err != nil {
			return "", err
		}
	}

	return cdn, nil
}

func CovertUrl2Inner(urlString string) (string, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	cosPattern := `\.umucdn\.cn`
	s3Pattern := `(umustatic\.com|amazonaws\.com|resource\.umu\.io)`
	var oss string
	if match, _ := regexp.MatchString(cosPattern, u.Host); match {
		// cos
		oss, err = CosCovertUrlCDN2OSS(urlString)
		if err != nil {
			return "", err
		}
	} else if match, _ := regexp.MatchString(s3Pattern, u.Host); match {
		// s3
		oss, err = S3CovertUrlCDN2OSS(urlString)
		if err != nil {
			return "", err
		}
	}

	return oss, nil
}
