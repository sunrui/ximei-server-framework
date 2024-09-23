/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-04-28 00:28:29
 */

package amqp

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
)

// Amqp 消息队列
type Amqp[P any, T any] struct {
	topic      string                  // 主题
	amqpConfig amqp.Config             // 配置
	logAdapter watermill.LoggerAdapter // 适配器
	publisher  *amqp.Publisher         // 发布者
}

// New 创建
func New[P any, T any](config Config, topic string) (mq *Amqp[P, T], err error) {
	mq = &Amqp[P, T]{
		topic: topic,
	}

	mq.amqpConfig = amqp.NewDurableQueueConfig(config.Uri)
	mq.logAdapter = watermill.NewStdLogger(false, false)

	if mq.publisher, err = amqp.NewPublisher(mq.amqpConfig, mq.logAdapter); err != nil {
		return nil, err
	}

	return
}

// Publish 发布
func (mq Amqp[P, T]) Publish(msgId string, t T) error {
	tBytes, _ := json.Marshal(t)
	return mq.publisher.Publish(mq.topic, message.NewMessage(msgId, tBytes))
}

// Subscriber 订阅
func (mq Amqp[P, T]) Subscriber(this P, subscriber Subscriber[P, T]) error {
	var amqpSubscriber *amqp.Subscriber
	var err error

	if amqpSubscriber, err = amqp.NewSubscriber(mq.amqpConfig, mq.logAdapter); err != nil {
		return err
	}

	var messages <-chan *message.Message
	messages, err = amqpSubscriber.Subscribe(context.Background(), mq.topic)
	if err != nil {
		return err
	}

	go func(messages <-chan *message.Message) {
		for msg := range messages {
			var t T
			if err = json.Unmarshal(msg.Payload, &t); err != nil {
				if subscriber.OnError(this, msg.UUID, err) {
					msg.Ack()
				}
			}

			if subscriber.OnProcess(this, msg.UUID, t) {
				msg.Ack()
			}
		}
	}(messages)

	return nil
}

// Stop 停止
func (mq Amqp[P, T]) Stop() error {
	return mq.publisher.Close()
}
