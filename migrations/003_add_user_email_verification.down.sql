-- Remove indexes (must be done before dropping columns)
DROP INDEX IF EXISTS idx_users_email_verification_token;
DROP INDEX IF EXISTS idx_users_unverified;

-- Remove email verification columns from users table
ALTER TABLE users
DROP COLUMN IF EXISTS is_email_verified,
DROP COLUMN IF EXISTS email_verification_token,
DROP COLUMN IF EXISTS email_verification_token_expires_at,
DROP COLUMN IF EXISTS email_verified_at;
