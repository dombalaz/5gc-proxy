FROM golang
WORKDIR /app
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /5gc-proxy

FROM alpine
COPY --from=0 /5gc-proxy /5gc-proxy
COPY api /api
COPY config.json /config.json
EXPOSE 8080
CMD ["/5gc-proxy"]