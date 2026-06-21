# WAVE-06 Codegen/Lint/Schema Evidence (TEST-W06-013/014/015)

## Commands Run
```bash
cd apps/api && go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate
cd apps/api && go run github.com/99designs/gqlgen generate --config atlas-gqlgen.yml
cd apps/api && go vet ./internal/atlas/...
cd apps/api && go test -count=1 ./internal/atlas/models/ ./internal/atlas/service/ ./internal/atlas/graph/resolver/
cd apps/api && go build -o /dev/null ./cmd/server
```

## Results
- **sqlc**: No drift (generated code clean)
- **gqlgen**: No drift (generated code clean)
- **go vet**: All packages clean
- **go test**: 3 packages pass (models: 1 test, service: 20 tests, resolver: 16 tests)
- **go build**: Clean build

## Total test count: 37 (1 model + 20 service + 16 resolver)
All pass.
