from golang:alpine3.19 as build
WORKDIR /usr/src/app
COPY . .
RUN apk add --no-cache make
RUN go mod tidy
RUN make build

from build as deploy
COPY --from=build /usr/src/app/dist ./dist
COPY --from=build /usr/src/app/makefile ./makefile
EXPOSE 3333

CMD ["make", "start"]
