# dnsgate-go

Go client library for [dnsgate](https://github.com/dnsgateway), the DNS Gateway. 
Use it to build your own client to a Gateway.

| Package | What it is |
|---|---|
| [`api/dnsgate/v1`](api/dnsgate/v1) | The Gateway's protobuf API (`gateway.proto`) and its generated Go and gRPC code. |
| [`client`](client) | A typed Gateway client: one method per operation, Gateway errors as `*client.Error`, and the Egress Agent side of Agent Egress. |
| [`gatewayadapter`](gatewayadapter) | The Gateway Adapter: a [`dnsprovider`](https://github.com/dnsgateway/dnsprovider) Provider backed by the Gateway, for DNS tools built on the shared Provider interface. |

## Usage

```go
conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})))
c := client.New(conn)
zones, err := c.ListZones(ctx, &dnsgatev1.ListZonesRequest{Token: serviceToken})
```

## License

Apache License 2.0, see [LICENSE](LICENSE).
