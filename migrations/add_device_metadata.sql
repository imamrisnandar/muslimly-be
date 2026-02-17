-- Migration: Add device metadata to user_devices table
-- Created: 2026-02-17
-- Purpose: Add device info and location columns for better analytics

ALTER TABLE user_devices 
ADD COLUMN IF NOT EXISTS device_model VARCHAR(100),
ADD COLUMN IF NOT EXISTS device_os_version VARCHAR(50),
ADD COLUMN IF NOT EXISTS app_version VARCHAR(20),
ADD COLUMN IF NOT EXISTS country_code VARCHAR(10),
ADD COLUMN IF NOT EXISTS timezone VARCHAR(50);

-- Optional: Add indexes for common queries
CREATE INDEX IF NOT EXISTS idx_user_devices_country ON user_devices(country_code);
CREATE INDEX IF NOT EXISTS idx_user_devices_platform_version ON user_devices(platform, device_os_version);
