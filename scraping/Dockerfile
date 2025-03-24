FROM node:22 AS ui-build

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY .git /app/.git
COPY ui /app
WORKDIR /app

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

RUN VITE_COMMIT_HASH="$(cat .git/$(cat .git/HEAD | cut -d' ' -f2) | cut -c1-12)" pnpm run build

FROM golang:1.23-alpine as build

RUN apk add --update --no-cache build-base

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

COPY --from=ui-build /app/dist /go/src/app/scraping/api/ui

WORKDIR /go/src/app
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -o out/bigcord-scraping \
    -v \
    -ldflags="-X main.CommitHash=$(cat .git/$(cat .git/HEAD | cut -d' ' -f2) | cut -c1-12)" \
    # delve
    -gcflags "all=-N -l" \
    github.com/K3das/bigcord/scraping

FROM alpine

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN mkdir -p /data/media && chown -R appuser:appgroup /data /data/media

COPY --from=build --chown=appuser:appgroup /go/bin/dlv /dlv
COPY --from=build --chown=appuser:appgroup /go/src/app/out/ /opt/app/

WORKDIR /opt/app

ENV SCRAPING_SQLITE_DSN=/data/scraping.db

USER appuser
ENTRYPOINT [ "/opt/app/bigcord-scraping" ]
