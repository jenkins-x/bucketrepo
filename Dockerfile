FROM gcr.io/distroless/base:nonroot

COPY --chown=nonroot:nonroot ./config/config.yaml ./
COPY --chown=nonroot:nonroot ./bin ./
ENTRYPOINT ["/home/nonroot/bucketrepo"]
