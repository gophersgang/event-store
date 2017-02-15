package eventpubsubhandlers

import (
	"github.com/vendasta/gosdks/logging"
	"golang.org/x/net/context"
)

func HandleEmailEvent(ctx context.Context, message []byte) error {
	logging.Infof(ctx, "Handling email event message: %s", message)
	return nil
}

func HandleCampaignCreated(ctx context.Context, message []byte) error {
	logging.Infof(ctx, "Handling campaign created message: %s", message)
	return nil
}
