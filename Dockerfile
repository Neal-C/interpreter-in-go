FROM golang:1.21.3 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /nealc-interpreter
FROM alpine:latest as production
WORKDIR /root/
COPY --from=build /nealc-interpreter ./
ENTRYPOINT ./nealc-interpreter