# Timestamp Oracle

Timestamp Oracle is a golang implementation of timestamp service. It only support simple TS request function, which returns an auto-increment logic timestamp. 

The TS oracle can handle crash recovery by writing WAL to durable disk.

The TS oracle has high performance by batching TS requests. Each client maintains only one in-flight TS request RPC. TS oracle also allocate TSs in batch, which reduces the WAL IO cost.

## Usage

* Start Oracle Server
```
    go run tso/examples/server.go --addr=":7070"
```

* Client stub

The client stub is thread-safe.

```go
    client, err := tso.NewClient(":7070")
    if err != nil {
        log.Fatalln(err)
    }
    derfer client.Close()
    if ts, err := client.TS(); err != nil {
        log.Println("ts error")
    } else {
        ...
    }
```
