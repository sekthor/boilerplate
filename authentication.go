package boilerplate

import (
	"context"
	"errors"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryJwtClaimsInterceptor[T jwt.Claims](jwksUrls []string, claimsFunc func() T) (grpc.UnaryServerInterceptor, error) {

	Keyfunc, err := keyfunc.NewDefault(jwksUrls)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		header := md["authorization"]

		if len(header) < 1 {
			return handler(ctx, req)
			//return nil, errMissingBearerToken
		}

		claims := claimsFunc()

		signed := strings.TrimPrefix(header[0], "Bearer ")
		token, err := jwt.ParseWithClaims(signed, claims, Keyfunc.Keyfunc)
		if err != nil || !token.Valid {
			return nil, errInvalidToken
		}

		span := trace.SpanFromContext(ctx)

		if sub, err := claims.GetSubject(); err != nil {
			span.SetAttributes(attribute.String("user.id", sub))
		}
		ctx = context.WithValue(ctx, "claims", claims)

		return handler(ctx, req)
	}, nil
}

func GetClaimsFromContext[T jwt.Claims](ctx context.Context) (claims T, err error) {
	var ok bool
	claims, ok = ctx.Value("claims").(T)
	if !ok {
		err = errors.New("no user in context")
	}
	return
}
