FROM docker.arvancloud.ir/golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build \
    -o /out/db-migrate \
    .

FROM scratch

COPY --from=build /out/db-migrate /db-migrate
COPY migrations /migrations

WORKDIR /

USER 65532:65532

ENTRYPOINT ["/db-migrate"]
CMD ["up"]
