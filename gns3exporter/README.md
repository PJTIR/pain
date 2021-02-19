# The GNS3 exporter

The program can export configuration files from Cisco switches and routers virtualized in GNS3, by utilizing its API. The whole thing runs in Docker to ensure maximum compatibility, as several pieces of software need to be installed. This program depends on the following:

- qemu-utils
- python

Docker is used to prevent cluttering the system with these programs.

## Instructions

First of all:
[Get Docker](https://docs.docker.com/get-docker/)

Simply follow those instructions for your system. If you're using Windows, just follow their instructions and install Docker desktop. Just use the command line!
If you're using stupidly complicated distro's such as Gentoo, find it out yourself. You must be good at that by now.  
_Tip: follow the [Post-installation steps for Linux](https://docs.docker.com/engine/install/linux-postinstall/)! If you don't all the following commands using Docker have to be preceded by sudo_

Then in this repository enter the `gns3exporter` and run the following command to build the docker image

```
$ docker build -t gns3exporter .
```

After having done so, you can view the image by running `docker image ls`. You will see something like this:

```
REPOSITORY                TAG           IMAGE ID       CREATED         SIZE
gns3exporter              latest        83907d4fdfcb   2 minutes ago   1.04GB
```

To run the container, the `docker run` command is used. For now, to get the output of the program (a config.zip file) you are required to open a second shell in which you can perform the copying. So open a second shell already 馬鹿！

Then to run the container run the following command in the first shell (or second doesn't matter):

```
$ docker run -it --name exporter \
    --cap-add SYS_ADMIN \
    --device /dev/fuse \
    --security-opt apparmor:unconfined \
    --env URL=http://172.21.185.1:3080 \
    --env PROJECTNAME={projectname} # use "testing" to test
    gns3exporter:latest
```

This will start the container and keep your shell attached to it's output, meaning you cannot run anything else in that shell until the container exits. This it won't do automatically, so you will have to end it at some point using CTRL+C.
Now wait for the program to finish doing it's thing. You will be notified by a friendly message telling you `Complete, copy the zip file and close the container!` that repeats every 5 seconds.

In your second shell, run the following command to copy. Also don't forget to change the last argument to change the destination where the file should be copied to:

```
$ docker cp exporter:/output/config.zip {destination path}
```

Well well, that's it. After the copy has been completed, you can stop the container by pressing CTRL+C in the shell with the container running. After you might want to clean that container up by performing:

```
$ docker rm exporter
```

This way, next time you copy my commands in your shell, you won't get an error saying that the name `exporter` is already taken.
