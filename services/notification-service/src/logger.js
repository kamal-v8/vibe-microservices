/**
 * Lightweight structured JSON logger for the Notification Service.
 *
 * Wraps console.log / console.error so that every log line is a
 * machine-parseable JSON object — perfect for log aggregators like
 * Fluentd, Loki, or CloudWatch.
 *
 * Format:
 *   {"level":"info","msg":"...","timestamp":"...","service":"notification-service",...extra}
 */

'use strict';

const SERVICE_NAME = 'notification-service';

/**
 * Build a JSON log line and write it to the appropriate console stream.
 *
 * @param {'info'|'warn'|'error'|'debug'} level
 * @param {string} msg   Human-readable message
 * @param {object} extra Arbitrary key/value pairs merged into the log entry
 */
function log(level, msg, extra = {}) {
  const entry = {
    level,
    msg,
    timestamp: new Date().toISOString(),
    service: SERVICE_NAME,
    ...extra,
  };

  const line = JSON.stringify(entry);

  if (level === 'error') {
    console.error(line);
  } else {
    console.log(line);
  }
}

const logger = {
  info:  (msg, extra = {}) => log('info',  msg, extra),
  warn:  (msg, extra = {}) => log('warn',  msg, extra),
  error: (msg, extra = {}) => log('error', msg, extra),
  debug: (msg, extra = {}) => log('debug', msg, extra),
};

module.exports = logger;
