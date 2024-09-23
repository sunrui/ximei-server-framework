/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-04-28 11:08:58
 */

package amqp

// Subscriber 订阅接口
type Subscriber[P any, T any] interface {
	OnProcess(this P, msgId string, t T) (ack bool)     // 成功
	OnError(this P, msgId string, err error) (ack bool) // 失败
}
