ALTER TABLE transactions ADD COLUMN IF NOT EXISTS recurring BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE transactions ADD COLUMN IF NOT EXISTS recurring_type SMALLINT NOT NULL DEFAULT 1;

-- A deleted budget must not delete an expense; it simply becomes unassigned.
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_budget_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_budget_id_fkey
    FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE SET NULL;
