package test

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

var serverAddr = "http://localhost:8080"
var testUserA = "douyinTestUserA"
var testUserB = "douyinTestUserB"

func newExpect(t *testing.T) *httpexpect.Expect {
	//创建一个基于指定配置的 httpexpect.Expect 对象，用于进行 HTTP 请求的测试和断言。
	//在测试中，可以使用该对象执行请求，并使用断言函数检查请求的响应结果是否符合预期。
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,              //提供默认HTTP客户端，用于发送HTTP请求
		BaseURL:  serverAddr,                      //设置服务器地址
		Reporter: httpexpect.NewAssertReporter(t), //用于记录测试结果
		Printers: []httpexpect.Printer{ //创建一个httpexpect.Printer类型的切片
			httpexpect.NewDebugPrinter(t, true), //用于输出调试信息
		},
	})
}

func getTestUserToken(user string, e *httpexpect.Expect) (int, string) {
	registerResp := e.POST("/douyin/user/register/").
		WithQuery("username", user).WithQuery("password", user).
		WithFormField("username", user).WithFormField("password", user).
		Expect().
		Status(http.StatusOK).
		JSON().Object()

	userId := 0
	token := registerResp.Value("token").String().Raw()
	if len(token) == 0 {
		loginResp := e.POST("/douyin/user/login/").
			WithQuery("username", user).WithQuery("password", user).
			WithFormField("username", user).WithFormField("password", user).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		loginToken := loginResp.Value("token").String()
		loginToken.Length().Gt(0)
		token = loginToken.Raw()
		userId = int(loginResp.Value("user_id").Number().Raw())
	} else {
		userId = int(registerResp.Value("user_id").Number().Raw())
	}
	return userId, token
}
