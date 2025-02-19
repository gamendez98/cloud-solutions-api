package rabbitConnection

type DocumentIndexingPublisher struct {
	RabbitMQPublisher
}

type DocumentIndexingMessage struct {
	DocumentId   int32  `json:"document_id"`
	DocumentText string `json:"document_text"`
}

const DOCUMENT_INDEXING_EXCHANGE = "document-indexing-exchange"
const DOCUMENT_INDEXING_QUEUE = "document-indexing-queue"

func NewDocumentIndexingPublisher() (*DocumentIndexingPublisher, error) {
	rabbitMQPublisher, err := NewRabbitMQPublisher(
		DOCUMENT_INDEXING_EXCHANGE,
		DOCUMENT_INDEXING_QUEUE,
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
