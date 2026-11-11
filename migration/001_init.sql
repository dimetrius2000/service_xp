create table users(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);

create or replace function get_user_by_email(p_email TEXT)
returns table (id INT, name TEXT, email TEXT)
as $$
BEGIN
    RETURN QUERY
        SELECT u.id,u.name,u.email
        FROM users u
        WHERE u.email = p_email;
END;
$$ LANGUAGE plpgsql;