package vstore

import (
	"errors"
	"github.com/vendasta/gosdks/pb/vstorepb"
)

// NewSchema creates a new schema for a vStore kind.
// Ex:
//	vstore.NewSchema(
//		"vbc",
//		"AccountGroup",
//		vstore.NewFieldBuilder().StringField(
//			"account_group_id", Required(),
//		).StructField(
//			"accounts", vstore.NewFieldBuilder().StringField("account", Required()).Build(),
//		).Build(),
//	)
//
func NewSchema(namespace string, kind string, primaryKey []string, properties []*Property, secondaryIndexes []*SecondaryIndex, backupConfig *BackupConfig) *Schema {
	return &Schema{
		Namespace: namespace,
		Kind: kind,
		PrimaryKey: primaryKey,
		Properties: properties,
		SecondaryIndexes: secondaryIndexes,
		BackupConfig: backupConfig,
	}
}

//NewPropertyBuilder returns a struct that stores a schema's property definitions
func NewPropertyBuilder() (*propertyBuilder) {
	return &propertyBuilder{}
}

//NewSecondaryIndexBuilder returns a struct that stores a schema's secondary index definitions
func NewSecondaryIndexBuilder() (*secondaryIndexBuilder) {
	return &secondaryIndexBuilder{}
}

//NewBackupConfigBuilder returns a struct that stores a schema's backup configuration details
func NewBackupConfigBuilder() (*backupConfigBuilder) {
	return &backupConfigBuilder{}
}

type fieldOption func(*Property)

//Required makes a property required.
func Required() fieldOption {
	return func(p *Property) {
		p.IsRequired = true
	}
}

//Repeated makes a property repeated.
func Repeated() fieldOption {
	return func(p *Property) {
		p.IsRepeated = true
	}
}

type elasticFieldOption func(*vstorepb.SecondaryIndexPropertyConfig_Elasticsearch, *Property)

func getElasticType(vstoreType FieldType) string {
	switch vstoreType {
	case StringType:
		return "string"
	case IntType:
		return "integer"
	case FloatType:
		return "double"
	case BoolType:
		return "boolean"
	case TimeType:
		return "date"
	case GeoPointType:
		return "geo_point"
	case StructType:
		return "nested"
	}
	panic(ErrUnhandledTypeConversion)
}

// Used for storing an elasticsearch property more then once with different index types.
// Example: vstore.ElasticsearchField("raw", "string", "not_analyzed")
// https://www.elastic.co/guide/en/elasticsearch/reference/2.4/multi-fields.html
func ElasticsearchField(name string, indexType string) elasticFieldOption {
	return func(c *vstorepb.SecondaryIndexPropertyConfig_Elasticsearch, p *Property) {
		c.Fields = append(c.Fields, &vstorepb.SecondaryIndexPropertyConfig_ElasticsearchField{
			Name: name,
			Type: getElasticType(p.FType),
			Index: indexType,
		})
	}
}

//Repeated makes a property repeated.
func ElasticsearchProperty(indexName string, indexType string, fieldOptions ...elasticFieldOption) fieldOption {
	return func(p *Property) {
		if p.SecondaryIndexConfigs == nil {
			p.SecondaryIndexConfigs = map[string]*vstorepb.SecondaryIndexPropertyConfig{}
		}
		config := &vstorepb.SecondaryIndexPropertyConfig_Elasticsearch{
			Type: getElasticType(p.FType),
			Index: indexType,
		}
		for _, opt := range fieldOptions {
			opt(config, p)
		}
		p.SecondaryIndexConfigs[indexName] = &vstorepb.SecondaryIndexPropertyConfig{
			&vstorepb.SecondaryIndexPropertyConfig_ElasticsearchPropertyConfig{
				config,
			},
		}
	}
}

type cloudSQLFieldOption func(*vstorepb.SecondaryIndexPropertyConfig_CloudSQL, *Property)

func CloudSQLFieldType(typeOverride string) cloudSQLFieldOption {
	return func(c *vstorepb.SecondaryIndexPropertyConfig_CloudSQL, p *Property) {
		c.Type = typeOverride
	}
}

