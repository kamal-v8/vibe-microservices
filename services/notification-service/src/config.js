/**
 * Configuration module for the Notification Service.
 *
 * Reads all tunables from environment variables (12-factor style)
 * and falls back to sensible development defaults.
 */

'use strict';

const config = {
  /** HTTP port for the Fastify server */
  port: parseInt(process.env.PORT, 10) || 8083,

  /** Redis connection string — used for both Pub/Sub and storage */
  redisUrl: process.env.REDIS_URL || 'redis://localhost:6379/0',

  /** Pino / application log level */
  logLevel: process.env.LOG_LEVEL || 'info',

  /** Redis Pub/Sub channel that carries platform-wide events */
  eventChannel: 'pulse:events',

  /** Prefix for per-user notification lists in Redis */
  notificationKeyPrefix: 'pulse:notifications:',

  /** Redis key for the global (all-users) notification list */
  allNotificationsKey: 'pulse:notifications:all',
};

module.exports = config;
