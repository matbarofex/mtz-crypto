CREATE TABLE "wallet_items" (
    "wallet_id" text NOT NULL,
    "symbol" text NOT NULL,
    "quantity" numeric NOT NULL DEFAULT 0.0,
    CONSTRAINT "pk_wallet_items" PRIMARY KEY ("wallet_id", "symbol")
);
