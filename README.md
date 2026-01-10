# Venio

<!-- TODO: Add project logo/banner -->

[![Build Status](https://github.com/lusoris/venio/workflows/CI/badge.svg)](https://github.com/lusoris/venio/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lusoris/venio)](https://go.dev/)
[![License](https://img.shields.io/github/license/lusoris/venio)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/lusoris/venio)](https://github.com/lusoris/venio/releases)

> **Unified Media Management System** - A comprehensive orchestration layer for Movies, TV Shows, Music, and Adult content.

Venio is a centralized management system that unifies Overseerr, Lidarr, Whisparr, and your media server into a single, powerful interface with Netflix-like UX, intelligent content lifecycle management, and community-driven features.

## ✨ Features

- **Unified Interface** - Single UI for Movies, TV, Music, and Adult content
- **Multi-User Support** - Advanced RBAC with OIDC integration
- **Smart Request System** - Auto-approval, merging, and community voting
- **Quality Management** - TRaSH Guides integration with dynamic profiles
- **Metadata Enrichment** - Multi-provider aggregation with conflict resolution
- **Content Lifecycle** - Automated retention and archive management
- **Parental Controls** - Complete isolation for adult content
- **Community Features** - Voting, collections, and watch parties

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/lusoris/venio.git
cd venio

# Copy environment template
cp .env.example .env

# Edit .env with your settings
# Then start with Docker Compose
docker compose up -d

# Access Venio at http://localhost:3690
```

## 📚 Documentation

- [Installation Guide](https://github.com/lusoris/venio/wiki/Installation)
- [Configuration Reference](https://github.com/lusoris/venio/wiki/Configuration)
- [User Manual](https://github.com/lusoris/venio/wiki/User-Manual)
- [API Documentation](https://github.com/lusoris/venio/wiki/API)
- [Development Guide](docs/development.md)
- [Architecture Overview](docs/architecture.md)

## 🗺️ Roadmap

### MVP (Current Phase)
- [x] Project Setup & Template
- [ ] Core User Management
- [ ] Overseerr Integration (Movies/TV)
- [ ] Lidarr Integration (Music)
- [ ] Whisparr Integration (Adult)
- [ ] Basic Request System
- [ ] Metadata Proxy

See the full [Roadmap](https://github.com/lusoris/venio/wiki/Roadmap) for planned features.

## 🛠️ Development

```bash
# Install dependencies
make

# Run development environment
make dev

# Run tests
make test

# Run linter
make lint

# Build
make build
```

See [Contributing Guide](CONTRIBUTING.md) for development setup details.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:
- Code of Conduct
- Development setup
- Coding standards
- Pull request process

## 📝 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Overseerr](https://github.com/sct/overseerr) - Media request management inspiration
- [Sonarr/Radarr/Lidarr](https://github.com/Sonarr) - The *arr ecosystem
- [TRaSH Guides](https://trash-guides.info/) - Quality profiles and guides
- All contributors who help make Venio better

## 📞 Support

- [GitHub Discussions](https://github.com/lusoris/venio/discussions) - Questions and community support
- [Issue Tracker](https://github.com/lusoris/venio/issues) - Bug reports and feature requests
- [Discord Server](#) - Real-time chat (coming soon)

---

**Made with ❤️ for the self-hosting community**
