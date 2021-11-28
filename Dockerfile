#
# BUILD
#
FROM golang:1.16-buster AS build
WORKDIR /app

COPY . /app/
RUN go mod download

RUN go build -o /jamz

#
# Deploy
#
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build jamz jamz

USER nonroot:nonroot

ENTRYPOINT ["/jamz"]


