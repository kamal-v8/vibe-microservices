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

RUN python -m venv /opt/venv

COPY requirements.txt .

RUN pip install --no-cache-dir --target=/opt/dependencies  -r requirements.txt 

COPY . . 

FROM python:3.11-slim

WORKDIR /app
RUN apt-get update && \
  apt-get upgrade -y && \
  rm -rf /var/lib/apt/lists/* && \
  pip uninstall -y pip setuptools wheel

COPY --from=builder /opt/dependencies /opt/dependencies
COPY . .
ENV PYTHONPATH=/opt/dependencies
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser
EXPOSE 8082

CMD ["python","-m","app.main"]
