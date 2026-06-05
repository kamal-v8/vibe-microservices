"""
Pydantic schemas (request / response models) for the Post Service.

These models drive automatic request validation, serialization, and
OpenAPI documentation generation in FastAPI.
"""

from datetime import datetime
from typing import Any, Optional
from uuid import UUID

from pydantic import BaseModel, Field


# ------------------------------------------------------------------
# Request models
# ------------------------------------------------------------------

class CreatePostRequest(BaseModel):
    """Payload for creating a new post."""

    user_id: UUID = Field(..., description="UUID of the post author")
    content: str = Field(
        ...,
        max_length=500,
        description="Post content (max 500 characters)",
    )


# ------------------------------------------------------------------
# Response models
# ------------------------------------------------------------------

class PostResponse(BaseModel):
    """Serialized representation of a single post."""

    id: UUID
    user_id: UUID
    content: str
    created_at: datetime

    class Config:
        from_attributes = True  # allow dict → model conversion


class HealthResponse(BaseModel):
    """Response from the /health endpoint."""

    status: str
    service: str
    timestamp: datetime


class APIResponse(BaseModel):
    """Consistent envelope for all API responses.

    Every endpoint returns this shape so clients can rely on a uniform
    contract:
        { "data": <payload | null>, "error": <message | null> }
    """

    data: Any = None
    error: Optional[str] = None
