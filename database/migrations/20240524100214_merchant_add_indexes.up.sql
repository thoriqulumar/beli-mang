-- Add indexes for merchant
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE INDEX "merchant_name_idx" ON merchant USING BTREE ((lower(name)), "name");
CREATE INDEX "merchant_location_idx" ON merchant USING GIST ((ll_to_earth(latitude, longitude)));
