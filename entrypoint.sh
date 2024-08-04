echo "Running in $APP_ENV environment"

CompileDaemon --build="go build -o /app/main cmd/server/main.go" --command="/app/main"