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
# cpus * 5 parallel calls(each 5 times) with only 1 SUCCESS
go test -run=none -bench=CreateUserUnique -count=0 .

# Benchmark: cpus * 2500 parallel calls(each 10 times) for [All roles]
go test -run=none -bench=UserRoles -count=0 .
```

### API Doc
see ./docs/swagger.yaml