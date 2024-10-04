FROM node:20-slim AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY frontend/ /app
WORKDIR /app

FROM base AS frontend_build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile
RUN BASE_URL=/aquarium/ pnpm run build

FROM golang:1.23-alpine3.20 AS build
WORKDIR /src
COPY ./server/. .

RUN go mod download

RUN go build -ldflags="-w -s" -o /src/bin/aquarium

RUN ["chmod", "+x", "/src/bin/aquarium"]

FROM scratch

WORKDIR /bin

COPY --from=build /src/bin/aquarium /bin/aquarium
COPY --from=frontend_build /app/dist ./assets/aquarium

ENV AQUARIUM_PORT=8080
EXPOSE 8080

CMD ["/bin/aquarium"]
