# Boilerplate

This project holds a Boilerplate server for 

- a grpc protobuf service
- with grpc-gateway
- opentelemetry

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
1. create a `boilerplate.GrpcRegisterFunc` and a `boilerplate.GatewayRegisterFunc` that wrap the generated register functions. These wrappers are called inbetween the creation of the servers and the start. They can also just be standalone functions (rather than being returned from a method on the ServiceImplementation).
    ```go
    func (i *ServiceImplementation) GrpcFunc() boilerplate.GrpcRegisterFunc {
        return func(s *grpc.Server) error {
            greeterv1.RegisterGreeterServiceServer(s, i)
            return nil
        }
    }

    func (i *ServiceImplementation) GatewayFunc() boilerplate.GatewayRegisterFunc {
        return func(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error {
            return greeterv1.RegisterGreeterServiceHandler(ctx, mux, cc)
        }
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
    server.RegisterGrpc(service.GrpcFunc())
    server.RegisterGateway(service.GatewayFunc())
    ```
1. Run the server
    ```go
    if err := server.Run(ctx); err != nil {
        log.Fatal(err)
    }
    ```
    
