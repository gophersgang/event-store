package eventrepository

import (
	"time"

	"github.com/vendasta/gosdks/vstore"
	"golang.org/x/net/context"
	"github.com/vendasta/event-store/pkg/event"
	"github.com/vendasta/event-store/pkg/utils"
)

// VStoreRepository partially implements OpportunityRepository
type VStoreRepository struct {
	client vstore.Interface
}

// CreateOpportunity creates an opportunity and store it
func (v *VStoreRepository) CreateEvent(ctx context.Context, aggregateType event.AggregateType, aggregateID string, payload string, timestamp time.Time) (*event.Event, error) {
	var e *event.Event
	eventID := event.GenerateEventID()

	err := v.client.Transaction(ctx, event.KeySet(eventID), func(t vstore.Transaction, m vstore.Model) error {
		if m == nil {
			e = &event.Event{
				EventID:        eventID,
				AggregateType:  aggregateType,
				AggregateID:    aggregateID,
				Payload:        payload,
				Timestamp:      timestamp,
			}
		} else {
			return utils.Error(utils.AlreadyExists, "An existing event with the given identifier %s already exists.", eventID)
		}
		return t.Save(e)
	})
	if err != nil {
		return nil, err
	}
	return e, nil
}

// NewVStoreRepository returns a new instance of a VStoreRepository
func NewVStoreRepository(v vstore.Interface) *VStoreRepository {
	return &VStoreRepository{v}
}
