![R.O.C](https://github.com/Happykat/R.O.C-CONTROLS/blob/master/assets/logo-roc-flat.png)


# R.O.C-CONTROLS

Repository of the ROC Controls

[![Build Status](https://travis-ci.org/Happykat/R.O.C-CONTROLS.svg?branch=development)](https://travis-ci.org/Happykat/R.O.C-CONTROLS)

Part of the [ROC project](Chttp://eip.epitech.eu/2017/remoteocculuscontroller/)

Here we aim to control any type of robots ! Be it your own one or some bigs shots robots.
 
# What is R.O.C Controls ?

R.O.C stands for Remote Oculus Controller. Our goal is to transcript the direct environment of a robot to distant user which controls it. It is made of three main part a client, a video-server and a control-server. [More info](http://http://eip.epitech.eu/2017/remoteocculuscontroller/) and [take look here](https://github.com/Remote-Oculus-Controller)

Here we are going to see the control part. It is designated be adaptive and evolve to the constraint of embedded system and user needs.

R.O.C is golang package build on top of [gobot.io](https://gobot.io/) that add network interaction and R.O.C protocols as well as a nice features.

# Getting started

The server is made entirely in [go](https://golang.org/) and we are using protobuff. The project (Controls) can be build on any platform as long as you meet this requirement. Nonetheless we recommend Linux as a development platform as the installation is easier and you will meet less trouble on the way.

###Todo list
- [ ] Smoothed movement in Motion example
- [ ] Add movement prediction in GPS module
- [ ] Look for a keep alive on websocket

#### Other List
- [ ] Finish 3D support for robot
- [ ] Complete Wiki

## Install

Required system packages installation.

### Golang:

Follow your system package manager instruction, be sure to have at least version 1.6. ex:

```bash
$ pacman -S golang
...
$ go verion
golang version go1.6 ....
```

Or you install via the source. [source](https://golang.org/dl/)

### Protobuff 3

Golang was supported in protobuff v2 so it is mandatory to have v3. Protobuff compile .proto and then create your .go source files.

It can be tricky to install depending on system and often the v3 is not present by official repository.

On Arch-based system use yaourt to access AUR packages:
```bash
yaourt -S protobuf-go
```

The installation will take some times. In the mean times you can pull the control repository or go through the wiki!

Once the installation is complete you should be able to do this:

```bash
$ protoc --version
libprotoc 3.X.X
```

Warning: The first number should be 3 and nothing else !

### PULL & BUILD !!!

Now that your system is up to the task you can build the controls.

```bash
go get Remote-Oculus-Controller/R.O.C-CONTROLS
```

On success you will now be able to build a robot that will be able to interact to other R.O.C projects !

Notes: R.O.C Controls use an other package "proto" which is a located [here](https://github.com/Remote-Oculus-Controller/proto)