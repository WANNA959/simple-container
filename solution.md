# Solution

## 基本要求

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

## 功能增强

> 创建容器无需指定network namespace，自动接入host bridge

```
# run centos-rootfs-based containers, auto join default bridge sc-br0(ref:docker0, auto assigned ip)
$ ./scadm run -it --name=centos1 --limits=cpu.shares=256,cpu.cfs_quota_us=10000,memory.limit_in_bytes=2097152 

$ ./scadm run -it --name=centos2 --limits=cpu.shares=512,cpu.cfs_quota_us=20000,memory.limit_in_bytes=4194304 
```

- host
  - simple-container首次启动的时候，默认创建sc-br0虚拟网卡
  - 使用sqlite轻量级数据库持久化记录sc-br0网段ip分配，确保给各容器分配unique ip
  - 容器关闭的时候clean工作
    - 删除对应的netns或是netns软链接
    - 删除sqlite db中持久化记录

![image-20220520151747738](https://tva1.sinaimg.cn/large/e6c9d24ely1h2ew8oqx0ij211c0cg40g.jpg)

![image-20220520151732647](https://tva1.sinaimg.cn/large/e6c9d24ely1h2ew7t05d0j20zq0kwadr.jpg)

- centos1

![image-20220520151647576](https://tva1.sinaimg.cn/large/e6c9d24ely1h2ew8j4qhdj20yt0u0dl8.jpg)

- centos2

![image-20220520151605945](https://tva1.sinaimg.cn/large/e6c9d24ely1h2ew61u3q9j211s0tu795.jpg)

> docker ps

- simple-container首次启动的时候，默认创建/var/run/container
  - 容器启动时以name为dir，即创建/var/run/container/{containerName}
    - config.json下存储container metadata（ContainerInfo struct），包含container name、pid、container id、创建时间等
      - json序列化方式存储到文件中
  - ./scdm ps命令从/var/run/container下读取各个container的metadata数据
  - 容器关闭的时候clean工作
    - 删除对应/var/run/container/{containerName}下的数据

```
$ ./scadm ps
ID          NAME        PID         STATUS      COMMAND     CREATED
78326536    centos1     8395        Running     unshare     2022-05-20 14:03:40
d820bf22    centos2     8345        Running     unshare     2022-05-20 14:03:05
```

![image-20220520140411177](https://tva1.sinaimg.cn/large/e6c9d24ely1h2eu3vwvynj2148076jsk.jpg)

## todo

- 状态记录 done
  - 轻量级db持久化（如sqlite3）
- 实现类似docker network资源(持久化资源) done
  - 容器创建时加入默认的bridge（docker0） or 手动指定network后自动加入bridge网段，无需创建容器后再搭建bridge
- 功能增强
  - docker ps(done)
  - docker --name(done)
  - 挂载
    - docekr -v 
    - docker images
  - ...