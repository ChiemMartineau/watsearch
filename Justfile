dev:
    go tool templ generate --watch --proxy="http://localhost:1234" --cmd="go run ./cmd/server"