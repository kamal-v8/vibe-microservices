"""
Configuration module for the Post Service.

Uses pydantic-settings to load configuration from environment variables,
following the 12-factor app methodology. A .env file is also supported
for local development.
"""

from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    # Server
    port: int = 8082

    # PostgreSQL connection string
    # Example: postgresql://user:pass@localhost:5432/pulse
    database_url: str = "postgresql://pulse:pulse@localhost:5432/pulse"

    # Redis connection string
    # Example: redis://localhost:6379/0
    redis_url: str = "redis://localhost:6379/0"

    # URL of the User Service for inter-service communication
    # Example: http://user-service:8081
    user_service_url: str = "http://localhost:8081"

    # Logging level
    log_level: str = "info"

    class Config:
        env_file = ".env"


# Singleton settings instance used throughout the application
settings = Settings()
