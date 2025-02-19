package rabbitConnection

type DocumentIndexingPublisher struct {
	RabbitMQPublisher
}

type DocumentIndexingMessage struct {
	DocumentId   int32  `json:"document_id"`
	DocumentText string `json:"document_text"`
}

const DocumentIndexingExchangeName = "document-indexing-exchange"
const DocumentIndexingQueueName = "document-indexing-queue"

func NewDocumentIndexingPublisher() (*DocumentIndexingPublisher, error) {
	rabbitMQPublisher, err := NewRabbitMQPublisher(
		DocumentIndexingExchangeName,
		DocumentIndexingQueueName,
	)
	if err != nil {
		return nil, err
	}
	return &DocumentIndexingPublisher{
		RabbitMQPublisher: *rabbitMQPublisher,
	}, nil
}

func (dip *DocumentIndexingPublisher) Publish(message DocumentIndexingMessage) error {
	return dip.PublishJSON(message)
}
