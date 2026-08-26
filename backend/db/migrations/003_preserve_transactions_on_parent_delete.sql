-- Historical transactions remain part of the ledger when a scenario or
-- property is removed; only their optional link is cleared.
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_budget_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_budget_id_fkey
    FOREIGN KEY (budget_id) REFERENCES budgets (id) ON DELETE SET NULL;

ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_property_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_property_id_fkey
    FOREIGN KEY (property_id) REFERENCES properties (id) ON DELETE SET NULL;
