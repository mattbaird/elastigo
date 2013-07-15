include_recipe "build-essential"
package "mercurial"

bash "Fetch golang source" do
  cwd "/home/vagrant"
  code <<-EOC
    hg clone #{node['golang']['repo']}
    cd /home/vagrant/go/
    hg update tip #{node['golang']['repo']}
  EOC
  creates "~/go"
end

bash "Build golang" do
  cwd "/home/vagrant/go/src"
  code <<-EOC
    ./all.bash
  EOC
end

bash "Export ENV Vars" do
  code <<-EOC
    echo 'export GOBIN=~/go/bin' >> /home/vagrant/.bashrc
    mkdir -p /home/vagrant/code/go/
    chown vagrant /home/vagrant/code/go/
    echo 'export GOPATH=/home/vagrant/code/go/' >> /home/vagrant/.bashrc
    echo 'export GOROOT=/usr/local/go/' >> /home/vagrant/.bashrc
    echo 'export PATH=$PATH:$GOBIN' >> /home/vagrant/.bashrc
    source /home/vagrant/.bashrc
  EOC

end