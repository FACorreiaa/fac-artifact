# Deployment Guide - Hetzner Cloud VPS

This guide will help you deploy your Go + Templ portfolio to a Hetzner Cloud VPS with Caddy for automatic HTTPS.

## Prerequisites

- [ ] Hetzner Cloud account
- [ ] Domain name (for HTTPS)
- [ ] SSH key pair

## Cost Estimate

**Hetzner Cloud CX22**: €4.15/month
- 2 vCPU AMD
- 4 GB RAM
- 40 GB SSD
- 20 TB traffic

Perfect for a static portfolio site!

---

## Step 1: Create Hetzner Cloud Server

1. Go to [Hetzner Cloud Console](https://console.hetzner.cloud/)
2. Create new project (e.g., "Portfolio")
3. Click "Add Server"
   - **Location**: Choose closest to your audience (Nuremberg, Helsinki, etc.)
   - **Image**: Ubuntu 24.04 LTS
   - **Type**: CX22 (€4.15/month)
   - **SSH Key**: Add your public key
   - **Name**: portfolio-server
4. Click "Create & Buy Now"
5. Note your server IP address (e.g., `123.45.67.89`)

---

## Step 2: Configure DNS (Point Domain to Server)

Add these DNS records at your domain registrar:

```
Type    Name    Value               TTL
A       @       123.45.67.89        3600
A       www     123.45.67.89        3600
```

Wait 5-10 minutes for DNS propagation.

---

## Step 3: Initial Server Setup

SSH into your server:

```bash
ssh root@YOUR_SERVER_IP
```

### Update system and install dependencies:

```bash
# Update packages
apt update && apt upgrade -y

# Install Caddy
apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt install caddy

# Create application directory
mkdir -p /var/www/goforge
chown -R www-data:www-data /var/www/goforge

# Enable firewall
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
```

---

## Step 4: Configure Caddy

Edit the Caddyfile on your server:

```bash
nano /etc/caddy/Caddyfile
```

Replace contents with (update `yourdomain.com` and email):

```caddyfile
{
    email your-email@example.com
}

yourdomain.com www.yourdomain.com {
    reverse_proxy localhost:8080

    encode gzip zstd

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
        X-XSS-Protection "1; mode=block"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }

    log {
        output file /var/log/caddy/access.log
        format json
    }
}
```

Reload Caddy:

```bash
systemctl reload caddy
```

---

## Step 5: Deploy Application

### Option A: Using Deployment Script (Recommended)

On your **local machine**, edit `deploy/deploy.sh`:

```bash
# Update these variables
SERVER_USER="root"
SERVER_HOST="YOUR_SERVER_IP"
```

Then deploy:

```bash
./deploy/deploy.sh
```

### Option B: Manual Deployment

1. **Build locally:**
   ```bash
   make build
   ```

2. **Upload to server:**
   ```bash
   scp ./bin/server root@YOUR_SERVER_IP:/var/www/goforge/
   rsync -avz ./assets/ root@YOUR_SERVER_IP:/var/www/goforge/assets/
   ```

3. **Install systemd service:**
   ```bash
   scp ./deploy/goforge.service root@YOUR_SERVER_IP:/etc/systemd/system/goforge.service

   ssh root@YOUR_SERVER_IP "systemctl daemon-reload && systemctl enable goforge && systemctl start goforge"
   ```

---

## Step 6: Verify Deployment

Check if services are running:

```bash
ssh root@YOUR_SERVER_IP

# Check Go app
systemctl status goforge

# Check Caddy
systemctl status caddy

# View application logs
journalctl -u goforge -f

# View Caddy logs
journalctl -u caddy -f
```

Visit your domain: `https://yourdomain.com`

✅ You should see your portfolio with automatic HTTPS!

---

## Updating Your Site

After making changes locally:

```bash
# Build and deploy
./deploy/deploy.sh
```

The script will:
1. Build the binary
2. Upload to server
3. Restart the service

---

## Monitoring & Maintenance

### View Logs
```bash
# Application logs
ssh root@YOUR_SERVER_IP 'journalctl -u goforge -f'

# Caddy logs
ssh root@YOUR_SERVER_IP 'journalctl -u caddy -f'
```

### Restart Services
```bash
ssh root@YOUR_SERVER_IP 'systemctl restart goforge'
ssh root@YOUR_SERVER_IP 'systemctl restart caddy'
```

### Check Resource Usage
```bash
ssh root@YOUR_SERVER_IP 'htop'
ssh root@YOUR_SERVER_IP 'df -h'
```

---

## Security Best Practices

### 1. Create Non-Root User (Recommended)

```bash
ssh root@YOUR_SERVER_IP

# Create new user
adduser deploy
usermod -aG sudo deploy

# Copy SSH keys
mkdir -p /home/deploy/.ssh
cp ~/.ssh/authorized_keys /home/deploy/.ssh/
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh
chmod 600 /home/deploy/.ssh/authorized_keys

# Update deploy script to use 'deploy' user instead of 'root'
```

### 2. Disable Root Login

```bash
nano /etc/ssh/sshd_config

# Set:
PermitRootLogin no
PasswordAuthentication no

# Restart SSH
systemctl restart sshd
```

### 3. Install Fail2Ban

```bash
apt install fail2ban -y
systemctl enable fail2ban
```

---

## Troubleshooting

### Port 8080 already in use
```bash
# Find what's using the port
lsof -i :8080

# Kill the process
kill -9 PID
```

### Certificate errors
```bash
# Check Caddy logs
journalctl -u caddy -n 100

# Verify DNS is pointing to server
dig yourdomain.com +short
```

### App won't start
```bash
# Check logs
journalctl -u goforge -n 100

# Check file permissions
ls -la /var/www/goforge

# Test binary directly
cd /var/www/goforge
GO_ENV=production PORT=8080 ./server
```

---

## Cost Summary

| Service | Monthly Cost |
|---------|-------------|
| Hetzner CX22 VPS | €4.15 |
| Domain name | ~€1-10 |
| **Total** | **~€5-15/month** |

SSL certificates are **FREE** via Caddy + Let's Encrypt!

---

## Alternative: Even Cheaper Options

If you want even cheaper hosting:

1. **Hetzner CX11** (€3.79/month) - 1 vCPU, 2GB RAM
   - Fine for low-traffic portfolio

2. **Oracle Cloud Free Tier** - FREE forever
   - 1 vCPU, 1GB RAM ARM instance
   - Limited but truly free

3. **Fly.io** - €0-3/month
   - 1 shared CPU, 256MB RAM free
   - Auto-scaling, global CDN
