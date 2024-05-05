CREATE TABLE IF NOT EXISTS grenades(
    id bigserial PRIMARY KEY,
    map varchar(30) NOT NULL,
    title text NOT NULL,
    description text,
    type varchar(30) NOT NULL,
    side varchar(30) NOT NULL,
    version integer NOT NULL DEFAULT 1
);