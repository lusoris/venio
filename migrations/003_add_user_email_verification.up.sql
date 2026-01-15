-- Add email verification columns to users table
ALTER TABLE users
ADD COLUMN is_email_verified BOOLEAN NOT NULL DEFAULT FALSE,
ADD COLUMN email_verification_token VARCHAR(255),
ADD COLUMN email_verification_token_expires_at TIMESTAMPTZ,
ADD COLUMN email_verified_at TIMESTAMPTZ;

-- Add index for verification token lookup (partial index for non-null tokens)
CREATE INDEX IF NOT EXISTS idx_users_email_verification_token 
    ON users(email_verification_token) 
    WHERE email_verification_token IS NOT NULL;

-- Add partial index for unverified users (for admin queries)
CREATE INDEX IF NOT EXISTS idx_users_unverified 
    ON users(email) 
    WHERE is_email_verified = FALSE;

-- Add comments for documentation
COMMENT ON COLUMN users.is_email_verified IS 'Whether user has verified their email address';
COMMENT ON COLUMN users.email_verification_token IS 'Token sent via email for verification (32-byte hex)';
COMMENT ON COLUMN users.email_verification_token_expires_at IS 'Expiration timestamp for verification token (typically 24-48 hours)';
COMMENT ON COLUMN users.email_verified_at IS 'Timestamp when email was successfully verified';
