# Quick Deployment Guide 🚀

## TL;DR - Get Your Portfolio Online in 30 Minutes

### Prerequisites
- [ ] Domain name (Namecheap, Cloudflare, etc.)
- [ ] SSH key generated (`ssh-keygen -t ed25519`)

---

## Step 1: Create Hetzner VPS (5 min)

1. Go to https://console.hetzner.cloud/
2. New Project → Add Server
3. Choose:
   - **Image**: Ubuntu 24.04
   - **Type**: CX22 (€4.15/mo)
   - **Location**: Closest to you
   - **SSH Key**: Upload `~/.ssh/id_ed25519.pub`
4. Copy your server IP: `123.45.67.89`

---

## Step 2: Point Domain to Server (2 min)

At your domain registrar, add:

```
Type: A     Name: @      Value: YOUR_SERVER_IP
Type: A     Name: www    Value: YOUR_SERVER_IP
```

---

## Step 3: Setup Server (10 min)

SSH into server:
```bash
ssh root@YOUR_SERVER_IP
```

Run setup commands:
```bash
# Update system
apt update && apt upgrade -y

# Install Caddy
apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update && apt install -y caddy

# Create app directory
mkdir -p /var/www/goforge
chown -R www-data:www-data /var/www/goforge

# Setup firewall
ufw allow 22/tcp && ufw allow 80/tcp && ufw allow 443/tcp
ufw --force enable
```

Configure Caddy:
```bash
nano /etc/caddy/Caddyfile
```

Paste this (replace `yourdomain.com` and email):
```caddyfile
{
    email your-email@example.com
}

yourdomain.com www.yourdomain.com {
    reverse_proxy localhost:8080
    encode gzip zstd
}
```

Save and reload:
```bash
systemctl reload caddy
```

---

## Step 4: Deploy Your App (10 min)

**On your local machine:**

1. Update deployment config:
```bash
nano deploy/deploy.sh
```

Change:
```bash
SERVER_USER="root"
SERVER_HOST="YOUR_SERVER_IP"  # Your actual IP
```

2. Deploy:
```bash
make deploy
```

---

## Step 5: Verify (1 min)

Visit: `https://yourdomain.com`

✅ **Done!** Your portfolio is live with HTTPS!

---

## Daily Workflow

### Make changes and deploy:
```bash
# Edit your .templ files
# ...

# Deploy updates
make deploy
```

### View logs:
```bash
ssh root@YOUR_SERVER_IP 'journalctl -u goforge -f'
```

---

## Troubleshooting

### Site not loading?
```bash
# Check if app is running
ssh root@YOUR_SERVER_IP 'systemctl status goforge'

# Check Caddy
ssh root@YOUR_SERVER_IP 'systemctl status caddy'
```

### HTTPS not working?
- Wait 2-3 minutes for Let's Encrypt
- Check DNS: `dig yourdomain.com +short` (should show your IP)
- Check Caddy logs: `journalctl -u caddy -n 50`

---

## Cost Breakdown

| Item | Cost |
|------|------|
| Hetzner CX22 | €4.15/mo |
| Domain | ~€10/year |
| SSL Cert | FREE (Caddy) |
| **Total** | **~€5/mo** |

---

## Next Steps

- [ ] Set up GitHub Actions for auto-deploy
- [ ] Add analytics (Plausible, Umami)
- [ ] Configure backups
- [ ] Add monitoring (Uptime Robot)

📖 Full docs: `deploy/DEPLOYMENT.md`
