package service

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/dao"
	"github.com/lllllan-fv/gateway-admin/internal/proxy/models"
	"github.com/lllllan-fv/gateway-admin/public/consts"
	"github.com/lllllan-fv/gateway-admin/public/resp"

	load_balance "github.com/lllllan-fv/gateway-admin/internal/proxy/service/load_balance"
)

func HTTPAccessMode(c *gin.Context) (*models.GatewayServiceInfo, error) {
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path

	for _, service := range dao.ListService(consts.HttpLoadType) {
		if service.RuleType == consts.DomainHTTPRuleType {
			if service.Rule == host {
				return service, nil
			}
		} else {
			if strings.HasPrefix(path, service.Rule) {
				return service, nil
			}
		}
	}

	return nil, errors.New("not matched service")
}

func NewLoadBalanceReverseProxy(c *gin.Context, lb load_balance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	// 请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		// todo 优化点3
		if err != nil || nextAddr == "" {
			panic("get next addr fail")
		}

		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}

		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	// 更改内容
	modifyFunc := func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		return nil
	}

	// 错误回调 ：关闭real_server时测试，错误回调
	// 范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		resp.Error(c, 999, err)
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
