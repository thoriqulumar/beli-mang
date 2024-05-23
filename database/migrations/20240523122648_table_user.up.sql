CREATE TYPE "role" AS ENUM (
  'admin',
  'user'
);


CREATE TABLE "user" (
     "id" uuid PRIMARY KEY,
     "username" varchar,
     "role" role,
     "email" varchar
     "password" varchar,
     "createdAt" timestamp
);