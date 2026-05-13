-- Up Migration
CREATE TABLE registration (
  id         UUID PRIMARY KEY,
  email      TEXT NOT NULL UNIQUE,
  full_name  TEXT NOT NULL,
  phone      TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT email_length     CHECK (char_length(email) BETWEEN 3 AND 254),
  CONSTRAINT email_format     CHECK (email ~ '^[^[:space:]@]+@[^[:space:]@]+\.[^[:space:]@]+$'),
  CONSTRAINT email_lower      CHECK (email = lower(email)),
  CONSTRAINT full_name_length CHECK (char_length(full_name) BETWEEN 2 AND 100),
  CONSTRAINT phone_length     CHECK (char_length(phone) BETWEEN 7 AND 20),
  CONSTRAINT phone_format     CHECK (phone ~ '^[+0-9 ()\-]+$')
);

-- Down Migration
DROP TABLE registration;
