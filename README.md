# Simple metrics collector 

Simple agent that regularly outputs: 

- Load average values
- Derived CPU percentage values
- Network interface statistics
- Disk partition usage in percent
- Memory usage in percent


This agent only work on Linux 2.6 or higher (3.14 or higher for memory
usage) as it uses the `/proc/` interface to query system data. 

## Build
```sh
go build  . 
```

## Run 

The program accepts the following options:
- `-i`: interval in seconds at which to poll
- `-p`: partition to poll
- `-n`: network interface to poll
```sh
./simple-metrics -i 15 -p /mnt/media -n eth0 
```
