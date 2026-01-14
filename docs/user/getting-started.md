# Venio User Guide

Welcome to Venio! This guide will help you get started with using Venio to manage your media automation stack.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Creating Your Account](#creating-your-account)
3. [Dashboard Overview](#dashboard-overview)
4. [Managing Users](#managing-users)
5. [Settings](#settings)
6. [Troubleshooting](#troubleshooting)

---

## Getting Started

Venio is a web-based orchestration platform for managing media automation tools (like Sonarr, Radarr, Lidarr, etc.). Access Venio through your web browser at the URL provided by your administrator.

**Default URL:** `http://localhost:3000` (development) or your configured domain in production.

---

## Creating Your Account

### First-Time Registration

1. Navigate to the Venio login page
2. Click **"Create Account"** or **"Register"**
3. Fill in the registration form:
   - **Email:** Your email address (must be unique)
   - **Username:** Choose a username (3-50 characters)
   - **Password:** Create a strong password (minimum 12 characters)
   - **First Name** and **Last Name** (optional)
4. Click **"Register"**
5. Check your email for a verification link (if email verification is enabled)
6. Once verified, log in with your credentials

### Password Requirements

For security, passwords must meet these criteria:
- Minimum 12 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character (!@#$%^&*)

---

## Dashboard Overview

After logging in, you'll see the main dashboard:

### Navigation Menu

- **Home:** Dashboard overview
- **Services:** Connected media automation services (coming soon)
- **Users:** User management (admin only)
- **Settings:** Your account settings
- **Logout:** Sign out of Venio

### Dashboard Widgets

(Coming in future releases)
- Service status indicators
- Recent activity feed
- Quick action buttons
- System notifications

---

## Managing Users

### Viewing Your Profile

1. Click your username in the top-right corner
2. Select **"Profile"**
3. View your account details:
   - Username
   - Email address
   - Account created date
   - Assigned roles

### Updating Your Profile

1. Navigate to **Settings** → **Profile**
2. Update your information:
   - First Name
   - Last Name
   - Email (requires verification)
3. Click **"Save Changes"**

### Changing Your Password

1. Go to **Settings** → **Security**
2. Click **"Change Password"**
3. Enter:
   - Current password
   - New password
   - Confirm new password
4. Click **"Update Password"**

---

## Settings

### Account Settings

- **Profile:** Update your name and contact information
- **Security:** Change password, enable two-factor authentication (coming soon)
- **Notifications:** Configure email and in-app notifications (coming soon)
- **Preferences:** Theme, language, timezone settings (coming soon)

### Privacy Settings

- **Data Export:** Request a copy of your data
- **Account Deletion:** Delete your Venio account permanently

---

## Troubleshooting

### Can't Log In?

**Problem:** "Invalid credentials" error

**Solutions:**
1. Verify your email/username is correct
2. Check for typos in your password
3. Click **"Forgot Password"** to reset
4. Contact your administrator if the account is locked

**Problem:** Account locked after multiple failed attempts

**Solutions:**
1. Wait 15 minutes for automatic unlock
2. Contact your administrator for manual unlock

### Email Not Received?

1. Check your spam/junk folder
2. Verify the email address is correct in your profile
3. Request a new verification email
4. Contact support if issues persist

### Page Not Loading?

1. Refresh the page (F5 or Ctrl+R)
2. Clear your browser cache
3. Try a different browser
4. Check your internet connection
5. Contact your administrator

### Feature Not Working?

1. Check if you have the required permissions
2. Log out and log back in
3. Try a different browser
4. Report the issue to your administrator

---

## Getting Help

### Contact Support

- **Administrator:** Contact your Venio administrator
- **Documentation:** [docs/](../INDEX.md)
- **GitHub Issues:** [Report a bug](https://github.com/lusoris/venio/issues)
- **Community:** [GitHub Discussions](https://github.com/lusoris/venio/discussions)

### Frequently Asked Questions

**Q: How do I reset my password?**
A: Click "Forgot Password" on the login page and follow the instructions sent to your email.

**Q: Can I use Venio on mobile?**
A: Yes! Venio is responsive and works on mobile browsers. A native mobile app may come in the future.

**Q: How do I enable dark mode?**
A: Dark mode is coming soon. For now, Venio uses your system's theme preference.

**Q: What browsers are supported?**
A: Modern versions of Chrome, Firefox, Safari, and Edge are fully supported.

---

## Document Revision

| Version | Date | Author | Changes | Source Version |
|---------|------|--------|---------|----------------|
| 1.0.0 | 2026-01-14 | AI Assistant | Initial user guide | - |

## Referenced Documentation

For more detailed information, see:
- [API Documentation](../dev/api.md) (developers)
- [Configuration Guide](../admin/configuration.md) (administrators)
- [Deployment Guide](../admin/deployment.md) (administrators)
