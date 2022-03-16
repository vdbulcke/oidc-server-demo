FROM alpine:latest  
RUN apk --no-cache add ca-certificates
## non privileged user
USER 1000
# EXPOSE 9000
WORKDIR /app/
COPY oidc-server /app/oidc-server

## For docker listen on all interfaces
ENTRYPOINT ["/app/oidc-server", "start", "--config", "/app/config.yaml", "--listen-addr", "0.0.0.0"]
