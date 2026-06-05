# TODO: Optimize with multi-stage build + distroless base
#
# This Dockerfile is INTENTIONALLY un-optimized — it uses a full golang base image,
# single-stage build, and includes the entire Go toolchain in the final image.
# Students will optimize it as a DevOps exercise (multi-stage, scratch/distroless,
# layer caching, .dockerignore, etc.).

# FROM golang:1.22
#
# WORKDIR /app
#
# # Copy dependency manifests first to leverage Docker layer caching.
# # (Even in this un-optimized version, this ordering habit is worth forming.)
# COPY go.mod go.sum ./
# RUN go mod download
#
# # Copy all source code and build the binary
# COPY . .
# RUN go build -o /app/user-service ./cmd/
# # RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/
#
#
# EXPOSE 8081
#
# CMD ["/app/user-service"]
#

FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/user-service ./cmd/
# RUN go build -o /app/user-service ./cmd/

FROM gcr.io/distroless/static-debian12

WORKDIR /app
USER nonroot:nonroot
COPY --from=builder /app/user-service /app/user-service

EXPOSE 8081
CMD ["/app/user-service"]


