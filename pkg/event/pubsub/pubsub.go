package eventpubsub

import (
	"cloud.google.com/go/pubsub"
	"github.com/vendasta/gosdks/logging"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

type pubsubWorker struct {
	subscriptionName string
	ctx              context.Context
	messageHandler   handlerFn
	pubsubClient     *pubsub.Client
}

func NewPubsubWorker(ctx context.Context, subscriptionName string, pubsubClient *pubsub.Client, messageHandler handlerFn) *pubsubWorker {
	return &pubsubWorker{
		subscriptionName: subscriptionName,
		pubsubClient:     pubsubClient,
		ctx:              ctx,
		messageHandler:   messageHandler,
	}
}

func (w *pubsubWorker) consume() error {
	// Create an iterator to pull messages via subscription1.
	it, err := w.pubsubClient.Subscription(w.subscriptionName).Pull(w.ctx)
	if err != nil {
		logging.Errorf(w.ctx, "Error pulling subscription %s: %s", w.subscriptionName, err)
		return err
	}

	defer it.Stop()

	// Consume all the messages
	for {
		msg, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logging.Errorf(w.ctx, "Error getting message: %s", err.Error())
			continue
		}
		logging.Infof(w.ctx, "Message data: %s", msg.Data)
		err = w.messageHandler(w.ctx, msg.Data) // Call the provided message handler
		if err != nil {
			logging.Errorf(w.ctx, "Error processing message: %s", err.Error())
			continue
		}
		msg.Done(true) // Acknowledge that we've consumed the message.
	}
	return nil
}

func (w *pubsubWorker) Work() error {
	for {
		err := w.consume()
		if err != nil {
			logging.Errorf(w.ctx, "Error consuming messages")
			return err
		}
	}
	return nil
}
