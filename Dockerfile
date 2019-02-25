FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY ./config/config.yaml /etc/nexus-minimal/
COPY ./bin/ /
ENTRYPOINT ["/nexus-minimal"]
