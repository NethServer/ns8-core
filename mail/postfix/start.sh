#!/bin/sh

envsubst '${MAILBOX_DOMAIN}' </etc/postfix/main.cf.template >/etc/postfix/main.cf

# Refresh file databases on startup
newaliases

exec /usr/sbin/postfix start-fg
