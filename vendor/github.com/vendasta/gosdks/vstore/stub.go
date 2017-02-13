package vstore

import (
	"golang.org/x/net/context"
	"github.com/vendasta/gosdks/pb/vstorepb"
	"sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strconv"
)

// NewHappyPathStub returns an in memory validation that has little to no validation but supports creating
// namespaces/kinds as well as getting/creating/updating entities in vstore. It does not do any secondary indexing.
func NewHappyPathStub() *VStoreHappyPathStub {
	return &VStoreHappyPathStub{
		entities: map[string][]*entity{},
		kinds: map[string][]string{},
	}
}


// Implements vstore.Interface
type VStoreHappyPathStub struct {
	entities   map[string][]*entity
	namespaces []string
	kinds      map[string][]string
	sync.Mutex
}

type entity struct {
	m  Model
	ks *KeySet
}

func (vs *VStoreHappyPathStub) GetMulti(ctx context.Context, keysets []*KeySet) ([]Model, error) {
	vs.Lock()
	defer vs.Unlock()

	toReturn := []Model{}
	for _, ks := range keysets {
		kspb := ks.ToKeySetPB()
		if !StringInSlice(kspb.Namespace, vs.namespaces) {
			return nil, grpc.Errorf(codes.InvalidArgument, "Namespace doesnt exist")
		}
		if !StringInSlice(kspb.Kind, vs.kinds[kspb.Namespace]) {
			return nil, grpc.Errorf(codes.InvalidArgument, "Kind doesnt exist")
		}
		found := false
		for _, existing := range vs.entities[kspb.Namespace + kspb.Kind] {
			if StringsEqual(existing.ks.keys, ks.keys) {
				found = true
				toReturn = append(toReturn, clone(existing.ks.namespace, existing.ks.kind, existing.m))
				break
			}
		}
		if !found {
			toReturn = append(toReturn, nil)
		}
	}
	return toReturn, nil
}

func (vs *VStoreHappyPathStub) Get(ctx context.Context, ks *KeySet) (Model, error) {
	r, err := vs.GetMulti(ctx, []*KeySet{ks}); if err != nil {
		return nil, err
	}
	return r[0], err
}

func (vs *VStoreHappyPathStub)  Lookup(ctx context.Context, namespace, kind string, opts ...LookupOption) (*LookupResult, error) {
	vs.Lock()
	defer vs.Unlock()

	if !StringInSlice(namespace, vs.namespaces) {
		return nil, grpc.Errorf(codes.InvalidArgument, "Namespace doesnt exist")
	}
	if !StringInSlice(kind, vs.kinds[namespace]) {
		return nil, grpc.Errorf(codes.InvalidArgument, "Kind doesnt exist")
	}
	var err error

	lo := &lookupOption{pageSize: 10, cursor: ""}
	for _, opt := range opts {
		opt(lo)
	}

	// set index to start scanning, based on cursor
	startIdx := 0
	if lo.cursor != "" {
		startIdx, err = strconv.Atoi(lo.cursor)
		if err != nil {
			return nil, grpc.Errorf(codes.InvalidArgument, "Invalid cursor")
		}
	}

	hasMore := false
	nextCursor := ""
	toReturn := []Model{}
	for i, e := range vs.entities[namespace + kind] {
		// skip entries until we reach the correct startIdx defined by the cursor
		if i < startIdx {
			continue
		}

		// if there are filters, we need to actually see if the keys of the entities match our filters, else just add the next entity
		if len(lo.filters) > 0 {
			match := true
			// for each prefix filter, check that the entity key's prefix is equal
			// if any of the key components are not equal to the defined prefixes in order, the entity is not a match
			for i, prefix := range lo.filters {
				if len(e.ks.keys) < i {
					match = false
					break
				} else if e.ks.keys[i] != prefix {
					match = false
					break
				}
			}
			if match {
				toReturn = append(toReturn, e.m)
			}
		} else {
			toReturn = append(toReturn, e.m)
		}

		// if this happens, there is more than a single page of results remaining, so set the information the client
		// will need to iterate over the next 1<=n<=lo.pageSize results
		if int64(len(toReturn)) > lo.pageSize {
			// remove the last item, we need to check to see if there are actually more results
			toReturn = toReturn[:len(toReturn)-1]
			hasMore = true
			nextCursor = strconv.FormatInt(int64(i), 10)
			break
		}
	}
	return &LookupResult{Results: toReturn, NextCursor: nextCursor, HasMore: hasMore}, nil
}

func (vs *VStoreHappyPathStub) Transaction(ctx context.Context, ks *KeySet, tx func(Transaction, Model) error) error {
	t := &transaction{}
	m, err := vs.Get(ctx, ks); if err != nil {
		return err
	}

	vs.Lock()
	defer vs.Unlock()

	err = tx(t, m); if err != nil {
		return err
	}
	index := -1
	for i, existing := range vs.entities[ks.namespace + ks.kind] {
		if StringsEqual(existing.ks.keys, ks.keys) {
			index = i
			break
		}
	}
	e := &entity{m: clone(ks.namespace, ks.kind, t.toSave), ks: ks}
	if index == -1 {
		vs.entities[ks.namespace + ks.kind] = append(vs.entities[ks.namespace + ks.kind], e)
	} else {
		vs.entities[ks.namespace + ks.kind][index] = e
	}

	return nil
}

