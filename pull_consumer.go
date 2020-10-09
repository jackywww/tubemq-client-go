/**
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tubeclient

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lubanproj/tubemq-client-go/client"
	"github.com/lubanproj/tubemq-client-go/constants"
	"github.com/lubanproj/tubemq-client-go/errs"
	"github.com/lubanproj/tubemq-client-go/model"
	"github.com/lubanproj/tubemq-client-go/protocol"
	"github.com/lubanproj/tubemq-client-go/rebalance"
	"github.com/lubanproj/tubemq-client-go/utils"
)

// PullConsumer implements Consumer with the pull consumption
type PullConsumer struct {
	config           *Config
	client           client.Client
	fetcher          *model.FetcherManager
	filterConds      []string
	rebalanceEvents  chan *rebalance.Event
	rebalanceResults chan *rebalance.Event

	shutdown        bool
	rebalanceClosed bool

	rmtDataCache *rmtDataCache
}

func (p *PullConsumer) Start() error {
	if err := checkConfig(p.config); err != nil {
		return err
	}

	for _, topic := range p.config.Topics {
		if err := p.SubscribeByFilterConds(topic, p.filterConds); err != nil {
			return err
		}
	}

	go p.HeartBeat()

	go p.rebalance()

	return nil

}

// Subscribe a topic
func (p *PullConsumer) Subscribe(topic string) error {
	return p.SubscribeByFilterConds(topic, nil)
}

// Subscribe a topic
func (p *PullConsumer) SubscribeByFilterConds(topic string, filterConds []string) error {
	if p.shutdown {
		return errors.New("subscribe error, client has been shutdown")
	}

	if topic == "" {
		return errors.New("subscribe error, topic is empty")
	}

	if len(filterConds) > constants.ConfigMaxFilterCondLength {
		return fmt.Errorf("subscribe error, filterConds length exceeds the maximum value : %d",
			constants.ConfigMaxFilterCondLength)
	}

	if err := p.Register(); err != nil {
		return err
	}

	return nil
}

func (p *PullConsumer) HeartBeat() error {
	d := p.config.HeartBeatPeriod
	if d == 0 {
		d = constants.ConfigHeartBeatPeriod
	}

	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		if err := p.heartbeatMaster(); err != nil {
			// TODO add log
		}
		if err := p.heartbeatBroker(); err != nil {
			// TODO add log
		}
	}
}

func (p *PullConsumer) heartbeatMaster() error {
	req := &protocol.HeartRequestC2M{
		GroupName: p.config.Group,
		ClientId:  p.config.ConsumerID,
	}
	rsp := &protocol.HeartResponseM2C{}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()


	opts := []client.Option{
		client.WithTarget(p.config.MasterAddress),
	}
	err := p.client.Invoke(ctx, req, rsp, opts...)
	if err != nil {
		return err
	}

	e := rsp.GetEvent()

	event := &rebalance.Event{
		RebalanceID:      e.GetRebalanceId(),
		Type:             rebalance.EventType(e.GetOpType()),
		Status:           rebalance.EventStatus(e.GetStatus()),
		SubscribeStrList: e.GetSubscribeInfo(),
	}

	p.rebalanceEvents <- event

	return nil
}

func (p *PullConsumer) heartbeatBroker() error {

	req := &protocol.HeartBeatRequestC2B{
		GroupName: p.config.Group,
		ClientId:  p.config.ConsumerID,
	}
	rsp := &protocol.HeartBeatResponseB2C{}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	err := p.client.Invoke(ctx, req, rsp)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return errors.New("heartbeat to broker fail")
	}

	return nil
}

func (p *PullConsumer) rebalance() {
	if p.shutdown || p.rebalanceClosed {
		return
	}

	for {
		event := <-p.rebalanceEvents
		switch event.Type {
		case rebalance.Connect:
		case rebalance.OnlyConnect:
			p.connectBroker(event, false)
			p.rebalanceResults <- event
		case rebalance.Disconnect:
		case rebalance.OnlyDisConnect:
			p.connectBroker(event, true)
			p.rebalanceResults <- event
		case rebalance.Report:
		case rebalance.StopRebalance:
		default:
		}
	}
}

func (p *PullConsumer) connectBroker(event *rebalance.Event, disconnect bool) {

	if p.shutdown || p.rebalanceClosed {
		return
	}

	brokerPartitionMap := make(map[*model.Broker][]*model.Partition)
	subStrList := event.SubscribeStrList
	subInfoList := utils.ConvertSubInfo(subStrList)

	for _, subInfo := range subInfoList {
		broker := subInfo.Partition.Broker
		partition := subInfo.Partition

		partitionList, ok := brokerPartitionMap[broker]
		if !ok {
			partitionList = []*model.Partition{}
			brokerPartitionMap[broker] = partitionList
		}

		if !utils.ContainsKey(partitionList, partition) {
			partitionList = append(partitionList, partition)
		}
	}

	// TODO filter unfinished partitions
	var unregisterPartitions []*model.Partition

	p.RegisterBroker(brokerPartitionMap, unregisterPartitions, disconnect)

}

func (p *PullConsumer) RegisterBroker(
	brokerPartitionMap map[*model.Broker][]*model.Partition,
	unregisteredPartitions []*model.Partition,
	unregister bool,
) {

	if len(brokerPartitionMap) == 0 {
		return
	}

	registerBrokerRetryTimes := 0

	for len(unregisteredPartitions) == 0 && registerBrokerRetryTimes < p.config.MaxRegisterRetryTimes {

		for broker, partitions := range brokerPartitionMap {
			p.registerBrokers(broker, partitions, unregister)
		}

	}

}

func (p *PullConsumer) registerBrokers(
	broker *model.Broker,
	partitions []*model.Partition,
	unregister bool,
) {

	// TODO handle error
	for _, partition := range partitions {
		p.registerBroker(broker, partition, unregister)
	}

}

func (p *PullConsumer) registerBroker(
	broker *model.Broker,
	partition *model.Partition,
	unregister bool,
) error {

	opType := constants.MsgOpTypeRegister

	if unregister {
		opType = constants.MsgOpTypeUnRegister
	}

	req := &protocol.RegisterRequestC2B{
		ClientId:    p.config.ConsumerID,
		GroupName:   p.config.Group,
		OpType:      int32(opType),
		PartitionId: partition.PartitionID,
	}

	rsp := &protocol.RegisterResponseB2C{}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()
	opts := []client.Option{
		client.WithTarget(fmt.Sprintf("%s:%d", broker.Host, broker.Port)),
	}

	if err := p.client.Invoke(ctx, req, rsp, opts...); err != nil {
		return err
	}

	if !rsp.Success || rsp.ErrCode != 0 {
		return errs.New(rsp.ErrCode, rsp.ErrMsg)
	}

	return nil
}

func (p *PullConsumer) PullMessage() (*PullContext, error) {

	partitionExt, err := p.rmtDataCache.selectPartition()
	if err != nil {
		return nil, err
	}
	partition := partitionExt.Partition

	pullContext := &PullContext{
		partition: partition,
	}
	req := &protocol.GetMessageRequestC2B{
		ClientId:         p.config.ConsumerID,
		PartitionId:      partition.PartitionID,
		GroupName:        p.config.Group,
		TopicName:        partition.Topic,
		LastPackConsumed: partitionExt.GetAndResetLastPackConsumed(),
	}

	rsp := &protocol.GetMessageResponseB2C{}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	broker := partition.Broker
	opts := []client.Option{
		client.WithTarget(fmt.Sprintf("%s:%d", broker.Host, broker.Port)),
	}

	err = p.client.Invoke(ctx, req, rsp, opts...)
	if err != nil {
		return nil, err
	}

	if !rsp.Success {
		switch rsp.GetErrCode() {
		case errs.Success:
			// TODO
			fillPullContext(rsp, pullContext)
		case errs.Hb_No_Node:
		case errs.Certificate_Failure:
		case errs.Duplicate_Partition:
			pullContext.setFailResult(rsp.GetErrCode(), rsp.GetErrMsg())
		case errs.ServerConsumeSpeedLimit:
			// TODO
		default:
			// TODO
		}
	}

	return pullContext, nil
}

func (p *PullConsumer) ConfirmConsume(pullContext *PullContext) error {
	confirmContextStrList := strings.Split(pullContext.confirmContext, constants.AttrSep)
	if len(confirmContextStrList) != 4 {
		return errors.New("confirmContext format error, value must be aaaa:bbbb:cccc:ddddd")
	}
	for index, confirmContextStr := range confirmContextStrList {
		if confirmContextStr == "" {
			return errors.New(fmt.Sprintf("ConfirmContext format error, the %d item is empty", index))
		}
	}

	partitionKey := strings.TrimSpace(confirmContextStrList[0]) + constants.AttrSep +
		strings.TrimSpace(confirmContextStrList[1]) + constants.AttrSep +
		strings.TrimSpace(confirmContextStrList[2])
	topicName := strings.TrimSpace(confirmContextStrList[1])
	timestamp, _ := strconv.ParseInt(strings.TrimSpace(confirmContextStrList[3]), 10, 64)

	if !p.rmtDataCache.isPartitionInUse(partitionKey, timestamp) {
		return errors.New("confirmContext's value invalid")
	}

	curPartition := p.rmtDataCache.getPartition(partitionKey)
	if curPartition == nil {
		return errors.New("partition not found by confirmContext")
	}

	req := &protocol.CommitOffsetRequestC2B{
		ClientId:         p.config.ConsumerID,
		TopicName:        topicName,
		PartitionId:      curPartition.Partition.PartitionID,
		GroupName:        p.config.Group,
		LastPackConsumed: curPartition.IsLastPackConsumed,
	}
	rsp := &protocol.CommitOffsetResponseB2C{}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	broker := curPartition.Partition.Broker
	opts := []client.Option{
		client.WithTarget(fmt.Sprintf("%s:%d", broker.Host, broker.Port)),
	}

	err := p.client.Invoke(ctx, req, rsp, opts...)
	if err != nil {
		return err
	}

	if !rsp.Success {
		return errs.New(rsp.ErrCode, rsp.ErrMsg)
	}

	return nil
}

func fillPullContext(rsp *protocol.GetMessageResponseB2C, pullContext *PullContext) {
	pullContext.messageList = rsp.GetMessages()
}

// Register a consumer
func (p *PullConsumer) Register() error {
	var err error

	registerRetryTimes := p.config.MaxRegisterRetryTimes
	if registerRetryTimes == 0 {
		registerRetryTimes = constants.ConfigRegisterRetryTimes
	}

	for i := 0; i < registerRetryTimes ; i++ {
		if err = p.registerMaster(); err == nil {
			fmt.Println(err.Error())
			break
		}
	}

	return err
}

func (p *PullConsumer) registerMaster() error {
	req := &protocol.RegisterRequestC2M{
		ClientId:  p.config.ConsumerID,
		HostName:  utils.GetLocalAddr(),
		GroupName: p.config.Group,
		TopicList: p.config.Topics,
	}

	rsp := &protocol.RegisterResponseM2C{}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.Timeout)
	defer cancel()

	err := p.client.Invoke(ctx, req, rsp)
	if err != nil {
		return err
	}

	if !rsp.Success {
		switch rsp.GetErrCode() {
		case errs.ConsumeGroupForbidden:
			return errors.New("Register fail, group forbidden")
		case errs.ConsumeConnectForbidden:
			return errors.New("Register fail, restricted consume content")
		default:
			return errors.New("Register to master failed, please check and retry later")
		}
	}

	return nil
}
