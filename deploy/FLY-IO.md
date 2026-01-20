# Deploy to Fly.io (FREE Tier)

Deploy your portfolio to Fly.io's free tier in under 10 minutes.

## Why Fly.io?

- ✅ **$0/month** - Free tier is generous
- ✅ **5 minute setup** - Easiest deployment
- ✅ **Global CDN** - Fast everywhere
- ✅ **Auto HTTPS** - Free SSL certificates
- ✅ **Auto-scaling** - Sleeps when idle, wakes instantly

## Prerequisites

- [ ] GitHub account
- [ ] Credit card (for verification, won't be charged)

---

## Step 1: Install Fly CLI (2 min)

### macOS:
```bash
brew install flyctl
```

### Linux:
```bash
curl -L https://fly.io/install.sh | sh
```

### Windows:
```powershell
pwsh -Command "iwr https://fly.io/install.ps1 -useb | iex"
```

---

## Step 2: Sign Up & Login (2 min)

```bash
# Sign up (opens browser)
fly auth signup

# Or login if you have an account
fly auth login
```

Add a credit card when prompted (for verification, **free tier won't charge**).

---

## Step 3: Configure Your App (1 min)

Edit `fly.toml`:

```toml
app = "my-portfolio"  # ⚠️ Change this to your unique app name
primary_region = "fra"  # Options: fra, ams, lhr, iad, syd, etc.
```

**Choose your region:**
- `fra` - Frankfurt (Europe)
- `ams` - Amsterdam (Europe)
- `lhr` - London (UK)
- `iad` - Virginia (US East)
- `lax` - Los Angeles (US West)
- `syd` - Sydney (Australia)
- `nrt` - Tokyo (Japan)

---

## Step 4: Deploy! (5 min)

From your project root:

```bash
# Launch app (first time)
fly launch --no-deploy

# Follow prompts:
# - Choose app name: my-portfolio
# - Choose region: fra (or closest to you)
# - Setup Postgres? NO
# - Setup Redis? NO
# - Deploy now? NO (we'll do it manually)

# Now deploy
fly deploy
```

**That's it!** ✨

Your app will be live at: `https://my-portfolio.fly.dev`

---

## Step 5: Add Custom Domain (Optional)

### Add your domain:
```bash
fly certs add yourdomain.com
fly certs add www.yourdomain.com
```

### Configure DNS at your registrar:

**For root domain (`yourdomain.com`):**
```
Type: A
Name: @
Value: (get from: fly ips list)
```

**For www subdomain:**
```
Type: CNAME
Name: www
Value: my-portfolio.fly.dev
```

**Get your IP addresses:**
```bash
fly ips list
```

Add both IPv4 and IPv6 addresses as A and AAAA records.

### Verify:
```bash
fly certs show yourdomain.com
```

Wait 5-10 minutes for SSL provisioning.

---

## Common Commands

```bash
# Deploy updates
fly deploy

# View logs
fly logs

# Check status
fly status

# SSH into machine
fly ssh console

# Scale (if you outgrow free tier)
fly scale memory 512  # Upgrade to 512MB

# View dashboard
fly dashboard
```

---

## Updating Your Site

Make changes and deploy:

```bash
# Edit your .templ files
# ...

# Deploy
fly deploy
```

**Deployment takes ~2 minutes.**

---

## Free Tier Limits

Fly.io free tier includes:

✅ **Included Free:**
- Up to 3 shared-cpu-1x VMs (256MB RAM)
- 3GB persistent storage
- 160GB outbound data transfer/month

**What this means for you:**
- Perfect for portfolio/CV site
- Handles ~100K+ pageviews/month
- Auto-sleeps when idle (wakes in <1 second)

**If you exceed free tier:**
- You'll get an email warning
- Automatic charges start at ~$2/month

---

## Performance Optimization

### 1. Enable Fly CDN (Free)
Already enabled! Your static assets are cached globally.

### 2. Add Multiple Regions (Paid)
```bash
# Add US East region
fly regions add iad

# Add Europe region
fly regions add ams
```

**Cost:** +$1.94/month per additional region

### 3. Prevent Sleep (Paid)
```bash
# Keep 1 machine always running
fly scale count 1 --max-per-region 1
```

Edit `fly.toml`:
```toml
min_machines_running = 1  # Keep awake
```

**Cost:** +$1.94/month

---

## Cost Calculator

| Configuration | Monthly Cost |
|---------------|--------------|
| **Free tier** (256MB, sleeps) | **$0** |
| Keep awake (256MB, 1 region) | $1.94 |
| Upgraded (512MB, 1 region) | $3.88 |
| Multi-region (256MB, 3 regions) | $5.82 |

**For a portfolio: FREE tier is perfect!**

---

## Monitoring

### View real-time logs:
```bash
fly logs -a my-portfolio
```

### View metrics:
```bash
fly dashboard
```

Or visit: https://fly.io/dashboard

### Set up alerts:
1. Go to https://fly.io/dashboard
2. Click your app
3. Go to "Monitoring"
4. Enable email alerts

---

## Comparing Deployment Methods

| Method | Setup Time | Monthly Cost | Complexity |
|--------|------------|--------------|------------|
| **Fly.io** | 10 min | $0 | ⭐ Easy |
| Hetzner + Caddy | 30 min | €4.15 | ⭐⭐ Medium |
| DigitalOcean | 30 min | $4-6 | ⭐⭐ Medium |

---

## Troubleshooting

### Deployment fails

**Check Docker build locally:**
```bash
docker build -t test .
docker run -p 8080:8080 test
```

Visit http://localhost:8080 to test.

### App crashes on startup

**View logs:**
```bash
fly logs
```

**Common issues:**
- Missing environment variables
- Port not set to 8080
- Health check failing

### Out of memory

**Check memory usage:**
```bash
fly vm status
```

**Upgrade to 512MB:**
```bash
fly scale memory 512
```

This moves you to paid tier ($3.88/month).

### Slow first request (cold start)

**This is normal for free tier** - app sleeps after inactivity.

**Solutions:**
1. Accept it (portfolio visitors won't notice <1s delay)
2. Keep awake with `min_machines_running = 1` (+$1.94/mo)
3. Use uptime monitoring (pings every 5 min to keep warm)

---

## Migration from Fly.io to Hetzner (Later)

If you outgrow Fly.io:

1. **Build locally:**
   ```bash
   make build
   ```

2. **Deploy to Hetzner:**
   ```bash
   ./deploy/deploy.sh
   ```

3. **Update DNS** to point to Hetzner IP

4. **Delete Fly.io app:**
   ```bash
   fly apps destroy my-portfolio
   ```

---

## Recommendation

**Start with Fly.io free tier** because:
- ✅ Zero cost to test
- ✅ 5-minute setup
- ✅ Easy to migrate later
- ✅ Perfect for portfolios

**Upgrade to Hetzner when:**
- You need a database (PostgreSQL)
- You have consistent high traffic
- You want more control
- You add real-time features

---

## Next Steps

After deploying:

- [ ] Add Google Analytics / Plausible
- [ ] Set up Uptime monitoring (UptimeRobot)
- [ ] Add meta tags for social sharing
- [ ] Submit to Google Search Console
- [ ] Add robots.txt and sitemap.xml

---

**Deploy now:**
```bash
fly launch
fly deploy
```

🎉 **Your portfolio will be live in 5 minutes!**
