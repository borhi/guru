FROM golang:1.14-alpine as build
WORKDIR /src/guru
COPY . .
RUN go mod vendor
RUN go build -mod=vendor

FROM alpine:latest
COPY --from=build /src/guru/.env .
COPY --from=build /src/guru/guru .
COPY --from=build /src/guru/swaggerui ./swaggerui
CMD ["./guru", "-port", "3000"]
EXPOSE 3000
