package vstore

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"github.com/vendasta/gosdks/pb/vstorepb"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"errors"
)

func TestVStoreTestSuite(t *testing.T) {
	suite.Run(t, new(VStoreTestSuite))
}

type VStoreTestSuite struct {
	suite.Suite
	clientMock *ClientMock
	adminClientMock *AdminClientMock
	v          Interface
}

type TestModel struct {
}

func (t *TestModel) Schema() *Schema {
	return NewSchema("vstore", "Song", []string{"song_id"}, nil, nil, nil)
}

var testEntity *vstorepb.Entity = &vstorepb.Entity{
	Namespace: "vstore",
	Kind: "Song",
	Version: 1,
	Values: &vstorepb.Struct{Values:map[string]*vstorepb.Value{"name": {&vstorepb.Value_StringValue{"Hello"}}}},
}

func (suite *VStoreTestSuite) SetupTest() {
	RegisterModel("vstore", "Song", (*Song)(nil))
	suite.clientMock = &ClientMock{}
	suite.adminClientMock = &AdminClientMock{}
	suite.v, _ = New(Client(suite.clientMock), AdminClient(suite.adminClientMock), Environment(Local))
}

// GetMulti Tests
func (suite *VStoreTestSuite) Test_GetMultiReturnsErrorOnClientError() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	suite.clientMock.On("Get", context.Background(), &vstorepb.GetRequest{KeySets:[]*vstorepb.KeySet{}}, mock.Anything).Return(nil, err)

	_, gotError := suite.v.GetMulti(context.Background(), nil)
	suite.Assert().Equal(err, gotError)
}

func (suite *VStoreTestSuite) Test_GetMultiReturnsModelOnSuccess() {
	r := vstorepb.GetResponse{
		Entities:[]*vstorepb.EntityResult{
			{
				Entity: testEntity,
			},
		},
	}
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(&r, nil)
	entities, err := suite.v.GetMulti(context.Background(), nil)
	suite.Assert().Nil(err)
	suite.Assert().Equal(1, len(entities))
	song := entities[0].(*Song)
	suite.Assert().Equal("Hello", song.Name)
}

func (suite *VStoreTestSuite) Test_GetMultiCallsWithExpectedRequest() {
	r := vstorepb.GetResponse{
		Entities:[]*vstorepb.EntityResult{
			{
				Entity: testEntity,
			},
		},
	}

	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(&r, nil)

	ks := []*KeySet{NewKeySet("vstore", "Song", []string{"k"})}
	suite.v.GetMulti(context.Background(), ks)

	expectedRequest := &vstorepb.GetRequest{[]*vstorepb.KeySet{{Namespace: "vstore", Kind: "Song", Keys: []string{"k"}}}}
	suite.clientMock.AssertCalled(suite.T(), "Get", context.Background(), expectedRequest, mock.Anything)
}

// Get Tests
func (suite *VStoreTestSuite) Test_GetReturnsModelOnSuccess() {
	r := vstorepb.GetResponse{
		Entities:[]*vstorepb.EntityResult{
			{
				Entity: testEntity,
			},
		},
	}
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(&r, nil)
	entity, err := suite.v.Get(context.Background(), NewKeySet("", "", []string{}))
	suite.Assert().Nil(err)
	song := entity.(*Song)
	suite.Assert().Equal("Hello", song.Name)
}

// Lookup Tests
func (suite *VStoreTestSuite) Test_LookupReturnsErrorOnClientError() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	suite.clientMock.On("Lookup", context.Background(), &vstorepb.LookupRequest{Namespace: "vstore", Kind: "Song", PageSize:10}, mock.Anything).Return(nil, err)

	_, gotError := suite.v.Lookup(context.Background(), "vstore", "Song")
	suite.Assert().Equal(err, gotError)
}

func (suite *VStoreTestSuite) Test_LookupReturnsExpectedResult() {
	r := vstorepb.LookupResponse{
		Entities:[]*vstorepb.EntityResult{
			{
				Entity: testEntity,
			},
		},
		NextCursor: "abc",
		HasMore: true,
	}
	suite.clientMock.On("Lookup", context.Background(), &vstorepb.LookupRequest{Namespace: "vstore", Kind: "Song", PageSize:10}, mock.Anything).Return(&r, nil)

	resp, err := suite.v.Lookup(context.Background(), "vstore", "Song")
	suite.Assert().Nil(err)
	suite.Assert().Equal(1, len(resp.Results))
	suite.Assert().Equal("abc", resp.NextCursor)
	suite.Assert().Equal(true, resp.HasMore)
	song := resp.Results[0].(*Song)
	suite.Assert().Equal("Hello", song.Name)
}

