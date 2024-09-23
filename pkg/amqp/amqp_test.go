/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-04-28 00:34:09
 */

package amqp

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type SmsSend struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type SmsSendSubscriber struct{}

func (smsSendSubscriber SmsSendSubscriber) OnProcess(_ any, msgId string, smsSend SmsSend) (ack bool) {
	smsSendBytes, _ := json.Marshal(smsSend)
	fmt.Println("OnProcess", msgId, string(smsSendBytes))
	return true
}

func (smsSendSubscriber SmsSendSubscriber) OnError(_ any, msgId string, err error) (ack bool) {
	fmt.Println("OnError", msgId, err.Error())
	return true
}

func TestMq_Publish(t *testing.T) {
	mq, err := New[any, SmsSend](Config{
		Uri: "amqp://guest:guest@localhost:5672/",
	}, "sms.send.topic")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("%d", i)
		if err = mq.Publish(id, SmsSend{
			Phone: fmt.Sprintf("130000000%02d", i),
			Code:  id,
		}); err != nil {
			t.Error(err)
		}
	}

	go func() {
		if err = mq.Subscriber(nil, SmsSendSubscriber{}); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(10 * time.Second)
	_ = mq.Stop()
}
