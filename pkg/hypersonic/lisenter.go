/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 00:12:01
 */

package hypersonic

import (
	"fmt"
	"log/slog"
)

// Listener 监听者
type Listener interface {
	OnLimit(limitType LimitType, req *Request)
	OnLog(req *Request, data *Data, err *Error)
}

// EchoListener 回显监听者
type EchoListener struct {
	requestModelPool ParameterPool[RequestModel] // 请求模型池
}

// NewEchoListener 创建适配器
func NewEchoListener() EchoListener {
	return EchoListener{
		requestModelPool: NewParameterPool[RequestModel](),
	}
}

// OnLimit 限制调用
func (echoListener EchoListener) OnLimit(limitType LimitType, req *Request) {
	msg := fmt.Sprintf("OnLimit Ip => %s, Uri => %s", req.GetIp(), req.GetUri())

	switch limitType {
	case LimitIp:
		slog.Error(msg)
	case LimitUserId:
		slog.Error(msg + fmt.Sprintf(", UserId => %s", req.Token.MustGetUserId()))
	}
}

// OnLog 日志
func (echoListener EchoListener) OnLog(req *Request, data *Data, err *Error) {
	// 请求模型
	requestModel := echoListener.requestModelPool.Acquire()
	defer echoListener.requestModelPool.Release(requestModel)

	response := func() string {
		if data != nil {
			return data.String()
		}

		if err != nil {
			return err.String()
		}

		return "<null>"
	}()

	requestModel.Filed(req, req.Token.GetUserId(), &response, req.GetElapsed())

	if data != nil {
		slog.Debug(requestModel.String())
	}

	if err != nil {
		slog.Error(requestModel.String())
	}
}
