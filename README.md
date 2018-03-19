# MemReflect
MemReflect is another implementation on the Memcached killswitch.
Unlike [Memfixed](https://github.com/649/Memfixed-Mitigation-Tool), which uses an active mitigation model, this one uses a passive model.
It sends back killswitch after receiving any UDP packet from 11211 port. This could mitigate the unknown vulnerable memcached servers.
Note: Some server does not support shutdown command, so there's no way to prevent them from being used anyway.
However by using flush_all, the amplify rate is limited.

## Usage
The program automantically sets iptables and routing to receive UDP packets from 11211.
TPROXY module and root permission is required.

### Arguments
-p    The port memreflect listen on (Can be any port rather than 11211)

-s    Use shutdown together with flush_all command

### Build and run
```
go get -t github.com/Max-Sum/memreflect/build
go build -o memreflect github.com/Max-Sum/memreflect/build
sudo ./memreflect -p 11211
```

## Docker
The program would set iptables and routing automantically, but you need to give the capability of net_admin to the docker.
### Tags
`latest` Contains program and source file
`binary` Contains only binary of the program

### Environment
MEMREFLECT_PORT        The port memreflect listen on

MEMREFLECT_SHUTDOWN    Use shutdown together with flush_all command if set

### Run
```
docker run --network=host -d -e MEMREFLECT_PORT=11211 --cap-add net_admin gzmaxsum/memreflect
```
or
```
docker run --network=host -d -e MEMREFLECT_PORT=11211 --privileaged=true gzmaxsum/memreflect
```

