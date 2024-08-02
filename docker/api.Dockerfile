FROM golang:1.22-alpine
WORKDIR /app

# add some necessary packages
RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

# prevent the re-installation of vendors at every change in the source code
COPY ./go.mod go.sum ./
RUN go mod download && go mod verify

# Download the Compile Daemon package and its dependencies
RUN go get -d -v github.com/githubnemo/CompileDaemon

# Install Compile Daemon for go
RUN go install -v github.com/githubnemo/CompileDaemon

# Make sure the Compile Daemon is in your PATH
ENV PATH=$PATH:/app/bin

# Copy and build the app
COPY . .
COPY ./entrypoint.sh /entrypoint.sh
COPY ./config.yml /config.yml

# # wait-for-it requires bash, which alpine doesn't ship with by default. Use wait-for instead
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

ENTRYPOINT [ "sh", "/entrypoint.sh" ]
