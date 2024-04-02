from golang:alpine3.19 as build
WORKDIR /usr/src/app
COPY . .
RUN apk add --no-cache make curl tar
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-386.tar.gz | tar xvz
RUN go mod tidy
RUN make build

from build as deploy
COPY --from=build /usr/src/app/dist ./dist
COPY --from=build /usr/src/app/makefile ./makefile
COPY --from=build /usr/src/app/migrate ./migrate
COPY --from=build /usr/src/app/.env.prod ./.env
EXPOSE 3333

CMD ["make", "migration_up"]
CMD ["make", "start"]
