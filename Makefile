.PHONY: test evalops evalops-gate

test:
	go test ./...

evalops:
	go test ./internal/eval ./cmd/evalops-target ./cmd/evalops-gate -count=1

evalops-gate:
	go run ./cmd/evalops-gate
