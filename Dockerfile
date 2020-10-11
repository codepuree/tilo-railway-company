FROM golang:1.15.2 AS go-builder
WORKDIR /go/src/github.com/codepuree/tilo-railway-company
COPY . /go/src/github.com/codepuree/tilo-railway-company
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -tags netgo -ldflags '-w' -o /go/bin/trc ./cmd/

FROM node:13.10.1-slim AS web-builder
RUN npm install -g esbuild-linux-64
WORKDIR /usr/src/trc/web/static
COPY ./web/static /usr/src/trc/web/static
RUN esbuild --bundle --minify ./js/*.js

FROM scratch
WORKDIR /opt/trc
COPY --from=go-builder /go/bin/trc /opt/trc/trc
# COPY --from=web-builder /usr/src/trc/web/static /opt/trc/web/static
COPY ./web/static /var/www/static
ENTRYPOINT [ "/opt/trc/trc" ]
EXPOSE 8080