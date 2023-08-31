# auth-service-sample

#### Language: Golang & Gin Framework

### Install:
```
go mod tidy
```

### Run:
```
# -p: serve port(default 8080)
# -tt: token lifetime(default 7200)

go run main.go -p 8080 -tt 7200
```

### Test:
```
# FullFlow Test: token lifetime 5 second
go test -run=FullFlow . -v

# Benchmark: unique for creating same user
# 10 times call with 1 SUCCESS and 9 FAIL
go test -run=none -bench=CreateUserUnique .

# Benchmark: 10000 times call for [All roles]
go test -run=none -bench=UserRoles -count=0 .
```

### API Doc
see ./docs/swagger.yaml