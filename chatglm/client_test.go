package chatglm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const TestAPIKey = "d7c2dfda469280506236c8ad1c2178b2.def7WeLWt4ebpdog"

func TestNew(t *testing.T) {
	cfg := Config{
		APIKey: TestAPIKey,
		Debug:  false,
	}

	_, err := New(cfg)
	assert.NoError(t, err, "New function should not return an error")
}

func TestSignToken(t *testing.T) {
	apiKey := TestAPIKey
	token, err := signToken(apiKey)

	assert.NoError(t, err, "signToken function should not return an error")
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestSyncInvoke(t *testing.T) {
	cfg := Config{
		APIKey: TestAPIKey,
		Debug:  true,
	}

	c, err := New(cfg)
	assert.NoError(t, err, "New function should not return an error")

	req := InvokeRequest{
		Prompt: []InvokeReqPrompt{
			{
				Role:    RoleUser,
				Content: "你好",
			},
		},
	}

	model := ModelChatGLMPro
	response, err := c.SyncInvoke(model, req)

	assert.NoError(t, err, "SyncInvoke function should not return an error")
	assert.NotNil(t, response, "Response should not be nil")
}

func TestAsyncInvoke(t *testing.T) {
	cfg := Config{
		APIKey: TestAPIKey,
		Debug:  false,
	}

	c, err := New(cfg)
	assert.NoError(t, err, "New function should not return an error")

	req := InvokeRequest{
		Prompt: []InvokeReqPrompt{
			{
				Role:    RoleUser,
				Content: "你好",
			},
		},
	}

	model := ModelChatGLMPro
	_, err = c.AsyncInvoke(model, req)

	assert.NoError(t, err, "AsyncInvoke function should not return an error")
}