func (suite *VStoreTestSuite) Test_CursorIsPassedInRequest() {
	r := vstorepb.LookupResponse{}
	suite.clientMock.On("Lookup", context.Background(), mock.Anything, mock.Anything).Return(&r, nil)

	suite.v.Lookup(context.Background(), "vstore", "Song", Cursor("cursor"))

	expectedRequest := &vstorepb.LookupRequest{
		Namespace: "vstore",
		Kind: "Song",
		PageSize: 10,
		Cursor: "cursor",
	}

	suite.clientMock.AssertCalled(suite.T(), "Lookup", context.Background(), expectedRequest, mock.Anything)
}

func (suite *VStoreTestSuite) Test_FilterIsPassedInRequest() {
	r := vstorepb.LookupResponse{}
	suite.clientMock.On("Lookup", context.Background(), mock.Anything, mock.Anything).Return(&r, nil)

	suite.v.Lookup(context.Background(), "vstore", "Song", Filter([]string{"filter"}))

	expectedRequest := &vstorepb.LookupRequest{
		Namespace: "vstore",
		Kind: "Song",
		PageSize: 10,
		Filter: &vstorepb.LookupFilter{Keys: []string{"filter"}},
	}

	suite.clientMock.AssertCalled(suite.T(), "Lookup", context.Background(), expectedRequest, mock.Anything)
}

func (suite *VStoreTestSuite) Test_PageSizeIsPassedInRequest() {
	r := vstorepb.LookupResponse{}
	suite.clientMock.On("Lookup", context.Background(), mock.Anything, mock.Anything).Return(&r, nil)

	suite.v.Lookup(context.Background(), "vstore", "Song", PageSize(1))

	expectedRequest := &vstorepb.LookupRequest{
		Namespace: "vstore",
		Kind: "Song",
		PageSize: 1,
	}

	suite.clientMock.AssertCalled(suite.T(), "Lookup", context.Background(), expectedRequest, mock.Anything)
}

// Transaction Tests
func (suite *VStoreTestSuite) TestTransaction_ErrorOnGetReturnsError() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(nil, err)

	ks := NewKeySet("vstore", "Song", []string{"k"})

	gotError := suite.v.Transaction(context.Background(), ks, func(Transaction, Model) error {
		return nil
	})
	suite.Assert().Equal(err, gotError)
}

func (suite *VStoreTestSuite) TestTransaction_CreatesEntityWhenGetReturnsEmptyResult() {
	toSave, _ := StructPBToModel(testEntity.Namespace, testEntity.Kind, testEntity.Values)
	getResp := &vstorepb.GetResponse{
		Entities: []*vstorepb.EntityResult{
			&vstorepb.EntityResult{},
		},
	}
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(getResp, nil)
	suite.clientMock.On("Create", context.Background(), mock.Anything, mock.Anything).Return(&google_protobuf.Empty{}, nil)

	ks := NewKeySet("vstore", "Song", []string{"k"})

	_ = suite.v.Transaction(context.Background(), ks, func(tx Transaction, m Model) error {
		tx.Save(toSave)
		return nil
	})
	suite.clientMock.AssertCalled(suite.T(), "Create", context.Background(), mock.AnythingOfType("*vstorepb.CreateRequest"), mock.AnythingOfType("[]grpc.CallOption"))
}

func (suite *VStoreTestSuite) TestTransaction_ReturnsErrorWhenCreateFails() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	toSave, _ := StructPBToModel(testEntity.Namespace, testEntity.Kind, testEntity.Values)
	getResp := &vstorepb.GetResponse{
		Entities: []*vstorepb.EntityResult{
			&vstorepb.EntityResult{},
		},
	}
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(getResp, nil)
	suite.clientMock.On("Create", context.Background(), mock.Anything, mock.Anything).Return(nil, err)

	ks := NewKeySet("vstore", "Song", []string{"k"})

	gotError := suite.v.Transaction(context.Background(), ks, func(tx Transaction, m Model) error {
		tx.Save(toSave)
		return nil
	})
	suite.Assert().Equal(err, gotError)
}

func (suite *VStoreTestSuite) TestTransaction_UpdatesEntityWhenGetReturnsNonEmptyResult() {
	toSave, _ := StructPBToModel(testEntity.Namespace, testEntity.Kind, testEntity.Values)
	getResp := &vstorepb.GetResponse{
		Entities: []*vstorepb.EntityResult{
			&vstorepb.EntityResult{
				Entity: testEntity,
			},
		},
	}
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(getResp, nil)
	suite.clientMock.On("Update", context.Background(), mock.Anything, mock.Anything).Return(&google_protobuf.Empty{}, nil)

	ks := NewKeySet("vstore", "Song", []string{"k"})

	_ = suite.v.Transaction(context.Background(), ks, func(tx Transaction, m Model) error {
		tx.Save(toSave)
		return nil
	})
	suite.clientMock.AssertCalled(suite.T(), "Update", context.Background(), mock.AnythingOfType("*vstorepb.UpdateRequest"), mock.AnythingOfType("[]grpc.CallOption"))
}

