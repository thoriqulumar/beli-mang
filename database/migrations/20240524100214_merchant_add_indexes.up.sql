CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE INDEX IF NOT EXISTS "merchant_name_idx" ON merchant USING BTREE ((lower(name)), "name");
CREATE INDEX IF NOT EXISTS "merchant_location_idx" ON merchant USING GIST ((ll_to_earth(latitude, longitude)));
