FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /grafana-dashboards-syncer

FROM debian:12-slim

COPY --from=build /grafana-dashboards-syncer /grafana-dashboards-syncer
COPY /dashboards /dashboards

ENTRYPOINT ["/grafana-dashboards-syncer"]
