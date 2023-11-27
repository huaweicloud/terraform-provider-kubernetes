Terraform Huawei Cloud Kubernetes Provider
==========================================

<!-- markdownlint-disable-next-line MD033 -->
<a href="https://www.huaweicloud.com/"><img width="450px" height="102px" src="https://console-static.huaweicloud.com/static/authui/20210202115135/public/custom/images/logo-en.svg"></a>

Requirements
------------

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x
* [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/huaweicloud/terraform-provider-k8s`

```sh
$ mkdir -p $GOPATH/src/github.com/huaweicloud; cd $GOPATH/src/github.com/huaweicloud
$ git clone https://github.com/huaweicloud/terraform-provider-k8s
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/huaweicloud/terraform-provider-k8s
$ make build
```

Using the provider
------------------

Please see the documentation at [provider usage](docs/index.md).

Or you can browse the documentation within this repo [here](https://github.com/huaweicloud/terraform-provider-k8s/tree/master/docs).
