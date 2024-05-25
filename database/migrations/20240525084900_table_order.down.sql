-- Drop indexes
DROP INDEX IF EXISTS idx_order_order_id;
DROP INDEX IF EXISTS idx_order_merchant_ids;
DROP INDEX IF EXISTS idx_order_merchant_categories;
DROP INDEX IF EXISTS idx_order_joined_merchant_name_gin;
DROP INDEX IF EXISTS idx_order_joined_items_name_gin;
DROP INDEX IF EXISTS idx_order_merchant_name_btree;
DROP INDEX IF EXISTS idx_order_user_location;

-- Drop the order table
DROP TABLE IF EXISTS "order";

-- Drop the orderStatus enum type
DROP TYPE IF EXISTS "orderStatus";

-- Drop the PostGIS extension if not needed anymore
-- (Note: Be cautious with this command as it will remove PostGIS extension from the entire database)
-- DROP EXTENSION IF EXISTS postgis;