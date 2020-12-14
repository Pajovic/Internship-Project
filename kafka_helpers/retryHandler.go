package kafka_helpers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type RetryHandler struct {
	Reader *kafka.Reader
	Writer *KafkaProducer
}

const (
	contextDeadlineExceeded = "context deadline exceeded"
)

func (handler *RetryHandler) TransferToMainTopic(w http.ResponseWriter, r *http.Request) {
	returnMessage := ""
	statusCode := 200

	numberOfTransferedProducts := 0

	for {
		m, err := handler.fetchWithTimeout(context.Background())
		if err != nil {
			kafkaError := strings.TrimSpace(err.Error())
			if kafkaError == contextDeadlineExceeded {
				returnMessage = "There are no messages to be read"
				break
			}
			fmt.Println("Error while fetching message from retry topic")
			continue
		}

		err = handler.Writer.WriteMessage(string(m.Value), string(m.Key))
		if err != nil {
			fmt.Println("Error while writing message to main topic")
			continue
		}

		numberOfTransferedProducts++

		err = handler.Reader.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("Failed to commit message")
			continue
		}
	}

	if numberOfTransferedProducts != 0 && returnMessage != "" {
		returnMessage = fmt.Sprintf("Managed to transfer %d products from deadletter queue.", numberOfTransferedProducts)
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(returnMessage))
}

func (handler *RetryHandler) fetchWithTimeout(ctx context.Context) (kafka.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Second)
	defer cancel()
	return handler.Reader.FetchMessage(ctx)
}
