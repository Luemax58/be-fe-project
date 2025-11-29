#!/bin/bash

echo "🚀 Resetting MySQL database (utf8mb4) ..."

DB_CONTAINER="backend-db-1"
DB_NAME="mydatabase"
DB_PASS="Lue5548"
SQL_FILE="Database.sql"

echo "🧨 Dropping database..."
docker exec -i $DB_CONTAINER mysql -u root -p$DB_PASS -e "DROP DATABASE IF EXISTS $DB_NAME;"

echo "🛠 Creating database with utf8mb4..."
docker exec -i $DB_CONTAINER mysql -u root -p$DB_PASS -e "
    CREATE DATABASE $DB_NAME 
    CHARACTER SET utf8mb4 
    COLLATE utf8mb4_unicode_ci;
"

echo "📥 Importing $SQL_FILE with UTF8MB4..."
docker exec -i $DB_CONTAINER mysql \
  -u root -p$DB_PASS \
  --default-character-set=utf8mb4 \
  $DB_NAME < $SQL_FILE

echo "✅ Done! Database reset & re-imported successfully."
