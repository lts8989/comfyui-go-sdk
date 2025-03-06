package test

import (
	"github.com/lts8989/sd_sdk/model"
	"github.com/lts8989/sd_sdk/sdk"
	"os"
	"testing"
)

const ClientId = "asdfasdf"

func init() {
	sdk.Setup("https://asdf.sdf.wer/", ClientId, 1, 2)

}
func TestWsServer(t *testing.T) {
	sdk.ConnectToWebSocket(receviedMsg)
	select {}
}

// type ReceivedMsgFun func(WsReceive) error
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
	promptResp, err := sdk.Prompt(ClientId, promptContent)
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
