# Security Policy

## Supported Versions

We actively maintain and support the latest version of knife. Security updates are applied to the latest stable release.

| Version | Supported          |
|---------|--------------------|
| Latest  | :white_check_mark: |
| Older   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in knife, we appreciate your help in disclosing it to us in a responsible manner.

**Please do NOT create a public GitHub issue for security vulnerabilities.**

### How to Report

To report a security vulnerability, please use GitHub's private vulnerability reporting:

1. Visit our security advisory page: <https://github.com/gantries/knife/security/advisories>
2. Click "Report a vulnerability"
3. Fill in the details:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if known)

This will create a private report visible only to the repository maintainers.

### What Happens Next

1. **Confirmation**: We will acknowledge receipt within 48 hours
2. **Verification**: We will verify and assess the vulnerability
3. **Fix Development**: We will develop a fix in a private branch
4. **Coordination**: We will coordinate the disclosure timeline with you
5. **Release**: We will release the security update
6. **Credit**: With your permission, we will credit you in the security advisory

### Security Best Practices

When using knife in production:

1. **Keep dependencies updated** - Regularly update to the latest version
2. **Review dependencies** - Audit your dependencies using `go list -json -m all`
3. **Enable security scanning** - Use tools like `gosec` or `golangci-lint`
4. **Follow secure coding practices** - Validate all inputs, use parameterized queries

### Security Announcements

Security announcements will be posted through:

- GitHub Security Advisories: <https://github.com/gantries/knife/security/advisories>
- Release notes for affected versions

### Dependency Vulnerabilities

We regularly scan our dependencies for known vulnerabilities. If you find a vulnerable dependency:

1. Check if there's an update available
2. Open an issue describing the vulnerability
3. Include the affected dependency and CVE reference

## Security Features

knife includes several security-focused features:

- **SQL Injection Prevention**: ORM layer uses parameterized queries
- **Input Validation**: Expression evaluation with sandboxed environment
- **Authentication**: JWT utilities with proper validation
- **Credential Management**: Secure handling of database credentials

## Contact

For general security questions or concerns:

- GitHub Security Advisories: <https://github.com/gantries/knife/security/advisories>

Thank you for helping keep knife safe!
