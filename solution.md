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

![image-20220516140501742](https://camo.githubusercontent.com/d5601f0e3fcf87f66b4a0d3c6f454fdf01d49b2e863176b5b302d324568d1fff/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f6536633964323465677931683261386833306767686a323138623075306a78652e6a7067)

![image-20220516140544819](https://camo.githubusercontent.com/eebb727de80bd325a56dd7ae2b908ec7f2e97f5bcf4ab52ecea3ebbb8560d224/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f6536633964323465677931683261386835787275666a323169763075303435302e6a7067)

- 能够为容器分配定量的CPU和内存资源

> 不同container的cgroup

![image-20220516141754766](https://camo.githubusercontent.com/6d4bcaf579279be3c6a4443217ec92b56e87a3b889b84f4c31db007c1c016906/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f6536633964323465677931683261386861796a68626a323179673064793738792e6a7067)

> demo

![image-20220516142001883](https://camo.githubusercontent.com/b6c034ceae74a5df55560f4ffd76adb50c5c9c6debd744113b4d22025b611b7b/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f653663396432346567793168326138686475376b6a6a32317475306d387161372e6a7067)

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

![image-20220520151747738](https://camo.githubusercontent.com/2ae6dacf4ffea40bbc91af6cbc7c05704b601b7db2c4b8ac36c59896aec0d061/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c793168326577386f717830696a323131633063673430672e6a7067)

![image-20220520151732647](https://camo.githubusercontent.com/6f5ca8832d3ec17c09f9083ac4ab08de7f93f67f5f507ed8c42367ea159c737e/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c7931683265773774303564306a32307a71306b776164722e6a7067)

- centos1

![image-20220520151647576](https://camo.githubusercontent.com/332dac253fde56a46c10d11bcc7cb9cf5baf15adf1f2116cd24b63891a8b0007/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c793168326577386a347168646a32307974307530646c382e6a7067)

- centos2

![image-20220520151605945](https://camo.githubusercontent.com/4b85d2a46f7323a0d2fc1a938d47d76218844c44b936f28d65d5854c1a3ad7aa/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c7931683265773631753371396a323131733074753739352e6a7067)

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

![image-20220520140411177](https://camo.githubusercontent.com/42ef7e3dcd296459b61809bb11eaed80e5d7ce95d7e557003bc27d261a306568/68747470733a2f2f747661312e73696e61696d672e636e2f6c617267652f65366339643234656c79316832657533767776796e6a323134383037366a736b2e6a7067)

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