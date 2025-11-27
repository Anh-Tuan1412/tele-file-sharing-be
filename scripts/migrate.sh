#!/bin/bash
# Migration script để chạy trên Fly.io
# Sử dụng: flyctl ssh console -> bash scripts/migrate.sh

set -e

echo "🔄 Bắt đầu chạy migrations..."

# Kiểm tra DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
    echo "❌ Lỗi: DATABASE_URL không được thiết lập"
    exit 1
fi

echo "📦 Kết nối database..."

# Chạy từng migration file
MIGRATIONS_DIR="migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "❌ Lỗi: Thư mục $MIGRATIONS_DIR không tìm thấy"
    exit 1
fi

# Chạy migrations theo thứ tự
for migration in $(ls -1 $MIGRATIONS_DIR/*.sql | sort); do
    echo "▶️  Chạy: $(basename $migration)"
    psql "$DATABASE_URL" -f "$migration"
    if [ $? -eq 0 ]; then
        echo "✅ Hoàn thành: $(basename $migration)"
    else
        echo "❌ Lỗi khi chạy: $(basename $migration)"
        exit 1
    fi
done

echo ""
echo "✅ Tất cả migrations đã chạy thành công!"
echo ""
echo "📊 Kiểm tra tables:"
psql "$DATABASE_URL" -c "\dt"