func (vs *VStoreHappyPathStub) CreateNamespace(ctx context.Context, namespace string, authorizedServiceAccounts []string) (error) {
	vs.Lock()
	defer vs.Unlock()

	if StringInSlice(namespace, vs.namespaces) {
		return grpc.Errorf(codes.AlreadyExists, "Namespace already exists")
	}
	vs.namespaces = append(vs.namespaces, namespace)
	return nil
}

func (vs *VStoreHappyPathStub) UpdateNamespace(ctx context.Context, namespace string, authorizedServiceAccounts []string) (error) {
	vs.Lock()
	defer vs.Unlock()

	if !StringInSlice(namespace, vs.namespaces) {
		return grpc.Errorf(codes.NotFound, "Namespace doesnt exist.")
	}
	return nil
}

func (vs *VStoreHappyPathStub) DeleteNamespace(ctx context.Context, namespace string) (error) {
	vs.Lock()
	defer vs.Unlock()
	if !StringInSlice(namespace, vs.namespaces) {
		return grpc.Errorf(codes.NotFound, "Namespace doesnt exist.")
	}
	namespaces := []string{}
	for _, ns := range vs.namespaces {
		if ns == namespace {
			continue
		}
		namespaces = append(namespaces, ns)
	}
	vs.namespaces = namespaces
	delete(vs.kinds, namespace)
	return nil
}

func (vs *VStoreHappyPathStub) CreateKind(ctx context.Context, schema *Schema) (error) {
	vs.Lock()
	defer vs.Unlock()
	if !StringInSlice(schema.Namespace, vs.namespaces) {
		return grpc.Errorf(codes.NotFound, "Namespace doesnt exist.")
	}

	if StringInSlice(schema.Kind, vs.kinds[schema.Namespace]) {
		return grpc.Errorf(codes.AlreadyExists, "Kind already exists")
	}
	vs.kinds[schema.Namespace] = append(vs.kinds[schema.Namespace], schema.Kind)
	return nil
}

func (vs *VStoreHappyPathStub) UpdateKind(ctx context.Context, schema *Schema) (error) {
	vs.Lock()
	defer vs.Unlock()
	if !StringInSlice(schema.Namespace, vs.namespaces) {
		return grpc.Errorf(codes.NotFound, "Namespace doesnt exist.")
	}

	if !StringInSlice(schema.Kind, vs.kinds[schema.Namespace]) {
		return grpc.Errorf(codes.NotFound, "Kind doenst exist.")
	}
	return nil
}

func (vs *VStoreHappyPathStub) GetKind(ctx context.Context, namespace string, kind string) (*vstorepb.GetKindResponse, error) {
	vs.Lock()
	defer vs.Unlock()
	return nil, nil
}

func (vs *VStoreHappyPathStub) DeleteKind(ctx context.Context, namespace, kind string) (error) {
	vs.Lock()
	defer vs.Unlock()
	if !StringInSlice(namespace, vs.kinds[namespace]) {
		return grpc.Errorf(codes.NotFound, "Kind doesnt exist.")
	}
	delete(vs.kinds, namespace)
	return nil
}

func (vs *VStoreHappyPathStub) RegisterKind(ctx context.Context, namespace, kind string, serviceAccounts []string, model Model) (*vstorepb.GetKindResponse, error) {
	RegisterModel(namespace, kind, model)
	err := vs.UpdateKind(ctx, model.Schema()); if err == nil {
		sch, err := vs.GetKind(ctx, namespace, kind); if err != nil {
			return nil, err
		}
		return sch, nil
	}; if grpc.Code(err) != codes.NotFound {
		return nil, err
	}

	// Kind doesn't exist yet; create it
	err = vs.CreateNamespace(ctx, namespace, serviceAccounts); if err != nil && grpc.Code(err) != codes.AlreadyExists {
		return nil, err
	}

	err = vs.CreateKind(ctx, model.Schema()); if err != nil {
		return nil, err
	}
	return &vstorepb.GetKindResponse{}, nil
}

func (vs *VStoreHappyPathStub) GetSecondaryIndexName(ctx context.Context, namespace, kind string, indexID string) (string, error) {
	vs.Lock()
	defer vs.Unlock()

	return "", nil
}

func StringInSlice(target string, list []string) bool {
	for _, candidate := range list {
		if candidate == target {
			return true
		}
	}
	return false
}

func StringsEqual(a []string, b[]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func clone(namespace, kind string, m Model) Model {
	s, err := ModelToStructPB(m); if err != nil {
		panic(err.Error())
	}
	mClone, err := StructPBToModel(namespace, kind, s); if err != nil {
		panic(err.Error())
	}
	return mClone
}
