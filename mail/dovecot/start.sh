#!/bin/sh

# Ensure mail homes root has proper ownership:
chown -c vmail:vmail /var/lib/vmail

exec dovecot -F
