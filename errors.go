package boilerplate

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata    = status.Errorf(codes.InvalidArgument, "missing metadata")
	errMissingBearerToken = status.Errorf(codes.InvalidArgument, "missing bearer token")
	errInvalidToken       = status.Errorf(codes.Unauthenticated, "invalid token")
)
