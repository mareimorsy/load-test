# load-test

## 1- Please install the following package

```
go get -u -v github.com/leesper/go_rng
```

### 2 - Start sleep service

```
go run limit.go
curl http://localhost:8081/sleep
```

### 3 - Start load test service

```
go run rps-poisson.go
curl http://localhost:8080/load
curl http://localhost:8080/graph
```
to show scatter chart please visit `http://localhost:8080`

