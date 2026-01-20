# GitHub Actions Auto-Deploy Setup

Automatically deploy your portfolio when you push to `main` branch.

## Setup Instructions

### 1. Generate SSH Key for GitHub Actions

On your **local machine**:

```bash
# Generate a new SSH key pair for deployments
ssh-keygen -t ed25519 -C "github-actions@deploy" -f ~/.ssh/github_actions_deploy

# Display the private key (you'll add this to GitHub)
cat ~/.ssh/github_actions_deploy

# Display the public key (you'll add this to your server)
cat ~/.ssh/github_actions_deploy.pub
```

### 2. Add Public Key to Server

SSH into your server:

```bash
ssh root@YOUR_SERVER_IP
```

Add the public key:

```bash
# Add to authorized_keys
echo "YOUR_PUBLIC_KEY_HERE" >> ~/.ssh/authorized_keys

# Verify
cat ~/.ssh/authorized_keys
```

### 3. Add Secrets to GitHub Repository

1. Go to your GitHub repository
2. Click **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**

Add these 3 secrets:

| Secret Name | Value |
|-------------|-------|
| `SSH_PRIVATE_KEY` | Contents of `~/.ssh/github_actions_deploy` (the **private** key) |
| `SERVER_HOST` | Your server IP (e.g., `123.45.67.89`) |
| `SERVER_USER` | Usually `root` or `deploy` |

**Example:**

```
SSH_PRIVATE_KEY:
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
...
(full private key contents)
...
-----END OPENSSH PRIVATE KEY-----

SERVER_HOST:
123.45.67.89

SERVER_USER:
root
```

### 4. Enable GitHub Actions

The workflow file is already in `.github/workflows/deploy.yml`.

It will automatically run when you:
- Push to `main` branch
- Manually trigger from GitHub Actions tab

### 5. Test the Workflow

Make a small change and push:

```bash
# Make a change
echo "# Test" >> README.md

# Commit and push
git add .
git commit -m "Test auto-deploy"
git push origin main
```

Go to GitHub → **Actions** tab → Watch the deployment!

---

## Workflow Overview

The GitHub Actions workflow does:

1. ✅ Checkout code
2. ✅ Setup Go + Templ + Tailwind
3. ✅ Build production binary
4. ✅ Upload to server via SSH
5. ✅ Restart systemd service
6. ✅ Verify deployment

**Build time:** ~1-2 minutes

---

## Manual Deployment vs Auto-Deploy

### Manual (using `make deploy`):
- ✅ Full control
- ✅ Build on your machine
- ❌ Requires manual action

### Auto (GitHub Actions):
- ✅ Automatic on push
- ✅ Consistent build environment
- ✅ Build logs in GitHub
- ❌ Requires GitHub secrets setup

**Recommendation:** Use both!
- Auto-deploy for quick changes
- Manual deploy for testing before push

---

## Troubleshooting

### Deployment fails with SSH error

**Check secrets:**
```bash
# In GitHub repo → Settings → Secrets
# Verify:
# - SSH_PRIVATE_KEY has the full private key (including header/footer)
# - SERVER_HOST is just the IP (no ssh:// prefix)
# - SERVER_USER is correct (usually 'root' or 'deploy')
```

### Build fails

**Check Go version:**
The workflow uses Go 1.23. Update `.github/workflows/deploy.yml` if needed:

```yaml
- name: Setup Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.23'  # Update this if needed
```

### Service restart fails

**Check permissions:**

```bash
ssh root@YOUR_SERVER_IP

# Ensure the user can restart services
# If using non-root user, add to sudoers:
echo "deploy ALL=(ALL) NOPASSWD: /bin/systemctl restart goforge" >> /etc/sudoers.d/deploy
```

---

## Monitoring Deployments

### View GitHub Actions logs:
1. Go to repo → **Actions**
2. Click on the workflow run
3. View detailed logs

### View server logs after deploy:
```bash
ssh root@YOUR_SERVER_IP 'journalctl -u goforge -n 50'
```

---

## Security Best Practices

### ✅ DO:
- Use a dedicated SSH key for GitHub Actions
- Limit key permissions on server
- Use GitHub's encrypted secrets
- Use a non-root user for deployments (optional but recommended)

### ❌ DON'T:
- Commit private keys to repository
- Use your personal SSH key
- Share secrets in commit messages
- Give unnecessary sudo permissions

---

## Advanced: Deploy to Multiple Environments

You can extend the workflow for staging/production:

```yaml
on:
  push:
    branches:
      - main        # Production
      - develop     # Staging

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      # ...
      - name: Set environment
        run: |
          if [ "${{ github.ref }}" == "refs/heads/main" ]; then
            echo "ENV=production" >> $GITHUB_ENV
            echo "SERVER_HOST=${{ secrets.PROD_SERVER_HOST }}" >> $GITHUB_ENV
          else
            echo "ENV=staging" >> $GITHUB_ENV
            echo "SERVER_HOST=${{ secrets.STAGING_SERVER_HOST }}" >> $GITHUB_ENV
          fi
```

Add `PROD_SERVER_HOST` and `STAGING_SERVER_HOST` secrets.

---

## Next Steps

- [ ] Set up deployment notifications (Slack, Discord)
- [ ] Add health checks after deployment
- [ ] Implement blue-green deployments
- [ ] Add automatic rollback on failure
