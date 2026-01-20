#!/bin/bash
# Deployment script for Hetzner VPS
# Run this on your local machine to deploy to production

set -e

# Configuration
SERVER_USER="root"  # Change to your VPS user
SERVER_HOST="your-server-ip"  # Change to your VPS IP
APP_NAME="fact-artifact"
DEPLOY_PATH="/var/www/${APP_NAME}"

echo "🚀 Starting deployment to ${SERVER_HOST}..."

# 1. Build the application locally
echo "📦 Building production binary..."
make build

# 2. Upload binary and assets to server
echo "📤 Uploading files to server..."
ssh ${SERVER_USER}@${SERVER_HOST} "mkdir -p ${DEPLOY_PATH}/assets"

# Upload binary
scp ./bin/server ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_PATH}/server

# Upload assets (CSS, JS, static files)
rsync -avz --delete ./assets/ ${SERVER_USER}@${SERVER_HOST}:${DEPLOY_PATH}/assets/

# 3. Upload and configure systemd service (first time only)
if ! ssh ${SERVER_USER}@${SERVER_HOST} "systemctl is-active --quiet ${APP_NAME}"; then
    echo "📝 Installing systemd service..."
    scp ./deploy/goforge.service ${SERVER_USER}@${SERVER_HOST}:/etc/systemd/system/${APP_NAME}.service
    ssh ${SERVER_USER}@${SERVER_HOST} "systemctl daemon-reload && systemctl enable ${APP_NAME}"
fi

# 4. Restart the service
echo "🔄 Restarting application..."
ssh ${SERVER_USER}@${SERVER_HOST} "systemctl restart ${APP_NAME}"

# 5. Check status
echo "✅ Checking application status..."
ssh ${SERVER_USER}@${SERVER_HOST} "systemctl status ${APP_NAME} --no-pager"

echo ""
echo "✅ Deployment complete!"
echo "🌐 Your app should be running at http://${SERVER_HOST}:8080"
echo "📊 Check logs with: ssh ${SERVER_USER}@${SERVER_HOST} 'journalctl -u ${APP_NAME} -f'"
