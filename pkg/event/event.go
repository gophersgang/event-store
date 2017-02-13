package event

import (
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pborman/uuid"
	"github.com/vendasta/gosdks/logging"
	"github.com/vendasta/gosdks/pb/vstorepb"
	"github.com/vendasta/gosdks/vstore"
	"golang.org/x/net/context"
	"github.com/vendasta/gosdks/pb/event-store/v1"
	"github.com/vendasta/event-store/pkg/config"
	"github.com/vendasta/event-store/pkg/utils"
)

const (
	// KIND for all subject entities.
	KIND = "DomainEvent"
)

// AggregateType is an enum of possible aggregate types
type AggregateType int32

const (
	CAMPAIGN AggregateType = 0
)

// Event vStore model
type Event struct {
	EventID  string `vstore:"event_id"`
	AggregateType AggregateType `vstore:"aggregate_type"`
	AggregateID string `vstore:"aggregate_id"`
	Payload string `vstore:"payload"`
	Timestamp    time.Time `vstore:"timestamp"`
}

// Proto converts Currency to Protobuf Currency
func (a AggregateType) Proto() (eventstore_v1.AggregateType, error) {
	switch a {
	case CAMPAIGN:
		return eventstore_v1.AggregateType_CAMPAIGN, nil
	}
	return -1, errors.New("Could not convert aggregate type to proto")
}

// AggregateTypeFromProto converts Protobuf AggregateType to AggregateType
func AggregateTypeFromProto(a eventstore_v1.AggregateType) (AggregateType, error) {
	switch a {
	case eventstore_v1.AggregateType_CAMPAIGN:
		return CAMPAIGN, nil
	}
	return -1, utils.Error(utils.InvalidArgument, "Could not convert aggregate type from proto")
}

// Proto converts Event to Protobuf Event
func (e Event) Proto() (*eventstore_v1.Event, error) {
	aggregateType, err := e.AggregateType.Proto()
	if err != nil {
		return nil, err
	}

	timestamp, err := ptypes.TimestampProto(e.Timestamp)
	if err != nil {
		return nil, err
	}

	return &eventstore_v1.Event{
		EventId:        e.EventID,
		AggregateType:  aggregateType,
		AggregateId:    e.AggregateID,
		Payload:        e.Payload,
		Timestamp:      timestamp,
	}, nil
}

// Schema representing an Event
func (e *Event) Schema() *vstore.Schema {
	event := vstore.NewPropertyBuilder().StringProperty(
		"event_id",
		vstore.Required(),
	).IntegerProperty(
		"aggregate_type",
		vstore.Required(),
	).StringProperty(
		"aggregate_id",
		vstore.Required(),
	).StringProperty(
		"payload",
		vstore.Required(),
	).TimeProperty(
		"timestamp",
		vstore.Required(),
	).Build()

	backupConfig := vstore.NewBackupConfigBuilder().PeriodicBackup(vstorepb.BackupConfig_WEEKLY).Build()
	return vstore.NewSchema(config.Namespace, KIND, []string{"event_id"}, event, nil, backupConfig)
}

// Initialize handles all bootstrapping required to be able to serve subject operations for this instance
func Initialize(ctx context.Context, vstoreClient vstore.Interface) error {
	logging.Debugf(ctx, "Attempting to Register Kind.")
	//vstoreClient.DeleteKind(ctx, config.Namespace, KIND)
	_, err := vstoreClient.RegisterKind(ctx, config.Namespace, KIND, []string{config.GetServiceAccount()}, (*Event)(nil))
	if err != nil {
		logging.Errorf(ctx, "Error Registering Kind: %s", err)
		return err
	}
	//indexName, err := vstoreClient.GetSecondaryIndexName(ctx, config.Namespace, KIND, config.ElasticIndexID)
	//if err != nil {
	//	logging.Errorf(ctx, "Error Getting Elastic Index Name: %s", err)
	//	return err
	//}
	//logging.Infof(ctx, "Elastic index for opportunities: %s", indexName)
	//config.OpportunityElasticIndex = indexName
	return nil
}

// GenerateEventID generates a new Event ID
func GenerateEventID() string {
	return "EVENT-" + uuid.New()
}

// KeySet builder for an Event
func KeySet(eventID string) *vstore.KeySet {
	return vstore.NewKeySet(config.Namespace, KIND, []string{eventID})
}
