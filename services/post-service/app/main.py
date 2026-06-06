"""
Pulse Post Service — application entry-point.

Responsibilities:
  • Create the FastAPI application with CORS + request-logging middleware
  • Wire up lifecycle events (startup / shutdown) to connect / disconnect
    from PostgreSQL, Redis, and the User Service HTTP client
  • Mount the /api/v1 router, health-check, and metrics stub
  • Configure structured JSON logging
  • Handle graceful shutdown on SIGTERM / SIGINT
"""

import asyncio
import json
import logging
import signal
import sys
import time
from contextlib import asynccontextmanager
from datetime import datetime, timezone

from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.database import Database
from app.routers.posts import router as posts_router
from app.schemas import APIResponse, HealthResponse
from app.services.event_publisher import EventPublisher
from app.services.user_client import UserServiceClient

# ------------------------------------------------------------------
# Structured JSON logging
# ------------------------------------------------------------------


class JSONFormatter(logging.Formatter):
    """Emit each log record as a single JSON object."""

    def format(self, record: logging.LogRecord) -> str:
        log_entry = {
            "level": record.levelname,
            "msg": record.getMessage(),
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "service": "post-service",
        }
        # Merge any extra keys added via `extra={...}`
        for key in (
            "user_id",
            "post_id",
            "status_code",
            "method",
            "path",
            "duration_ms",
            "attempt",
            "base_url",
            "channel",
        ):
            value = getattr(record, key, None)
            if value is not None:
                log_entry[key] = value
        return json.dumps(log_entry)


def _configure_logging() -> None:
    """Set up structured JSON logging for the entire application."""
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(JSONFormatter())

    logger = logging.getLogger("post-service")
    logger.setLevel(settings.log_level.upper())
    logger.addHandler(handler)
    logger.propagate = False


_configure_logging()
logger = logging.getLogger("post-service")


# ------------------------------------------------------------------
# Application lifespan (startup + shutdown)
# ------------------------------------------------------------------


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Manage the application lifecycle.

    On startup:
      1. Connect to PostgreSQL (with retries).
      2. Connect to Redis.
      3. Initialise the User Service HTTP client.

    On shutdown:
      Tear down all connections gracefully.
    """
    # -- Startup -------------------------------------------------------
    logger.info("Starting Post Service on port %d", settings.port)

    # PostgreSQL
    db = Database(settings.database_url)
    await asyncio.to_thread(db.connect)
    app.state.db = db

    # Redis / Event Publisher
    event_publisher = EventPublisher(settings.redis_url)
    await event_publisher.connect()
    app.state.event_publisher = event_publisher

    # User Service client
    user_client = UserServiceClient(settings.user_service_url)
    await user_client.start()
    app.state.user_client = user_client

    logger.info("Post Service startup complete")

    yield  # ← application runs here

    # -- Shutdown ------------------------------------------------------
    logger.info("Shutting down Post Service …")
    await user_client.close()
    await event_publisher.disconnect()
    await asyncio.to_thread(db.disconnect)
    logger.info("Post Service shutdown complete")


# ------------------------------------------------------------------
# FastAPI app
# ------------------------------------------------------------------

app = FastAPI(
    title="Pulse Post Service",
    description="Manages posts for the Pulse social platform",
    version="1.0.0",
    lifespan=lifespan,
)

# -- Middleware --------------------------------------------------------

# CORS — allow everything during development
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.middleware("http")
async def request_logging_middleware(request: Request, call_next):
    """Log every incoming request with method, path, status, and duration."""
    start = time.perf_counter()
    response: Response = await call_next(request)
    duration_ms = round((time.perf_counter() - start) * 1000, 2)

    logger.info(
        "%s %s → %d (%.2fms)",
        request.method,
        request.url.path,
        response.status_code,
        duration_ms,
        extra={
            "method": request.method,
            "path": request.url.path,
            "status_code": response.status_code,
            "duration_ms": duration_ms,
        },
    )
    return response


# -- Routers -----------------------------------------------------------

app.include_router(posts_router)


# -- Health check ------------------------------------------------------


@app.get("/api/v1/health", response_model=HealthResponse, tags=["health"])
async def health_check():
    """Health endpoint — checks DB and Redis connectivity."""
    db_ok = await asyncio.to_thread(app.state.db.ping)
    redis_ok = await app.state.event_publisher.ping()

    overall = "healthy" if (db_ok and redis_ok) else "degraded"

    return HealthResponse(
        status=overall,
        service="post-service",
        timestamp=datetime.now(timezone.utc),
    )


# -- Metrics stub ------------------------------------------------------


@app.get("/metrics", tags=["metrics"])
async def metrics():
    """Stub metrics endpoint (Prometheus / OpenTelemetry placeholder)."""
    return {"message": "metrics endpoint — not yet implemented"}


# ------------------------------------------------------------------
# Graceful shutdown signal handlers
# ------------------------------------------------------------------


def _handle_signal(sig, _frame):
    """Log the received signal and let uvicorn handle shutdown."""
    logger.info(
        "Received signal %s — initiating graceful shutdown", signal.Signals(sig).name
    )
    sys.exit(0)


signal.signal(signal.SIGTERM, _handle_signal)
signal.signal(signal.SIGINT, _handle_signal)


# ------------------------------------------------------------------
# Run with uvicorn when executed directly
# ------------------------------------------------------------------

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=settings.port,
        reload=False,
    )
