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

type Claims struct {
	jwt.RegisteredClaims
}

func UnaryJwtInterceptor(jwksUrls []string) (grpc.UnaryServerInterceptor, error) {

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

		signed := strings.TrimPrefix(header[0], "Bearer ")
		claims := Claims{}
		token, err := jwt.ParseWithClaims(signed, &claims, Keyfunc.Keyfunc)
		if err != nil || !token.Valid {
			return nil, errInvalidToken
		}

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("user.id", claims.Subject))
		ctx = context.WithValue(ctx, "claims", claims)

		return handler(ctx, req)
	}, nil
}

func GetClaimsFromContext(ctx context.Context) (claims Claims, err error) {
	var ok bool
	claims, ok = ctx.Value("claims").(Claims)
	if !ok {
		err = errors.New("no user in context")
	}
	return
}

func GetSubjectFromContext(ctx context.Context) (string, error) {
	claims, err := GetClaimsFromContext(ctx)
	return claims.Subject, err
}
