/**
 * EventSubscriber — listens on the Redis Pub/Sub channel `pulse:events`
 * and converts incoming platform events into stored notifications.
 *
 * Architecture note:
 *   Redis requires a dedicated client for SUBSCRIBE — once a client enters
 *   subscriber mode it cannot run regular commands.  That's why this class
 *   receives TWO clients: one for subscribing and one for writing.
 */

'use strict';

const { v4: uuidv4 } = require('uuid');
const config = require('./config');
const logger = require('./logger');

/** Maximum notifications to keep per list (global + per-user) */
const MAX_NOTIFICATIONS = 100;

class EventSubscriber {
  /**
   * @param {import('ioredis').Redis} redisSubscriberClient  Dedicated SUBSCRIBE client
   * @param {import('ioredis').Redis} redisStorageClient     Client for LPUSH / LTRIM
   */
  constructor(redisSubscriberClient, redisStorageClient) {
    this.subscriber = redisSubscriberClient;
    this.storage = redisStorageClient;
    this._messageHandler = null; // stored so we can cleanly remove it
  }

  /**
   * Subscribe to the platform event channel and begin processing messages.
   */
  async start() {
    // Register the message handler BEFORE subscribing so we never miss an event
    this._messageHandler = async (channel, message) => {
      await this._handleMessage(channel, message);
    };
    this.subscriber.on('message', this._messageHandler);

    await this.subscriber.subscribe(config.eventChannel);
    logger.info(`Subscribed to Redis channel: ${config.eventChannel}`);
  }

  /**
   * Unsubscribe and remove the message listener.
   */
  async stop() {
    try {
      await this.subscriber.unsubscribe(config.eventChannel);
      if (this._messageHandler) {
        this.subscriber.removeListener('message', this._messageHandler);
      }
      logger.info(`Unsubscribed from Redis channel: ${config.eventChannel}`);
    } catch (err) {
      logger.error('Error during unsubscribe', { error: err.message });
    }
  }

  // -------------------------------------------------------------------------
  // Internal helpers
  // -------------------------------------------------------------------------

  /**
   * Process a single Pub/Sub message.
   *
   * @param {string} channel  The channel name (should be pulse:events)
   * @param {string} message  Raw JSON string
   */
  async _handleMessage(channel, message) {
    try {
      const event = JSON.parse(message);

      if (event.event === 'new_post') {
        await this._handleNewPost(event);
      } else {
        logger.warn(`Unknown event type received: ${event.event}`, {
          channel,
          eventType: event.event,
        });
      }
    } catch (err) {
      logger.error('Failed to process event message', {
        error: err.message,
        rawMessage: message,
      });
    }
  }

  /**
   * Turn a `new_post` event into a notification and persist it in Redis.
   *
   * @param {object} event  Parsed event payload
   */
  async _handleNewPost(event) {
    const { id: postId, user_id: userId } = event.data;

    // Build the notification object
    const notification = {
      id: uuidv4(),
      type: 'new_post',
      message: `User ${userId} created a new post`,
      post_id: postId,
      user_id: userId,
      created_at: new Date().toISOString(),
      read: false,
    };

    const serialized = JSON.stringify(notification);

    // Key for the user-specific notification list
    const userKey = `${config.notificationKeyPrefix}${userId}`;

    // Store in both the global list and the per-user list, then trim
    // Using a pipeline keeps round-trips to a minimum
    const pipeline = this.storage.pipeline();
    pipeline.lpush(config.allNotificationsKey, serialized);
    pipeline.ltrim(config.allNotificationsKey, 0, MAX_NOTIFICATIONS - 1);
    pipeline.lpush(userKey, serialized);
    pipeline.ltrim(userKey, 0, MAX_NOTIFICATIONS - 1);
    await pipeline.exec();

    logger.info(`Notification stored for user ${userId}`, {
      notificationId: notification.id,
      postId,
      userId,
    });
  }
}

module.exports = EventSubscriber;
