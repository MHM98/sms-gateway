FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN  go build \
    -o /out/sms-gateway \
    .

FROM scratch

COPY --from=build /out/sms-gateway /sms-gateway

USER 65532:65532

EXPOSE 3000

STOPSIGNAL SIGTERM

ENTRYPOINT ["/sms-gateway"]
