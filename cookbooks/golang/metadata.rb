maintainer        "Michael Klishin"
maintainer_email  "michael@defprotocol.org"
license           "Apache 2.0"
description       "Installs go language from Go Language Gophers Ubuntu PPA"
long_description  IO.read(File.join(File.dirname(__FILE__), 'README.md'))
version           "1.0.0"
recipe            "golang", "Installs go"

depends "apt"

supports "ubuntu"
