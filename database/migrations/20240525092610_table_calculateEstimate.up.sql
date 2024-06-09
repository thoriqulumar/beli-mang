CREATE TABLE IF NOT EXISTS "calculatedEstimate" (
      "calculatedEstimateId" UUID NOT NULL PRIMARY KEY,
      "totalPrice" INTEGER NOT NULL,
      "estimatedDeliveryTimeInMinutes" INTEGER NOT NULL,
      "orderId" UUID NOT NULL,
      "createdAt" TIMESTAMP NOT NULL,
      CONSTRAINT fk_calculatedEstimate_orderId
          FOREIGN KEY("orderId")
              REFERENCES "order"("orderId")
              ON DELETE CASCADE
);

