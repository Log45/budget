-- Budget scenarios, category catalog, and transaction links.
ALTER TABLE budgets ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT 'Untitled budget';
ALTER TABLE budgets ADD COLUMN IF NOT EXISTS start_date DATE;
ALTER TABLE budgets ADD COLUMN IF NOT EXISTS end_date DATE;

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    transaction_type SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT categories_custom_name_unique UNIQUE (user_id, name)
);
CREATE UNIQUE INDEX IF NOT EXISTS categories_default_name_unique ON categories (name) WHERE user_id IS NULL;

INSERT INTO categories (name, transaction_type) VALUES
    ('Groceries', 1), ('Rent', 1), ('Utilities', 1), ('Dining', 1),
    ('Transportation', 1), ('Insurance', 1), ('Healthcare', 1),
    ('Savings', 1), ('Investments', 1), ('Entertainment', 1),
    ('Salary', 0), ('Rental income', 0)
ON CONFLICT DO NOTHING;

CREATE TABLE IF NOT EXISTS budget_categories (
    budget_id BIGINT NOT NULL REFERENCES budgets (id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories (id),
    planned_amount BIGINT NOT NULL DEFAULT 0 CHECK (planned_amount >= 0),
    PRIMARY KEY (budget_id, category_id)
);

ALTER TABLE transactions ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES categories (id);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS property_id BIGINT REFERENCES properties (id);
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS account_id BIGINT REFERENCES accounts (id);
CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions (category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_property_id ON transactions (property_id);
