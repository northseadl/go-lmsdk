package chatglm

type PromptRole string

const (
	RoleUser      PromptRole = "user"
	RoleAssistant PromptRole = "assistant"
)

const (
	TaskStatusProcessing = "PROCESSING"
	TaskStatusSuccess    = "SUCCESS"
	TaskStatusFailed     = "FAILED"
)

type InvokeRequest struct {
	// 调用对话模型时，将当前对话信息列表作为提示输入给模型; 按照 Role => Content 的键值对形式进行传参; 总长度超过模型最长输入限制后会自动截断，需按时间由旧到新排序
	Prompt []InvokeReqPrompt `json:"prompt"`
	// 采样温度，控制输出的随机性，必须为正数, 取值范围是：(0.0,1.0]，不能等于 0,默认值为 0.95
	Temperature *float64 `json:"temperature"`
	// 用温度取样的另一种方法，称为核取样, 取值范围是：(0.0, 1.0) 开区间，不能等于 0 或 1，默认值为 0.7, 模型考虑具有 top_p 概率质量tokens的结果
	TopP *float64 `json:"top_p"`
	// 由用户端传参，需保证唯一性；用于区分每次请求的唯一标识，用户端不传时平台会默认生成
	RequestId string `json:"request_id"`
	// [程序自动设定] SSE接口调用时，用于控制每次返回内容方式是增量还是全量
	Incremental bool `json:"incremental"`
	// [程序自动设定] json_string: 返回标准的 JSON 字符串; text: 返回原始的文本内容
	ReturnType string `json:"return_type"`
	// 用于控制请求时的外部信息引用，目前用于控制是否引用外部信息
	Ref *InvokeRef `json:"ref"`
}

type InvokeReqPrompt struct {
	Role    PromptRole `json:"role"`
	Content string     `json:"content"`
}

type InvokeRef struct {
	// 是否启用外部信息引用
	Enable bool `json:"enable"`
	// 搜索时的query词，不指定则默认按照prompt信息进行搜索
	SearchQuery string `json:"search_query"`
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

type TaskInfo struct {
	RequestID  string `json:"request_id"`
	TaskID     string `json:"task_id"`
	TaskStatus string `json:"task_status"`
}

type Choice struct {
	Role    PromptRole `json:"role"`
	Content string     `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type InvokeSyncResponse struct {
	BaseResponse
	Data struct {
		Choices []Choice `json:"choices"`
		TaskInfo
	}
	Usage Usage `json:"usage"`
}

type InvokeAsyncTaskResponse struct {
	BaseResponse
	Data TaskInfo `json:"data"`
}

type InvokeAsyncTaskQueryResponse struct {
	BaseResponse
	Data struct {
		TaskInfo
		Choices []Choice `json:"choices"`
		Usage   Usage    `json:"usage"`
	}
}
