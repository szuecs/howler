#!/bin/sh
set -e

### BEGIN INIT INFO
# Provides:           howler
# Required-Start:
# Required-Stop:
# Should-Start:
# Should-Stop:
# Default-Start:
# Default-Stop:
# Short-Description:  Howler is an service which is intended to be an endpoint to receive events from the Marathon Event Bus and process them in arbitrary backends.
# Description:
#  https://github.com/zalando-techmonkeys/howler
### END INIT INFO

export PATH=/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin:/opt/bin:

BASE=$(basename $0)

# modify these in /etc/default/$BASE (/etc/default/howler)
HOWLER=/usr/bin/$BASE
# This is the pid file managed by howler itself
HOWLER_PIDFILE=/var/run/$BASE.pid
HOWLER_LOGDIR=/var/log/$BASE/
HOWLER_DESC="Howler: waits to hear something in the Marathon Event Bus and shouts it to the other monkeys"
HOWLER_OPTS=""

# Get lsb functions
. /lib/lsb/init-functions

if [ -f /etc/default/$BASE ]; then
	. /etc/default/$BASE
fi

# see also init_is_upstart in /lib/lsb/init-functions (which isn't available in Ubuntu 12.04, or we'd use it)
if false && [ -x /sbin/initctl ] && /sbin/initctl version 2>/dev/null | grep -q upstart; then
	log_failure_msg "$HOWLER_DESC is managed via upstart, try using service $BASE $1"
	exit 1
fi

# Check howler is present
if [ ! -x $HOWLER ]; then
	log_failure_msg "$HOWLER not present or not executable"
	exit 1
fi

fail_unless_root() {
	if [ "$(id -u)" != '0' ]; then
		log_failure_msg "$HOWLER_DESC must be run as root"
		exit 1
	fi
}

create_logdir() {
    if [ ! -d $HOWLER_LOGDIR ]; then
        mkdir -p $HOWLER_LOGDIR
    fi
}

HOWLER_START="start-stop-daemon \
--start \
--background \
--quiet \
--exec $HOWLER \
--make-pidfile \
--pidfile $HOWLER_PIDFILE \
-- $HOWLER_OPTS"

HOWLER_STOP="start-stop-daemon \
--stop \
--pidfile $HOWLER_PIDFILE"

case "$1" in
	start)
		fail_unless_root
        create_logdir
		log_begin_msg "Starting $HOWLER_DESC: $BASE"
        $HOWLER_START
		log_end_msg $?
		;;

	stop)
		fail_unless_root
		log_begin_msg "Stopping $HOWLER_DESC: $BASE"
        $HOWLER_STOP
		log_end_msg $?
		;;

	restart | force-reload)
		fail_unless_root
        create_logdir
		log_begin_msg "Restarting $HOWLER_DESC: $BASE"
        $HOWLER_STOP
        $HOWLER_START
		log_end_msg $?
		;;

	status)
		status_of_proc -p "$HOWLER_PIDFILE" "$HOWLER" "$HOWLER_DESC"
		;;

	*)
		echo "Usage: $0 {start|stop|restart|status}"
		exit 1
		;;
esac
