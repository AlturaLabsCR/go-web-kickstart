CREATE VIEW IF NOT EXISTS user_meta AS
SELECT
    u.id,
    u.created_at AS created,
    COALESCE(un.name, '') AS name
FROM users u
LEFT JOIN user_names un ON u.id = un."user";
