# Hosting Comparison for Go + Templ Portfolio

Complete comparison of hosting options for your static portfolio site.

## TL;DR - Quick Recommendation

**For Portfolio/CV Site:** Use **Fly.io Free Tier** ($0/month)

---

## Detailed Comparison

### 1. Fly.io 🥇 (RECOMMENDED)

| Feature | Free Tier | Paid Tier |
|---------|-----------|-----------|
| **Cost** | **$0/month** | $1.94+/month |
| **RAM** | 256MB | 256MB - 8GB |
| **Storage** | 3GB | Unlimited |
| **Traffic** | 160GB/month | Unlimited |
| **Setup Time** | 10 minutes | 10 minutes |
| **HTTPS** | ✅ Auto (free) | ✅ Auto (free) |
| **CDN** | ✅ Global | ✅ Global |
| **Scaling** | Auto (0-3 VMs) | Custom |
| **Domains** | Unlimited | Unlimited |

**Pros:**
- ✅ **FREE** for small sites
- ✅ Fastest deployment (one command)
- ✅ Global CDN included
- ✅ Auto HTTPS
- ✅ Auto-scaling (sleeps when idle)
- ✅ Built for Go apps

**Cons:**
- ❌ Sleeps after inactivity (~1s wake time)
- ❌ Limited to 256MB RAM on free tier
- ❌ Requires credit card for signup

**Perfect for:**
- ✅ Portfolios
- ✅ CV/Resume sites
- ✅ Landing pages
- ✅ Low-traffic sites (<100K views/month)

**Setup:** See [FLY-IO.md](./FLY-IO.md)

---

### 2. Hetzner Cloud VPS 🥈

| Feature | CX11 | CX22 (Recommended) |
|---------|------|---------------------|
| **Cost** | €3.79/month | **€4.15/month** |
| **RAM** | 2GB | 4GB |
| **vCPU** | 1 AMD | 2 AMD |
| **Storage** | 20GB SSD | 40GB SSD |
| **Traffic** | 20TB | 20TB |
| **Setup Time** | 30 minutes | 30 minutes |
| **HTTPS** | ✅ Caddy (free) | ✅ Caddy (free) |
| **CDN** | ❌ No (single location) | ❌ No |

**Pros:**
- ✅ **Best value** - 4GB RAM for €4
- ✅ Always on (no sleep)
- ✅ Full VPS control
- ✅ Can add PostgreSQL/Redis later
- ✅ EU-based (GDPR compliant)
- ✅ Predictable performance

**Cons:**
- ❌ Costs money (€4.15/month)
- ❌ More setup required
- ❌ Single location (not global)
- ❌ Manual server management

**Perfect for:**
- ✅ High-traffic sites
- ✅ Need database
- ✅ Want full control
- ✅ Consistent load

**Setup:** See [DEPLOYMENT.md](./DEPLOYMENT.md)

---

### 3. DigitalOcean 🥉

| Feature | Basic ($4) | Basic ($6) |
|---------|------------|------------|
| **Cost** | $4/month | $6/month |
| **RAM** | 512MB | 1GB |
| **vCPU** | 1 | 1 |
| **Storage** | 10GB | 25GB |
| **Traffic** | 500GB | 1TB |
| **Setup Time** | 30 minutes | 30 minutes |

**Pros:**
- ✅ Good documentation
- ✅ Managed databases available
- ✅ Simple UI
- ✅ App Platform option (PaaS)

**Cons:**
- ❌ **Worse value** than Hetzner
- ❌ Same setup complexity as Hetzner
- ❌ More expensive for same specs

**Verdict:** Hetzner offers better value

---

### 4. Oracle Cloud Free Tier 💎

| Feature | Always Free |
|---------|-------------|
| **Cost** | **$0/month FOREVER** |
| **RAM** | 1GB |
| **vCPU** | 1 ARM |
| **Storage** | 50GB |
| **Traffic** | 10TB/month |
| **Setup Time** | 60+ minutes |

**Pros:**
- ✅ **FREE FOREVER**
- ✅ Very generous limits
- ✅ Full VPS control
- ✅ Great for learning

**Cons:**
- ❌ **Complex signup** (often rejects cards)
- ❌ ARM architecture (different from x86)
- ❌ Can terminate accounts without warning
- ❌ Confusing UI
- ❌ Not reliable for production

**Verdict:** Great if you can sign up, but risky

---

### 5. Railway

| Feature | Free | Hobby |
|---------|------|-------|
| **Cost** | $0 | $5/month |
| **Free Credits** | $5/month | Unlimited |
| **Setup Time** | 5 minutes | 5 minutes |

**Pros:**
- ✅ Very easy deployment
- ✅ $5 free credits monthly
- ✅ Good for testing

**Cons:**
- ❌ Free credits run out quickly
- ❌ $5/month for production
- ❌ Less generous than Fly.io

**Verdict:** Use Fly.io instead (better free tier)

---

### 6. Render

| Feature | Free | Starter |
|---------|------|---------|
| **Cost** | $0 | $7/month |
| **RAM** | 512MB | 512MB |
| **Setup Time** | 5 minutes | 5 minutes |

**Pros:**
- ✅ Easy deployment
- ✅ Auto HTTPS
- ✅ Good UI

**Cons:**
- ❌ Free tier **spins down** after 15min inactivity
- ❌ **Slow cold starts** (30+ seconds)
- ❌ $7/month to stay awake
- ❌ More expensive than alternatives

