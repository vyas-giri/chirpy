-- +goose Up
CREATE TABLE "user" (
    id uuid PRIMARY KEY,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    email text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE "user";