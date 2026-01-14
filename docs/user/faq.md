# Frequently Asked Questions (FAQ)

Common questions and answers about using Venio.

## General Questions

### What is Venio?

Venio is an orchestration platform for managing media automation services like Sonarr, Radarr, Lidarr, and more. It provides a unified interface for controlling and monitoring your entire media stack.

### Is Venio free?

Yes, Venio is open-source software released under the MIT License. You can use, modify, and distribute it freely.

### What operating systems does Venio support?

Venio runs on:
- Linux (recommended for production)
- macOS (development)
- Windows (development and production with Docker)
- Any platform that supports Docker containers

---

## Account & Authentication

### How do I create an account?

1. Navigate to the Venio login page
2. Click "Register" or "Create Account"
3. Fill in your details and submit
4. Verify your email (if enabled)
5. Log in with your credentials

See [Getting Started Guide](getting-started.md#creating-your-account) for detailed steps.

### I forgot my password. How do I reset it?

1. Click "Forgot Password" on the login page
2. Enter your email address
3. Check your email for a reset link
4. Follow the link and create a new password

### My account is locked. What should I do?

Accounts are temporarily locked after 5 failed login attempts. Wait 15 minutes for automatic unlock, or contact your administrator for immediate assistance.

### Can I change my username?

Currently, usernames cannot be changed after registration. Please choose carefully during signup.

### Can I change my email address?

Yes, you can update your email address in **Settings** → **Profile**. You'll need to verify the new email address.

---

## Features & Capabilities

### What services can Venio manage?

(Coming in future releases)
- Sonarr (TV shows)
- Radarr (Movies)
- Lidarr (Music)
- Readarr (Books)
- Prowlarr (Indexers)
- Overseerr/Jellyseerr (Requests)
- Media servers (Plex, Jellyfin, Emby)

### Can I use Venio on mobile?

Yes! Venio's web interface is fully responsive and works on mobile browsers. A native mobile app may be developed in the future.

### Does Venio support dark mode?

Dark mode support is coming soon. Currently, Venio respects your system's theme preference.

### Can I customize the dashboard?

Dashboard customization features are planned for future releases. Stay tuned!

---

## Security & Privacy

### How does Venio protect my data?

Venio uses industry-standard security practices:
- Passwords hashed with bcrypt (cost 12)
- JWT tokens with secure HTTP-only cookies
- HTTPS encryption (in production)
- SQL injection prevention
- Rate limiting on sensitive endpoints
- Regular security audits

See [Security Hardening Guide](../dev/security-hardening.md) for technical details.

### Who can see my data?

- **You:** Full access to your profile and activity
- **Administrators:** Can view user accounts and manage permissions
- **Other users:** Cannot see your data unless explicitly shared

### Can I export my data?

Yes, you can request a data export in **Settings** → **Privacy** → **Export Data**. You'll receive a JSON file with all your account information.

### How do I delete my account?

1. Go to **Settings** → **Privacy**
2. Click **"Delete Account"**
3. Confirm the deletion
4. Your account and associated data will be permanently removed

**Warning:** Account deletion is irreversible!

---

## Technical Questions

### What technologies does Venio use?

**Backend:**
- Go 1.25
- Gin web framework
- PostgreSQL 18.1 database
- Redis 8.4 caching

**Frontend:**
- Next.js 15
- React 19
- TypeScript 5.7
- Tailwind CSS 4

### What browsers are supported?

Venio supports modern browsers:
- ✅ Chrome/Edge 120+
- ✅ Firefox 120+
- ✅ Safari 17+
- ❌ Internet Explorer (not supported)

### Can I self-host Venio?

Yes! Venio is designed for self-hosting. See the [Deployment Guide](../admin/deployment.md) for instructions.

### What are the system requirements?

**Minimum:**
- 1 CPU core
- 512 MB RAM
- 1 GB disk space
- Docker support

**Recommended:**
- 2+ CPU cores
- 2 GB RAM
- 10 GB disk space
- SSD storage

### Does Venio support multiple users?

Yes! Venio has full multi-user support with role-based access control (RBAC). Administrators can create users, assign roles, and manage permissions.

---

## Troubleshooting

### Venio won't load in my browser

Try these steps:
1. Refresh the page (F5)
2. Clear browser cache (Ctrl+Shift+Delete)
3. Try incognito/private mode
4. Try a different browser
5. Check if the server is running
6. Contact your administrator

### I get a "Connection refused" error

This means the Venio server isn't running or isn't reachable:
1. Verify the server is started
2. Check the URL is correct
3. Ensure your firewall allows connections
4. Contact your administrator

### Features are missing or disabled

Some features may require specific permissions:
1. Check your assigned roles in **Profile**
2. Contact an administrator to request permissions
3. Some features may be in development (coming soon)

### The page is slow or unresponsive

1. Check your internet connection
2. Try a different browser
3. Clear browser cache
4. Contact your administrator (server may be overloaded)

---

## Getting More Help

### Where can I find more documentation?

- **User Guide:** [getting-started.md](getting-started.md)
- **Admin Guide:** [../admin/deployment.md](../admin/deployment.md)
- **Developer Docs:** [../dev/](../dev/)
- **Main Index:** [../INDEX.md](../INDEX.md)

### How do I report a bug?

1. Check if it's already reported: [GitHub Issues](https://github.com/lusoris/venio/issues)
2. If not, create a new issue with:
   - Description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots (if applicable)
   - Browser/OS version

### How do I request a feature?

1. Check existing requests: [GitHub Discussions](https://github.com/lusoris/venio/discussions)
2. Create a new discussion with:
   - Feature description
   - Use case
   - Why it would be valuable

### Where can I get community support?

- **GitHub Discussions:** [lusoris/venio/discussions](https://github.com/lusoris/venio/discussions)
- **GitHub Issues:** [lusoris/venio/issues](https://github.com/lusoris/venio/issues)
- **Discord:** (coming soon)

---

## Document Revision

| Version | Date | Author | Changes | Source Version |
|---------|------|--------|---------|----------------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial FAQ | - |

## Referenced Documentation

- [User Guide](getting-started.md)
- [Admin Guide](../admin/configuration.md)
- [Security Guide](../dev/security-hardening.md)
