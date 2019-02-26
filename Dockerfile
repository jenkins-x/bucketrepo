FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY ./config/config.yaml /
COPY ./bin/ /
ENTRYPOINT ["/bucketrepo"]
