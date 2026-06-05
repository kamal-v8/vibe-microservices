# TODO: Optimize with multi-stage build + distroless base
# This Dockerfile is intentionally un-optimized for a DevOps exercise.
# Improvements to consider:
#   - Multi-stage build to separate build and runtime layers
#   - Use node:20-alpine or distroless for a smaller final image
#   - Run as a non-root user
#   - Leverage .dockerignore to exclude unnecessary files
#
# FROM node:20
#
# WORKDIR /app
#
# # Copy dependency manifests first (layer caching for npm install)
# COPY package.json package-lock.json* ./
#
# RUN npm install
#
# # Copy the rest of the application source
# COPY . .
#
# EXPOSE 8083
#
# CMD ["node", "src/index.js"]
#

FROM node:20 AS builder

WORKDIR /app

COPY package.json package.lock.json* ./


RUN npm install ci --only=production

COPY . .

FROM gcr.io/distroless/nodejs20-debian12

WORKDIR /app

COPY --from=builder /app /app

USER nonroot:nonroot

EXPOSE 8083

CMD [ "src/index.js" ]