func CloudSQLProperty(indexName string, fieldOptions ...cloudSQLFieldOption) fieldOption {
	return func(p *Property) {
		if p.SecondaryIndexConfigs == nil {
			p.SecondaryIndexConfigs = map[string]*vstorepb.SecondaryIndexPropertyConfig{}
		}
		config := &vstorepb.SecondaryIndexPropertyConfig_CloudSQL{}
		for _, opt := range fieldOptions {
			opt(config, p)
		}
		p.SecondaryIndexConfigs[indexName] = &vstorepb.SecondaryIndexPropertyConfig{
			&vstorepb.SecondaryIndexPropertyConfig_CloudsqlPropertyConfig{
				config,
			},
		}
	}
}

//Schema defines how a VStore Kind is both stored and replicated
type Schema struct {
	Namespace        string
	Kind             string
	PrimaryKey       []string
	Properties       []*Property
	SecondaryIndexes []*SecondaryIndex
	BackupConfig     *BackupConfig
}

//BackupConfig defines the kind's backup strategy
type BackupConfig struct {
	BackupConfigPb *vstorepb.BackupConfig
}

type backupConfigBuilder struct {
	config *BackupConfig
}

//PeriodicBackup adds a backup policy that causes the kind to be backed up based on a defined period
func (b *backupConfigBuilder) PeriodicBackup(frequency vstorepb.BackupConfig_BackupFrequency) *backupConfigBuilder {
	b.config = &BackupConfig{
		&vstorepb.BackupConfig{
			Frequency: frequency,
		},
	}
	return b
}

//Build returns the backup config
func (b *backupConfigBuilder) Build() *BackupConfig {
	return b.config
}

//SecondaryIndex is a construct that defines a replication destination for a kind, including any configuration relevant for that index type
type SecondaryIndex struct {
	SecondaryIndexPb *vstorepb.SecondaryIndex
}

type ElasticsearchIndexOption func(*vstorepb.ElasticsearchConfig)

//ElasticsearchNumberOfShards controls how many shards the index is created with.
func ElasticsearchNumberOfShards(n int64) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		c.NumberOfShards = n
	}
}

//ElasticsearchNumberOfReplicas controls how many replicas the index is created with.
func ElasticsearchNumberOfReplicas(n int64) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		c.NumberOfReplicas = n
	}
}

//ElasticsearchRefreshInterval controls how often your index is refreshed
func ElasticsearchRefreshInterval(s string) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		c.RefreshInterval = s
	}
}

//ElasticsearchAnalyzer creates a new ElasticsearchAnalyzer
func ElasticsearchAnalyzer(name, analyzerType, tokenizer string, stemExclusion, stopWords, charFilter, filter []string) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		a := &vstorepb.ElasticsearchAnalyzer{
			Name: name,
			Type: analyzerType,
			Tokenizer: tokenizer,
			StemExclusion: stemExclusion,
			StopWords: stopWords,
			CharFilter: charFilter,
			Filter: filter,
		}
		if c.Analysis == nil {
			c.Analysis = &vstorepb.ElasticsearchAnalysis{}
		}
		c.Analysis.Analyzers = append(c.Analysis.Analyzers, a)
	}
}

//ElasticsearchFilter creates a new ElasticsearchFilter
func ElasticsearchFilter(name, analyzerType, pattern, replacement string, synonyms []string) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		f := &vstorepb.ElasticsearchFilter{
			Name: name,
			Type: analyzerType,
			Pattern: pattern,
			Replacement: replacement,
			Synonyms: synonyms,
		}
		if c.Analysis == nil {
			c.Analysis = &vstorepb.ElasticsearchAnalysis{}
		}
		c.Analysis.Filters = append(c.Analysis.Filters, f)
	}
}

