package model

type SystemStatsResp struct {
	System struct {
		OS             string   `json:"os"`
		RamTotal       int64    `json:"ram_total"`
		RamFree        int64    `json:"ram_free"`
		ComfyuiVersion string   `json:"comfyui_version"`
		PythonVersion  string   `json:"python_version"`
		PytorchVersion string   `json:"pytorch_version"`
		EmbeddedPython bool     `json:"embedded_python"`
		Argv           []string `json:"argv"`
	} `json:"system"`
	Devices []struct {
		Name           string `json:"name"`
		Type           string `json:"type"`
		Index          int    `json:"index"`
		VramTotal      int64  `json:"vram_total"`
		VramFree       int64  `json:"vram_free"`
		TorchVramTotal int64  `json:"torch_vram_total"`
		TorchVramFree  int64  `json:"torch_vram_free"`
	} `json:"devices"`
}

type PromptReq struct {
	ClientId string `json:"client_id"`
	Prompt   any    `json:"prompt"`
}

type PromptResp struct {
	PromptId string `json:"prompt_id"`
	Error    struct {
		Type      string `json:"type"`
		Message   string `json:"message"`
		Details   string `json:"details"`
		ExtraInfo any    `json:"extra_info"`
	}
}

type ViewReq struct {
	Filename  string `json:"filename"`  // 文件名
	Subfolder string `json:"subfolder"` // 子文件夹
	Type      string `json:"type"`      // 类型
}

const (
	ReceiveTypeStart     = "execution_start"
	ReceiveTypeExecuting = "executing"
	ReceiveTypeSuccess   = "executed"
)

var ReceiveTypeMap = map[string]int8{
	ReceiveTypeStart:     1,
	ReceiveTypeExecuting: 2,
	ReceiveTypeSuccess:   3,
}

var ReceiveTypeDescMap = map[int8]string{
	1: "开始执行",
	2: "执行中",
	3: "执行完成",
}

type WsReceive struct {
	Type string `json:"type"`
	Data struct {
		Node        string `json:"node"`
		DisplayNode string `json:"display_node"`
		Output      struct {
			Images []struct {
				Filename  string `json:"filename"`
				Subfolder string `json:"subfolder"`
				Type      string `json:"type"`
			} `json:"images"`
		} `json:"output"`
		PromptID string `json:"prompt_id"`
	} `json:"data"`
}

type ReceivedMsgFun func(WsReceive) error
