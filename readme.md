# ncadm

scadm, a commond-line tool to control simple container

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
   create-netns    create network namespace
   connect-pair    connect two netns with veth pair
   connect-bridge  connect to host bridge
   delete-netns    delete network namespace
   create-cgroup   create cgroup
   delete-cgroup   delete cgroup
   set-cgroup      set cgroup limits
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ip value        leader host ip (default: "127.0.0.1")
   --port value      network grpc control port (default: "6440")
   --bootport value  network grpc bootstrap control port (default: "6439")
   --cacert value    ca cert filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/ca.pem")
   --cert value      client cert filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/client.pem")
   --key value       client key filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/client-key.pem")
   --help, -h        show help (default: false)
   --version, -v     print the version (default: false)
  
  
  
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

![image-20220513214223391](https://tva1.sinaimg.cn/large/e6c9d24ely1h289xb089pj21080n6jvs.jpg)

#### Create bridge connect

> assign ip for bridge+netns

```shell
$ ./scadm connect-bridge --name=netns1 --subnet=10.0.2.2/24
$ ./scadm connect-bridge --name=netns2 --subnet=10.0.2.3/24

------------------------------------------------
simple-container controller:
    master bridge name: master-br0
    master bridge subnet: 10.99.2.1/24
    netns name: netns2
    netns subnet: 10.99.2.3/24
------------------------------------------------
```

- host bridge

![image-20220513223642186](https://tva1.sinaimg.cn/large/e6c9d24ely1h289xbfnk5j213c0sqwk5.jpg)

- Host vs netns1

![image-20220513223540694](https://tva1.sinaimg.cn/large/e6c9d24ely1h289xh2pyhj21300k6adw.jpg)

- netns1 vs netns2

![image-20220513223621075](https://tva1.sinaimg.cn/large/e6c9d24ely1h289xdbwvlj213q0lcjvo.jpg)
