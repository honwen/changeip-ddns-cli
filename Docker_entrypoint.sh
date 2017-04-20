#!/bin/sh

cmdArgs="$*"
if [ -n "$cmdArgs" ]; then
  /opt/changeip $cmdArgs
  exit 0
fi

Username=${Username:-1234567890}
Password=${Password:-abcdefghijklmn}
Domain=${Domain:-ddns.changeip.com}
Redo=${Redo:-0}

cat > /opt/supervisord.conf <<EOF
[supervisord]
nodaemon=true

[program:changeip]
command=/opt/changeip --username ${Username} --password ${Password} auto-update --domain ${Domain} --redo ${Redo}
autorestart=true
redirect_stderr=true
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0

EOF

/usr/bin/supervisord -c /opt/supervisord.conf
