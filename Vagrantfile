Vagrant::Config.run do |config|
  config.vm.box = "lucid64"
  config.vm.box_url = "http://files.vagrantup.com/lucid64.box"
  config.vm.forward_port  80, 8080
  config.vm.forward_port  9300, 9300

   config.vm.provision :chef_solo do |chef|
     chef.cookbooks_path = "cookbooks"
     chef.add_recipe("apt")
     chef.add_recipe("java")
     chef.add_recipe("elasticsearch")
   end
end