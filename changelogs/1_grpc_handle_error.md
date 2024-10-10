
Необходимо заменить стандартную fatal на fatal с задержкой
А именно, заменить:
```go
s.log.Fatal(grpcServer.Serve(l))
```
На
```go
// start serving grpc
err := grpcServer.Serve(l)
s.log.Warn(err)
// wait 3 second, when grpc is ended to server. Is needed for log other fatals, if exists
time.Sleep(time.Second * 3)
// and crash error
s.log.Fatal(err)
```