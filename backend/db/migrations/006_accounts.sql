ALTER TABLE accounts ADD COLUMN IF NOT EXISTS account_type TEXT NOT NULL DEFAULT 'checking';
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS balance BIGINT NOT NULL DEFAULT 0;
ALTER TABLE accounts ADD COLUMN IF NOT EXISTS credit_limit BIGINT CHECK (credit_limit IS NULL OR credit_limit >= 0);
CREATE INDEX IF NOT EXISTS idx_accounts_user_created ON accounts (user_id, created_at DESC);

ALTER TABLE transactions DROP CONSTRAINT IF EXISTS transactions_account_id_fkey;
ALTER TABLE transactions ADD CONSTRAINT transactions_account_id_fkey
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE SET NULL;
