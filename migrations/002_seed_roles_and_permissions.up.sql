-- Insert default roles
INSERT INTO roles (name, description, created_at) VALUES
  ('admin', 'Administrator with full access', NOW()),
  ('moderator', 'Moderator with moderation capabilities', NOW()),
  ('user', 'Regular user with basic access', NOW()),
  ('guest', 'Guest with read-only access', NOW())
ON CONFLICT (name) DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (name, description, created_at) VALUES
  ('users:read', 'Read user information', NOW()),
  ('users:write', 'Create and edit users', NOW()),
  ('users:delete', 'Delete users', NOW()),
  ('roles:read', 'Read roles', NOW()),
  ('roles:write', 'Create and edit roles', NOW()),
  ('roles:delete', 'Delete roles', NOW()),
  ('permissions:read', 'Read permissions', NOW()),
  ('permissions:write', 'Create and edit permissions', NOW()),
  ('permissions:delete', 'Delete permissions', NOW()),
  ('content:read', 'Read content', NOW()),
  ('content:write', 'Create and edit content', NOW()),
  ('content:delete', 'Delete content', NOW()),
  ('content:moderate', 'Moderate content', NOW()),
  ('settings:read', 'Read application settings', NOW()),
  ('settings:write', 'Modify application settings', NOW()),
  ('audit:read', 'Read audit logs', NOW())
ON CONFLICT (name) DO NOTHING;

-- Assign permissions to admin role (all permissions)
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign permissions to moderator role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'moderator'
  AND p.name IN ('users:read', 'content:read', 'content:moderate', 'audit:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign permissions to user role
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'user'
  AND p.name IN ('users:read', 'content:read', 'content:write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign permissions to guest role (read-only)
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT r.id, p.id, NOW()
FROM roles r, permissions p
WHERE r.name = 'guest'
  AND p.name IN ('users:read', 'content:read', 'settings:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;
