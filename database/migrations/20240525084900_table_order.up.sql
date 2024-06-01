CREATE TYPE "orderStatus" AS ENUM (
  'DRAFT',
  'CREATED'
);


CREATE TABLE "order" (
     "orderId" uuid NOT NULL PRIMARY KEY,
     "orderStatus" "orderStatus" NOT NULL,
     detail jsonb NOT NULL,
     "merchantIds" uuid[] NOT NULL,
     "joinedMerchantName" varchar NOT NULL,
     "merchantCategories" "merchantCategory"[] NOT NULL,
     "joinedItemsName" varchar NOT NULL,
     "userId" uuid NOT NULL,
     "userLatitude" DOUBLE PRECISION NOT NULL, -- Separate column for latitude
     "userLongitude" DOUBLE PRECISION NOT NULL, -- Separate column for longitude
     "createdAt" timestamp NOT NULL
);

-- Index on orderId
CREATE INDEX idx_order_order_id ON "order" ("orderId");

-- Index on merchantId
CREATE INDEX idx_order_merchant_ids ON "order" ("merchantIds");

-- Index on merchantCategory
CREATE INDEX idx_order_merchant_categories ON "order" ("merchantCategories");

-- GIN index on joinedMerchantName for fast full-text search
CREATE INDEX idx_order_joined_merchant_name_gin ON "order" USING GIN (to_tsvector('english', "joinedMerchantName"));

-- GIN index on joinedItemsName for fast full-text search
CREATE INDEX idx_order_joined_items_name_gin ON "order" USING GIN (to_tsvector('english', "joinedItemsName"));

-- BTREE index on merchantName for fast text search
CREATE INDEX idx_order_merchant_name_btree ON "order" ("joinedMerchantName");

-- Create a spatial index for efficient spatial queries
CREATE INDEX idx_order_user_location ON "order" USING GIST ((ll_to_earth("userLatitude", "userLongitude")));