FROM golang AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main 

# second stage distroless image

FROM gcr.io/distroless/base

COPY --from=base /app/main .

EXPOSE 8080

CMD [ "./main" ]