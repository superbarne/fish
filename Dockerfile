ARG VERSION=1.23-alpine3.20
FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY frontend/ /app
WORKDIR /app

FROM base AS frontend_build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN pnpm run build

FROM golang:$VERSION AS build
WORKDIR /src
COPY ./server/. .

WORKDIR /src/aquarium
RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /src/bin/aquarium

RUN ["chmod", "+x", "/src/bin/aquarium"]

FROM scratch

COPY --from=build /src/bin/aquarium /bin/aquarium
COPY --from=frontend_build /app/dist ./assets/aquarium

WORKDIR /bin
EXPOSE 5555

ENV AQUARIUM_PORT=8080
EXPOSE 8080

CMD ["/bin/aquarium"]
