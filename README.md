# MemReflect
MemReflect is another implementation on the Memcached killswitch.
Unlike [Memfixed](https://github.com/649/Memfixed-Mitigation-Tool), which uses an active mitigation model, this one uses a passive model.
It sends back killswitch after receiving any UDP packet from 11211 port. This could kill the unknown vulnerable memcached servers.
Your server might be flooded before you send few killswitch, but it does help killing some servers.

## Usage
The program automantically sets iptables and routing to receive UDP packets from 11211.
Therefor, TPROXY module and root permission is required.

### Arguments
-p    The port memreflect listen on
-s    Use shutdown instead of flush_all command

### Build and run
```
go get -t github.com/Max-Sum/memreflect/build
go build -o memreflect github.com/Max-Sum/memreflect/build
sudo ./memreflect -p 11211
```

## Docker
### Environment
MEMREFLECT_PORT        The port memreflect listen on
MEMREFLECT_SHUTDOWN    Use shutdown instead of flush_all command if set

### Run
The program would set iptables and routing automantically, but you need to give the capability of net_admin to the docker.
```
docker run --network=host -d -e MEMREFLECT_PORT=11211 --cap-add net_admin gzmaxsum/memreflect
```
