FROM alpine:latest
WORKDIR /app
COPY bin/* /app/caco
ENTRYPOINT ["/app/caco"]