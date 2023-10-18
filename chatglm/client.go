package chatglm

import (
	"fmt"
	"github.com/artisancloud/httphelper"
	"github.com/artisancloud/httphelper/client"
	"github.com/artisancloud/httphelper/dataflow"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

const (
	endpoint = "https://open.bigmodel.cn"
	// /api/paas/v3/model-api/{model}/{invoke_method}
	modelUriFormat     = "/api/paas/v3/model-api/%s/%s"
	asyncTaskUriFormat = "/api/paas/v3/model-api/-/async-invoke/%s"
	sseUriFormat       = "/api/paas/v3/model-api/%s/sse-invoke"
)

const (
	ModelChatGLMPro  = "chatglm_pro"
	ModelChatGLMStd  = "chatglm_std"
	ModelChatGLMLite = "chatglm_lite"
)

type ChatGLM struct {
	helper httphelper.Helper
}

func New(cfg Config) (*ChatGLM, error) {
	conf := &httphelper.Config{
		Config: &client.Config{
			Timeout: time.Second * 300,
		},
		BaseUrl: endpoint,
	}

	token, err := signToken(cfg.APIKey)
	if err != nil {
		return nil, err
	}

	helper, err := httphelper.NewRequestHelper(conf)
	helper.WithMiddleware(func(handle dataflow.RequestHandle) dataflow.RequestHandle {
		return func(request *http.Request, response *http.Response) (err error) {
			request.Header.Set("Authorization", token)
			return handle(request, response)
		}
	})
	if cfg.Debug {
		helper.WithMiddleware(httphelper.HttpDebugMiddleware(true))
	}
	if err != nil {
		return nil, fmt.Errorf("connect to %s failed: %s", conf.BaseUrl, err)
	}
	return &ChatGLM{
		helper: helper,
	}, nil
}

// signToken 签名
func signToken(apiKey string) (token string, err error) {
	parts := strings.Split(apiKey, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid api key: %s", apiKey)
	}
	id := parts[0]
	secret := parts[1]

	payload := jwt.MapClaims{
		"api_key":   id,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
		"timestamp": time.Now().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	jwtToken.Header["sign_type"] = "SIGN"

	signedToken, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// SyncInvoke HTTP同步调用模型
func (c *ChatGLM) SyncInvoke(model string, req InvokeRequest) (response *InvokeSyncResponse, err error) {
	err = c.helper.Df().Method(http.MethodPost).
		Uri(fmt.Sprintf(modelUriFormat, model, "invoke")).
		Json(req).
		Result(&response)
	return response, err
}

// AsyncInvoke HTTP 异步调用模型, 返回 TaskInfo
func (c *ChatGLM) AsyncInvoke(model string, req InvokeRequest) (response *InvokeAsyncTaskResponse, err error) {
	err = c.helper.Df().Method(http.MethodPost).
		Uri(fmt.Sprintf(modelUriFormat, model, "async-invoke")).
		Json(req).
		Result(&response)
	return response, err
}

// AsyncInvokeTaskQuery HTTP 查询异步调用任务结果
func (c *ChatGLM) AsyncInvokeTaskQuery(model string, taskID string) (response *InvokeAsyncTaskQueryResponse, err error) {
	err = c.helper.Df().Method(http.MethodGet).
		Uri(fmt.Sprintf(asyncTaskUriFormat, taskID)).
		Result(&response)
	return response, err
}

// SSEInvoke HTTP SSE调用模型
func (c *ChatGLM) SSEInvoke(model string, req InvokeRequest) (response *InvokeSyncResponse, err error) {
	err = c.helper.Df().Method(http.MethodPost).
		Uri(fmt.Sprintf(sseUriFormat, model)).
		Json(req).
		Result(&response)
	return response, err
}
