maintainer        "Josh Holt"
maintainer_email  "holt.josh@gmail.com"
license           "MIT"
description       "Installs/Configures GoLang"
long_description  IO.read(File.join(File.dirname(__FILE__), 'README.md'))
version           "0.0.1"

recipe "golang", "Installs GoLang based on the default install method"
recipe "golang::local", "Installs golang in vagrant's home directory"
recipe "golang::global", "Installs golang in /usr/local/go"

depends "build-essential"

%w{ debian ubuntu centos redhat }.each do |os|
  supports os
end