-- Remove role-permission assignments
DELETE FROM role_permissions
WHERE role_id IN (SELECT id FROM roles WHERE name IN ('admin', 'moderator', 'user', 'guest'));

-- Remove default permissions
DELETE FROM permissions
WHERE name IN (
  'users:read', 'users:write', 'users:delete',
  'roles:read', 'roles:write', 'roles:delete',
  'permissions:read', 'permissions:write', 'permissions:delete',
  'content:read', 'content:write', 'content:delete', 'content:moderate',
  'settings:read', 'settings:write',
  'audit:read'
);

-- Remove default roles
DELETE FROM roles WHERE name IN ('admin', 'moderator', 'user', 'guest');
