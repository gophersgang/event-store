# Vendasta API Interfaces
This repo will be the source of truth for our proto files. Your work on a new microservice or API will usually start here. PR's into this repo are a great opportunity to design your API interface ahead of time and have people vet your syntax and parameter naming, without being watered down in unrelated implementation details.

## Directory Structure
Do yourself a huge favour and start with a versioned folder structure to begin with. This saves us from having to store commit hashes or tags when we check out a proto for compilation, but rather we can just point to the folder. Of course semantic versioning applies, no breaking changes without creating a new version.

Your folder tree will then look as follows:
```
- Listings
- CHANGELOG.md
-- V1
--- Listings
---- listings.proto
-- V2
--- Listings
---- listings.proto
```
As you can see, the listings subfolder is duplicated under the version name. This folder structure is optomized for the code generation... in certain languages, what your folder is named is going to make a difference on your generated code (for example, the name of the module). You don't want your module to be named `V1` instead of `Listings`.

Do not use hyphens in your path. Protoc python grpc generation can not handle hyphenated filenames or filepaths at this time: https://github.com/grpc/grpc/issues/5226

## Proto Version/Package Name
ALWAYS use proto V3. 

Your package name should be in the format \<microservice\>proto, if it is an internal only proto. For external protos use vendasta.\<microservice\>proto

```
syntax = "proto3";

package datalakeproto;
```

## Proto Syntax Standards

### Naming Conventions

#### Services
Use the microservice's name, since the microservice's name will directly reflect the domain that it is encapsulating. DO NOT postfix the service with "Service", this will result in stutter in auto-generated code.
```
service Datalake {...}
```

When you have other services that represent different levels of authentication, use a postfix that describes the user.
```
service DatalakeAdmin {...}
service DatalakeInternal {...}
```
#### RPC's
Since the microservice should already be striving to work in a single domain, there should not be a large number of domain-specific services defined from a single microservice. However, there are two other reasons why you might have multiple services in a single microservice: 1) services requiring different levels of authentication (datalake, datalake admin), or 2) when it makes sense to split up your connection pooling (TODO: provide more context as to when this is needed).

There is a limited number of verbs to use when considering an RPC name. Your RPC will fall under one of the following categories:
- Get
- List
- Search
- Create
- Update
- Delete
- Intent or Specific Action
- Replace

Since we are favouring less services based on the microservice name as opposed to more services based on the microservice's models, we will need to include the model name in the RPC name.

So when you combine the verb with your model name, your service will end up looking like this:
```
service Datalake {
  rpc GetListing (...) returns (...) {};
  rpc ListListings (...) returns (...) {}; // Returned in some defined order, ex: created, alphabetical
  rpc SearchListings (...) returns (...) {}; // Returned in scored order, the better the match to the search, the higher
  rpc CreateListing (...) returns (...) {};
  rpc UpdateListing (...) returns (...) {};
  rpc DeleteListing (...) returns (...) {};
  rpc RescrapeListing (...) returns (...) {}; // Intent/specific action
  rpc ReplaceListing (...) returns (...) {}; // Replaces a listing if it exists, otherwise creates it
}
```
#### Messages
There are two types of messages:
- Messages defining generic, reuseable types/objects
- Messages defining RPC request/response schemas

In the first case of a generic type or object, the model name alone (without no prefix or postfix) should be used:
```
message Listing {...}
message ListingStats {...}
```

In the second case, when you're defining messages for the purpose of being a request or response, you WILL use both a prefix and a postfix. Prefix with the rpc verb from the RPC section, postfix with `Request` or `Response`.
```
message GetListingRequest {
  string listing_id = 1;
}

message Listing {...}

message ListListingsResponse {
  repeated Listing listings = 1;
  TODO: page size, offset, total listings, etc (how do we do default page size and offset?)
}

service Datalake {
  rpc GetListing (GetListingRequest) returns (GetListingResponse) {};
  rpc ListListings (ListListingsRequest) returns (ListListingsResponse) {};
  rpc SearchListings (SearchListingsRequest) returns (SearchListingsResponse) {};
  rpc CreateListing (CreateListingRequest) returns (CreateListingResponse) {};
  rpc UpdateListing (UpdateListingRequest) returns (UpdateListingResponse) {};
  rpc DeleteListing (DeleteListingRequest) returns (DeleteListingResponse) {};
  rpc RescrapeListing (RescrapeListingRequest) returns (RescrapeListingResponse) {};
  rpc ReplaceListing (ReplaceListingRequest) returns (ReplaceListingResponse) {};
}
```

