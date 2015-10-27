Cool tricks to be a bytemark badass
===============================

* Output the uptime for all your machines in the "critical" group:
  `for i in $(bytemark list vms critical); do echo "${i%%.*}:"; ssh $i uptime; done`

* Add 10GB of space to each archive grade disk in your "storage" vm:
  `for disc in $(bytemark list discs storage | grep "archive grade"); do bytemark resize disc --size +10G $machine $(awk '{print $2}'); done`
