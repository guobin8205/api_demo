package kafka

import (
	"github.com/Shopify/sarama"
	log "github.com/thinkboy/log4go"
	"encoding/json"
	"github.com/axgle/mahonia"
)

const (
	KafkaPushsTopic = "sj_chat"
)

var (
	producer sarama.AsyncProducer
)

type KafkaMsg struct {

}

func InitKafka(kafkaAddrs []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err = sarama.NewAsyncProducer(kafkaAddrs, config)
	go handleSuccess()
	go handleError()
	return
}

func handleSuccess() {
	var (
		pm *sarama.ProducerMessage
	)
	for {
		pm = <-producer.Successes()
		if pm != nil {
			log.Info("producer message success, partition:%d offset:%d key:%v valus:%s", pm.Partition, pm.Offset, pm.Key, pm.Value)
		}
	}
}

func handleError() {
	var (
		err *sarama.ProducerError
	)
	for {
		err = <-producer.Errors()
		if err != nil {
			log.Error("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
		}
	}
}

func MSGpushKafka(msg interface{}) (err error) {
	var b_msg []byte
	if b_msg, err = json.Marshal(msg); err != nil {
		return
	}
	send_msg := string(b_msg)
	enc:=mahonia.NewEncoder("gbk")
	content_gbk := enc.ConvertString(send_msg)
	producer.Input() <- &sarama.ProducerMessage{Topic: KafkaPushsTopic, Value: sarama.StringEncoder(content_gbk)}
	return
}