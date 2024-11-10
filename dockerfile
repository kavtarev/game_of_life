FROM golang:1.23.0-alpine3.20 as build
WORKDIR /app
COPY . .
WORKDIR /app/back
RUN go mod download
RUN go build -o ./bin/game_of_life ./cmd/main
CMD [ "./main" ]

FROM alpine:3
COPY --from=build app/back/bin/game_of_life /game/back/game_of_life
COPY --from=build app/front /game/front
ENTRYPOINT ["/game/back/game_of_life"]