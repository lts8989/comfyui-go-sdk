# ComfyUI SDK

ComfyUI SDK 是 [ComfyUI](https://github.com/comfyanonymous/ComfyUI) 服务的开发工具包，旨在为开发者提供一套完整的工具和接口，以便快速构建和扩展
`ComfyUI` 的功能。

## 主要功能

绘图任务下发、任务执行结果查询，任务图片下载，以及接收任务进度推送。

## 提前准备

可访问的 `ComfyUI` 服务，记下域名，后面要用到。

## 初始化

* 调用 Setup 方法初始化参数。

| 参数名          | 类型     | 默认值 | 说明                                                |
|--------------|--------|-----|---------------------------------------------------|
| domain       | string |     | ComfyUI 的域名，需要带有协议，http or https 。                |
| clientId     | string |     | 客户端id，当前服务的唯一标识符。websocket，只会接收当前服务下发的任务通知。       |
| pingSec      | uint   | 30  | 单位（秒）。心跳间隔时间，用于维持当前服务与 ComfyUI 服务的websocket 连接状态。 |
| reconnectSec | uint   | 5   | 单位（秒）。如果 websocket 与 ComfyUI 服务断开链接，重连时间间隔。       |

* 调用 [log.go](log/log.go) 文件中的 `InitLogger` 方法初始化日志模块，默认使用 `fmt` 输出。
## 实现的 HTTP 接口

都在 [sdk.go](sdk/sdk.go) 文件中

### 绘图任务下发

* Prompt

**请求参数**

| 参数名        | 类型     | 说明                     |
|------------|--------|------------------------|
| promptByte | []byte | ComfyUI 的流程图，json格式字符串 |

**返回值**

| 属性名       | 类型     | 说明   |
|-----------|--------|------|
| PromptId  | string | 任务id |
| Error     | struct | 错误相关 |
| Type      | string | 错误类型 |
| Message   | string | 错误信息 |        
| Details   | string | 错误详情 |
| ExtraInfo | any    | 扩展信息 |

ComfyUI 的流程图请在web中调试无误后导出，见下图。

![](asdfasdf)

### 查询任务执行结果

* History

**请求参数**

| 参数名      | 类型     | 默认值 | 说明                |
|----------|--------|-----|-------------------|
| promptId | string |     | 任务id，Prompt接口的返回值 |

**返回值**

返回值为数组，每个 item 的属性如下

| 属性名       | 类型     | 说明   |
|-----------|--------|------|
| Filename  | string | 文件名  |
| Subfolder | string | 子文件夹 |
| Type      | string | 类型   |

### 获取图片

* View

**请求参数**

与 History 方法的返回值一致

| 参数名       | 类型     | 默认值 | 说明   |
|-----------|--------|-----|------|
| Filename  | string |     | 文件名  |
| Subfolder | string |     | 子文件夹 |
| Type      | string |     | 类型   |

**返回值**

返回值为图片的字节流

### 服务状态

* GetSystemStats

获取 ComfyUI 服务相关信息，可用于ComfyUI 服务探活。

**返回值**

| 属性名            | 类型       | 说明             |
|----------------|----------|----------------|
| System         | struct   | 系统参数           |
| OS             | string   | 操作系统           |
| RamTotal       | int64    | 内存总量           |
| RamFree        | int64    | 内存剩余           |
| ComfyuiVersion | string   | ComfyUI版本号     |        
| PythonVersion  | string   | Python版本号      |
| PytorchVersion | string   | PyTorch版本号     |
| EmbeddedPython | bool     | 是否嵌入Python环境   |
| Argv           | []string | ComfyUI版本号启动参数 |
| Devices        | struct   | 显卡参数           |
| Name           | string   | 名称             |
| Type           | string   | 类型             |
| Index          | int      | 索引             |
| VramTotal      | int64    | 显存总量           |
| VramFree       | int64    | 显存剩余           |
| TorchVramTotal | int64    |                |    
| TorchVramFree  | int64    |                |

## websocket 接口

`ConnectToWebSocket` 方法允许客户端通过 `WebSocket` 连接到服务器。此方接收任务执行实时更新的消息。请自定义
`model.ReceivedMsgFun` 方法用于接收 `ComfyUI` 发来的消息。

接收消息对象结构如下

| 属性名         | 类型       | 说明                                                                    |
|-------------|----------|-----------------------------------------------------------------------|
| Type        | string   | 任务状态。execution_start、executing、executed。只有这3种状态有价值，其他状态都表示任务正在执行中或无效。 |
| Data        | struct   |                                                                       |
| Node        | string   | 服务器node编号                                                             |
| DisplayNode | string   | node名称                                                                |
| Output      | struct   |                                                                       |
| images      | []struct | 执行完成时，返回图像结果列表，与 History 接口返回值一致                                      |
| Filename    | string   | 文件名                                                                   |
| Subfolder   | string   | 子文件夹                                                                  |
| Type        | string   | 类型                                                                    |
| PromptID    | string   | 任务id                                                                  |

`websocket` 连接每30s会 `ping` 一次 `ComfyUI` 服务，如果断开连接，每 5s 尝试重连。如果因为 `websocket`
断开连接导致错过了任务的执行结果的推送，可以调用 `History` 接口获取任务执行状态以及结果。

## 示例项目

[ComfyUI API](https://github.com/lts8989/sd_api)

## 安装

    go get github.com/lts8989/sd_sdk

## ps

* 方法名、参数和返回值的属性名保持与 `ComfyUI` 的接口一致。
* 项目大部分代码使用 `LLM` 生成