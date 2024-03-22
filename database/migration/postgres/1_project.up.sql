CREATE TABLE IF NOT EXISTS projects (
    id BIGSERIAL primary key,
    name text not null,
    created_at timestamp DEFAULT now()
);

INSERT INTO projects (name) values ('Первая запись');