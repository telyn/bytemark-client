#
# Regular cron jobs for the bytemark-client package
#
0 4	* * *	root	[ -x /usr/bin/bytemark-client_maintenance ] && /usr/bin/bytemark-client_maintenance