//ElasticsearchCharFilter creates a new ElasticsearchCharFilter
func ElasticsearchCharFilter(name, analyzerType, pattern, replacement string) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		f := &vstorepb.ElasticsearchCharFilter{
			Name: name,
			Type: analyzerType,
			Pattern: pattern,
			Replacement: replacement,
		}
		if c.Analysis == nil {
			c.Analysis = &vstorepb.ElasticsearchAnalysis{}
		}
		c.Analysis.CharFilters = append(c.Analysis.CharFilters, f)
	}
}

//ElasticsearchTokenizer creates a new ElasticsearchTokenizer
func ElasticsearchTokenizer(name, analyzerType, pattern, delimiter string) ElasticsearchIndexOption {
	return func(c *vstorepb.ElasticsearchConfig) {
		f := &vstorepb.ElasticsearchTokenizer{
			Name: name,
			Type: analyzerType,
			Pattern: pattern,
			Delimiter: delimiter,
		}
		if c.Analysis == nil {
			c.Analysis = &vstorepb.ElasticsearchAnalysis{}
		}
		c.Analysis.Tokenizers = append(c.Analysis.Tokenizers, f)
	}
}

type secondaryIndexBuilder struct {
	indexes []*SecondaryIndex
}

//Elasticsearch adds an elasticsearch index to the definition with the specified settings
func (s *secondaryIndexBuilder) RawElasticsearch(name string, MappingJSON string, SettingsJSON string) *secondaryIndexBuilder {
	i := &vstorepb.SecondaryIndex{
		Name: name,
		Index: &vstorepb.SecondaryIndex_EsRawConfig{
			EsRawConfig: &vstorepb.ElasticsearchRawConfig{
				MappingJson: MappingJSON,
				SettingsJson: SettingsJSON,
			},
		},
	}
	s.indexes = append(s.indexes, &SecondaryIndex{i})
	return s
}

//Elasticsearch adds an elasticsearch index to the definition with the specified settings
func (s *secondaryIndexBuilder) Elasticsearch(name string, opts ...ElasticsearchIndexOption) *secondaryIndexBuilder {
	c := &vstorepb.ElasticsearchConfig{
		NumberOfShards: 5,
		NumberOfReplicas: 1,
		RefreshInterval: "1s",
	}
	for _, opt := range opts {
		opt(c)
	}

	i := &vstorepb.SecondaryIndex{
		Name: name,
		Index: &vstorepb.SecondaryIndex_EsConfig{
			EsConfig: c,
		},
	}
	s.indexes = append(s.indexes, &SecondaryIndex{i})
	return s
}

//CloudSQL adds a CloudSQL index to the definition with the specified settings
func (s *secondaryIndexBuilder) CloudSQL(name, instanceIP, userName, password, projectID, instanceName string, clientKey, clientCert, serverCertificateAuthority []byte) *secondaryIndexBuilder {
	i := &vstorepb.SecondaryIndex{
		Name: name,
		Index: &vstorepb.SecondaryIndex_CloudSqlConfig{
			CloudSqlConfig: &vstorepb.CloudSQLConfig{
				InstanceIp: instanceIP,
				UserName: userName,
				Password: password,
				ClientKey: clientKey,
				ClientCert: clientCert,
				ServerCertificateAuthority: serverCertificateAuthority,
				ProjectId: projectID,
				InstanceName: instanceName,
			},
		},
	}
	s.indexes = append(s.indexes, &SecondaryIndex{i})
	return s
}

//Build returns the secondary index definitions
func (s *secondaryIndexBuilder) Build() []*SecondaryIndex {
	return s.indexes
}

type propertyBuilder struct {
	properties []*Property
}

//StringProperty adds a string property to the list of property definitions
func (s *propertyBuilder) StringProperty(name string, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: StringType,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}
	return s
}

//IntegerProperty adds an integer property to the list of property definitions
func (s *propertyBuilder) IntegerProperty(name string, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: IntType,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}

	return s
}

//IntegerProperty adds an integer property to the list of property definitions
func (s *propertyBuilder) FloatProperty(name string, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: FloatType,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}

	return s
}

