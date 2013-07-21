# Go Language Chef Cookbook

This is an OpsCode Chef cookbook for [Go, the programming language](http://golang.org).

It uses the ["Go Language Gophers" Ubuntu PPA](https://launchpad.net/~gophers/+archive/go)
and allows you to tweak version using Chef node attributes.

It is released under the [Apache Public License 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).


## Recipes

Main recipe is `golang::default`.


## Attributes

* `[:golang][:version]` (default: "stable"): Go version to install, either `stable` or `tip` (the bleeding edge)


## Supported OSes

Ubuntu 10.10 to 12.04, will likely work just as well on Debian unstable.


## Dependencies

None.


## Copyright & License

Michael S. Klishin, Travis CI Development Team, 2012.

Released under the [Apache Public License 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
