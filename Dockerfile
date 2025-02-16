FROM golang:1.23-alpine AS build
WORKDIR /app
# install deps
COPY go.mod go.sum ./
RUN go mod download
# copy source files
COPY . .
# run build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app
# run migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV GOOSE_DRIVER="postgres" \
    GOOSE_DBSTRING="$POSTGRES_DB_URL" \
    GOOSE_MIGRATION_DIR="./migrations"
RUN goose up

FROM node:22-alpine AS node-build
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app/frontend
# install deps
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
# copy source files
COPY frontend/ .
COPY templates/ /app/templates
# run build
RUN pnpm run build

FROM alpine:latest
WORKDIR /app
ENV NODE_ENV="production"
COPY --from=build /app/bin/app .
COPY --from=build /app/templates/ templates/
COPY --from=node-build /app/frontend/dist/ frontend/dist/
EXPOSE 8080

# run
CMD ["./app"]
