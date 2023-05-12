CREATE TABLE if not exists email_verification (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL ,
    token VARCHAR(100) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
