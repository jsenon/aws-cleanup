FROM golang:1.14 as build

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN make build_local

FROM gcr.io/distroless/base
COPY --from=build /app/aws-cleanup /
WORKDIR /app/data
CMD ["/aws-cleanup"]
