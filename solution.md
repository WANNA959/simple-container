# Solution

```shell
# create network namespace
$ ./scadm create-netns --name=netns1
$ ./scadm create-netns --name=netns2

# run centos-rootfs-based containers
$ ./scadm run -it --net=netns2 --limits=cpu.shares=256,cpu.cfs_quota_us=10000,memory.limit_in_bytes=2097152 

$ ./scadm run -it --net=netns3 --limits=cpu.shares=512,cpu.cfs_quota_us=20000,memory.limit_in_bytes=4194304 

# build network bridge 
$ ./scadm connect-bridge --name=netns2 --subnet=10.0.2.3/24
$ ./scadm connect-bridge --name=netns3 --subnet=10.0.2.4/24
```

- 实现进程、用户、文件系统、网络等方面的隔离

```shell
cd /proc/{pid}/ns

ll -h
```

> 进程

```shell
ps -ef
```

> 用户

```shell
echo '1 0 1' > /proc/5850/uid_map
echo '1 0 1' > /proc/5850/gid_map

id
```

> 文件系统

```shell
ls /root
```

> 网络

```shell
ifconfig
ip addr
```

- 能够在ubuntu系统运行centos环境

- 能够实现同一操作系统两个容器之间的网络通信

![image-20220516140501742](https://tva1.sinaimg.cn/large/e6c9d24egy1h2a8h30gghj218b0u0jxe.jpg)

![image-20220516140544819](https://tva1.sinaimg.cn/large/e6c9d24egy1h2a8h5xrufj21iv0u0450.jpg)

- 能够为容器分配定量的CPU和内存资源

> 不同container的cgroup

![image-20220516141754766](https://tva1.sinaimg.cn/large/e6c9d24egy1h2a8hayjhbj21yg0dy78y.jpg)

> demo

![image-20220516142001883](https://tva1.sinaimg.cn/large/e6c9d24egy1h2a8hdu7kjj21tu0m8qa7.jpg)

## todo
- 状态记录（如内存记录subnet ip分配）
  - 轻量级db持久化（如sqlite3）
- 实现类似docker network资源(持久化资源)
  - 容器创建时加入默认的bridge（docker0） or 手动指定network后自动加入bridge网段，无需创建容器后再搭建bridge
- 功能增强
  - docker ps
  - ...