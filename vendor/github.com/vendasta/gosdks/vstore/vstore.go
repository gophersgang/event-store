package vstore

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"errors"
	"github.com/vendasta/gosdks/pb/vstorepb"
	"google.golang.org/grpc/codes"
)

//ClientOption modifies the settings and behaviour of the vstore client.
type ClientOption func(*clientOption)

type clientOption struct {
	client      vstorepb.VStoreClient
	adminClient vstorepb.VStoreAdminClient
	environment *env
	dialOptions []grpc.DialOption
}

// Environment sets the environment vStore should connect to. Either `Environment` or `Client` must be supplied.
func Environment(e env) ClientOption {
	return func(c *clientOption) {
		c.environment = &e
	}
}

// Client manually sets the client for vStore to use. Either `Environment` or `Client` must be supplied.
func Client(client vstorepb.VStoreClient) ClientOption {
	return func(c *clientOption) {
		c.client = client
	}
}

// AdminClient manually sets the admin client for vStore to use.
func AdminClient(client vstorepb.VStoreAdminClient) ClientOption {
	return func(c *clientOption) {
		c.adminClient = client
	}
}

// GRPCDialOptions manually sets the dial options for vStore to use when establishing connections.
func GRPCDialOptions(dialOptions ...grpc.DialOption) ClientOption {
	return func(c *clientOption) {
		c.dialOptions = dialOptions
	}
}

// New returns a new vStore client
func New(cs ...ClientOption) (Interface, error) {
	co := clientOption{}
	for _, c := range cs {
		c(&co)
	}
	if co.environment == nil {
		co.environment = Env()
	}
	if co.client == nil {
		var err error
		co.client, co.adminClient, err = newClient(*co.environment, co.dialOptions...); if err != nil {
			return nil, err
		}
	}

	return &vStore{client: co.client, adminClient: co.adminClient}, nil
}

type lookupOption struct {
	pageSize int64
	cursor   string
	filters  []string
}

//LookupOption augments the behavior of the lookup API
type LookupOption func(*lookupOption)

//PageSize sets a maximum page size on the lookup request
func PageSize(pageSize int64) LookupOption {
	return func(l *lookupOption) {
		l.pageSize = pageSize
	}
}

//Cursor sets a cursor on the lookup request, enabling easy paging
func Cursor(cursor string) LookupOption {
	return func(l *lookupOption) {
		l.cursor = cursor
	}
}

//Filter sets a prefix filter on the lookup
// Entities in VStore:
// ['LIS-101', 'RVW-223']
// ['LIS-101', 'RVW-676']
// ['LIS-555', 'RVW-444']
// Applying a filter like Filter([]string{'LIS-101'}) would mean that the lookup api would return ['LIS-101', 'RVW-223'] and ['LIS-101', 'RVW-676'] but not ['LIS-555', 'RVW-444']
func Filter(filters []string) LookupOption {
	return func(l *lookupOption) {
		l.filters = filters
	}
}

// Interface is the complete VStore client interface
type Interface interface {
	GetMulti(context.Context, []*KeySet) ([]Model, error)
	Get(context.Context, *KeySet) (Model, error)
	Lookup(ctx context.Context, namespace, kind string, opts ...LookupOption) (*LookupResult, error)
	Transaction(context.Context, *KeySet, func(Transaction, Model) error) error

	CreateNamespace(ctx context.Context, namespace string, authorizedServiceAccounts []string) (error)
	UpdateNamespace(ctx context.Context, namespace string, authorizedServiceAccounts []string) (error)
	DeleteNamespace(ctx context.Context, namespace string) (error)

	CreateKind(ctx context.Context, schema *Schema) (error)
	UpdateKind(ctx context.Context, schema *Schema) (error)
	GetKind(ctx context.Context, namespace string, kind string) (*vstorepb.GetKindResponse, error)
	DeleteKind(ctx context.Context, namespace, kind string) (error)

	RegisterKind(ctx context.Context, namespace, kind string, serviceAccounts []string, model Model) (*vstorepb.GetKindResponse, error)
	GetSecondaryIndexName(ctx context.Context, namespace, kind string, indexID string) (string, error)
}

// Implements vstore.Interface
type vStore struct {
	client      vstorepb.VStoreClient
	adminClient vstorepb.VStoreAdminClient
	user        *UserInfo
}

// GetMulti returns a set of rows for the given set of keysets.
func (v *vStore) GetMulti(ctx context.Context, keysets []*KeySet) ([]Model, error) {
	kspbs := make([]*vstorepb.KeySet, len(keysets))
	for i, ks := range keysets {
		kspbs[i] = ks.ToKeySetPB()
	}

	res, err := v.client.Get(ctx, &vstorepb.GetRequest{KeySets:kspbs}, grpc.FailFast(false)); if err != nil {
		return nil, err
	}
	var models []Model
	for _, e := range res.Entities {
		var m Model
		if e.Entity != nil {
			m, err = StructPBToModel(e.Entity.Namespace, e.Entity.Kind, e.Entity.Values); if err != nil {
				return nil, err
			}
		}
		models = append(models, m)
	}
	return models, nil
}

