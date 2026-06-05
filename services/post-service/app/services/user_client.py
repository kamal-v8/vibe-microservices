"""
HTTP client for the User Service.

The Post Service needs to verify that a user exists before allowing
a post to be created.  This module encapsulates that inter-service
call using httpx's async client.
"""

import logging

import httpx

from app.config import settings

logger = logging.getLogger("post-service")


class UserServiceClient:
    """Async HTTP client that talks to the User Service."""

    def __init__(self, base_url: str | None = None) -> None:
        self.base_url = (base_url or settings.user_service_url).rstrip("/")
        self.client: httpx.AsyncClient | None = None

    async def start(self) -> None:
        """Create the underlying httpx.AsyncClient."""
        self.client = httpx.AsyncClient(
            base_url=self.base_url,
            timeout=httpx.Timeout(10.0),
        )
        logger.info("UserServiceClient initialised", extra={"base_url": self.base_url})

    async def close(self) -> None:
        """Shut down the httpx.AsyncClient."""
        if self.client:
            await self.client.aclose()
            logger.info("UserServiceClient closed")

    async def validate_user(self, user_id: str) -> bool:
        """Check whether a user exists by calling the User Service.

        Args:
            user_id: The UUID of the user to validate.

        Returns:
            True if the User Service responds with 200, False otherwise.
        """
        url = f"/api/v1/users/{user_id}"
        try:
            response = await self.client.get(url)

            if response.status_code == 200:
                logger.info("User validated successfully", extra={"user_id": user_id})
                return True

            if response.status_code == 404:
                logger.info("User not found", extra={"user_id": user_id})
                return False

            # Unexpected status
            logger.warning(
                "Unexpected status from User Service",
                extra={"user_id": user_id, "status_code": response.status_code},
            )
            return False

        except httpx.ConnectError as exc:
            logger.warning(
                "Could not connect to User Service: %s",
                str(exc),
                extra={"user_id": user_id},
            )
            return False
        except httpx.HTTPError as exc:
            logger.warning(
                "HTTP error calling User Service: %s",
                str(exc),
                extra={"user_id": user_id},
            )
            return False
