package kafkala

import (
	"context"

	"github.com/segmentio/kafka-go"
)

const (
	createCompanyTopic = "create-company"
	updateCompanyTopic = "update-company"
)

// CompanyProducer interface
type CompanyProducer interface {
	PublishCreate(ctx context.Context, msgs ...kafka.Message) error
	PublishUpdate(ctx context.Context, msgs ...kafka.Message) error
	Close()
	Run()
	GetNewKafkaWriter(topic string) *kafka.Writer
}

type CompaniesProducer struct {
	createWriter *kafka.Writer
	updateWriter *kafka.Writer
}

// NewCompanyProducer constructor
func NewCompanyProducer() *CompaniesProducer {
	cp := CompaniesProducer{}
	cp.Run()
	return &cp
}

// GetNewKafkaWriter Create new kafka writer
func (p *CompaniesProducer) GetNewKafkaWriter(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: topic,
	}
	return w
}

// Run init producers writers
func (p *CompaniesProducer) Run() {
	p.createWriter = p.GetNewKafkaWriter(createCompanyTopic)
	p.updateWriter = p.GetNewKafkaWriter(updateCompanyTopic)
}

// Close close writers
func (p CompaniesProducer) Close() {
	p.createWriter.Close()
	p.updateWriter.Close()
}

// PublishCreate publish messages to create topic
func (p *CompaniesProducer) PublishCreate(ctx context.Context, msgs ...kafka.Message) error {
	return p.createWriter.WriteMessages(ctx, msgs...)
}

// PublishUpdate publish messages to update topic
func (p *CompaniesProducer) PublishUpdate(ctx context.Context, msgs ...kafka.Message) error {
	return p.updateWriter.WriteMessages(ctx, msgs...)
}
