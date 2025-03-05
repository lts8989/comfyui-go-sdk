package sd_sdk

import (
	"github.com/lts8989/sd_sdk/model"
	"github.com/lts8989/sd_sdk/sdk"
	"testing"
)

func init() {
	sdk.Setup("https://asdf.sdf.wer/", "", 1, 2)

}
func TestWsServer(t *testing.T) {
	sdk.ConnectToWebSocket(aaa)
	select {}
}

// type ReceivedMsgFun func(WsReceive) error
func aaa(aaddd model.WsReceive) error {
	return nil
}
