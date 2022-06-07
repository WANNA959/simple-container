# scadm

scadm, a commond-line tool to control simple-container

**The function of simple-container is implemented by exec.command on Ubuntu-20.04**

## build & installation

```shell
go build -o scadm main.go
cp scadm /usr/sbin/
```

## usage

### help

```shell
$ ./scadm -h

NAME:
   scadm - scadm, a commond-line tool to control simple container

USAGE:
   scadm [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   run             run a container
   create-netns    create network namespace
   connect-pair    connect two netns with veth pair
   connect-bridge  connect to host bridge
   delete-netns    delete network namespace
   create-cgroup   create cgroup
   delete-cgroup   delete cgroup
   set-cgroup      set cgroup limits
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
  
  
$ ./scadm --version
scadm version 0.1.0
go version go1.17.6


$ ./scadm create-netns -h
NAME:
   scadm create-netns - create network namespace

USAGE:
   scadm [global options] create-netns [options]

OPTIONS:
   --name value  name of netns
   --help, -h    show help (default: false)
```

### Network

#### Create veth pair

```shell
$ ./scadm create-netns --name=netns1

------------------------------------------------
simple-container controller:
    netns name: netns1
------------------------------------------------

$ ./scadm create-netns --name=netns2

------------------------------------------------
simple-container controller:
    netns name: netns2
------------------------------------------------

$ ./scadm connect-pair --netns=netns1,netns2 --subnets=10.99.1.1/24,10.99.1.2/24

------------------------------------------------
simple-container controller:
    netns1
        name: netns1
        subnet: 10.99.1.1/24
        veth name: veth4c598a51
    netns2
        name: netns2
        subnet: 10.99.1.2/24
        veth name: veth59e318ea
------------------------------------------------
```

![image-20220513214223391](https://camo.githubusercontent.com/bb271f586d67726e6ef53322fb606206eadd3945d609b12f8e26a6ef862419c1/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c7931683238397862303839706a32313038306e366a76732e6a7067)

#### Create bridge connect

> assign ip for bridge+netns

```shell
$ ./scadm connect-bridge --name=netns1 --subnet=10.0.2.2/24
$ ./scadm connect-bridge --name=netns2 --subnet=10.0.2.3/24
$ ./scadm connect-bridge --name=netns3 --subnet=10.0.2.4/24

------------------------------------------------
simple-container controller:
    master bridge name: master-br0
    master bridge subnet: 10.99.2.1/24
    netns name: netns2
    netns subnet: 10.99.2.3/24
------------------------------------------------
```

- host bridge

![image-20220513223642186](https://camo.githubusercontent.com/a0d7da9d51c6bee92950bc34d5ad44ab9d11d2413247e14a2ff2d6c685df9fe7/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c7931683238397862666e6b356a32313363307371776b352e6a7067)

- Host vs netns1

![image-20220513223540694](https://camo.githubusercontent.com/e77660c30e2b65c707d4bb27b0dbb121b5c53ad6df929d16b8b53a6c0afd407a/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c7931683238397868327079686a32313330306b366164772e6a7067)

- netns1 vs netns2

![image-20220513223621075](https://camo.githubusercontent.com/5a56fe6dc6f44044e2886022b4101da0bac5220afbed8f5eaa9f2ba783c9f828/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c79316832383978646277766c6a32313371306c636a766f2e6a7067)

### Cgroups

```shell
create-cgroup
delete-cgroup
set-cgroup
```

### Container

```shell
$ ./scadm run -it --net=netns3 --limits=cpu.shares=512,cpu.cfs_quota_us=10000,memory.limit_in_bytes=2097152 
```

## Solution

[solution for course report](./solution.md)
