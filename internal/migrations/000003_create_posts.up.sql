CREATE TABLE IF NOT EXISTS posts (
                                     id bigserial PRIMARY KEY,
                                     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    description text,
    author_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    type text NOT NULL,
    version integer NOT NULL DEFAULT 1
    );
