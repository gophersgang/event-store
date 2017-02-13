package utils

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ErrorType is an enum encapsulating the spectrum of all possible types of errors raised by the application
type ErrorType int64

const (
	// NotFound corresponds to errors caused by missing entities
	NotFound ErrorType = 1 + iota
	// InvalidArgument corresponds to errors caused by missing or malformed arguments supplied by a client
	InvalidArgument
	// AlreadyExists corresponds to errors caused by an entity already existing
	AlreadyExists
	// PermissionDenied corresponds to a user not having permission to access a resource.
	PermissionDenied
	// Unauthenticated indicates the request does not have valid authentication credentials for the operation.
	Unauthenticated
	// Unimplemented corresponds to a function that is unimplemented
	Unimplemented
	// Unknown Error occurred
	Unknown
)

// EventStoreError is an error that can be translated to a GRPC-compliant error
type EventStoreError struct {
	msg     string
	errType ErrorType
}

// Error returns the message associated with this error
func (v EventStoreError) Error() string {
	return v.msg
}

// ErrorType returns the ErrorType associated with this error
func (v EventStoreError) ErrorType() ErrorType {
	return v.errType
}

// GRPCError returns an error in a format such that it can be consumed by GRPC
func (v EventStoreError) GRPCError() error {
	if v.errType == NotFound {
		return grpc.Errorf(codes.NotFound, v.msg)
	} else if v.errType == InvalidArgument {
		return grpc.Errorf(codes.InvalidArgument, v.msg)
	} else if v.errType == AlreadyExists {
		return grpc.Errorf(codes.AlreadyExists, v.msg)
	} else if v.errType == PermissionDenied {
		return grpc.Errorf(codes.PermissionDenied, v.msg)
	} else if v.errType == Unauthenticated {
		return grpc.Errorf(codes.Unauthenticated, v.msg)
	} else if v.errType == Unimplemented {
		return grpc.Errorf(codes.Unimplemented, v.msg)
	}
	return grpc.Errorf(codes.Unknown, "Unknown server error.")
}

// Error returns a SalesOpportunitiesError
func Error(errorType ErrorType, format string, a ...interface{}) error {
	return EventStoreError{msg: fmt.Sprintf(format, a...), errType: errorType}
}

// FromError given an error tries to return a proper EventStoreError.
func FromError(err error) EventStoreError {
	eventStoreError, ok := err.(EventStoreError)
	if !ok {
		return Error(Unknown, err.Error()).(EventStoreError)
	}
	return eventStoreError
}

// IsError returns true/false if the given err matches the errorType type.
func IsError(errorType ErrorType, err error) bool {
	eventStoreError, ok := err.(EventStoreError)
	if !ok {
		return false
	}
	return eventStoreError.errType == errorType
}
