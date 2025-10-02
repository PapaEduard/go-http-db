#!/bin/bash
set -e

BACKUP_DIR="./db_backups"
DB_CONTAINER="go-http-db_db_1"
DB_NAME="app_db"
DB_USER="postgres"

mkdir -p $BACKUP_DIR

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/backup_${DB_NAME}_$TIMESTAMP.sql"

docker exec $DB_CONTAINER pg_dump -U $DB_USER $DB_NAME > $BACKUP_FILE

echo "Backup completed: $BACKUP_FILE"
