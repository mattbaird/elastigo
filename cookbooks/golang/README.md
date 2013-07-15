###Installs GoLang

1. config.add_recipe "golang::global" --> /usr/local/go
2. config.add_recipe "golang::local"  --> /home/vagrant/go

When using this cookbook you should decide if you want a global install
or a local user install. You can do so by choosing one of the two options
above.