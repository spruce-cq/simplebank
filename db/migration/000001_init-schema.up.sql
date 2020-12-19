CREATE TABLE "accounts" (
  "id" SERIAL PRIMARY KEY,
  "owner" varchar UNIQUE NOT NULL,
  "balance" int NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" SERIAL PRIMARY KEY,
  "account_id" int NOT NULL,
  "amount" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transfers" (
  "id" SERIAL PRIMARY KEY,
  "from_account_id" int NOT NULL,
  "to_account_id" int NOT NULL,
  "amount" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT 'now()'
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

CREATE INDEX "idx_owner" ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
