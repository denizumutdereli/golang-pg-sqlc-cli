-- CREATE TYPE "ORDER_STATUS" AS ENUM (
--   'Open',
--   'Filled',
--   'Cancelled',
--   'Invistigate',
--   'Error'
-- );

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" decimal NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "token_id" bigint NOT NULL,
  "amount" decimal NOT NULL,
  "side" varchar NOT NULL DEFAULT 'BUY',
  "price" decimal NOT NULL,
  "status" varchar NOT NULL DEFAULT 'OPEN',
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" decimal NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "tokens" (
  "id" bigserial PRIMARY KEY,
  "shortName" varchar NOT NULL,
  "name" decimal NOT NULL,
  "details_etc" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "orders" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

CREATE INDEX ON "tokens" ("shortName");

CREATE INDEX ON "tokens" ("name");

COMMENT ON COLUMN "orders"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "orders"."price" IS 'can not be negative';

COMMENT ON COLUMN "transfers"."amount" IS 'it must be positive';

ALTER TABLE "orders" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("token_id") REFERENCES "tokens" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
