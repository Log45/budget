-- Numeric values 3 and 4 previously represented Weekly and Daily. Shift them
-- once to make room for the new Semi-Monthly value at 3. The former combined
-- value at 2 is retained as Bi-Weekly.
CREATE TABLE IF NOT EXISTS app_migrations (
    key TEXT PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM app_migrations WHERE key = '005_split_biweekly_semimonthly') THEN
        UPDATE budgets SET type = type + 1 WHERE type >= 3;
        UPDATE transactions SET recurring_type = recurring_type + 1 WHERE recurring_type >= 3;
        INSERT INTO app_migrations (key) VALUES ('005_split_biweekly_semimonthly');
    END IF;
END $$;
