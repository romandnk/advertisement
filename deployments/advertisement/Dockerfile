FROM golang:1.20 as build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o "./bin/advertisement" ./cmd/advertisement

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/advertisement ./bin/

COPY ./configs/config.yaml ./configs/
COPY ./configs/.env ./configs/
COPY ./static/images ./static/images

CMD ["./bin/advertisement"]