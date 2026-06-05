/**
 * Notification REST routes — Fastify plugin.
 *
 * Endpoints:
 *   GET /api/v1/notifications           — list all notifications (global)
 *   GET /api/v1/notifications/:userId   — list notifications for a specific user
 *
 * Notifications are stored as JSON strings in Redis lists and parsed on read.
 */

'use strict';

const config = require('../config');
const logger = require('../logger');

/** Default and maximum page size */
const DEFAULT_LIMIT = 50;

/**
 * Fastify route plugin.  The `redis` storage client is expected to be
 * available via `fastify.redis` (decorated in index.js).
 *
 * @param {import('fastify').FastifyInstance} fastify
 * @param {object} _options
 */
async function notificationRoutes(fastify, _options) {
  const redis = fastify.redis;

  // --------------------------------------------------------------------------
  // GET /api/v1/notifications — Global notification feed
  // --------------------------------------------------------------------------
  fastify.get('/api/v1/notifications', async (request, reply) => {
    try {
      const limit = Math.min(
        parseInt(request.query.limit, 10) || DEFAULT_LIMIT,
        DEFAULT_LIMIT,
      );
      const offset = parseInt(request.query.offset, 10) || 0;

      // LRANGE is 0-indexed and inclusive on both ends
      const rawItems = await redis.lrange(
        config.allNotificationsKey,
        offset,
        offset + limit - 1,
      );

      const notifications = rawItems.map((item) => {
        try {
          return JSON.parse(item);
        } catch {
          logger.warn('Skipping malformed notification entry', { raw: item });
          return null;
        }
      }).filter(Boolean);

      return reply.code(200).send({ data: notifications, error: null });
    } catch (err) {
      logger.error('Failed to fetch notifications', { error: err.message });
      return reply.code(500).send({
        data: null,
        error: 'Internal server error while fetching notifications',
      });
    }
  });

  // --------------------------------------------------------------------------
  // GET /api/v1/notifications/:userId — Per-user notification feed
  // --------------------------------------------------------------------------
  fastify.get('/api/v1/notifications/:userId', async (request, reply) => {
    try {
      const { userId } = request.params;
      const userKey = `${config.notificationKeyPrefix}${userId}`;

      const limit = Math.min(
        parseInt(request.query.limit, 10) || DEFAULT_LIMIT,
        DEFAULT_LIMIT,
      );
      const offset = parseInt(request.query.offset, 10) || 0;

      const rawItems = await redis.lrange(userKey, offset, offset + limit - 1);

      const notifications = rawItems.map((item) => {
        try {
          return JSON.parse(item);
        } catch {
          logger.warn('Skipping malformed notification entry', { raw: item });
          return null;
        }
      }).filter(Boolean);

      // An empty list is perfectly valid — the user just has no notifications yet
      return reply.code(200).send({ data: notifications, error: null });
    } catch (err) {
      logger.error('Failed to fetch user notifications', {
        error: err.message,
        userId: request.params.userId,
      });
      return reply.code(500).send({
        data: null,
        error: 'Internal server error while fetching user notifications',
      });
    }
  });
}

module.exports = notificationRoutes;
