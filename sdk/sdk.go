package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lts8989/sd_sdk/log"
	"github.com/lts8989/sd_sdk/model"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type method string

const (
	Post method = "POST"
	Get  method = "GET"
)

var _domain, _clientId string
var _apiProtocol, _wsProtocol string
var _pingSec, _reconnectSec uint

func Setup(domain, clientId string, pingSec, reconnectSec uint) error {
	// 根据domain判断API请求使用的协议
	if len(domain) > 0 {
		if domain[:5] == "https" {
			_apiProtocol = "https"
			_wsProtocol = "wss"
		} else if domain[:4] == "http" {
			_apiProtocol = "http"
			_wsProtocol = "ws"
		}
	}

	if len(_apiProtocol) == 0 {
		log.Error("setup的domain请携带协议，http 或者 https")
		return errors.New("域名没有携带协议，http 或者 https")
	}

	// 过滤掉domain的协议部分和结尾斜杠
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimSuffix(domain, "/")

	_domain = domain
	_clientId = clientId

	if pingSec > 0 {
		_pingSec = pingSec
	} else {
		_pingSec = 30
	}

	if reconnectSec > 0 {
		_reconnectSec = reconnectSec
	} else {
		_reconnectSec = 5
	}

	return nil
}

// callSDAPI 调用 sdapi 的统一方法
func callSDAPI(apiPath string, method method, data interface{}) ([]byte, error) {
	// 构造完整的URL
	apiPath = strings.TrimPrefix(apiPath, "/")
	apiUrl := fmt.Sprintf("%s://%s/%s", _apiProtocol, _domain, apiPath)

	// 创建请求
	var body io.Reader
	if data != nil {
		if method == Post {
			// 将 data 参数序列化为 JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(jsonData)
		} else if method == Get {
			queryString := structToQuery(data)
			apiUrl = fmt.Sprintf("%s?%s", apiUrl, queryString)
		}

	}
	req, err := http.NewRequest(string(method), apiUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// structToQuery 将结构体转换为 URL 查询参数
func structToQuery(p any) string {
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)

	query := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("json")
		if tag != "" {
			query.Set(tag, fmt.Sprintf("%v", field.Interface()))
		}
	}
	return query.Encode()
}

// GetSystemStats 服务状态
func GetSystemStats() (*model.SystemStatsResp, error) {
	responseBody, err := callSDAPI("/system_stats", Get, nil)
	if err != nil {
		log.Errorf("获取服务状态出错:%+v", err)
		return nil, err
	}

	var stats model.SystemStatsResp
	err = json.Unmarshal(responseBody, &stats)
	if err != nil {
		log.Errorf("获取服务状态解析出错:%+v", err)
		return nil, err
	}

	return &stats, nil
}

// View 获取图片
func View(req model.ViewReq) ([]byte, error) {
	responseBody, err := callSDAPI("/view", Get, req)
	if err != nil {
		log.Errorf("获取图片出错:%+v", err)
		return nil, err
	}
	return responseBody, nil
}

// Prompt 绘图任务下发
func Prompt(clientId string, promptByte []byte) (*model.PromptResp, error) {
	paramsObj := make(map[string]interface{})
	if err := json.Unmarshal(promptByte, &paramsObj); err != nil {
		log.Errorf("提示词转json出错，%+v", err)
		return nil, err
	}
	req := model.PromptReq{
		ClientId: clientId,
		Prompt:   paramsObj,
	}
	responseBody, err := callSDAPI("/prompt", Post, req)
	if err != nil {
		log.Errorf("绘图任务下发出错:%+v", err)
		return nil, err
	}
	var resp model.PromptResp
	err = json.Unmarshal(responseBody, &resp)
	if err != nil {
		log.Errorf("绘图任务下发解析json出错:%+v", err)
		return nil, err
	}
	return &resp, nil
}

