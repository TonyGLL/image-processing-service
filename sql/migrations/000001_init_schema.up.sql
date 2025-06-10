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

CREATE TABLE image_processing_schema.images (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE image_processing_schema.images_options (
    id SERIAL PRIMARY KEY,
    resize_width integer NOT NULL,
    resize_height integer NOT NULL,
    crop_width integer NOT NULL,
    crop_height integer NOT NULL,
    crop_x integer NOT NULL,
    crop_y integer NOT NULL,
    rotate integer NOT NULL,
    format varchar(10) NOT NULL,
    grayscale boolean NOT NULL DEFAULT false,
    sepia boolean NOT NULL DEFAULT false,
    image_id integer NOT NULL
);

-- Add foreign keys
ALTER TABLE image_processing_schema.passwords
  ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES image_processing_schema.users (id);

ALTER TABLE image_processing_schema.images_options
  ADD CONSTRAINT fk_image_id FOREIGN KEY (image_id) REFERENCES image_processing_schema.images (id);