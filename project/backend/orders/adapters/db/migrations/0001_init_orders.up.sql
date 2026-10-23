-- todo: implement
BEGIN;

CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE orders.customers (
    customer_uuid uuid NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    address json NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    PRIMARY KEY (customer_uuid)
);

COMMIT;