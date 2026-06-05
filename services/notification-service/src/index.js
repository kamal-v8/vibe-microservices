/**
 * Main entry point for the Pulse Notification Service.
 *
 * Responsibilities:
 *   1. Boot a Fastify HTTP server
 *   2. Connect TWO Redis clients (subscriber + storage)
 *   3. Start the EventSubscriber to listen on pulse:events
 *   4. Expose REST API, health-check, and metrics endpoints
 *   5. Handle graceful shutdown on SIGTERM / SIGINT
 */

'use strict';

const Fastify = require('fastify');
const Redis = require('ioredis');
const config = require('./config');
const logger = require('./logger');
const EventSubscriber = require('./subscriber');
const notificationRoutes = require('./routes/notifications');

// ---------------------------------------------------------------------------
// Fastify instance
// ---------------------------------------------------------------------------
const fastify = Fastify({
  logger: {
    level: config.logLevel,
  },
});

// ---------------------------------------------------------------------------
// Redis clients
// ---------------------------------------------------------------------------

/**
 * Helper: create an ioredis client with standard event handlers.
 *
 * @param {string} name  Friendly label used in log messages
 * @returns {import('ioredis').Redis}
 */
function createRedisClient(name) {
  const client = new Redis(config.redisUrl, {
    // ioredis auto-reconnects by default — we just log the attempts
    retryStrategy(times) {
      const delay = Math.min(times * 200, 5000);
      logger.warn(`Redis ${name} reconnecting — attempt ${times}, delay ${delay}ms`);
      return delay;
    },
    maxRetriesPerRequest: null, // let ioredis retry indefinitely
    lazyConnect: false,         // connect immediately on instantiation
  });

  client.on('connect', () => {
    logger.info(`Redis ${name} client connected`, { redisUrl: config.redisUrl });
  });

  client.on('error', (err) => {
    logger.error(`Redis ${name} client error`, { error: err.message });
  });

  client.on('close', () => {
    logger.warn(`Redis ${name} client connection closed`);
  });

  return client;
}

// Dedicated client for SUBSCRIBE (cannot run other commands while subscribed)
const subscriberClient = createRedisClient('subscriber');
// General-purpose client for LPUSH / LRANGE / LTRIM / PING
const storageClient = createRedisClient('storage');

// ---------------------------------------------------------------------------
// Event subscriber
// ---------------------------------------------------------------------------
const eventSubscriber = new EventSubscriber(subscriberClient, storageClient);

// ---------------------------------------------------------------------------
// Decorate Fastify so route plugins can access the storage client
// ---------------------------------------------------------------------------
fastify.decorate('redis', storageClient);

// ---------------------------------------------------------------------------
// Health check — GET /api/v1/health
// ---------------------------------------------------------------------------
fastify.get('/api/v1/health', async (_request, reply) => {
  try {
    await storageClient.ping();
    return reply.code(200).send({
      status: 'ok',
      service: 'notification-service',
      timestamp: new Date().toISOString(),
      redis: 'connected',
    });
  } catch (err) {
    logger.error('Health check — Redis ping failed', { error: err.message });
    return reply.code(200).send({
      status: 'degraded',
      service: 'notification-service',
      timestamp: new Date().toISOString(),
      redis: 'disconnected',
    });
  }
});

// ---------------------------------------------------------------------------
// Metrics stub — GET /metrics
// ---------------------------------------------------------------------------
fastify.get('/metrics', async (_request, reply) => {
  return reply.code(200).send({
    message: 'metrics endpoint - integrate with Prometheus',
  });
});

// ---------------------------------------------------------------------------
// Register route plugins
// ---------------------------------------------------------------------------
fastify.register(notificationRoutes);

// ---------------------------------------------------------------------------
// Graceful shutdown
// ---------------------------------------------------------------------------

/**
 * Perform an orderly teardown of all resources.
 *
 * @param {string} signal  The signal that triggered shutdown
 */
async function shutdown(signal) {
  logger.info(`Received ${signal} — starting graceful shutdown`);

  try {
    // 1. Stop subscribing to events
    await eventSubscriber.stop();

    // 2. Disconnect Redis clients
    subscriberClient.disconnect();
    storageClient.disconnect();

    // 3. Close the HTTP server (finishes in-flight requests)
    await fastify.close();

    logger.info('Graceful shutdown complete');
  } catch (err) {
    logger.error('Error during shutdown', { error: err.message });
  } finally {
    process.exit(0);
  }
}

process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));

// ---------------------------------------------------------------------------
// Start
// ---------------------------------------------------------------------------
async function main() {
  try {
    // Start listening for platform events on Redis Pub/Sub
    await eventSubscriber.start();

    // Start the HTTP server on all interfaces
    await fastify.listen({ port: config.port, host: '0.0.0.0' });

    logger.info(`Notification service listening on port ${config.port}`, {
      port: config.port,
    });
  } catch (err) {
    logger.error('Failed to start notification service', { error: err.message });
    process.exit(1);
  }
}

main();
