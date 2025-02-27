#!/bin/bash

if [ -z "$DB_USERNAME" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_DATABASE" ]; then
  echo "Ошибка: Не все переменные окружения установлены."
  exit 1
fi

echo "Applying database migrations..."
/migrate -path=/migrations -database="postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable" up

if [ $? -eq 0 ]; then
  echo "Migrations applied successfully."
else
  echo "Failed to apply migrations."
  exit 1
fi

# Запускаем основное приложение
exec "$@"