func (suite *VStoreTestSuite) TestTransaction_ReturnsErrorWhenUpdateFails() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	toSave, _ := StructPBToModel(testEntity.Namespace, testEntity.Kind, testEntity.Values)
	getResp := &vstorepb.GetResponse{
		Entities: []*vstorepb.EntityResult{
			&vstorepb.EntityResult{
				Entity: testEntity,
			},
		},
	}
	suite.clientMock.On("Get", context.Background(), mock.Anything, mock.Anything).Return(getResp, nil)
	suite.clientMock.On("Update", context.Background(), mock.Anything, mock.Anything).Return(nil, err)

	ks := NewKeySet("vstore", "Song", []string{"k"})

	gotError := suite.v.Transaction(context.Background(), ks, func(tx Transaction, m Model) error {
		tx.Save(toSave)
		return nil
	})
	suite.Assert().Equal(err, gotError)
}

// Register tests
func (suite *VStoreTestSuite) TestRegister_DoesNotCreateAfterSuccessfulUpdate() {
	res := &vstorepb.GetKindResponse{PrimaryKey: []string{"Test"}}
	suite.adminClientMock.On("UpdateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, nil)
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	kindRes, err := suite.v.RegisterKind(context.Background(), "vstore", "Song", []string{"datapond@vendasta.com"}, (*TestModel)(nil))

	suite.Assert().Equal(kindRes.PrimaryKey, []string{"Test"})
	suite.Assert().Nil(err)
	suite.adminClientMock.AssertNotCalled(suite.T(), "CreateKind")
}

func (suite *VStoreTestSuite) TestRegister_CreatesNamespaceAndKindAfterFailedUpdate() {
	res := &vstorepb.GetKindResponse{PrimaryKey: []string{"Test"}}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	suite.adminClientMock.On("UpdateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, grpc.Errorf(codes.NotFound, "Not Found"))
	suite.adminClientMock.On("CreateNamespace", context.Background(), mock.Anything, mock.Anything).Return(nil, nil)
	suite.adminClientMock.On("CreateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, nil)

	kindRes, err := suite.v.RegisterKind(context.Background(), "vstore", "Song", []string{"datapond@vendasta.com"}, (*TestModel)(nil))

	suite.Assert().Equal(kindRes.PrimaryKey, []string{"Test"})
	suite.Assert().Nil(err)
}

func (suite *VStoreTestSuite) TestRegister_RaisesErrorIfUpdateFailsWithNonNotFoundCode() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	suite.adminClientMock.On("UpdateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, err)

	kindRes, gotErr := suite.v.RegisterKind(context.Background(), "vstore", "Song", []string{"datapond@vendasta.com"}, (*TestModel)(nil))

	suite.Assert().Equal(err, gotErr)
	suite.Assert().Nil(kindRes)
}

func (suite *VStoreTestSuite) TestRegister_RaisesErrorIfCreateNamespaceFailsWithNoneAlreadyExistsError() {
	err := grpc.Errorf(codes.Unavailable, "Broked")
	suite.adminClientMock.On("UpdateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, grpc.Errorf(codes.NotFound, "Not Found"))
	suite.adminClientMock.On("CreateNamespace", context.Background(), mock.Anything, mock.Anything).Return(nil, err)

	kindRes, gotErr := suite.v.RegisterKind(context.Background(), "vstore", "Song", []string{"datapond@vendasta.com"}, (*TestModel)(nil))

	suite.Assert().Nil(kindRes)
	suite.Assert().Equal(err, gotErr)
	suite.adminClientMock.AssertNotCalled(suite.T(), "CreateKind")
}

func (suite *VStoreTestSuite) TestRegister_DoesNotRaiseErrorIfCreateNamespaceFailsWithAlreadyExistsError() {
	res := &vstorepb.GetKindResponse{PrimaryKey: []string{"Test"}}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	suite.adminClientMock.On("UpdateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, grpc.Errorf(codes.NotFound, "Not Found"))
	suite.adminClientMock.On("CreateNamespace", context.Background(), mock.Anything, mock.Anything).Return(nil, grpc.Errorf(codes.AlreadyExists, "Already Exists"))
	suite.adminClientMock.On("CreateKind", context.Background(), mock.Anything, mock.Anything).Return(nil, nil)

	kindRes, err := suite.v.RegisterKind(context.Background(), "vstore", "Song", []string{"datapond@vendasta.com"}, (*TestModel)(nil))

	suite.Assert().Equal(kindRes.PrimaryKey, []string{"Test"})
	suite.Assert().Nil(err)
}

// GetSecondaryIndexName tests
func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsErrorWhenGetKindFails() {
	expectedErr := errors.New("GetKind blew the f up")
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(nil, expectedErr)

	n, actualErr := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "elasticsearch1")
	suite.Assert().Zero(n)
	suite.Assert().Equal(expectedErr, actualErr)
}

func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsErrorWhenSchemaHasNoSecondaryIndexes() {
	res := &vstorepb.GetKindResponse{PrimaryKey: []string{"Test"}, SecondaryIndexes: []*vstorepb.SecondaryIndex{}}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	n, err := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "elasticsearch1")

	suite.Assert().Zero(n)
	suite.Assert().EqualError(err, "Could not find any secondary indexes on the schema.")
}

func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsErrorWhenSpecificIndexIsMissingFromSchema() {
	res := &vstorepb.GetKindResponse{
		PrimaryKey: []string{"Test"},
		SecondaryIndexes: []*vstorepb.SecondaryIndex{
			&vstorepb.SecondaryIndex{Name: "my-sql-index"},
		},
	}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	n, err := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "elasticsearch1")

	suite.Assert().Zero(n)
	suite.Assert().EqualError(err, "Could not find the specified secondary index on the schema.")
}

func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsCloudSQLIndexName() {
	res := &vstorepb.GetKindResponse{
		PrimaryKey: []string{"Test"},
		SecondaryIndexes: []*vstorepb.SecondaryIndex{
			&vstorepb.SecondaryIndex{
				Name: "cloudsql",
				Index: &vstorepb.SecondaryIndex_CloudSqlConfig{
					CloudSqlConfig: &vstorepb.CloudSQLConfig{
						IndexName: "vstor-bbass-gibberish",
					},
				},
			},
		},
	}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	n, err := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "cloudsql")

	suite.Assert().Equal("vstor-bbass-gibberish", n)
	suite.Assert().Nil(err)
}

func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsElasticsearchIndexName() {
	res := &vstorepb.GetKindResponse{
		PrimaryKey: []string{"Test"},
		SecondaryIndexes: []*vstorepb.SecondaryIndex{
			&vstorepb.SecondaryIndex{
				Name: "elasticsearch1",
				Index: &vstorepb.SecondaryIndex_EsConfig{
					EsConfig: &vstorepb.ElasticsearchConfig{
						IndexName: "vstor-bbass-gibberish",
					},
				},
			},
		},
	}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	n, err := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "elasticsearch1")

	suite.Assert().Equal("vstor-bbass-gibberish", n)
	suite.Assert().Nil(err)
}

