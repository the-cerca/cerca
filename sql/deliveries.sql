CREATE TABLE IF NOT EXISTS deliveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_id INTEGER REFERENCES services(id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    freelance_id UUID REFERENCES users(id) ON DELETE CASCADE,
    price FLOAT NOT NULL,
    "start_date" TIMESTAMP NOT NULL,
    "end_date" TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP )
