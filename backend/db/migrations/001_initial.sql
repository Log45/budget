-- Initial schema for the budget app prototype.
-- Amounts are stored in cents (BIGINT) to avoid floating-point errors.

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    username      TEXT        NOT NULL,
    email         TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS accounts (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT      NOT NULL REFERENCES users (id),
    name             TEXT        NOT NULL,
    type             TEXT        NOT NULL DEFAULT 'manual',
    plaid_account_id TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS loans (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT             NOT NULL REFERENCES users (id),
    name            TEXT               NOT NULL,
    principal       BIGINT             NOT NULL,
    current_balance BIGINT             NOT NULL,
    rate            DOUBLE PRECISION   NOT NULL,
    term            INT                NOT NULL,
    start_date      DATE               NOT NULL,
    created_at      TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS budgets (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users (id),
    type       SMALLINT    NOT NULL DEFAULT 1,
    net_income BIGINT      NOT NULL,
    balance    BIGINT      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id          TEXT        PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (id),
    budget_id   BIGINT      REFERENCES budgets (id),
    amount      BIGINT      NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    type        SMALLINT    NOT NULL,
    source      TEXT        NOT NULL DEFAULT '',
    destination TEXT        NOT NULL DEFAULT '',
    date        DATE        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS properties (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT      NOT NULL REFERENCES users (id),
    name            TEXT        NOT NULL,
    address         TEXT        NOT NULL DEFAULT '',
    type            TEXT        NOT NULL DEFAULT 'house',
    purchase_price  BIGINT,
    current_value   BIGINT,
    loan_id         BIGINT      REFERENCES loans (id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts (user_id);
CREATE INDEX IF NOT EXISTS idx_loans_user_id ON loans (user_id);
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (date);
CREATE INDEX IF NOT EXISTS idx_properties_user_id ON properties (user_id);
