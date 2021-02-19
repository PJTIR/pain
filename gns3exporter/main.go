package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jhoonb/archivex"
)

// Hints for a node being a Cisco router or switch
const(
	switchDiskImage string = "vios_l2-adventerprisek9-m.vmdk.SSA.152-4.0.55.E"
	routerDiskImage string = "c7200-adventerprisek9-mz.124-24.T5.image"
)

func main(){
	// Checking all environment variables to be present
	//projectName := os.Getenv("PROJECTNAME")
	//url := os.Getenv("URL")

	projectName := "testing"
	url := "http://10.0.0.9:3080"

	if projectName == ""{
		log.Fatalf("The project name should not be empty\n")
	}
	
	project, err := getProjectByName(url, projectName)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("Project name: %s\n", project.Name)
	
	routers := make([]*Node, 0)
	switches := make([]*Node, 0)

	for i := 0; i < len(project.Nodes); i++ {
		if project.Nodes[i].Properties["hda_disk_image"] == switchDiskImage{
			switches = append(switches, &project.Nodes[i])
		}else if project.Nodes[i].Properties["image"] == routerDiskImage{
			routers = append(routers, &project.Nodes[i])
		}
	}

	// Check if the slices actually work
	/*for i, s := range routers{
		fmt.Printf("Router %d: %s\n", i, s.Name)
	}

	for i, s := range switches{
		fmt.Printf("Switch %d: %s\n", i, s.Name)
	}*/

	err = createFolder("configs")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chdir("configs")
	if err != nil {
		log.Fatal(err)
	}
	for _, router := range routers{
		err := createFolder(router.Name)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(router.Name)
		if err != nil {
			log.Fatal(err)
		}
		err = router.getFile("configs/i1_startup-config.cfg", "./" + router.Name + ".cfg")
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir("..")
		if err != nil {
			log.Fatal()
		}
	}	
	

	for _, switchy := range switches{
		err := createFolder(switchy.Name)
		if err != nil {
			log.Fatal(err)
		}
		img := switchy.Properties["hda_disk_image"].(string)
		if _, err := os.Stat(img); os.IsNotExist(err){
			err := switchy.getNodeImage("./" + img)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Dir at this point is ./config
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)
		err = os.Chdir(switchy.Name)
		if err != nil {
			log.Fatal(err)
		}
		err = switchy.getFile("hda_disk.qcow2", "./hda_disk.qcow2")
		if err != nil {
			log.Fatal(err)
		}

		// Now mount the image
		// First rebase the image to a Cisco IOS file
		cmd := exec.Command("qemu-img", "rebase", "-f", "qcow2", "-u", "-b", "../" + img, "hda_disk.qcow2")
		//fmt.Println(cmd.String())
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Oh no bad pp %s\n", err)
		}
		err = createFolder("./mount")
		if err != nil {
			log.Fatalf("Can't create folder ./mount: %s\n", err)
		}
		// Then mount
		cmd = exec.Command("guestmount", "-a", "hda_disk.qcow2", "-m" , "/dev/sda1", "./mount")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			wd, _ := os.Getwd()
			log.Fatalf("command \"guestmount -a hda_disk.qcow2 -m /dev/sda1 ./mount\" failed.\nThe working directory was %s\n%s\n", wd, err)
		}
		
		/*cmd = exec.Command("mount", "/dev/nbd0p1", "./mount") // Creates a folder ./configs/switchname/mount
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}*/
		// Now in ./config/switchname/mount
		// Now we can also copy the nvram and other files like vlan.dat
		err = os.Chdir("mount")
		if err != nil {
			log.Fatalf("Changing dirs to mount failed, %s\n", err)
		}
		// Now we copy the nvram
		err = copyFunction("../nvram", "nvram")
		if err != nil {
			log.Fatalf("copy nvram failed, %s\n", err)
		}
		// If the vlan.dat file exists, copy it
		if _, err = os.Stat("vlan.dat"); !os.IsNotExist(err){
			err = copyFunction("../vlan.dat", "vlan.dat")
			if err != nil {
				log.Fatalf("copy vlan.dat failed, %s\n", err)
			}
		}
		// Return to the switchname folder
		err = os.Chdir("..")
		if err != nil {
			log.Fatalf("changing dirs to .. failed, %s\n", err)
		}
		// Unmount the folder and disconnect using qemu-nbd
		cmd = exec.Command("umount", "./mount")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			wd, _ := os.Getwd()
			log.Fatalf("unmounting dir ./mount failed\nWorking dir was %s, %s\n", wd, err)
		}
		/*cmd = exec.Command("sudo", "qemu-nbd", "-d", "/dev/nbd0")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err != nil {
			log.Fatalf("6, %s", err)
		}*/

		// The mount has been disconnected
		// Now export the startupconfig from the nvram and save it
		// We do that using the iou_export program found in the project root
		// This is obtained from https://git.b-ehlers.de/ehlers/IOUtools
		cmd = exec.Command("../../iou_export", "nvram", "startupconfig.cfg", "privateconfig.cfg")
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalf("iou_export, %s\n",err)
		}
		// Remove mount folder and nvram file
		err = os.Remove("mount")
		if err != nil {
			log.Fatalf("Removing mount folder failed %s\n", err)
		}
		err = os.Remove("nvram")
		if err != nil {
			log.Fatalf("Removing nvram failed %s\n", err)
		}
		err = os.Remove("hda_disk.qcow2")
		if err != nil {
			log.Fatalf("Removing hda_disk.qcow2 failed %s\n", err)
		}
		

		// Change back to the configs folder and repeat until no switches are left
		err = os.Chdir("..")
		if err != nil {
			wd, _ := os.Getwd()
			log.Fatalf("changing dir to .. failed\nWorking dir was %s, %s\n", wd, err)
		}
	}

	// Alright, finally done...
	// All the configs should now be in the following structure:
	// - configs/
	//  	routerName/
	// 			routerName.cfg
	//		...etc
	// 		switchName/
	//			startupconfig.cfg
	//			privateconfig.cfg
	//			vlan.dat (if it exists on nvram)

	// The following step is to zip it all up and copy it somewhere yet to be decided
	// Change back to the root folder
	err = os.Chdir("..")
	if err != nil {
		log.Fatalf("8, %s", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	zipFile := new(archivex.ZipFile)
	err = zipFile.Create("config.zip")
	if err != nil {
		log.Fatal(err)
	}
	err = zipFile.AddAll(currentDir + "/configs", false)
	if err != nil {
		log.Fatal(err)
	}
	err = zipFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	
	err = createFolder("/output")
	if err != nil {
		log.Fatal(err)
	}
	err = copyFunction("output/config.zip", "config.zip")

	for{
		fmt.Printf("Complete, copy the zip file and close the container!\n")
		time.Sleep(time.Second * 5)
	}



}	 

func copyFunction(dest string, src string) error{
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	_, err = io.Copy(destfile, srcfile)
	if err != nil {
		return err
	}
	destfile.Close()
	srcfile.Close()
	return nil
}

func createFolder(name string) error{
	err := os.MkdirAll(name, 0777)
	if err != nil {
		return err
	}
	return nil
}	

/* 
Docker starts the go program
	go reads the ip and port of the server
	go then downloads all the config files for routers

	go then downloads all hda_disk.qcow2 images in order and performs the following
		mount this disk
		perform iou_export on the nvram (https://www.b-ehlers.de/blog/posts/2017-10-26-inspect-modify-qemu-images/)
		copy the config to a folder named after the node name
		remove the hda_disk file
		repeat until all switches are done
	zip all of these files
	copy this zip to a special folder that's mounted to the host
	wait forever until the host stops the container
*/