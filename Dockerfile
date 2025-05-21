# syntax=docker/dockerfile:1
ARG GO_VERSION=1.24.0
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server ./

FROM oven/bun:1.2-alpine AS assets
WORKDIR /app

COPY package.json bun.lock webpack.config.js ./
RUN bun install --frozen-lockfile
COPY assets/ assets
RUN bun run build

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
    ca-certificates \
    tzdata \
    ffmpeg \
    && \
    update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY .env .
COPY public/ public

COPY --from=build /bin/server /bin/
COPY --from=assets /app/public/build.js /app/public/1.build.js public/

EXPOSE 8000

ENTRYPOINT [ "/bin/server" ]
