FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build \
    -o /out/sms-worker \
    .

FROM scratch

COPY --from=build /out/sms-worker /sms-worker

USER 65532:65532

ENTRYPOINT ["/sms-worker"]
