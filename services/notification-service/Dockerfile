FROM node:21 AS builder
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci --only=production
COPY . .

FROM node:26-alpine
WORKDIR /app
RUN rm -rf /usr/local/lib/node_modules/npm /usr/local/bin/npm /usr/local/bin/npx
COPY --from=builder /app /app
RUN chown -R node:node /app
USER node
EXPOSE 8083
CMD ["node", "src/index.js"]
