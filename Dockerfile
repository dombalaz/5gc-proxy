FROM golang
WORKDIR /app
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /5gc-proxy

FROM scratch
COPY --from=0 /5gc-proxy /bin/5gc-proxy
COPY api /
EXPOSE 8080
CMD ["/bin/5gc-proxy"]