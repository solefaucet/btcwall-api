from fabric.api import env, run
from fabric.operations import run, local, put, sudo
from fabric.decorators import runs_once, parallel
from os import popen, remove
from tempfile import NamedTemporaryFile

env.shell = '/bin/bash -lc'
env.user = 'd'
env.roledefs.update({
    'production': [
        '192.168.193.149', # api-1
        '192.168.194.49', # api-2
    ]
})

# Heaven will execute fab -R production deploy:branch_name=master
def deploy(branch_name):
    print("Executing on %s as %s" % (env.host, env.user))

    local_task(branch_name)
    remote_setup()
    restart()

@runs_once
def local_task(branch_name):
    # checkout the branch, it can also be a tag
    local('git checkout %s' % branch_name)

    # build btcwall-api with go1.7.3
    goversion = '1.7.3'
    gitcommit = popen('git rev-parse --short HEAD').read()
    buildtime = popen("date '+%Y-%m-%d_%I:%M:%S%p'").read()
    flags = '-X main.goVersion=%s -X main.buildTime=%s -X main.gitCommit=%s' % (goversion, buildtime, gitcommit)
    project = 'github.com/solefaucet/btcwall-api'
    directory = '/go/src/%s' % project
    local('docker run --rm -v "$PWD":%s -w %s golang:%s go build -ldflags "%s" -v -o btcwall-api' % (directory, directory, goversion, flags))

@parallel
def remote_setup():
    put('/etc/btcwall/api.supervisor.conf', '/etc/supervisor/conf.d/btcwall-api.conf', mode=0644, use_sudo=True)
    put('btcwall-api', '/usr/local/bin/btcwall-api', mode=0755, use_sudo=True)
    put('/opt/GeoLite2-City.mmdb', '/opt/GeoLite2-City.mmdb', mode=0644, use_sudo=True)
    put('/usr/local/share/swagger', '/opt', mode=0755, use_sudo=True)
    put('apidoc', '/opt', mode=0755, use_sudo=True)
    conf = popen('cat /etc/btcwall/api.yml').read()
    conf = replace_macro(conf, { 'internal_ip': env.host })
    with NamedTemporaryFile(delete=False) as f:
        f.write(conf)
        f.close()
        put(f.name, '/etc/btcwall-api.yml', mode=0644, use_sudo=True)
        remove(f.name)

def replace_macro(s, macro_map):
    for key, val in macro_map.iteritems():
        s = s.replace("%"+key+"%", val)
    return s

def restart():
    # restart btcwall-api
    sudo('supervisorctl update all')
    sudo('supervisorctl restart btcwall-api')
