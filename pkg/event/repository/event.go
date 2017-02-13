package eventrepository

// EventRepository is the interface for the event domain, implemented by vstore components
type EventRepository struct {
	*VStoreRepository
}

// NewEventRepository returns a new instance of a EventRepository
func NewEventRepository(v *VStoreRepository) *EventRepository {
	return &EventRepository{v}
}
