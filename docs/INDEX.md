# Documentation Overview

Complete documentation for the Venio project. Start here to find what you need.

> **Documentation Structure:**
> - **[user/](user/)** - End-user documentation (✅ in `main` branch)
> - **[admin/](admin/)** - Administrator guides (✅ in `main` branch)
> - **[dev/](dev/)** - Developer documentation (⚠️ `develop` branch only)

## 🚀 Getting Started

### For End Users
- [User Guide](user/getting-started.md) - How to use Venio
- [FAQ](user/faq.md) - Common questions and answers

### For Administrators
- [Configuration Guide](admin/configuration.md) - Environment variables and settings
- [Deployment Guide](admin/deployment.md) - Docker, Kubernetes, production setup

### For Developers
1. **Windows Users:** Start with [Windows Setup Guide](dev/windows-setup.md)
   - Automated setup with `setup-windows-dev.ps1` (installs everything including GNU Make)
   - Manual step-by-step instructions
   - Troubleshooting guide for common issues

2. **Linux/macOS Users:** See [Development Guide](dev/development.md)
   - Environment setup
   - Running locally
   - Using build tools

3. **Everyone:** Read [Quick Start](#quick-start) below

### Quick Start (5 minutes)

#### Windows
```powershell
# 1. Run automated setup (Administrator)
.\scripts\setup-windows-dev.ps1

# 2. Start services
docker compose up postgres redis -d

# 3. Run backend (Terminal 1)
go run cmd/venio/main.go

# 4. Run frontend (Terminal 2)
cd web
npm run dev

# 5. Access at http://localhost:3000
```

#### Linux/macOS
```bash
# 1. Install dependencies (see development.md)
make install

# 2. Start services
docker compose up postgres redis -d

# 3. Run backend
make dev

# 4. Run frontend
cd web && npm run dev

# 5. Access at http://localhost:3000
```

---

## 📚 Documentation Structure

### 👥 User Documentation ([user/](user/))
Available in both `main` and `develop` branches.

| Document | Audience | Time | Content |
|----------|----------|------|---------|
| **[Getting Started](user/getting-started.md)** | End Users | 10 min | Account setup, dashboard overview, basic usage |
| **[FAQ](user/faq.md)** | End Users | 5 min | Common questions, troubleshooting, getting help |

### 🔧 Administrator Documentation ([admin/](admin/))
Available in both `main` and `develop` branches.

| Document | Audience | Time | Content |
|----------|----------|------|---------|
| **[Configuration Reference](admin/configuration.md)** | Admins/DevOps | 10 min | Environment variables, options, best practices |
| **[Deployment Guide](admin/deployment.md)** | Admins/DevOps | 20 min | Docker, Kubernetes, production setup, monitoring |

### 💻 Developer Documentation ([dev/](dev/))
**Only available in `develop` branch** (not synced to `main`).

#### Essential Reading
| Document | Audience | Time | Content |
|----------|----------|------|---------|
| **[README.md](../README.md)** | Everyone | 5 min | Project overview, features, quick links |
| **[PROJECT_STATUS.md](dev/PROJECT_STATUS.md)** | Developers | 10 min | Complete status, what's done, what's next |
| **[Architecture Overview](dev/architecture.md)** | Developers | 15 min | System design, components, data flow |

#### Setup & Development
| Document | Audience | Time | Content |
|----------|----------|------|---------|
| **[Windows Setup Guide](dev/windows-setup.md)** | Windows Developers | 20-30 min | Automated/manual setup, Make installation, troubleshooting |
| **[Development Guide](dev/development.md)** | All Developers | 20 min | General setup, running locally, configuration |
| **[API Documentation](dev/api.md)** | Backend/Frontend Dev | 15 min | All endpoints, request/response examples, testing |

#### Best Practices & Standards
| Document | Audience | Time | Content |
|----------|----------|------|---------|
| **[Best Practices](dev/best-practices.md)** | All Developers | 30 min | **NEW:** Gin, pgx, JWT, Next.js, PostgreSQL, Redis patterns |
| **[Security Hardening](dev/security-hardening.md)** | All Developers | 30 min | **NEW:** OWASP Top 10, secrets management, container security |
| **[Project Guidelines](dev/project-guidelines.md)** | All Developers + AI | 20 min | Coding standards, commit conventions, AI instructions |
| **[Project Standards](dev/project-standards.md)** | All Developers | 10 min | Tech stack versions, CalVer, dependency policy |
| **[Contributing Guide](../CONTRIBUTING.md)** | Contributors | 10 min | Code of conduct, PR process, development workflow |

---

## 🎯 Find What You Need

### I want to...

#### Use Venio as an end user
- [Getting Started Guide](user/getting-started.md) - Account setup and basic usage
- [FAQ](user/faq.md) - Common questions and troubleshooting

#### Deploy and manage Venio
- [Configuration Guide](admin/configuration.md) - Environment variables and settings
- [Deployment Guide](admin/deployment.md) - Docker, Kubernetes, production

#### Develop Venio features
- **Windows:** [Windows Setup Guide](dev/windows-setup.md) → Run `setup-windows-dev.ps1`
- **Linux/macOS:** [Development Guide](dev/development.md) → Follow "Initial Setup"
- **Quick Start:** [PROJECT_STATUS.md - Quick Start](dev/PROJECT_STATUS.md#quick-start-for-new-developers)

#### Understand how Venio works
1. [Project Overview](../README.md#-features)
2. [Architecture Overview](dev/architecture.md)
3. [API Documentation](dev/api.md)

#### Build/modify the backend (Go)
1. [Best Practices - Go, Gin, pgx](dev/best-practices.md) - **NEW:** Framework-specific patterns
2. [Security Hardening](dev/security-hardening.md) - **NEW:** OWASP Top 10 coverage
3. [Development Guide](dev/development.md#running-locally)
4. [Project Guidelines](dev/project-guidelines.md#development-standards)
5. [Architecture Overview](dev/architecture.md#package-structure)

#### Build/modify the frontend (React/TypeScript)
1. [Best Practices - Next.js, React, TypeScript](dev/best-practices.md) - **NEW:** Frontend patterns
2. [Development Guide](dev/development.md)
3. Navigate to `web/` directory
4. [Project Guidelines - TypeScript section](dev/project-guidelines.md#typescriptreact)

#### Deploy Venio
1. [Deployment Guide](admin/deployment.md)
2. [Configuration Reference](admin/configuration.md)
3. [Architecture Overview](dev/architecture.md#deployment)

#### Understand the project status
- [PROJECT_STATUS.md](dev/PROJECT_STATUS.md) - Everything about what's done, bugs, roadmap

#### Write AI-assisted code
- [Project Guidelines - AI Instructions](dev/project-guidelines.md#ai-assistant-instructions)
- [Testing Guidelines](.github/instructions/testing-guidelines.instructions.md) - **NEW:** Testing best practices
- Follow the [Security Guidelines](dev/project-guidelines.md#security-guidelines)

#### Debug/troubleshoot issues
1. [Windows Setup - Troubleshooting](dev/windows-setup.md#troubleshooting) (Windows)
2. [Development Guide - Configuration](dev/development.md#configuration)
3. [PROJECT_STATUS - Known Issues](dev/PROJECT_STATUS.md#known-issues--workarounds)

---

## 📋 Document Checklist

### Completeness by Document

#### ✅ Complete
- [x] README.md - Project overview, quick start, links
- [x] User Guide - End-user documentation (getting started, FAQ)
- [x] Admin Guide - Configuration and deployment
- [x] Development Guide - All platforms setup, running, config
- [x] Windows Setup Guide - Automated setup, manual steps, troubleshooting
- [x] Architecture Overview - System design, components, interactions
- [x] API Documentation - All endpoints with examples
- [x] **Best Practices** - **NEW:** Gin, pgx, JWT, Next.js, PostgreSQL, Redis patterns
- [x] **Security Hardening** - **NEW:** OWASP Top 10, secrets management, container security
- [x] Project Guidelines - Coding standards, security, AI instructions
- [x] PROJECT_STATUS.md - Complete implementation status
- [x] Project Structure - Folder organization

#### 🟡 Partial
- [ ] Frontend component library documentation
- [ ] Service integration examples

#### 🔴 TODO
- [ ] Database schema diagram (visual)
- [ ] Performance tuning guide
- [ ] Monitoring & logging setup guide

---

## 🔐 Security Documentation

Security is integrated throughout documentation:

- **[Project Guidelines - Security Section](project-guidelines.md#security-guidelines)**
  - Authentication & Authorization
  - Data Protection
  - API Security
  - Code Review Checklist

- **[API Documentation - Security Notes](api.md)**
  - Authentication requirements
  - Authorization checks
  - Input validation

- **[Architecture - Security Layer](architecture.md)**
  - Middleware chain
  - Token validation
  - Error handling

---

## 🛠️ Tools & Scripts Reference

### Automated Setup
- **Windows:** `scripts/setup-windows-dev.ps1` - One-command environment setup
- **All Platforms:** See [Development Guide](development.md#initial-setup)

### Build Tools
- **Makefile:** `make help` - View all available commands
- **PowerShell:** `.\build.ps1 help` - Windows alternative
- See [Development Guide - Build Tools](development.md#build-tools-and-commands)

### Git Hooks
- **Lefthook:** Pre-commit hooks for linting, formatting
- See [Project Guidelines - Commit Messages](project-guidelines.md#commit-messages)

### Development Tools
- **Air:** Hot reload - `air`
- **Delve:** Debugging - `dlv debug ./cmd/venio`
- **golangci-lint:** Linting - `golangci-lint run ./...`
- **goimports:** Formatting - `goimports -w ./internal`

---

## 📖 Reading Order (Recommended)

### For New Contributors
1. [README.md](../README.md) - 5 min
2. [PROJECT_STATUS.md](PROJECT_STATUS.md) - 10 min (focus on MVP section)
3. [Development Guide](development.md) - 20 min
4. [Architecture Overview](architecture.md) - 15 min
5. [Project Guidelines](project-guidelines.md) - 20 min
6. **Total:** ~70 minutes to be fully oriented

### For Code Reviews
1. [Project Guidelines](project-guidelines.md) - Check coding standards
2. [API Documentation](api.md) - Verify endpoints
3. [Security Guidelines](project-guidelines.md#security-guidelines) - Review security aspects

### For Debugging
1. [Configuration Reference](configuration.md) - Verify settings
2. [Architecture Overview](architecture.md) - Understand components
3. [Windows Setup - Troubleshooting](windows-setup.md#troubleshooting) or [Development Guide](development.md)

---

## 🔗 External Links

### Technologies
- [Go Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Next.js Documentation](https://nextjs.org/docs)
- [React Documentation](https://react.dev)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)

### Tools
- [Docker Docs](https://docs.docker.com/)
- [Git Documentation](https://git-scm.com/doc)
- [VSCode Docs](https://code.visualstudio.com/docs)

### Best Practices
- [OWASP Security Guidelines](https://owasp.org/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

## ❓ FAQ

### Where do I find environment variable options?
→ [Configuration Reference](configuration.md)

### How do I report a bug?
→ [Contributing Guide](../CONTRIBUTING.md)

### What should I do before making a commit?
→ [Project Guidelines - Review Checklist](project-guidelines.md#review-checklist)

### How do I debug the backend?
→ [Development Guide - Debugging](development.md#debugging)

### What's the current project status?
→ [PROJECT_STATUS.md](PROJECT_STATUS.md)

### How do I integrate a new service?
→ See Phase 3 in [PROJECT_STATUS.md - Roadmap](PROJECT_STATUS.md#next-steps-roadmap)

---

## 📞 Support

### If you need help
1. Check the relevant documentation section above
2. Search for your issue in [Windows Setup - Troubleshooting](windows-setup.md#troubleshooting)
3. Review [PROJECT_STATUS - Known Issues](PROJECT_STATUS.md#known-issues--workarounds)
4. Ask in GitHub Discussions: https://github.com/lusoris/venio/discussions

### Documentation Issues
- Found outdated information? [Open an issue](https://github.com/lusoris/venio/issues)
- Have a suggestion? [Start a discussion](https://github.com/lusoris/venio/discussions)

---

## 📝 Last Updated

- **Date:** January 14, 2026
- **Status:** Complete & Current
- **Next Update:** As new features are added

---

**Total Documentation:** ~50,000 words across 10 comprehensive guides covering setup, development, architecture, API, security, and AI instructions.
