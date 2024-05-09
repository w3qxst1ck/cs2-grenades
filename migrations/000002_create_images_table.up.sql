CREATE TABLE IF NOT EXISTS images (
    id bigserial PRIMARY KEY,
    name varchar(30) NOT NULL,
    grenade_id bigint NOT NULL REFERENCES grenades on DELETE CASCADE
);