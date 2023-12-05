package uhttp

import (
	"golang.org/x/net/publicsuffix"
	"net/url"
	"path"
	"strings"
)

type TLD struct {
	Domain string
	Port   string
	ICANN  bool
}

func JoinPath(baseUrl string, elem ...string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	paths := make([]string, 0, len(elem)+1)
	paths = append(paths, u.Path)
	for _, e := range elem {
		paths = append(paths, e)
	}
	u.Path = path.Join(paths...)
	s := u.String()

	return s, nil
}

// TLDParser 解析URL，返回一个TLD结构，其中包含额外的字段。
func TLDParser(s string) (*TLD, error) {
	// 解析URL
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	// 如果Host为空，则返回一个空的TLD结构
	if len(u.Host) == 0 {
		return &TLD{Domain: s}, nil
	}
	var domain string
	for i := len(u.Host) - 1; i >= 0; i-- {
		if u.Host[i] == ':' {
			domain = u.Host[:i]
			break
		} else if u.Host[i] < '0' || u.Host[i] > '9' {
			domain = u.Host
			break
		}
	}
	// 获取EffectiveTLDPlusOne
	domain, err = publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return nil, err
	}
	// 获取域名和端口
	dom, port := parsePort(u.Host)
	// 获取公共后缀和ICANN状态
	_, icann := publicsuffix.PublicSuffix(strings.ToLower(dom))

	return &TLD{
		Domain: domain,
		Port:   port,
		ICANN:  icann,
	}, nil
}

// domainPort 从主机字符串中提取域名和端口。
func parsePort(host string) (string, string) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:]
		} else if host[i] < '0' || host[i] > '9' {
			return host, ""
		}
	}
	// 只有当字符串全为数字时才会执行到这里，
	// net/url 应该防止这种情况发生
	return host, ""
}
