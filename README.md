# Boilerplate

🚨 **This is still a very experimental project** API-breaking changes at will 🚨

This project holds a Boilerplate server for 

- a grpc protobuf service
- with grpc-gateway
- opentelemetry

## Feature Status

- gPRC Server
    - ✅ insecure
    - ⭕ TLS (currently hard coded server name)
    - ⭕ mTLS (currently only supports 1 CA)
- gRPC Gateway
    - ✅ insecure
- JWT Authentication
    - ✅ multiple issuers (supply *n* jwks endpoints used to check jwt signatures)
    - ✅ access token claims from request context
    - ❌ API for accessing claims (e.g. `ClaimsFromContext(context.Context) (Claims, error)`)
- Opentelemetry
    - ✅ Tracing Exporter
    - ✅ Metrics Exporter
    - ✅ Logger Exporter (with logrus bridge)
    - ❌ customizable default resources/attributes

## Usage

```
go get -u github.com/sekthor/boilerplate
```

1. Create your `.proto` files and generate the code using [buf]() cli.
1. Implement your protobuf service.
    ```go
    type ServiceImplementation struct {
        greeterv1.UnimplementedGreeterServiceServer
    }

    func (i *ServiceImplementation) SayHello(ctx context.Context, req *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
        name := req.GetName()

        return &greeterv1.SayHelloResponse{
            Message: fmt.Sprintf("Hello %s!", name),
        }, nil
    }
    ```
1. create a `boilerplate.GrpcRegisterFunc` and a `boilerplate.GatewayRegisterFunc` that wrap the generated register functions. These wrappers are called inbetween the creation of the servers and the start.
    ```go
    var grpcFunc = func(s *grpc.Server) error {
        greeterv1.RegisterGreeterServiceServer(s, i)
        return nil
    }

    var gatewayFunc = func(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error {
        return greeterv1.RegisterGreeterServiceHandler(ctx, mux, cc)
    }
    ```
1. add the `boilerplate.BoilerplateServer` interface in your service implementation and create the instance
    ```go
    type ServiceImplementation struct {
        greeterv1.UnimplementedGreeterServiceServer
        server boilerplate.BoilerplateServer
    }
    ```
    ```go
    server := boilerplate.Default()

    service := ServiceImplementation{}
    service.server = server
    ```
1. Register the service
    ```go
    server.RegisterGrpc(grpcFunc)
    server.RegisterGateway(gatewayFunc)
    ```
1. Run the server
    ```go
    if err := server.Run(ctx); err != nil {
        log.Fatal(err)
    }
    ```
    
