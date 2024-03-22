CREATE TABLE IF NOT EXISTS goods (
    id BIGSERIAL,
    project_id BIGSERIAL,
    name text not null,
    description text,
    priority integer not null,
    removed BOOLEAN not null DEFAULT FALSE,
    created_at timestamp DEFAULT now(),
    CONSTRAINT id_project_id_pk PRIMARY KEY (id, project_id)
);

create OR REPLACE function get_priority() returns int language plpgsql as
$$
    Declare
    max_priority integer;
    Begin
        select max(priority) into max_priority from goods;
        IF max_priority is NULL THEN
            max_priority = 0;
        END IF;
        return max_priority + 1;
    End;
$$;

ALTER TABLE goods ALTER COLUMN priority SET DEFAULT get_priority();

CREATE EXTENSION IF NOT EXISTS btree_gin;

CREATE INDEX IF NOT EXISTS idx_name_gin ON goods USING GIN (name);
