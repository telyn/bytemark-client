Cool tricks to be a BigV badass
===============================

* Output the uptime for all your machines in the "critical" group:
  `for i in $(bigv list vms critical); do echo "${i%%.*}:"; ssh $i uptime; done`

* Add 10GB of space to each archive grade disk in your "storage" vm:
  `for disc in $(bigv list discs storage | grep "archive grade"); do bigv resize disc --size +10G $machine $(awk '{print $2}'); done`
