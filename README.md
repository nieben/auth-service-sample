# auth-service-sample

Language: Golang & Gin Framework

Install:
```
go mod tidy
```

Run:
```
# -p: serve port(default 8080)
# -tt: token lifetime(default 7200)

go run main.go -p 8080 -tt 7200
```

Test:
```
# token lifetime: 5 second

go test . -v
```

API Doc: see ./docs/swagger.yaml 