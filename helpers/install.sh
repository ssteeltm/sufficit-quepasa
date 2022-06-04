#!/bin/bash

echo UPDATING LOGGING
ln -sf /opt/quepasa-source/helpers/quepasa-syslog.conf /etc/rsyslog.d/10-quepasa.conf

echo UPDATING LOGROTATE
ln -sf /opt/quepasa-source/helpers/quepasa.logrotate.d /etc/logrotate.d/quepasa

/bin/mkdir -p /var/log/quepasa
/bin/chmod 755 /var/log/quepasa
/bin/chown syslog:adm /var/log/quepasa

echo RESTARTING SERVICES
systemctl restart rsyslog

echo UPDATING SYSTEMD SERVICE
ln -sf /opt/quepasa-source/helpers/quepasa.service /etc/systemd/system/quepasa.service
systemctl daemon-reload

adduser --disabled-password --gecos "" -home /opt/quepasa quepasa
chown -R quepasa /opt/quepasa-source

echo UPDATING GO PATHS
mkdir -p /usr/local/go/bin
ln -sf /usr/lib/go/bin/go /usr/local/go/bin/go

cp /opt/quepasa-source/helpers/.env /opt/quepasa/.env

systemctl enable quepasa.service
systemctl start quepasa
exit 0