# Watch your .go files and invoke go build if the files changed.
CompileDaemon --build="go build -o main cmd/server/main.go"  --command=./main 
