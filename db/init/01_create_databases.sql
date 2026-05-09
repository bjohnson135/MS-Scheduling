-- Bootstrap script for the dev MySQL container.
-- Compose mounts ./db/init/*.sql into /docker-entrypoint-initdb.d/, which
-- the official mysql image runs in alphabetical order on first boot only.
-- Subsequent boots reuse the persistent volume; if you need to rerun this,
-- `make reset` drops the volume.

CREATE DATABASE IF NOT EXISTS `account`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

CREATE DATABASE IF NOT EXISTS `company`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_0900_ai_ci;

-- Grant the dev root user access to both. The default mysql:8.4 image already
-- gives root all privileges on all databases, so this is a no-op safeguard
-- in case the image changes.
GRANT ALL PRIVILEGES ON `account`.* TO 'root'@'%';
GRANT ALL PRIVILEGES ON `company`.* TO 'root'@'%';
FLUSH PRIVILEGES;
