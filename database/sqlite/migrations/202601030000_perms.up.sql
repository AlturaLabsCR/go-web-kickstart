CREATE TABLE IF NOT EXISTS permissions (
  name TEXT PRIMARY KEY,
  description TEXT NOT NULL
);

INSERT INTO permissions (name, description) VALUES
  ('perm.manage_users', 'perm.manage_users.desc'),
  ('perm.change_name', 'perm.change_name.desc')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS roles (
  name TEXT PRIMARY KEY,
  description TEXT NOT NULL
);

INSERT INTO roles (name, description) VALUES
  ('role.admin', 'role.admin.desc'),
  ('role.default', 'role.default.desc')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS role_permissions (
  role TEXT NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
  permission TEXT NOT NULL REFERENCES permissions(name) ON DELETE CASCADE,
  PRIMARY KEY (role, permission)
);

INSERT INTO role_permissions (role, permission) VALUES
  ('role.admin', 'perm.manage_users'),
  ('role.admin', 'perm.change_name'),
  ('role.default', 'perm.change_name')
ON CONFLICT (role, permission) DO NOTHING;

CREATE TABLE IF NOT EXISTS user_roles (
  user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role TEXT NOT NULL REFERENCES roles(name) ON DELETE CASCADE,
  PRIMARY KEY (user, role)
);
