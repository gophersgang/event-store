package eventpubsub

import (
	"cloud.google.com/go/pubsub"
	"github.com/vendasta/gosdks/logging"
	"golang.org/x/net/context"
)

type handlerFn func(context.Context, []byte) error

type PubsubSubscription struct {
	name    string
	handler handlerFn
}

func NewPubsubSubscription(name string, handler handlerFn) *PubsubSubscription{
	return &PubsubSubscription{
		name: name,
		handler: handler,
	}
}

type PubsubManager struct {
	ctx           context.Context
	subscriptions []*PubsubSubscription
	projectId     string
}

func NewPubsubManager(ctx context.Context, subscriptions []*PubsubSubscription, projectId string) *PubsubManager {
	return &PubsubManager{
		ctx:           ctx,
		subscriptions: subscriptions,
		projectId:     projectId,
	}
}

func (m *PubsubManager) Start() error {
	client, err := pubsub.NewClient(m.ctx, m.projectId)
	if err != nil {
		logging.Errorf(m.ctx, "Error creating pubsub client: %s", err)
		return err
	}

	for _, subscription := range m.subscriptions {
		go NewPubsubWorker(m.ctx, subscription.name, client, subscription.handler).Work()
	}

	return nil
}
