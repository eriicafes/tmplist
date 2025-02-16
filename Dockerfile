FROM golang:1.23-alpine AS build
WORKDIR /app
# install deps
COPY go.mod go.sum ./
RUN go mod download
# copy source files
COPY . .
# run build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app

FROM node:22-alpine AS node-build
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app
# install deps
COPY frontend/package.json frontend/pnpm-lock.yaml frontend/
RUN cd frontend && pnpm install --frozen-lockfile
# copy source files
COPY frontend/ templates/ ./
# run build
RUN cd frontend && pnpm run build

FROM alpine:latest
WORKDIR /app
ENV NODE_ENV="production"
COPY --from=build /app/bin/app .
COPY --from=build /app/templates templates
COPY --from=node-build /app/frontend/dist frontend/dist
EXPOSE 8080

# run
CMD ["./app"]
