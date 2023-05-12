-- Starter beta 
CREATE EXTENSION "uuid-ossp";

CREATE TABLE IF NOT EXISTS categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS image_services (
  id SERIAL PRIMARY KEY,
  filename VARCHAR(255) NOT NULL,
  path VARCHAR(255) NOT NULL,
  alt_text VARCHAR(255) DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS services (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  starting_at float NOT NULL,
  description VARCHAR(200) NOT NULL,
  image_service_id INT,
  FOREIGN KEY (image_service_id) REFERENCES image_services(id) ON DELETE CASCADE,
  category_id INT REFERENCES categories(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  firstname varchar(50) NULL DEFAULT '',
  lastname varchar(50) NULL DEFAULT '',
  username varchar(50) UNIQUE NOT NULL,
  password_hashed text NOT NULL,
  email text NOT NULL UNIQUE,
  birthday date,
  is_verified bool DEFAULT false,
  is_freelance bool DEFAULT false,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP

);
CREATE TABLE IF NOT EXISTS email_verification (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id uuid NOT NULL,
  token VARCHAR(100) NOT NULL UNIQUE,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS deliveries (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL REFERENCES services(id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    freelance_id UUID REFERENCES users(id) ON DELETE CASCADE,
    price FLOAT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    complete bool DEFAULT false, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID,
  token VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
