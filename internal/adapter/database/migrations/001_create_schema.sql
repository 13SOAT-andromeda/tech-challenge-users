-- 001_create_schema.sql
-- usersdb schema: identity and profile domains extracted from garagedb

CREATE TABLE IF NOT EXISTS "Person" (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL,
    email          TEXT NOT NULL,
    contact        TEXT,
    document       TEXT NOT NULL,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    -- Address embedded (vo_address)
    address        TEXT,
    address_number TEXT,
    neighborhood   TEXT,
    city           TEXT,
    country        TEXT,
    zip_code       TEXT,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "User" (
    id         BIGSERIAL PRIMARY KEY,
    password   TEXT NOT NULL,
    role       TEXT NOT NULL CHECK (role IN ('customer', 'attendant', 'mechanic', 'administrator')),
    person_id  BIGINT NOT NULL REFERENCES "Person"(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "Employee" (
    id         BIGSERIAL PRIMARY KEY,
    position   TEXT NOT NULL,
    person_id  BIGINT NOT NULL REFERENCES "Person"(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "Customer" (
    id         BIGSERIAL PRIMARY KEY,
    type       TEXT NOT NULL CHECK (type IN ('pf', 'pj')),
    person_id  BIGINT NOT NULL REFERENCES "Person"(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "Vehicle" (
    id         BIGSERIAL PRIMARY KEY,
    plate      TEXT NOT NULL,
    name       TEXT NOT NULL,
    year       INT  NOT NULL,
    brand      TEXT NOT NULL,
    color      TEXT NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "Customer_Vehicle" (
    id          BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES "Customer"(id),
    vehicle_id  BIGINT NOT NULL REFERENCES "Vehicle"(id),
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS "Company" (
    id             BIGSERIAL PRIMARY KEY,
    name           TEXT NOT NULL,
    email          TEXT NOT NULL,
    document       TEXT NOT NULL,
    contact        TEXT NOT NULL,
    -- Address embedded (vo_address)
    address        TEXT,
    address_number TEXT,
    neighborhood   TEXT,
    city           TEXT,
    country        TEXT,
    zip_code       TEXT,
    created_at     TIMESTAMPTZ,
    updated_at     TIMESTAMPTZ,
    deleted_at     TIMESTAMPTZ
);

-- Partial unique indexes (soft-delete compatible: only enforce on non-deleted rows)
CREATE UNIQUE INDEX IF NOT EXISTS idx_person_email
    ON "Person"(email) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_person_document
    ON "Person"(document) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_vehicle_plate
    ON "Vehicle"(plate) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_cv_customer
    ON "Customer_Vehicle"(customer_id) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_company_document
    ON "Company"(document) WHERE deleted_at IS NULL;

-- Non-unique index for common lookups
CREATE INDEX IF NOT EXISTS idx_user_person_id
    ON "User"(person_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_customer_person_id
    ON "Customer"(person_id) WHERE deleted_at IS NULL;
