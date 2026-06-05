"""
Redis event publisher for the Post Service.

Publishes domain events (e.g. "new_post") to a Redis Pub/Sub channel
so that other services (feed-service, notification-service, etc.) can
react in near-real-time.
"""

import json
import logging
from datetime import datetime, timezone

import redis.asyncio as aioredis

logger = logging.getLogger("post-service")

# Channel that all Pulse platform events are published to
EVENTS_CHANNEL = "pulse:events"


class EventPublisher:
    """Publishes JSON-encoded events to a Redis Pub/Sub channel."""

    def __init__(self, redis_url: str) -> None:
        self.redis_url = redis_url
        self.redis: aioredis.Redis | None = None

    async def connect(self) -> None:
        """Establish the Redis connection."""
        self.redis = aioredis.from_url(
            self.redis_url,
            decode_responses=True,
        )
        # Verify connectivity
        await self.redis.ping()
        logger.info("EventPublisher connected to Redis")

    async def disconnect(self) -> None:
        """Close the Redis connection."""
        if self.redis:
            await self.redis.close()
            logger.info("EventPublisher disconnected from Redis")

    async def publish_new_post(self, post_data: dict) -> None:
        """Publish a 'new_post' event to the pulse:events channel.

        Args:
            post_data: Dict representation of the newly created post.
        """
        event = {
            "event": "new_post",
            "data": post_data,
            "timestamp": datetime.now(timezone.utc).isoformat(),
        }
        payload = json.dumps(event, default=str)

        await self.redis.publish(EVENTS_CHANNEL, payload)
        logger.info(
            "Published new_post event",
            extra={"post_id": post_data.get("id"), "channel": EVENTS_CHANNEL},
        )

    async def ping(self) -> bool:
        """Return True if Redis is reachable."""
        try:
            return await self.redis.ping()
        except Exception:
            return False
