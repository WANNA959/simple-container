package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"simple-container/pkg/utils"
	"strings"
)

const (
	RootUrl       string = "/root/.simple-container/images"
	MntUrl        string = "/root/.simple-container/mnt/"
	WriteLayerUrl string = "/root/.simple-container/writeLayer"
	ReadLayerUrl  string = "/root/.simple-container/readLayer"
	WorkLayerUrl  string = "/root/.simple-container/workLayer"
)

// 创建一个Overlay系统作为容器的根目录
func NewWorkSpace(volume string, imageName string, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateWorkLayer(containerName)
	CreateMountPoint(containerName, imageName)

	// 判断是否执行挂载数据卷操作
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs, containerName)
			log.Printf("NewWorkSpace volume urls %q", volumeURLs)
		} else {
			log.Printf("Volume parameter input is not correct.")
		}
	}
}

// 解析volume字符串
func volumeUrlExtract(volume string) []string {
	var volumeURLs []string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs
}

// 创建只读层
func CreateReadOnlyLayer(imageName string) error {
	unTarFolderUrl := filepath.Join(ReadLayerUrl, imageName)
	imageUrl := filepath.Join(RootUrl, imageName+".tar")
	exist := utils.Exists(unTarFolderUrl)
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl, 0622); err != nil {
			log.Fatalf("Mkdir read layer dir %s error. %v", unTarFolderUrl, err)
			return err
		}
		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			log.Fatalf("Untar read dir %s error %v", unTarFolderUrl, err)
			return err
		}
	}
	return nil
}

// 创建读写层
func CreateWriteLayer(containerName string) {
	writeURL := filepath.Join(WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Fatalf("Mkdir write layer dir %s error. %v", writeURL, err)
	}
}

// 创建工作层
func CreateWorkLayer(containerName string) {
	workURL := filepath.Join(WorkLayerUrl, containerName)
	if err := os.MkdirAll(workURL, 0777); err != nil {
		log.Fatalf("Mkdir work layer dir %s error. %v", workURL, err)
	}
}

// 挂载数据卷
func MountVolume(volumeURLs []string, containerName string) error {
	// 在宿主机创建宿主目录
	parentUrl := volumeURLs[0]
	if err := os.MkdirAll(parentUrl, 0777); err != nil {
		log.Printf("Mkdir parent dir %s error. %v", parentUrl, err)
	}
	CreateVolumePoint(parentUrl)

	// 在容器里创建挂载目录
	containerUrl := volumeURLs[1]
	mntURL := filepath.Join(MntUrl, containerName)
	containerVolumeURL := filepath.Join(mntURL, containerUrl)
	if err := os.MkdirAll(containerVolumeURL, 0777); err != nil {
		log.Printf("Mkdir container dir %s error. %v", containerVolumeURL, err)
	}

	// 把宿主机文件目录挂载到容器挂载点
	lower := "lowerdir=" + parentUrl + "/readLayer/"
	upper := "upperdir=" + parentUrl + "/writeLayer/"
	work := "workdir=" + parentUrl + "/workLayer/"
	parm := lower + "," + upper + "," + work
	_, err := exec.Command("mount", "-t", "overlay", "overlay", "-o", parm, containerVolumeURL).CombinedOutput()
	if err != nil {
		log.Fatalf("Mount volume failed. %v", err)
		return err
	}
	return nil
}

// 创建宿主机挂载卷的read/write.work目录
func CreateVolumePoint(parentUrl string) {
	readURL := parentUrl + "/readLayer/"
	writeURL := parentUrl + "/writeLayer/"
	workURL := parentUrl + "/workLayer/"

	if err := os.MkdirAll(readURL, 0777); err != nil {
		log.Fatalf("Mkdir dir %s error. %v", readURL, err)
	}

	if err := os.MkdirAll(writeURL, 0777); err != nil {
		log.Fatalf("Mkdir dir %s error. %v", writeURL, err)
	}

	if err := os.MkdirAll(workURL, 0777); err != nil {
		log.Fatalf("Mkdir dir %s error. %v", workURL, err)
	}
}

// 创建容器目录
func CreateMountPoint(containerName string, imageName string) error {
	mntUrl := filepath.Join(MntUrl, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		log.Fatalf("Mkdir mountpoint dir %s error. %v", mntUrl, err)
	}
	tmpReadLayer := filepath.Join(ReadLayerUrl, imageName)
	tmpWriteLayer := filepath.Join(WriteLayerUrl, containerName)
	tmpwork := filepath.Join(WorkLayerUrl, containerName)
	/*
		lowerdir read-only
		upperdir read-write cow
		workdir help upperdir
		merged dir
	*/
	lower := "lowerdir=" + tmpReadLayer
	upper := "upperdir=" + tmpWriteLayer
	work := "workdir=" + tmpwork
	parm := lower + "," + upper + "," + work
	//cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	scmds := fmt.Sprintf("mount -t overlay overlay -o %s %s", parm, mntUrl)
	log.Printf("exec command:%+v", scmds)
	_, err := exec.Command("mount", "-t", "overlay", "overlay", "-o", parm, mntUrl).CombinedOutput()
	//_, err := exec.Command("bash", "-c", scmds).CombinedOutput()
	if err != nil {
		log.Fatalf("Run command for creating mount point failed %v", err)
		return err
	}
	return nil
}

// 当容器退出时删除Overlay系统
func DeleteWorkSpace(volume string, containerName string) {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(volumeURLs, containerName)
		} else {
			DeleteMountPoint(containerName)
		}
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
	DeleteWorkLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := filepath.Join(MntUrl, containerName)
	_, err := exec.Command("bash", "-c", "umount "+mntURL).CombinedOutput()
	if err != nil {
		log.Fatalf("Unmount %s error %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Fatalf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteMountPointWithVolume(volumeURLs []string, containerName string) error {
	// 卸载容器里的volume挂载点
	mntURL := filepath.Join(MntUrl, containerName)
	containerUrl := mntURL + "/" + volumeURLs[1]
	_, err := exec.Command("bash", "-c", "umount "+containerUrl).CombinedOutput()
	if err != nil {
		log.Fatalf("Umount volume %s failed. %v", containerUrl, err)
	}

	// 卸载整个容器挂载点
	if err != nil {
		log.Fatalf("Unmount %s error %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		log.Fatalf("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteWriteLayer(containerName string) {
	writeURL := filepath.Join(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Fatalf("Remove writeLayer dir %s error %v", writeURL, err)
	}
}

func DeleteWorkLayer(containerName string) {
	workURL := filepath.Join(WorkLayerUrl, containerName)
	if err := os.RemoveAll(workURL); err != nil {
		log.Fatalf("Remove dir %s error %v", workURL, err)
	}
}
