package event

import (
	"time"

	"golang.org/x/net/context"
)

// Repository Interface for managing the storage of Opportunities
type Repository interface {
	CreateEvent(ctx context.Context, aggregateType AggregateType, aggregateID string, payload string, timestamp time.Time) (*Event, error)
	//ListEvents(ctx context.Context, accountGroupID string, cursor string, pageSize int64) ([]*Opportunity, string, bool, error)
}

//// ElasticRepository Interface for operations implemented by elastic
//type ElasticRepository interface {
//	ListOpportunityKeyIDsForAccountGroup(ctx context.Context, accountGroupID string, cursor string, pageSize int64) ([]*KeyIDs, string, bool, error)
//	GetOpportunityCountByAccountGroups(ctx context.Context, accountGroupIDs []string) (*AccountGroupOpportunityCounts, error)
//}
