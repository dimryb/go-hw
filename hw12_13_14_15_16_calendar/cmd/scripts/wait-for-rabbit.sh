#!/bin/sh
set -e

echo "Waiting for RabbitMQ at ${RABBIT_HOST}:${RABBIT_PORT}..."

while ! nc -z "${RABBIT_HOST}" "${RABBIT_PORT}"; do
  echo "RabbitMQ is not ready yet..."
  sleep 1
done

echo "RabbitMQ is ready!"