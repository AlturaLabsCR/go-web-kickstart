-- name: GetRoles :many
SELECT r.name
FROM roles r
JOIN user_roles ur ON r.name = ur.role
WHERE ur.user = $1;

-- name: GetPermissions :many
SELECT DISTINCT p.name
FROM permissions p
JOIN role_permissions rp ON p.name = rp.permission
JOIN user_roles ur ON rp.role = ur.role
WHERE ur.user = $1;

-- name: SetRole :exec
INSERT INTO user_roles ("user", role)
VALUES ($1, $2)
ON CONFLICT ("user", role) DO NOTHING;
