CREATE SCHEMA userbalancedb;

SET search_path TO userbalancedb;

CREATE TABLE "users" (
  "user_id" integer PRIMARY KEY,
  "real_balance" numeric(12,3) NOT NULL DEFAULT 0,
  "reserved_balance" numeric(12,3) NOT NULL DEFAULT 0
);

CREATE TABLE "orders" (
  "order_id" integer PRIMARY KEY,
  "user_id" integer NOT NULL,
  "service_id" integer NOT NULL,
  "activated_at" timestamp NOT NULL DEFAULT now()::timestamp,
  "price" numeric(12,3)

  -- PRIMARY KEY (order_id, service_id)
);

CREATE TABLE "payment_history" (
  "order_id" integer,
  "payment_date" timestamp NOT NULL DEFAULT now()::timestamp,
  "amount" numeric(12,3),

  PRIMARY KEY (order_id, payment_date)
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "payment_history" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("order_id");

ALTER TABLE "users" ADD CONSTRAINT "real_balance_check" CHECK (real_balance >= 0);

ALTER TABLE "users" ADD CONSTRAINT "reserved_balance_check" CHECK (reserved_balance >= 0);


set timezone to 'Europe/Moscow';