// Get returns a single row from vStore by its KeySet.
func (v *vStore) Get(ctx context.Context, keyset *KeySet) (Model, error) {
	entities, err := v.GetMulti(ctx, []*KeySet{keyset}); if err != nil {
		return nil, err
	}
	return entities[0], nil
}

// Lookup supports fetching a page of rows from vStore for a single namespace/kind.
func (v *vStore) Lookup(ctx context.Context, namespace, kind string, opts ...LookupOption) (*LookupResult, error) {
	options := lookupOption{pageSize: 10}
	for _, opt := range opts {
		opt(&options)
	}
	var filter *vstorepb.LookupFilter
	if len(options.filters) > 0 {
		filter = &vstorepb.LookupFilter{
			Keys: options.filters,
		}
	}
	req := &vstorepb.LookupRequest{
		Namespace: namespace,
		Kind: kind,
		PageSize: options.pageSize,
		Cursor: options.cursor,
		Filter: filter,
	}
	r, err := v.client.Lookup(ctx, req, grpc.FailFast(false)); if err != nil {
		return nil, err
	}
	var m = make([]Model, len(r.Entities))
	for i, e := range r.Entities {
		mod, err := StructPBToModel(e.Entity.Namespace, e.Entity.Kind, e.Entity.Values); if err != nil {
			return nil, err
		}
		m[i] = mod
	}
	return &LookupResult{Results: m, NextCursor: r.NextCursor, HasMore: r.HasMore}, nil
}

// Transaction allows updating a single row in vStore.
func (v *vStore) Transaction(ctx context.Context, ks *KeySet, f func(Transaction, Model) error) error {

	tx := &transaction{}
	res, err := v.client.Get(ctx, &vstorepb.GetRequest{KeySets:[]*vstorepb.KeySet{ks.ToKeySetPB()}}, grpc.FailFast(false)); if err != nil {
		return err
	}
	entity := res.Entities[0];
	var m Model
	if entity.Entity != nil {
		m, err = StructPBToModel(entity.Entity.Namespace, entity.Entity.Kind, entity.Entity.Values); if err != nil {
			return err
		}
	}
	err = f(tx, m); if err != nil {
		return err
	}
	if tx.toSave == nil {
		return nil
	}
	s, err := ModelToStructPB(tx.toSave); if err != nil {
		return err
	}
	e := vstorepb.Entity{
		Namespace: ks.namespace,
		Kind: ks.kind,
		Values: s,
	}
	var rpcError error
	if m == nil {
		e.Version = 1
		_, rpcError = v.client.Create(ctx, &vstorepb.CreateRequest{Entity: &e})
		return rpcError
	}
	e.Version = entity.Entity.Version
	_, rpcError = v.client.Update(ctx, &vstorepb.UpdateRequest{Entity: &e})
	return rpcError
}


// CreateNamespace
func (v *vStore) CreateNamespace(ctx context.Context, namespace string, authorizedServiceAccounts []string) (error) {
	if (v.adminClient == nil) {
		return errors.New("Admin client must be initialized.")
	}
	_, err := v.adminClient.CreateNamespace(ctx,
		&vstorepb.CreateNamespaceRequest{Namespace: namespace, AuthorizedServiceAccounts: authorizedServiceAccounts},
	)
	return err
}


// UpdateNamespace allows the updating of authorized service accounts.
func (v *vStore) UpdateNamespace(ctx context.Context, namespace string, authorizedServiceAccounts []string) (error) {
	if (v.adminClient == nil) {
		return errors.New("Admin client must be initialized.")
	}
	_, err := v.adminClient.UpdateNamespace(ctx,
		&vstorepb.UpdateNamespaceRequest{Namespace: namespace, AuthorizedServiceAccounts: authorizedServiceAccounts},
	)
	return err
}


// DeleteNamespace removes the given namespace from vStore and all of its kinds.  This is a permanent process and
// can not be reversed.
func (v *vStore) DeleteNamespace(ctx context.Context, namespace string) (error) {
	if (v.adminClient == nil) {
		return errors.New("Admin client must be initialized.")
	}
	_, err := v.adminClient.DeleteNamespace(ctx,
		&vstorepb.DeleteNamespaceRequest{Namespace: namespace},
	)
	return err
}

