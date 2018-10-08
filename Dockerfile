FROM golang:1.11.1-alpine
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .cgo -o exporter . && ls

FROM scratch
COPY --from=0 /go/src/app/exporter /bin/exporter
USER 65534:65534
ENTRYPOINT ["/bin/exporter"]
