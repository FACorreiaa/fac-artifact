# Deployment Documentation

Complete deployment setup for your Go + Templ portfolio on Hetzner Cloud VPS with Caddy.

## 📚 Documentation Index

| File | Description | When to Use |
|------|-------------|-------------|
| **[HOSTING-COMPARISON.md](./HOSTING-COMPARISON.md)** | Compare Fly.io vs Hetzner vs others | 🤔 Choosing where to host |
| **[FLY-IO.md](./FLY-IO.md)** | Deploy to Fly.io (FREE tier) | 🆓 Fastest & cheapest option |
| **[QUICK-START.md](./QUICK-START.md)** | Deploy to Hetzner in 30 minutes | 🚀 Traditional VPS setup |
| **[DEPLOYMENT.md](./DEPLOYMENT.md)** | Detailed Hetzner guide | 📖 Full reference and troubleshooting |
| **[GITHUB-ACTIONS.md](./GITHUB-ACTIONS.md)** | Automated CI/CD setup | 🤖 Set up auto-deploy on push |

## 🎯 Quick Links

### First Time Setup
👉 **[Start Here: Quick Start Guide](./QUICK-START.md)**

### Files in This Directory

```
deploy/
├── README.md              # This file
├── QUICK-START.md         # 30-min deployment guide
├── DEPLOYMENT.md          # Detailed documentation
├── GITHUB-ACTIONS.md      # Auto-deploy setup
├── deploy.sh              # Deployment script
├── goforge.service        # Systemd service file
└── ../Caddyfile           # Caddy configuration
```

## 📋 Deployment Checklist

- [ ] **Prerequisites**
  - [ ] Domain name purchased
  - [ ] SSH key generated
  - [ ] Hetzner Cloud account created

- [ ] **Server Setup** (see [QUICK-START.md](./QUICK-START.md))
  - [ ] VPS created (CX22, €4.15/mo)
  - [ ] DNS configured (A records)
  - [ ] Caddy installed
  - [ ] Firewall configured

- [ ] **First Deployment**
  - [ ] Update `deploy.sh` with server IP
  - [ ] Run `make deploy`
  - [ ] Verify at https://yourdomain.com

- [ ] **Optional: Auto-Deploy** (see [GITHUB-ACTIONS.md](./GITHUB-ACTIONS.md))
  - [ ] GitHub secrets configured
  - [ ] Test auto-deploy on push

## 🚀 Common Commands

```bash
# First time setup
make deploy-setup          # Show setup instructions

# Deploy to production
make deploy                # Build and deploy

# Manual steps
make build                 # Build production binary
ssh root@SERVER 'systemctl restart goforge'  # Restart app
ssh root@SERVER 'journalctl -u goforge -f'   # View logs
```

## 💰 Hosting Options Comparison

| Provider | Plan | Monthly Cost | Setup Complexity |
|----------|------|--------------|------------------|
| **Hetzner Cloud** | CX22 | €4.15 | ⭐⭐ Medium |
| Hetzner Cloud | CX11 | €3.79 | ⭐⭐ Medium |
| Oracle Cloud | Free Tier | €0 | ⭐⭐⭐ Complex |
| Fly.io | Hobby | €0-3 | ⭐ Easy |
| Railway | Hobby | $5 | ⭐ Easy |
| Render | Free | €0 | ⭐ Easy |

**Recommendation:** Hetzner CX22 for best value + performance.

## 🛠️ Architecture

```
Internet
    ↓
Caddy (Port 443/80)
    ↓ HTTPS ↓
    ↓ Automatic Let's Encrypt
    ↓ Reverse Proxy
    ↓
Go App (localhost:8080)
    ↓
Serves Templ templates + static assets
```

**Benefits:**
- ✅ Automatic HTTPS certificates
- ✅ HTTP/2 and HTTP/3 support
- ✅ Gzip compression
- ✅ Security headers
- ✅ Zero-downtime reloads

## 📊 Performance Expectations

**Hetzner CX22 (€4.15/mo):**
- Handles: ~10,000+ requests/day
- Response time: <50ms
- Uptime: 99.9%+

Perfect for portfolios and small business sites!

## 🔐 Security Features

- ✅ HTTPS by default (Let's Encrypt)
- ✅ Firewall (UFW)
- ✅ Security headers (HSTS, XSS, etc.)
- ✅ systemd hardening (NoNewPrivileges, PrivateTmp)
- ✅ Non-root execution (www-data user)

## 📈 Next Steps After Deployment

### Essential
- [ ] Set up monitoring (UptimeRobot, Uptime.com)
- [ ] Configure backups (Hetzner snapshots)
- [ ] Add analytics (Plausible, Umami)

### Optional
- [ ] CDN (Cloudflare free tier)
- [ ] Email notifications (health checks)
- [ ] Staging environment
- [ ] Database (PostgreSQL)

## 🆘 Need Help?

### Common Issues

**Site not loading?**
1. Check DNS: `dig yourdomain.com +short`
2. Check app: `ssh root@SERVER 'systemctl status goforge'`
3. Check Caddy: `ssh root@SERVER 'systemctl status caddy'`

**HTTPS not working?**
1. Wait 2-3 minutes for Let's Encrypt
2. Check Caddy logs: `journalctl -u caddy -n 50`
3. Verify email in Caddyfile is correct

**Deployment fails?**
1. Check `deploy.sh` has correct IP
2. Verify SSH key access: `ssh root@SERVER`
3. Check disk space: `ssh root@SERVER 'df -h'`

### Documentation

- [Hetzner Cloud Docs](https://docs.hetzner.com/cloud/)
- [Caddy Documentation](https://caddyserver.com/docs/)
- [Templ Documentation](https://templ.guide/)
- [Full Deployment Guide](./DEPLOYMENT.md)

## 💡 Pro Tips

1. **Use staging first**: Test on a CX11 (€3.79) before production
2. **Enable Hetzner backups**: €0.80/mo for daily snapshots
3. **Monitor your app**: Set up Uptime Robot (free tier)
4. **Use GitHub Actions**: Auto-deploy on every push to main
5. **Cloudflare free tier**: Add CDN + DDoS protection

## 📞 Support

- 📧 Hetzner Support: https://console.hetzner.cloud/support
- 💬 Caddy Community: https://caddy.community/
- 🐛 Project Issues: GitHub Issues tab

---

**Happy Deploying! 🚀**
