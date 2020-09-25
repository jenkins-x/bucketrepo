FROM gcr.io/distroless/base:nonroot
COPY ./config/config.yaml /home/nonroot/
COPY ./bin /home/nonroot/
CMD ["/home/nonroot/bucketrepo"]
