CREATE TYPE "merchantCategory" AS ENUM (
  'SmallRestaurant',
  'MediumRestaurant',
  'LargeRestaurant',
  'MerchandiseRestaurant',
  'BoothKiosk',
  'ConvenienceStore'
);


CREATE TABLE "merchant" (
     "id" uuid PRIMARY KEY,
     "name" varchar,
     "category" "merchantCategory",
     "imageUrl" varchar,
     "latitude" float,
     "longitude" float,
     "createdAt" timestamp
);