// CreateKind makes a new kind in a specific namespace.
func (v *vStore) CreateKind(ctx context.Context, schema *Schema) (error) {
	if (v.adminClient == nil) {
		return errors.New("Admin client must be initialized.")
	}
	properties, err := PropertiesToPb(schema.Properties...); if err != nil {
		return err
	}
	_, err = v.adminClient.CreateKind(ctx,
		&vstorepb.CreateKindRequest{
			Namespace: schema.Namespace,
			Kind: schema.Kind,
			PrimaryKey: schema.PrimaryKey,
			Properties: properties,
			SecondaryIndexes: SecondaryIndexesToPb(schema.SecondaryIndexes...),
			BackupConfig: BackupConfigToPb(schema.BackupConfig),
		})
	return err
}

// GetKind returns the kind by its namespace/name pair.
func (v *vStore) GetKind(ctx context.Context, namespace string, kind string) (*vstorepb.GetKindResponse, error) {
	if (v.adminClient == nil) {
		return nil, errors.New("Admin client must be initialized.")
	}
	r, err := v.adminClient.GetKind(ctx,
		&vstorepb.GetKindRequest{
			Namespace: namespace,
			Kind: kind,
		}); if err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateKind supports the addition of new fields and updates any supported settings.
func (v *vStore) UpdateKind(ctx context.Context, schema *Schema) (error) {
	if (v.adminClient == nil) {
		return errors.New("Admin client must be initialized.")
	}
	properties, err := PropertiesToPb(schema.Properties...); if err != nil {
		return err
	}
	_, err = v.adminClient.UpdateKind(ctx,
		&vstorepb.UpdateKindRequest{
			Namespace: schema.Namespace,
			Kind: schema.Kind,
			Properties: properties,
			SecondaryIndexes: SecondaryIndexesToPb(schema.SecondaryIndexes...),
		})
	return err
}

// DeleteKind removes a kind from vStore and deletes all of its associated data and secondary indexes. This is a
// permanent process and can not be reversed.
func (v *vStore) DeleteKind(ctx context.Context, namespace, kind string) (error) {
	if (v.adminClient == nil) {
		return errors.New("Admin client must be initialized.")
	}
	_, err := v.adminClient.DeleteKind(ctx,
		&vstorepb.DeleteKindRequest{
			Namespace: namespace,
			Kind: kind,
		})
	return err
}

//LookupResult holds all of the information returned by the Lookup API that is relevant for a client
type LookupResult struct {
	Results    []Model
	NextCursor string
	HasMore    bool
}

// Register handles the registration of a certain kind with VStore, returning the kind as it exists in VStore
func (v *vStore) RegisterKind(ctx context.Context, namespace, kind string, serviceAccounts []string, model Model) (*vstorepb.GetKindResponse, error) {
	RegisterModel(namespace, kind, model)
	err := v.UpdateKind(ctx, model.Schema()); if err == nil {
		sch, err := v.GetKind(ctx, namespace, kind); if err != nil {
			return nil, err
		}
		return sch, nil
	}; if grpc.Code(err) != codes.NotFound {
		return nil, err
	}

	// Kind doesn't exist yet; create it
	err = v.CreateNamespace(ctx, namespace, serviceAccounts); if err != nil && grpc.Code(err) != codes.AlreadyExists {
		return nil, err
	}

	err = v.CreateKind(ctx, model.Schema()); if err != nil {
		return nil, err
	}

	sch, err := v.GetKind(ctx, namespace, kind); if err != nil {
		return nil, err
	}

	return sch, nil
}

// GetSecondaryIndexName will tell you the name of table on the secondary index that VStore has created for the
// secondary index specified by indexID. The possible valid values for indexID are the same as the identifiers for
// the secondary indexes in your model's Schema.
func (v *vStore) GetSecondaryIndexName(ctx context.Context, namespace, kind string, indexID string) (string, error) {
	sch, err := v.GetKind(ctx, namespace, kind); if err != nil {
		return "", err
	}

	si := sch.GetSecondaryIndexes(); if si == nil  {
		return "", errors.New("Could not find any secondary indexes on the schema.")
	}; if len(si) < 1 {
		return "", errors.New("Could not find any secondary indexes on the schema.")
	}

	var r *vstorepb.SecondaryIndex
	for _, i := range si {
		if i.Name == indexID {
			r = i
			break
		}
	}
	if r == nil {
		return "", errors.New("Could not find the specified secondary index on the schema.")
	}

	sql := r.GetCloudSqlConfig(); if sql != nil {
		return sql.GetIndexName(), nil
	}
	es := r.GetEsConfig(); if es != nil {
		return es.GetIndexName(), nil
	}
	esRaw := r.GetEsRawConfig(); if esRaw != nil {
		return esRaw.GetIndexName(), nil
	}
	return "", errors.New("Could not determine the type of secondary index")
}
