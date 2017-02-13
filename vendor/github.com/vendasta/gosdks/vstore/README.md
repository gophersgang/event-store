# VStore Client

## Installation

* YOLO development? `go get github.com/vendasta/gosdks`
* Using [glide](https://github.com/Masterminds/glide) package manager? `glide get -s -v github.com/vendasta/gosdks`

## Visualization

The vstore ui is available at https://vstore-`env`.vendasta-internal.com/

* Local: https://vstore-test.vendasta-internal.com/ (no this is not a typo, see why below)
* Test: https://vstore-test.vendasta-internal.com/
* Demo: https://vstore-demo.vendasta-internal.com/
* Production: https://vstore-production.vendasta-internal.com/

From this interface you can look at each table's entities, activity (in terms of job manage by VStore) and configure your email to be alerted by VStore when something bad happens with your table.

### About the local environment

Your local data is actually stored on VStore test, hidden behind a namespace calculated based on your gcloud auth email. On the test environment, the `data-lake` namespace refers to the data-lake app running on test, but the `data-lake-dwalker@vendasta.com` namespace refers to the table that belongs to MY locally running data-lake app. You can not see other people's local namespaces and you can not take any actions against someone else's local namespace.

## Usage

### Introduction

This article is very opinionated and tries to demonstrate concepts with real world examples. These examples are inspired by real code but have been edited for brevity in many places.

If it's TLDR, try jumping around the client library code in your IDE - that's probably a better way to understand what is available, but it might not give you much prescription on how to use it.

Lastly, the hardest/most tedious part is by far the initialization required. Once you get past that, the rest will become more obvious.

If you're just maintaining a system that is already being bootstrapped correctly, then you don't necessarily need to read the sections dealing with initialization.

### Client Initialization

You want to initialize the client at the top level entry point of your application and reuse references to it.
```
var (
	vstoreClient vstore.Interface
)

func main() {
	err = vstore.New(); if err != nil {
		log.Fatalf("Error initializing vStore client. %s", err.Error())
	}
	// now you can inject vstoreClient into structs that need it
}
```

The initialization function also supports passing a variety of options, to customize or extend its behavior, or for testing.

#### Examples
Using a DialOptions option to add an interceptor for logging. **(BEST PRACTICE)**
```
vstoreClient, err = vstore.New(vstore.DialOptions(grpc.WithUnaryInterceptor(logging.ClientInterceptor()))); if err != nil {
    return fmt.Errorf("Error initializing vstore client %s", err.Error())
}
```

Using an Environment option to set the env explicitly (useful for testing)
```
err = vstore.New(vstore.Environment(vstore.Local)); if err != nil {
    log.Fatalf("Error initializing vStore client. %s", err.Error())
}
```

Using a Client option to replace the underlying client with a mock client or an extended, experimental client.
```
clientMock = &vstore.VStoreClientMock{}
err = vstore.New(vstore.Client(clientMock)); if err != nil {
    log.Fatalf("Error initializing vStore client. %s", err.Error())
}
```

Using multiple options in conjunction with one another.
```
clientMock = &vstore.VStoreClientMock{}
err = vstore.New(vstore.Client(clientMock), vstore.Environment(vstore.Local)); if err != nil {
    log.Fatalf("Error initializing vStore client. %s", err.Error())
}
```

### Namespace and Kind initialization
Generally, several replicas of any app you write could be running at any time.
Your app needs to be capable of registering (*in an idempotent manner*) some information with VStore to enable properly scoped communication.

There are 2 concepts in VStore that you need to understand:

1. Namespaces
  * Namespaces provide a way for VStore to understand who a particular Kind is relevant for.
  * For example, if a service account 'repcore-prod' registers a Namespace called 'cs-prod', then only that service account is capable of reading from or writing to any Kinds within the 'cs-prod' namespace, including the creation or deletion of any Kinds.
  * When your app starts, you need to make sure it has a namespace registered with VStore.
  * Normally this only has to happen once, but it is best practice to make this registration happen in the code upon bootstrap, so that if you are deploying to a new environment or running the app locally for the first time, the registration happens automatically.
  * Multiple service accounts can be authorized for a single namespace, but you should have a "good reason" for doing so.
2. Kinds
  * Kinds can be thought of as tables in a typical database. Resource access is authorized on a namespace-by-namespace basis, not kind by kind.
  * In order for you to read or write any entities to a Kind, you first need to make sure it is initialized in VStore.
  * Like namespace registration, you likely only need to create a Kind once, but it is best practice for that logic to live inside your application for portability.
  * Kinds can also be deleted, but let us hope you are not programmatically deleting Kinds. Again, a "good reason" is probably required.

Generally, registering both your namespace(s) and kind(s) can be done within the same stream of execution.

#### Examples

`RawListing` is a struct that implements `vstore.Model`. Think of this type's purpose like you would a class that subclasses `ndb.Model`. It is the domain representation of your data, not necessarily the API representation - it defines what and how you will save that struct's data to VStore.

You can use the RegisterKind method on the client to idempotently register a specific kind and namespace with vstore. This should be done every time an instance is initialized. The operation performed here is equivalent to an "Upsert", so it should also be used to update your kind when you make schema changes.
```
schemaResponse, err := vstoreClient.RegisterKind(ctx, "datalake", "RawListing", []string{"data-lake-local@vendasta.com"}, (*RawListing)(nil))
```

What if you need to know about a secondary index? You need to use `vstore.GetSecondaryIndexName` to have VStore tell you what the "table name" of the index is called on that secondary index.
This will need to be done for every instance on every bootstrap, since secondary indexes and their table names are managed by VStore.
```
// we want to be able to issue all Elasticsearch requests for this kind against this index
myElasticIndexName, err := vstoreClient.GetSecondaryIndexName(ctx, "datalake", "RawListing", "elasticsearch1")
```

Think of it like this:
You've told vstore that for the datalake-RawListing table, you want it to be replicated to elasticsearch, and you chose to call that specific elasticsearch configuration `elasticsearch1`.
VStore is the one who will make that index for you on Elasticsearch and manage it though, so you don't get a direct reference to the name of the table it has created, which is why you should be asking for it using the `GetSecondaryIndexName` api.

### KeySets

A `vstore.KeySet` fulfills the same purpose that `ndb.Key` does in ndb. It is how you construct an identifier for an entity.

#### Examples

A localized wrapper might look like this:
```
func BuildKey(listingID string) *vstore.KeySet {
	return vstore.NewKeySet("datalake", "RawListing", []string{listingID})
}
```

You'll notice the final parameter is a list.
This means that, like in datastore, keys can be composed of multiple parts (hence the 'Set' in KeySet), and you can utilize some special behaviour with the Lookup API if you take advantage of this.
```
func BuildKey(listingID, ReviewID string) *vstore.KeySet {
	return vstore.NewKeySet("datalake", "RawReview", []string{listingID, reviewID})
}
```

### Get API

`vstore.Get` will return either an instance of the struct that you registered with the namespace-kind associated with your provided keyset if the entity exists, or nil if it does not.

`vstore.Get` is strongly consistent. This means that as soon as a `vstore.Update` operation completes, you could call `vstore.Get` on the keyset corresponding to the entity you just updated and get the latest version back, guaranteed.

#### Examples

Since the KeySet contains the information about the namespace, kind, and entity ID, all that you pass to get is a KeySet, constructing one if need be.

```
type vStoreListingService struct {
	vstore vstore.Interface
}

func BuildKey(listingID string) *vstore.KeySet {
	return vstore.NewKeySet("datalake", "RawListing", []string{listingID})
}

func (v *vStoreListingService) Get(ctx context.Context, listingID string) (*RawListing, error) {
	m, err := v.vstore.Get(ctx, KeySet(listingID)); if err != nil {
	    // If an error happens here, it's unexpected. This indicates something is wrong with your vstore configuration or vstore itself.
		logging.Errorf(ctx, "Error getting raw listing from vstore. %s", err.Error())
		return nil, err
	}
	if m == nil {
	    // this means the entity is missing. depending on your use case, you may not want to throw an error here
		return nil, ErrNotFound
	}
	// m at this point is a vstore.Model. That means we can't reference stuff like m.RawListingId
	// or any other properties specific to the RawListing type. We need to cast the vstore.Model
	// back to the type we associated with our namespace and kind previously. Behold the pedantry of static typing!
	l, ok := m.(*RawListing); if !ok {
		logging.Errorf(ctx, "Got unexpected model: %v", m)
		return nil, errors.New("Got unexpected model from vstore.")
	}
	return l, nil
}
```
### GetMulti API

`vstore.GetMulti` will return a list of found results for the provided KeySets. This means that if you provide 10 KeySets to `vstore.GetMulti` but only 8 are found, the list that `vstore.GetMulti` returns will only have 8 items in it.

`vstore.GetMulti` is strongly consistent.

#### Examples

You should pass a list of KeySets to GetMulti.
```
func (v *vstoreListingService) GetMultipleListings(ctx context.Context, listingIds []string) ([]*RawListing, error) {
    var keySets []*vstore.KeySet
    for id := range listingIds {
        keySets = append(keySets, BuildKey(id)
    }

    models, err := v.vstore.GetMulti(ctx, keySets); if err != nil {
	    // If an error happens here, it's unexpected. This indicates something is wrong with your vstore configuration or vstore itself.
		logging.Errorf(ctx, "Error getting raw listings from vstore. %s", err.Error())
		return nil, err
	}

	// convert vstore.Model instances back to our domain objects
    var results []*RawReview
    for _, item := range models {
        e, ok := item.(*RawReview); if !ok {
            logging.Errorf(ctx, "Got unexpected model: %v", e)
            return nil, errors.New("Got unexpected model from vstore.")
        }
        results = append(results, e)
    }

    return results, nil
}
```

### Lookup API

`vstore.Lookup` will, in all expected circumstances, return a list of results, a boolean value indicating whether more results are available on the next page, and a cursor that can be used the get the next page of results if there are more results.

This lookup is strongly consistent, but does not support searching or partial matching on fields.

#### Examples
Using the lookup endpoint to serve just the first page of results. Notice that although we pass the cursor back to the caller, they'd have no way to get the second page since we don't put a cursor into the request.
```
func (v *vStoreReviewService) LookupReviews(ctx context.Context) ([]*RawReview, string, bool, error) (
    r, err := v.vstore.Lookup(ctx, "datalake", "RawReview"); if err != nil {
        // something is wrong with our configuration or vstore, we can't offer any results
        return nil, "", false, err
    }

    // convert vstore.Model instances back to our domain objects
    // hopefully you use a single method for doing this that is common to this layer of operations
    var results []*RawReview
    for _, item := range r.Results {
        e, ok := m.(*RawReview); if !ok {
            logging.Errorf(ctx, "Got unexpected model: %v", e)
            return nil, "", false, errors.New("Got unexpected model from vstore.")
        }
        results = append(results, e)
    }
    return results, r.NextCursor, r.HasMore, nil
```

We can modify our function to accept a cursor param and optionally append a lookup option.
```
func (v *vStoreReviewService) LookupReviews(ctx context.Context, cursor string) ([]*RawReview, string, bool, error) (
    var options []*vstore.LookupOption
    if cursor != "" {
        options = append(options, vstore.Cursor(cursor))
    }
    r, err := v.vstore.Lookup(ctx, "datalake", "RawReview", options...); if err != nil {
        // something is wrong with our configuration or vstore, we can't offer any results
        return nil, "", false, err
    }

    // convert vstore.Model instances back to our domain objects (truncated)
    ...

    return results, r.NextCursor, r.HasMore, nil
```

This means that you can page over the results:
```
firstPageOfResults, cursor, hasMore, err := v.LookupReviews(ctx, "")
nextPageOfResults, cursor, hasMore, err := v.LookupReviews(ctx, cursor)
```

Hopefully the options mechanism makes it obvious how to make your apis more flexible. Let's add an option for page size, just because we can.
```
func (v *vStoreReviewService) LookupReviews(ctx context.Context, cursor string, pageSize int64) ([]*RawReview, string, bool, error) (
    var options []*vstore.LookupOption
    if cursor != "" {
        options = append(options, vstore.Cursor(cursor))
    }
    if pageSize == 0 {
        pageSize = 10 // set a default size of 10 if the parameter is ignored by the caller
    }
    options = append(options, vstore.PageSize(pageSize))
    r, err := v.vstore.Lookup(ctx, "datalake", "RawReview", options...); if err != nil {
        // something is wrong with our configuration or vstore, we can't offer any results
        return nil, "", false, err
    }

    // convert vstore.Model instances back to our domain objects (truncated)
    ...

    return results, r.NextCursor, r.HasMore, nil
```

The Filter `vstore.LookupOption` can be used to execute prefix filters on entity key fragments. This is the same mechanism that Datastore exposes through ancestor queries.

You can imagine that a Review belongs to a Listing. This concept is fairly prevalent throughout the data-oriented parts of our system. A Listing's `KeySet` consists of only a single fragment: `['LIS-XXX']`.

Since a Review belongs to a Listing, you could configure the Review `Kind` to have a `KeySet` consisting of two or more fragments: `['LIS-XXX','RVW-YYY']`. This is, within the scope of BigTable anyway, equivalent to defining a parent-child relationship in NDB.

Using the Lookup API and the Filter option, VStore supports filtering on key prefixes, and in this case we mean you can look up Reviews by their listing id.
```
func (v *vStoreReviewService) LookupReviews(ctx context.Context, cursor string, pageSize int64, listingId string) ([]*RawReview, string, bool, error) (
    var options []*vstore.LookupOption
    if listingId != "" {
        options = append(options, vstore.Filter([]string{listingId}))
    }
    if cursor != "" {
        options = append(options, vstore.Cursor(cursor))
    }
    if pageSize == 0 {
        pageSize = 10 // set a default size of 10 if the parameter is ignored by the caller
    }
    options = append(options, vstore.PageSize(pageSize))
    r, err := v.vstore.Lookup(ctx, "datalake", "RawReview", options...); if err != nil {
        // something is wrong with our configuration or vstore, we can't offer any results
        return nil, "", false, err
    }

    // convert vstore.Model instances back to our domain objects (truncated)
    ...

    return results, r.NextCursor, r.HasMore, nil
```
All of the results returned by the call `results, cursor, hasMore, err := v.LookupReviews(ctx, "", 0, "LIS-225")` will be such that their keys begin with `LIS-225`. This is equivalent to an ancestor query in NDB and possesses the same quality of being strongly consistent. You can page on these results like any other, but you should keep re-applying the filter in subsequent calls - the cursor itself does not encode or map to any filter data.

### Transaction API
Writes to VStore are accomplished using the `vstore.Transaction` API. Non transactional writes are not supported, nor are multi-row transactions, even if the keys involved have a common prefix. You can not read any entity apart from the entity you are writing within your transaction.

(TODO: Add this back) Conflicting transactions will automatically retry and will raise an error if a transaction fails repeatedly. Failed transactions are caused by contention on the entity you're trying to write. Always consider the structure of your data with contention in mind - entities that are prone to being subject to contention are the same ones that experience rapid writes (such as an AccountGroup during the creation process, when we are filling in NAP data from inference or rapidly adding accounts and services in separate write operations).

#### Examples

This method takes a proto representation of the listing and makes sure it is saved with its Modified time adjusted to the current time. This is about the simplest kind of write possible.
```
func (v *vStoreListingService) Replace(ctx context.Context, listing *datalakeproto.RawListing) error {
	// The keyset parameter is the entity or row id (in the case of a create) that you will be performing the transaction upon
	err := v.vstore.Transaction(ctx, KeySet(listing.RawListingId), func(t vstore.Transaction, m vstore.Model) error {
		now, err := ptypes.TimestampProto(time.Now().UTC()); if err != nil {
			return err
		}
		listing.Modified = now

        // Note that we need to convert our listing object to a vstore.Model
        // In this case we reference an adapter function ToModel(l *datalakeproto.RawListing) vstore.Model which does this for us
        // t.Save(YourModel) does the work of actually enforcing that the commit happens or tries to happen.
		return t.Save(ToModel(listing))
	})
	if err != nil {
		logging.Errorf(ctx, "Error committing transaction %s", err.Error())
	}
	return err
}
```

You'll notice that we completely ignored `m vstore.Model` in our inner function, despite it being part of the signature required by the Transaction API. So what's `m`?

`m` is the entity that your `KeySet` parameter is currently pointing to. It will be `nil` if your `KeySet` does not reference an entity (in this case you are creating a new entity). Since our API was simply called `Replace` we didnt care about the presence or absence of `m` in our transaction, it had no bearing on our logic. It should be noted that versions are managed for you by VStore internally and you should not have to worry about them in your application code.

So let's inspect `m` and figure out whether we should be setting a created time or not.
```
func (v *vStoreListingService) Replace(ctx context.Context, listing *datalakeproto.RawListing) error {
	err := v.vstore.Transaction(ctx, KeySet(listing.RawListingId), func(t vstore.Transaction, m vstore.Model) error {
		now, err := ptypes.TimestampProto(time.Now().UTC()); if err != nil {
			return err
		}
		listing.Modified = now
        if m == nil {
		    listing.Created = now
		}

		return t.Save(ToModel(listing))
	})
	if err != nil {
		logging.Errorf(ctx, "Error committing transaction %s", err.Error())
	}
	return err
}
```

## Schemas
A VStore schema is the low level definition of your kind as it will be understood by VStore. The client application is responsible for defining and registering their schema.

Let's look at an example kind and its schema as it exists in the client application:
```
type Song struct {
	ArtistId        string `vstore:"artist_id"`
	AlbumId         string `vstore:"album_id"`
	TrackNumber     string `vstore:"track_number"`
	TrackTitle      string `vstore:"track_title"`
	Tags            []string `vstore:"tags"`
	DurationSeconds int64 `vstore:"duration_seconds"`
	WentGold        bool `vstore:"went_gold"`
	ReleaseDate		time.Time `vstore:"release_date"`
}

func (s *Song) Schema() *vstore.Schema {
	fields := vstore.NewPropertyBuilder().StringProperty(
		"artist_id",
		vstore.Required(),
		vstore.ElasticsearchProperty("elasticsearch", "analyzed", vstore.ElasticsearchField("raw", "not_analyzed")),
		vstore.ElasticsearchProperty("elasticsearch-v2", "analyzed"),
	).StringProperty(
		"album_id",
		vstore.Required(),
	).StringProperty(
		"track_number",
		vstore.Required(),
	).StringProperty(
		"track_title",
	).StringProperty(
		"tags",
		vstore.Repeated(),
	).IntegerProperty(
		"duration_seconds",
	).BooleanProperty(
		"went_gold",
	).TimeProperty(
		"release_date",
	).Build()
	secondaryIndexes := vstore.NewSecondaryIndexBuilder().
	Elasticsearch("elasticsearch", vstore.ElasticsearchNumberOfReplicas(1), vstore.ElasticsearchNumberOfShards(5), vstore.ElasticsearchRefreshInterval("1s")).
	Elasticsearch("elasticsearch-v2", vstore.ElasticsearchNumberOfReplicas(2), vstore.ElasticsearchNumberOfShards(5), vstore.ElasticsearchRefreshInterval("2s")).
	CloudSQL("cloud-sql", "104.154.165.235", "root", "password", "repcore-prod", "vstore", sqlClientKey, sqlClientCert, sqlCA).
	Build()
	backupConfig := vstore.NewBackupConfigBuilder().PeriodicBackup(vstorepb.BackupConfig_DAILY).Build()
	schema := vstore.NewSchema(calculatedNamespace, *kind, []string{"artist_id", "album_id", "track_number"}, fields, secondaryIndexes, backupConfig)
	return schema
}
```

This schema isn't trivial and actually utilizes most of VStore's major features, so let's examine one piece at a time.
```
type Song struct {
	ArtistId        string `vstore:"artist_id"`
	...
}
```
The critical thing to understand here is the field "tag": `vstore:"artist_id"`

If you're more of a read the docs guy: https://golang.org/pkg/reflect/#StructTag

Otherwise, what it means is that this field has some metadata that is important to the VStore client's internals. `artist_id` is the field name that VStore will store this field as, so when you call Elasticsearch or CloudSQL and look for this property, you'll want to look for `artist_id`, because that's how VStore is going to store it in Bigtable and replicate it to all the secondary indexes. Use snake case. You must provide this tag for the VStore client to serialize from or deserialize properties into this struct.

```
fields := vstore.NewPropertyBuilder().StringProperty(
		"artist_id",
		vstore.Required(),
		vstore.ElasticsearchProperty("elasticsearch", "analyzed", vstore.ElasticsearchField("raw", "not_analyzed")),
		vstore.ElasticsearchProperty("elasticsearch-v2", "analyzed"),
	)...
```
This is what a single property on your schema might look like. You're saying that one of your fields is a StringProperty named `"artist_id"`. Also note, `Song.ArtistId` is a `string`, and it is important you don't try to define a StringProperty as an `int64` on your struct. These much match, hopefully the mapping is obvious.

`vstore.Required()` is an option saying this field is required, but not necessarily a primary key. All properties that are primary keys must be required. Be careful with required fields that aren't part of the key, and don't try to add new required field(s) to existing schemas. Since `artist_id` is the first part of our 3 piece primary key, we make it required.

`vstore.ElasticsearchProperty("elasticsearch", "analyzed", vstore.ElasticsearchField("raw", "not_analyzed")),` is a configuration for a single property for a single secondary index. We are saying that we want the `artist_id` field in the index named `"elasticsearch"` to be an `"analyzed"` field, and that furthermore, with the `vstore.ElasticsearchField("raw", "not_analyzed")` option, we also want this field's `"raw"` property to be `"not_analyzed"`. Why would you want this? You can configure a field like this in order to be able to search on the analyzed version of the field by performing your query on the term `"artist_id"`, but you can also perform a different search on the unanalyzed version of the field by searching against the `"artist_id.raw"` term.

It's worth noting that VStore will always provide sensible defaults for property configurations on each kind of secondary index. Even if you don't specify an ElasticsearchProperty on this field, it will still end up in Elasticsearch with a default configuration.

`vstore.ElasticsearchProperty("elasticsearch-v2", "analyzed"),` is actually talking about a DIFFERENT elasticsearch index that you've chosen to call `"elasticsearch-v2"`, and so can have a completely different configuration for this property on that index.

```
secondaryIndexes := vstore.NewSecondaryIndexBuilder().
	Elasticsearch("elasticsearch", vstore.ElasticsearchNumberOfReplicas(1), vstore.ElasticsearchNumberOfShards(5), vstore.ElasticsearchRefreshInterval("1s")).
	Elasticsearch("elasticsearch-v2", vstore.ElasticsearchNumberOfReplicas(2), vstore.ElasticsearchNumberOfShards(5), vstore.ElasticsearchRefreshInterval("2s")).
	CloudSQL("cloud-sql", "104.154.165.235", "root", "password", "repcore-prod", "vstore", sqlClientKey, sqlClientCert, sqlCA).
	Build()
```
You can see that we've defined 3 different secondary indexes for this single kind, and that the identifiers ("elasticsearch-v2") match the property-level configuration identifiers.

Without getting into specifics, the major difference between `CloudSQL` and `Elasticsearch` indexes right now is that ALL elasticsearch indexes use the same instance (the same one we use for everything everywhere), whereas CloudSQL indices require you to provide the instance yourself, meaning you must supply the configuration with the instance's IP, login info, and credentials. All of this information is stored in VStore and all communication between VStore and CloudSQL happens over TLS.

```
backupConfig := vstore.NewBackupConfigBuilder().PeriodicBackup(vstorepb.BackupConfig_DAILY).Build()
```

The backup configuration right now is simple and should be straight forward. You can choose between daily, weekly or monthly backup policies to cloud storage, all of which expire after 90 days. If you plan on restoring a backup, please have a quick chat with someone listed at the bottom of this readme first and/or ensure you understand the consequences.

The last thing we need to do is put all the pieces together and define our primary keys (`"artist_id", "album_id", "track_number"`):
```
schema := vstore.NewSchema(calculatedNamespace, *kind, []string{"artist_id", "album_id", "track_number"}, fields, secondaryIndexes, backupConfig)
```


## History

See the [changelog](https://github.com/vendasta/gosdks/blob/master/vstore/CHANGELOG.md)

## Credits

These people know answers to questions:
* [Braden Bassingthwaite](https://github.com/bbassingthwaite-va)
* [Dustin Walker](https://github.com/dwalker-va)