**Verdict:** Use Fly.io instead (better free tier)

---

## Side-by-Side Comparison

| Provider | Monthly Cost | RAM | Always On? | Setup | Global CDN | Best For |
|----------|--------------|-----|------------|-------|------------|----------|
| **Fly.io** | **$0** | 256MB | ❌ (1s wake) | ⭐ Easy | ✅ Yes | 🏆 Portfolios |
| **Hetzner CX22** | **€4.15** | 4GB | ✅ Yes | ⭐⭐ Medium | ❌ No | 🏆 High traffic |
| Hetzner CX11 | €3.79 | 2GB | ✅ Yes | ⭐⭐ Medium | ❌ No | Budget VPS |
| DigitalOcean | $4-6 | 512MB-1GB | ✅ Yes | ⭐⭐ Medium | ❌ No | - |
| Oracle Free | $0 | 1GB | ✅ Yes | ⭐⭐⭐ Hard | ❌ No | Experimental |
| Railway | $5 | 512MB | ✅ Yes | ⭐ Easy | ❌ No | - |
| Render | $7 | 512MB | ✅ Yes | ⭐ Easy | ❌ No | - |

---

## My Recommendation by Use Case

### Portfolio / CV Site (Your Case)
**Winner: Fly.io Free Tier** ($0/month)

**Why:**
- Zero cost
- Auto-scaling
- Global CDN
- Perfect for sporadic traffic

**Deploy with:**
```bash
make deploy-fly
```

### High Traffic / Business Site
**Winner: Hetzner CX22** (€4.15/month)

**Why:**
- Always on
- 4GB RAM
- Can add database
- Best value

**Deploy with:**
```bash
make deploy
```

### Experimental / Learning
**Winner: Oracle Cloud Free Tier** ($0/month)

**Why:**
- Free forever
- Full VPS
- Learn server management

### Need Database
**Winner: Hetzner CX22** (€4.15/month)

**Why:**
- Enough RAM for PostgreSQL
- Full control
- No vendor lock-in

---

## Cost Projection (1 Year)

| Provider | Setup | Monthly | Annual Total |
|----------|-------|---------|--------------|
| **Fly.io Free** | $0 | $0 | **$0** ✨ |
| Fly.io Paid | $0 | $1.94 | $23.28 |
| **Hetzner CX22** | $0 | €4.15 | **~$54** |
| Hetzner CX11 | $0 | €3.79 | ~$49 |
| DigitalOcean | $0 | $4-6 | $48-72 |
| Railway Hobby | $0 | $5 | $60 |
| Render Starter | $0 | $7 | $84 |

**Add domain:** ~$10-15/year (any registrar)

---

## Performance Comparison

### Load Time (First Visit)

| Provider | Cold Start | Warm Response |
|----------|------------|---------------|
| **Fly.io** | ~800ms | <50ms |
| **Hetzner** | N/A | <30ms |
| DigitalOcean | N/A | <40ms |
| Render Free | **30+ seconds** | <100ms |

**Note:** Fly.io sleeps = fast wake. Render sleeps = slow wake.

---

## Migration Path

### Start Free → Scale Later

**Phase 1: Launch (Month 1-6)**
- Use **Fly.io free tier**
- $0 cost while validating

**Phase 2: Growth (Month 6-12)**
- Upgrade to **Hetzner CX22** if needed
- Add PostgreSQL database
- €4.15/month

**Phase 3: Scale (Year 2+)**
- Stay on Hetzner or
- Move to managed services (Render, Railway)
- ~$20-50/month

---

## Final Verdict

### For YOUR Portfolio/CV Site:

**Start with Fly.io Free Tier** ✅

**Reasons:**
1. $0 cost (can't beat free!)
2. 5-minute setup
3. Global CDN (fast everywhere)
4. Easy to migrate later
5. Perfect for portfolio traffic

### Upgrade to Hetzner when:
- ❌ You add a database
- ❌ Constant high traffic
- ❌ Need more control
- ❌ Want no cold starts

---

## Quick Deploy Commands

### Deploy to Fly.io (FREE):
```bash
fly launch
fly deploy
```

### Deploy to Hetzner (€4.15/mo):
```bash
./deploy/deploy.sh
```

---

## Questions to Ask Yourself

**Choose Fly.io if:**
- ✅ Budget is $0
- ✅ Traffic is sporadic
- ✅ Don't need database (yet)
- ✅ Want easiest setup

**Choose Hetzner if:**
- ✅ Need database
- ✅ High consistent traffic
- ✅ Want full control
- ✅ Prefer traditional VPS
- ✅ €4/month is okay

**Choose DigitalOcean if:**
- ✅ Already familiar with DO
- ❌ (Otherwise, use Hetzner - better value)

---

## My Personal Recommendation

**For a portfolio site in 2026:**

1. **Start:** Fly.io free tier ($0)
2. **Monitor:** Track traffic for 3 months
3. **Upgrade:** Only if you consistently exceed free tier
4. **Database?** Then migrate to Hetzner CX22 (€4.15)

**Most portfolios never need to upgrade from free tier.**

---

📖 **Detailed guides:**
- [Fly.io Deployment](./FLY-IO.md)
- [Hetzner Deployment](./DEPLOYMENT.md)
- [GitHub Actions Auto-Deploy](./GITHUB-ACTIONS.md)
