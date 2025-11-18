package publisher

import (
	"Paprika/models"
	"Paprika/publisher/pb"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type Topic string

type Publisher struct {
	pb.UnimplementedPaprikaServer

	listener   net.Listener
	grpcServer *grpc.Server

	mu         sync.RWMutex
	subscibers map[string]pb.Paprika_GetSpreadStreamServer
}

func New(port string) (*Publisher, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("cannot create gRPC server: %v", err)
	}

	pub := &Publisher{
		listener:   listener,
		subscibers: make(map[string]pb.Paprika_GetSpreadStreamServer),
	}

	s := grpc.NewServer()
	pb.RegisterPaprikaServer(s, pub)
	pub.grpcServer = s

	return pub, nil
}

// Creates a gRPC server
//
// 'port' argument should not be with ':', just "1234"
func (this *Publisher) Listen() error {
	return this.grpcServer.Serve(this.listener)
}

func (this *Publisher) GetSpreadStream(req *pb.StreamRequest, stream pb.Paprika_GetSpreadStreamServer) error {
	this.mu.Lock()
	if _, exist := this.subscibers[req.GetTopic()]; exist {
		return fmt.Errorf("failed to add stream with topic '%s', because stream with this topic already exists", req.GetTopic())
	}

	this.subscibers[req.GetTopic()] = stream
	log.Printf("New Client with topic: %s", req.GetTopic())
	this.mu.Unlock()

	defer func() {
		this.mu.Lock()
		delete(this.subscibers, req.GetTopic())
		this.mu.Unlock()
	}()

	<-stream.Context().Done()
	return stream.Context().Err()
}

func (this *Publisher) PublishSpread(msgTopic string, spread *models.Spread) {
	spreadProto := this.spreadToProto(msgTopic, spread)

	this.mu.RLock()
	defer this.mu.RUnlock()

	for streamTopic, stream := range this.subscibers {
		log.Printf("Pattern: %s | Topic: %s", streamTopic, msgTopic)
		if this.matchTopic(msgTopic, streamTopic) {
			go stream.Send(spreadProto)
		}
	}
}

// cex.futures.mexc.BTCUSDT | cex.*.futures | cex.futures.>
//
// Wildcast pattern like in NATS
func (this *Publisher) matchTopic(topic, pattern string) bool {
	t := strings.Split(topic, ".")
	p := strings.Split(pattern, ".")

	ti := 0
	pi := 0

	for pi < len(p) {
		pp := p[pi]

		if pp == ">" {
			// only allowed at the end
			return pi == len(p)-1
		}

		if ti >= len(t) {
			return false
		}

		if pp != "*" && pp != t[ti] {
			return false
		}

		pi++
		ti++
	}

	return ti == len(t)
}

func (this *Publisher) spreadToProto(topic string, spread *models.Spread) *pb.MsgSpread {
	return &pb.MsgSpread{
		Topic:  topic,
		To:     this.tickerToProto(spread.To),
		From:   this.tickerToProto(spread.From),
		Spread: spread.Value,
	}
}

func (this *Publisher) tickerToProto(tick *models.Ticker) *pb.MsgTicker {
	return &pb.MsgTicker{
		Exchange:     tick.Exchange,
		Market:       tick.Market,
		Symbol:       tick.Symbol,
		Ask:          tick.Ask,
		Bid:          tick.Bid,
		Volume:       tick.Volume,
		IsVolumeUsdt: tick.IsVolumeUsdt,
		Timestamp:    tick.Timestamp,
	}
}

func (this *Publisher) Stop() error {
	this.grpcServer.GracefulStop()
	if err := this.listener.Close(); err != nil {
		return err
	}

	return nil
}
