#
# Cookbook Name:: golang
# Recipe:: gvm
#
# Copyright 2012, Michael S. Klishin, Travis CI Development Team
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

include_recipe "gvm"

log "Default Go version will be #{node.golang.multi.default_version}"

home    = node.travis_build_environment.home
gvm     = "source #{home}/.gvm/scripts/gvm && gvm"
env     = {'HOME' => home}
user    = node.travis_build_environment.user
aliases = node.golang.multi.aliases || {}

setup = lambda do |bash|
  bash.user user
  bash.environment env
end

node.golang.multi.versions.each do |v|
  bash "golang::multi: installing #{v}" do
    setup.call(self)
    code "#{gvm} install #{v}"
    not_if { File.exists?(File.join(home, ".gvm", "gos", v)) }
  end
end

bash "set #{node.golang.multi.default_version} to be the default Go runtime version" do
  setup.call(self)
  code "#{gvm} use #{node.golang.multi.default_version} --default"
end

aliases.each do |new_name, existing_name|
  bash "alias #{existing_name} => #{new_name}" do
    setup.call(self)
    code "#{gvm} alias create #{new_name} #{existing_name}"

    ignore_failure true # alias creation is not idempotent. @fd.
  end
end
