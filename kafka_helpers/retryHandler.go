package kafka_helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

type RetryHandler struct {
	Reader *kafka.Reader
	Writer *KafkaProducer
}

const (
	contextDeadlineExceeded = "context deadline exceeded"
	badRunsLimit            = 4
)

func (handler *RetryHandler) TransferToMainTopic(w http.ResponseWriter, r *http.Request) {
	returnMessage := ""
	statusCode := 200

	numberOfTransferedProducts := 0

	badRunsFetching := 0
	badRunsTransfering := 0
	badRunsCommiting := 0

	for {
		returnMessage, statusCode = determineBadRunError(badRunsFetching, badRunsTransfering, badRunsCommiting)
		if returnMessage != "" {
			break
		}

		m, err := handler.fetchWithTimeout(context.Background())
		if err != nil {
			if err == errors.New(contextDeadlineExceeded) {
				returnMessage = "There are no messages to be read"
				break
			}
			fmt.Println("Error while fetching message from retry topic ")
			badRunsFetching++
			continue
		}

		err = handler.Writer.WriteMessage(string(m.Value), string(m.Key))
		if err != nil {
			fmt.Println("Error while writing message to main topic")
			badRunsTransfering++
			continue
		}

		numberOfTransferedProducts++

		err = handler.Reader.CommitMessages(context.Background(), m)
		if err != nil {
			log.Println("Failed to commit message")
			badRunsCommiting++
			continue
		}
	}

	if numberOfTransferedProducts != 0 && returnMessage != "" {
		returnMessage = fmt.Sprintf("Managed to transfer %d products from deadletter queue.", numberOfTransferedProducts)
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(returnMessage))
}

func determineBadRunError(fetchNum int, transferNum int, commitNum int) (string, int) {
	returnMessage := ""
	statusCode := 200

	if fetchNum > badRunsLimit {
		returnMessage = "Failed to fetch too many times"
		statusCode = 500
	}

	if transferNum > badRunsLimit {
		returnMessage = "Failed to transfer too many times"
		statusCode = 500
	}

	if commitNum > badRunsLimit {
		returnMessage = "Failed to commit too many times"
		statusCode = 500
	}

	return returnMessage, statusCode
}

func (handler *RetryHandler) fetchWithTimeout(ctx context.Context) (kafka.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Second)
	defer cancel()
	return handler.Reader.FetchMessage(ctx)
}
