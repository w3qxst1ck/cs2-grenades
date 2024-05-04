CREATE TABLE IF NOT EXISTS grenades(
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    description text,
    type text NOT NULL,
    side text NOT NULL,
    version integer NOT NULL DEFAULT 1
);