//BooleanProperty adds a boolean property to the list of property definitions
func (s *propertyBuilder) BooleanProperty(name string, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: BoolType,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}

	return s
}

//TimeProperty adds a time property to the list of property definitions
func (s *propertyBuilder) TimeProperty(name string, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: TimeType,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}

	return s
}

//GeoPointProperty adds a geo point property to the list of property definitions
func (s *propertyBuilder) GeoPointProperty(name string, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: GeoPointType,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}

	return s
}

//StructProperty adds a struct property to the list of property definitions
func (s *propertyBuilder) StructProperty(name string, properties []*Property, opts ...fieldOption) *propertyBuilder {
	f := &Property{
		Name: name,
		FType: StructType,
		Properties: properties,
	}
	s.properties = append(s.properties, f)
	for _, opt := range opts {
		opt(f)
	}

	return s
}

//Build returns the property definitions
func (s *propertyBuilder) Build() []*Property {
	return s.properties
}

//Property is the definition of a single field/column and how it is stored, validated and structured.
type Property struct {
	Name                  string
	FType                 FieldType
	IsRequired            bool
	IsRepeated            bool
	Properties            []*Property
	SecondaryIndexConfigs map[string]*vstorepb.SecondaryIndexPropertyConfig
}

//ToPb return the protobuf representation of a Property
func (p *Property) ToPb() (property *vstorepb.Property, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error); if !ok || err != ErrUnhandledTypeConversion {
				panic(r)
			}
		}
	}()
	var properties []*vstorepb.Property
	if len(p.Properties) > 0 {
		properties, err = PropertiesToPb(p.Properties...); if err != nil {
			return
		}
	}
	property = &vstorepb.Property{
		Name: p.Name,
		Required: p.IsRequired,
		Repeated: p.IsRepeated,
		Type: p.FType.ToPb(),
		Properties: properties,
		SecondaryIndexConfigs: p.SecondaryIndexConfigs,
	}
	return
}

//PropertiesToPb returns a list of protobuf serialized properties
func PropertiesToPb(properties ...*Property) ([]*vstorepb.Property, error) {
	var err error
	propertiesPb := make([]*vstorepb.Property, len(properties))
	for i, p := range properties {
		propertiesPb[i], err = p.ToPb(); if err != nil {
			return nil, err
		}
	}
	return propertiesPb, nil
}

type FieldType int64

//The complete spectrum of supported field types
const (
	StringType FieldType = iota
	IntType
	FloatType
	BoolType
	TimeType
	GeoPointType
	StructType
)

//ErrUnhandledTypeConversion is thrown when an unsupported field type is detected
var ErrUnhandledTypeConversion = errors.New("Unhandled type conversion to pb.")

func (f FieldType) ToPb() vstorepb.Property_Type {
	switch f {
	case StringType:
		return vstorepb.Property_STRING
	case IntType:
		return vstorepb.Property_INT64
	case FloatType:
		return vstorepb.Property_DOUBLE
	case BoolType:
		return vstorepb.Property_BOOL
	case TimeType:
		return vstorepb.Property_TIMESTAMP
	case GeoPointType:
		return vstorepb.Property_GEOPOINT
	case StructType:
		return vstorepb.Property_STRUCT
	}
	panic(ErrUnhandledTypeConversion)
}

//SecondaryIndexesToPb transforms a list of secondary indexes into protobuf format
func SecondaryIndexesToPb(secondaryIndexes ...*SecondaryIndex) []*vstorepb.SecondaryIndex {
	secondaryIndexesPb := make([]*vstorepb.SecondaryIndex, len(secondaryIndexes))
	for i, secondaryIndex := range secondaryIndexes {
		secondaryIndexesPb[i] = secondaryIndex.SecondaryIndexPb
	}
	return secondaryIndexesPb
}

//BackupConfigToPb transforms a backup configuration into protobuf format
func BackupConfigToPb(backupConfig *BackupConfig) *vstorepb.BackupConfig {
	if backupConfig == nil {
		return nil
	}
	return backupConfig.BackupConfigPb
}
