#!/bin/bash
# PostgreSQL Backup with Auto Telegram Upload via Bot API
# Prerequisites: Telegram Bot configured in /root/muslimly-be/.telegram.env

set -e

BACKUP_DIR="/root/muslimly-be/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/muslimly_db_$TIMESTAMP.sql.gz"
CONTAINER_NAME="muslimly_db"
RETENTION_DAYS=3
TELEGRAM_CONFIG="/root/muslimly-be/.telegram.env"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Load Telegram config
if [ -f "$TELEGRAM_CONFIG" ]; then
    source "$TELEGRAM_CONFIG"
else
    echo "ERROR: Telegram config not found at $TELEGRAM_CONFIG"
    exit 1
fi

# Load database credentials
if [ -f /root/muslimly-be/.secrets ]; then
    source /root/muslimly-be/.secrets
else
    echo "ERROR: .secrets file not found"
    exit 1
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Step 1: Create database backup
echo "Creating backup..."
docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" "$DB_NAME" | gzip > "$BACKUP_FILE"

# Check if backup was successful
if [ -f "$BACKUP_FILE" ]; then
    SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
    echo -e "${GREEN}✅ Backup created: $BACKUP_FILE ($SIZE)${NC}"
else
    echo -e "${RED}❌ Backup failed!${NC}"
    exit 1
fi

# Step 2: Upload to Telegram via Bot API
echo "Uploading to Telegram..."

MESSAGE="📊 DATABASE BACKUP REPORT
Date: $(date '+%Y-%m-%d %H:%M:%S %Z')
STATUS: ✅ SUCCESS
File: muslimly_db_$TIMESTAMP.sql.gz
Size: $SIZE
Retention: $RETENTION_DAYS days"

# Send message with file attachment
curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendDocument" \
  -F "chat_id=${TELEGRAM_CHAT_ID}" \
  -F "caption=${MESSAGE}" \
  -F "document=@${BACKUP_FILE}" \
  > /tmp/telegram_response.json

# Check if upload was successful
if grep -q '"ok":true' /tmp/telegram_response.json; then
    echo -e "${GREEN}✅ Upload to Telegram successful!${NC}"
else
    echo -e "${RED}❌ Upload to Telegram failed!${NC}"
    exit 1
fi

# Step 3: Cleanup old backups
echo "Cleaning up old backups..."
find "$BACKUP_DIR" -name "muslimly_db_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete
BACKUP_COUNT=$(ls -1 "$BACKUP_DIR"/muslimly_db_*.sql.gz 2>/dev/null | wc -l)
echo "Local backups: $BACKUP_COUNT files"

# Done!
echo ""
echo "=== Backup completed at $(date) ==="
echo -e "${GREEN}📊 STATUS: ✅ OK${NC}"
echo "File sent to Telegram ✓"
