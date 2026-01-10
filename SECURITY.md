# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

### Preferred Method: GitHub Security Advisories

1. Go to the [Security Advisories](https://github.com/lusoris/venio/security/advisories) page
2. Click "Report a vulnerability"
3. Fill out the form with as much detail as possible

### Alternative: Email

Send an email to: **security@example.com** (TODO: Update with actual email)

Please include the following information:
- Type of vulnerability
- Full path to source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the vulnerability

## Response Timeline

- **Initial Response:** Within 48 hours of report
- **Status Update:** Within 7 days with assessment
- **Fix Timeline:** Depends on severity
  - **Critical:** Within 7 days
  - **High:** Within 30 days
  - **Medium/Low:** Next regular release

## Disclosure Policy

- Security issues will be disclosed **after** a fix is available
- We will credit the reporter (unless they prefer to remain anonymous)
- CVE numbers will be requested for significant vulnerabilities
- Advisory will be published on GitHub Security Advisories

## Security Best Practices for Deployment

When deploying Venio in production:

1. **Use HTTPS** - Always use TLS/SSL
2. **Strong Secrets** - Generate strong passwords and API keys
3. **Network Isolation** - Use firewalls and network segmentation
4. **Regular Updates** - Keep Venio and dependencies updated
5. **Access Control** - Implement proper RBAC and OIDC
6. **Audit Logs** - Enable and monitor audit logging
7. **Backup** - Regular backups of databases and configs
8. **Container Security** - Use non-root users, scan images

## Known Security Considerations

### Adult Content Isolation

Venio implements complete isolation for adult content:
- Separate metadata filtering
- Parental controls with PIN
- Module-level access control

Ensure proper parental controls are configured if minors have access.

### API Keys

Venio handles API keys for multiple services (Arrs, Overseerr, etc.):
- Keys are stored encrypted in the database
- Never log API keys
- Rotate keys regularly

### Multi-User Environment

In multi-user deployments:
- Review permission settings carefully
- Use OIDC for centralized authentication
- Enable audit logging
- Regular permission audits

## Security Update Notifications

Subscribe to security updates:
- GitHub Watch → Custom → Security alerts
- GitHub Releases
- Discord Server (coming soon)

## Bug Bounty Program

**Status:** Not currently active

We may introduce a bug bounty program in the future for significant vulnerabilities.

---

**Thank you for helping keep Venio and its users safe!**