func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsElasticsearchRawConfigIndexName() {
	res := &vstorepb.GetKindResponse{
		PrimaryKey: []string{"Test"},
		SecondaryIndexes: []*vstorepb.SecondaryIndex{
			&vstorepb.SecondaryIndex{
				Name: "elasticsearch1",
				Index: &vstorepb.SecondaryIndex_EsRawConfig{
					EsRawConfig: &vstorepb.ElasticsearchRawConfig{
						IndexName: "vstor-bbass-gibberish",
					},
				},
			},
		},
	}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	n, err := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "elasticsearch1")

	suite.Assert().Equal("vstor-bbass-gibberish", n)
	suite.Assert().Nil(err)
}

func (suite *VStoreTestSuite) TestGetSecondaryIndexName_ReturnsCorrectNameFromMultipleIndexConfigs() {
	res := &vstorepb.GetKindResponse{
		PrimaryKey: []string{"Test"},
		SecondaryIndexes: []*vstorepb.SecondaryIndex{
			&vstorepb.SecondaryIndex{
				Name: "es-2",
				Index: &vstorepb.SecondaryIndex_EsRawConfig{
					EsRawConfig: &vstorepb.ElasticsearchRawConfig{
						IndexName: "vstore-this-is-not-what-you-want",
					},
				},
			},
			&vstorepb.SecondaryIndex{
				Name: "elasticsearch1",
				Index: &vstorepb.SecondaryIndex_EsRawConfig{
					EsRawConfig: &vstorepb.ElasticsearchRawConfig{
						IndexName: "vstor-bbass-gibberish",
					},
				},
			},
		},
	}
	suite.adminClientMock.On("GetKind", context.Background(), mock.Anything, mock.Anything).Return(res, nil)

	n, err := suite.v.GetSecondaryIndexName(context.Background(), "vstore", "Song", "elasticsearch1")

	suite.Assert().Equal("vstor-bbass-gibberish", n)
	suite.Assert().Nil(err)
}
