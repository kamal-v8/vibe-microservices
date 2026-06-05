# TODO: Optimize with multi-stage build + distroless base
#
# This Dockerfile is INTENTIONALLY un-optimized.
# It uses the full python:3.12 image (≈1 GB) and a single stage.
# As a DevOps exercise, consider:
#   - Multi-stage builds to separate build and runtime
#   - Slim or distroless base images to shrink the final image
#   - Layer caching optimizations (requirements before source)
#   - Non-root user for security
#   - .dockerignore to reduce context size
#
# FROM python:3.12
#
# WORKDIR /app
#
# # Install Python dependencies
# COPY requirements.txt .
# RUN pip install --no-cache-dir -r requirements.txt
#
# # Copy application source
# COPY . .
#
# EXPOSE 8082
#
# CMD ["python", "-m", "app.main"]
#

FROM python:3.11-slim AS builder 

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt 

COPY . . 

FROM gcr.io/distroless/python3-debian12 

WORKDIR /app

# COPY --from=builder /app /app

COPY --from=builder /usr/local/lib/python3.11/site-packages /usr/lib/python3/dist-packages
COPY . .

USER nonroot:nonroot

EXPOSE 8082

CMD ["-m","app.main"]


