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
  service_id INT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP,
  FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);
