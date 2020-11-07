FROM alpine:latest
WORKDIR /app
COPY ./cmd/simple_dns/simple_dns /app
#COPY ./cert/ /app/cert
COPY ./configs/config.yaml /app/configs/config.yaml

CMD ["/app/simple_dns"]