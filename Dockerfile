from golang:alpine3.19 as build
WORKDIR /usr/src/app
COPY . .
RUN apk add --no-cache make
RUN go mod tidy
RUN make build

from golang:alpine3.19 as deploy
COPY from=build /usr/src/app/dist ./dist
EXPOSE 8080

CMD ["./dist/main"]
