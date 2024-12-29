FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /unraid-monitoring-operator

FROM debian:12-slim

COPY --from=build /unraid-monitoring-operator /unraid-monitoring-operator
COPY /dashboards /dashboards

ENTRYPOINT ["/unraid-monitoring-operator"]
