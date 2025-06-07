-- Create schema
CREATE SCHEMA image_processing_schema;

-- Tables creation
CREATE TABLE image_processing_schema.users (
    id SERIAL PRIMARY KEY,
    username varchar(20) UNIQUE NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    deleted boolean NOT NULL DEFAULT false
);

CREATE TABLE image_processing_schema.passwords (
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    user_id integer NOT NULL
);

-- Add foreign keys
ALTER TABLE image_processing_schema.passwords
  ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES image_processing_schema.users (id);