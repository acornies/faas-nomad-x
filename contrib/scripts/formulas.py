#!/usr/bin/env python

from StringIO import StringIO
from zipfile import ZipFile
from urllib import urlopen
import shutil
import os

formulas = {}

# Key   = github URL name to formula; ie: https://github.com/${org}/${key}/archive/master.zip
# Value/List = directory/name of formulas within repo (these get copied to salt states root, and is therefore the name referenced in top.sls)

formulas["vault-formula"] = ["vault"]
formulas["nomad-formula"] = ["nomad"]
formulas["consul-formula"] = ["consul"]

if __name__ == "__main__":
    for formula in formulas:
        print "processing {0}".format(formula)
        print "downloading https://github.com/saltstack-formulas/{0}/archive/master.zip".format(formula)
        resp = urlopen('https://github.com/saltstack-formulas/{0}/archive/master.zip'.format(formula))
        zipfile = ZipFile(StringIO(resp.read()))
        if os.path.exists('/vagrant/contrib/salt/states/{0}'.format(formula)):
            print "deleting existing /vagrant/contrib/salt/states/{0}".format(formula)
            shutil.rmtree('/vagrant/contrib/salt/states/{0}'.format(formula))
        zipfile.extractall("/vagrant/contrib/salt/states/{0}".format(formula))
        for state in formulas[formula]:
            if os.path.exists('/vagrant/contrib/salt/states/{0}'.format(state)):
                print "deleting existing /vagrant/contrib/salt/states/{0}".format(state)
                shutil.rmtree('/vagrant/contrib/salt/states/{0}'.format(state))
            print "moving /vagrant/contrib/salt/states/{0}/{1}-master/{2} -> /vagrant/contrib/salt/states/{3}".format(formula,formula,state,state)
            shutil.move("/vagrant/contrib/salt/states/{0}/{1}-master/{2}".format(formula,formula,state), "/vagrant/contrib/salt/states/")
        print "deleting /vagrant/contrib/salt/states/{0}".format(formula)
        shutil.rmtree('/vagrant/contrib/salt/states/{0}'.format(formula))