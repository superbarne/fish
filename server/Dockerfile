ARG VERSION=1.23-alpine3.20
FROM golang:$VERSION AS build

WORKDIR /src
COPY . .

WORKDIR /src/aquarium
RUN go mod download
RUN go build -o /src/bin/aquarium

FROM scratch

COPY --from=build /src/bin/aquarium /bin/aquarium

WORKDIR /bin
EXPOSE 5555

ENV AQUARIUM_PORT=8080
EXPOSE 8080

CMD ["aquarium"]
