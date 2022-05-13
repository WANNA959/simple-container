# ncadm

scadm, a commond-line tool to control simple container

## build & installation

```shell
go build -o scadm main.go
cp scadm /usr/sbin/
```

## usage

```shell
$ ./scadm -h

NAME:
   ncadm - ncadm, a commond-line tool to control node join to litekube network-controller

USAGE:
   ncadm [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   create-bootstrap-token  create network bootstrap token info
   get-token               get grpc server ip/port/certs
   check-conn-state        check network conn state
   unregister              close network connection unregister bind ip
   check-health            check health of control and bootstrap grpc
   help, h                 Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --ip value        leader host ip (default: "127.0.0.1")
   --port value      network grpc control port (default: "6440")
   --bootport value  network grpc bootstrap control port (default: "6439")
   --cacert value    ca cert filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/ca.pem")
   --cert value      client cert filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/client.pem")
   --key value       client key filepath of network grpc server (default: "/root/.litekube/nc/certs/grpc/client-key.pem")
   --help, -h        show help (default: false)
   --version, -v     print the version (default: false)
```

