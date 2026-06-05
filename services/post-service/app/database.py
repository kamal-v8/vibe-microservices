"""
Database module — thin wrapper around psycopg2 for direct PostgreSQL access.

This deliberately does NOT use SQLAlchemy or any ORM.  All queries are raw SQL,
and results come back as plain Python dicts thanks to RealDictCursor.

Connection retry logic is built in so the service can start before the
database is fully ready (common in container orchestration).
"""

import logging
import time
from typing import Any

import psycopg2
import psycopg2.extras

logger = logging.getLogger("post-service")


class Database:
    """Simple PostgreSQL connection wrapper with retry and dict cursors."""

    def __init__(self, dsn: str) -> None:
        """
        Args:
            dsn: PostgreSQL connection string (DATABASE_URL).
        """
        self.dsn = dsn
        self.connection = None

    # ------------------------------------------------------------------
    # Lifecycle
    # ------------------------------------------------------------------

    def connect(self, retries: int = 5, backoff: float = 3.0) -> None:
        """Establish a database connection with retry logic.

        Args:
            retries:  Number of connection attempts before giving up.
            backoff:  Seconds to wait between retries.
        """
        for attempt in range(1, retries + 1):
            try:
                self.connection = psycopg2.connect(
                    self.dsn,
                    cursor_factory=psycopg2.extras.RealDictCursor,
                )
                # Enable autocommit so each statement runs in its own
                # transaction — keeps things simple for a microservice.
                self.connection.autocommit = True
                logger.info(
                    "Database connection established",
                    extra={"attempt": attempt},
                )
                return
            except psycopg2.OperationalError as exc:
                logger.warning(
                    "Database connection attempt %d/%d failed: %s",
                    attempt,
                    retries,
                    str(exc),
                )
                if attempt < retries:
                    time.sleep(backoff)
                else:
                    logger.error("All %d database connection attempts exhausted", retries)
                    raise

    def disconnect(self) -> None:
        """Close the database connection gracefully."""
        if self.connection and not self.connection.closed:
            self.connection.close()
            logger.info("Database connection closed")

    # ------------------------------------------------------------------
    # Query helpers
    # ------------------------------------------------------------------

    def execute(self, query: str, params: tuple | None = None) -> Any:
        """Execute a query and return the cursor (useful for INSERT/UPDATE/DELETE).

        Args:
            query:  SQL query string with %s placeholders.
            params: Tuple of parameters to bind.

        Returns:
            The cursor after execution.
        """
        cursor = self.connection.cursor()
        try:
            cursor.execute(query, params)
            return cursor
        except Exception:
            cursor.close()
            raise

    def fetchone(self, query: str, params: tuple | None = None) -> dict | None:
        """Execute a query and return a single row as a dict, or None.

        Args:
            query:  SQL query string.
            params: Tuple of parameters.

        Returns:
            A dict representing one row, or None if no rows matched.
        """
        cursor = self.execute(query, params)
        try:
            return cursor.fetchone()
        finally:
            cursor.close()

    def fetchall(self, query: str, params: tuple | None = None) -> list[dict]:
        """Execute a query and return all rows as a list of dicts.

        Args:
            query:  SQL query string.
            params: Tuple of parameters.

        Returns:
            A (possibly empty) list of dicts.
        """
        cursor = self.execute(query, params)
        try:
            return cursor.fetchall()
        finally:
            cursor.close()

    # ------------------------------------------------------------------
    # Health check
    # ------------------------------------------------------------------

    def ping(self) -> bool:
        """Return True if the database is reachable."""
        try:
            self.fetchone("SELECT 1")
            return True
        except Exception:
            return False
