Terraform Provisioner for Online.net server rescue mode
========================================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

About The Provisioner
---------------------

This provisioner is built to enable the [Online.net rescue mode](https://documentation.online.net/en/dedicated-server/rescue/rescue-mode)
on a server and then execute actions to it via [remote-exec](https://www.terraform.io/docs/provisioners/remote-exec.html) when it is being
created.  
This can be used to for example format the disks or install a specific OS.

Building The Provisioner
---------------------

Clone repository to: `$GOPATH/src/github.com/src-d/terraform-provisioner-online-rescue`

```sh
$ mkdir -p $GOPATH/src/github.com/scr-d; cd $GOPATH/src/github.com/src-d
$ git clone git@github.com:src-d/terraform-provisioner-online-rescue
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/src-d/terraform-provisioner-online-rescue
$ make build
```

Using the provisioner
------------------

```
resource "online_server" "dedibox" {
    provisioner "online-rescue" {
        "enabled" = true
        "server" = "${online_server.dedibox.id}"
        "image" = "ubuntu-18.04_amd64"
        "credentials_dir" = "/tmp/${online_server.dedibox.id}"
    }

    // the following is a hack as provisioners can not output data
    connection {
        host = "${file("/tmp/${online_server.dedibox.id}/ip")}"
        type     = "ssh"
        user     = "${file("/tmp/${online_server.dedibox.id}/username")}"
        password = "${file("/tmp/${online_server.dedibox.id}/password")}"
    }

    provisioner "remote-exec" {
        inline = [
        "mkfs.ext4 /dev/sda",
        ]
    }

    provisioner "online-rescue" {
        "enabled" = true
        "server" = "${online_server.dedibox.id}"
    }
}
```

Developing the Provisioner
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provisioner-online-rescue
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
