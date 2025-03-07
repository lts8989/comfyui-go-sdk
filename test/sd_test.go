package test

import (
	"fmt"
	"github.com/lts8989/comfyui-go-sdk/model"
	"github.com/lts8989/comfyui-go-sdk/sdk"
	"os"
	"testing"
)

const ClientId = "533ef3a3-39c0-4e39-9ced-37d290f371f8"

func init() {
	if err := sdk.Setup("https://rnuix-34-133-62-132.a.free.pinggy.link/", ClientId, 1, 2); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func TestWsServer(t *testing.T) {
	sdk.ConnectToWebSocket(receviedMsg)
	select {}
}

func receviedMsg(receive model.WsReceive) error {
	return nil
}

func TestGetSystemStats(t *testing.T) {
	stats, err := sdk.GetSystemStats()
	if err != nil {
		t.Errorf("GetSystemStats api err:%+v", err)
		return
	}
	if len(stats.Devices) == 0 {
		t.Errorf("stats.Devices is empty")
	} else {
		t.Log("stats.Devices len is ", len(stats.Devices))
	}
}

func TestPrompt(t *testing.T) {
	promptContent, err := os.ReadFile("default_temp.txt")
	if err != nil {
		t.Errorf("提示词文件打开错误，err:%+v", err)
		return
	}
	promptResp, err := sdk.Prompt(promptContent)
	if err != nil {
		t.Errorf("绘图任务下发api错误，err:%+v", err)
		return
	}
	if len(promptResp.PromptId) == 0 {
		t.Errorf("绘图任务下发错误，msg:%s", promptResp.Error.Message)
		return
	}
	t.Logf("prompt_id:%s", promptResp.PromptId)
}

func TestHistory(t *testing.T) {
	history, err := sdk.History("PromptId")
	if err != nil {
		t.Errorf("获取任务执行结果错误，err:%+v", err)
		return
	}
	if len(history) == 0 {
		t.Errorf("任务结果为空，%+v", history)
		return
	}
	t.Logf("获取任务执行结果成功，%+v", history)
}

func TestView(t *testing.T) {
	req := model.ViewReq{
		Filename:  "",
		Subfolder: "",
		Type:      "",
	}
	_, err := sdk.View(req)
	if err != nil {
		t.Errorf("获取任务执行结果错误，err:%+v", err)
		return
	}

	t.Log("获取任务执行结果成功")
}
