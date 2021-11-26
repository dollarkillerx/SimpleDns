FROM alpine:latest
WORKDIR /app
COPY ./cmd/simple_dns/simple_dns /app
COPY ./config.yaml /app/config.yaml
COPY ./routing_table.csv /app/configs/routing_table.csv

CMD ["/app/simple_dns"]
