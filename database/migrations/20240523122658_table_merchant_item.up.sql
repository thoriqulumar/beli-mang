CREATE TYPE "itemCategory" AS ENUM (
    'Beverage',
    'Food',
    'Snack',
    'Condiments',
    'Additions'
);


CREATE TABLE "merchantItem" (
     "id" uuid PRIMARY KEY,
     "merchantId" uuid,
     "name" varchar,
     "category" "itemCategory",
     "imageUrl" varchar,
     "price" int,
     "createdAt" timestamp
);