// History 调用history接口，获取任务执行结果
func History(promptId string) ([]model.ViewReq, error) {
	responseBody, err := callSDAPI("/history/"+promptId, Get, nil)
	if err != nil {
		log.Errorf("获取任务执行结果出错，promptId：%s,err:%+v", promptId, err)
		return nil, err
	}

	var value map[string]interface{}
	err = json.Unmarshal(responseBody, &value)
	if err != nil {
		log.Errorf("获取任务执行结果出错，promptId：%s,err:%+v,resp:%s", promptId, err, responseBody)
		return nil, err
	}

	//region 解析并断言返回值中数据类型
	vAny, ok := value[promptId]
	if !ok {
		return nil, errors.New("没有找到prompt_id")
	}

	vValue, ok := vAny.(map[string]interface{})
	if !ok {
		return nil, errors.New("json数据类型错误")
	}

	promptAny := vValue["prompt"]
	promptValue, ok := promptAny.([]interface{})
	if !ok {
		return nil, errors.New("prompt数据类型错误")
	}

	servIds := make([]string, 0)
	for _, v := range promptValue {
		servIdsAny, ok := v.([]interface{})
		if ok {
			for _, servIdAny := range servIdsAny {
				servId, ok := servIdAny.(string)
				if !ok {
					return nil, errors.New("servid数据类型错误")
				}
				servIds = append(servIds, servId)
			}
		}
	}

	outputsAny := vValue["outputs"]
	outputsValue, ok := outputsAny.(map[string]interface{})
	if !ok {
		return nil, errors.New("outputs数据类型错误")
	}

	list := make([]model.ViewReq, 0)
	for _, servId := range servIds {
		servAny := outputsValue[servId]
		servValue, ok := servAny.(map[string]interface{})
		if !ok {
			return nil, errors.New("serv数据类型错误")
		}

		imagesAny := servValue["images"]
		by, _ := json.Marshal(imagesAny)
		imageList := make([]model.ViewReq, 0)
		_ = json.Unmarshal(by, &imageList)
		list = append(list, imageList...)
	}
	//endregion

	return list, nil

}

var (
	conn *websocket.Conn
)

// connectToWebSocket 连接到第三方 WebSocket 服务
func ConnectToWebSocket(recFun model.ReceivedMsgFun) {
	wsUrl := fmt.Sprintf("%s://%s/ws?clientId=%s", _wsProtocol, _domain, _clientId)
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Errorf("Error connecting to WebSocket,url:%s,%s", wsUrl, err.Error())
		return
	}
	defer conn.Close()

	log.Infof("Connected to WebSocket:%s" + wsUrl)

	go receiveMessages(recFun)

	go monitorConnection(wsUrl, recFun)

}

// receiveMessages 接收 WebSocket 消息
func receiveMessages(receivedMsgFun model.ReceivedMsgFun) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Infof("Error while reading message:%v", err)
			break
		}
		log.Infof("Received message: %s\n", msg)
		var receivedData model.WsReceive
		err = json.Unmarshal(msg, &receivedData)
		if err != nil {
			log.Errorf("Received data json error:%v", err)
		}

		err = receivedMsgFun(receivedData)
		if err != nil {
			log.Errorf("rece err:%v", err)
		}
	}
}

// monitorConnection 定期检查 WebSocket 连接状态
func monitorConnection(wsURL string, receivedMsgFun model.ReceivedMsgFun) {
	for {
		// 发送心跳消息
		err := conn.WriteMessage(websocket.PingMessage, []byte("ping"))
		if err != nil {
			log.Infof("WebSocket connection lost:%v", err)
			reconnect(wsURL, receivedMsgFun) // 连接丢失，尝试重连
			break
		}

		// 等待一段时间再发送下一个心跳
		time.Sleep(time.Duration(_pingSec) * time.Second) // 每 10 秒发送一次心跳
	}
}

// reconnect 尝试重新建立 WebSocket 连接
func reconnect(wsURL string, receivedMsgFun model.ReceivedMsgFun) {
	for {
		log.Info("Attempting to reconnect...")
		time.Sleep(time.Duration(_reconnectSec) * time.Second) // 等待 5 秒后重试连接

		var err error
		conn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
		if err == nil {
			log.Infof("Reconnected to WebSocket:%s", wsURL)
			go receiveMessages(receivedMsgFun)          // 重新启动接收消息的 goroutine
			go monitorConnection(wsURL, receivedMsgFun) // 重新启动监控连接的 goroutine
			return
		}
		log.Infof("Reconnect failed:%v", err)
	}
}
