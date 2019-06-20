# -*- coding: utf-8 -*-
# vim: set softtabstop=2 tabstop=2 shiftwidth=2 expandtab autoindent ft=sls syntax=yaml:

{% from "nomad/map.jinja" import nomad with context %}

nomad-bin-dir:
  file.directory:
   - name: {{ nomad.bin_dir }}
   - makedirs: True

nomad-config-dir:
  file.directory:
    - name: {{ nomad.config_dir }}
    - makedirs: True

nomad-data-dir:
  file.directory:
    - name: {{ nomad.config.data_dir }}
    - makedirs: True

{% if nomad.build %}
nomad-prepare-build:
  file.directory:
    - name: /tmp/nomad-v{{ nomad.version }}/go/src/github.com/hashicorp
    - user: nobody
    - makedirs: True
    - unless:
      - '{{ nomad.bin_dir }}/nomad -v && {{ nomad.bin_dir }}/nomad -v | grep {{ nomad.version }}'

nomad-check-build-packages:
  pkg.installed:
    - names:
      - {{ nomad.git_pkg }}
      - {{ nomad.make_pkg }}
      - {{ nomad.gcc_pkg }}
    - onlyif:
      - 'command -v go'
    - onchanges:
      - nomad-prepare-build

nomad-checkout-repository:
  git.latest:
    - name: https://github.com/hashicorp/nomad.git
    - target: /tmp/nomad-v{{ nomad.version }}/go/src/github.com/hashicorp/nomad
    - rev: v{{ nomad.version }}
    - user: nobody
    - require:
      - nomad-check-build-packages
    - onchanges:
      - nomad-prepare-build
    - unless:
      - 'test -d /tmp/nomad-v{{ nomad.version }}/go/src/github.com/hashicorp/nomad'

nomad-build:
  cmd.run:
    - names: 
      - 'make bootstrap'
      - 'make dev'
    - cwd: '/tmp/nomad-v{{ nomad.version }}/go/src/github.com/hashicorp/nomad'
    - runas: nobody
    - env:
      - GOPATH: /tmp/nomad-v{{ nomad.version }}/go
      - PATH: '/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/bin:/usr/local/sbin:/tmp/nomad-v{{ nomad.version }}/go/bin'
    - require:
      - nomad-checkout-repository
    - onchanges:
      - nomad-prepare-build

nomad-install-binary:
  file.copy:
    - name: {{ nomad.bin_dir }}/nomad
    - source: /tmp/nomad-v{{ nomad.version }}/go/bin/nomad
    - force: True
    - require:
      - service: nomad-install-binary
      - nomad-bin-dir
    - onchanges:
      - nomad-build
  cmd.run:
    - name: '/usr/bin/strip {{ nomad.bin_dir }}/nomad'
    - onchanges:
      - file: nomad-install-binary
  service.dead:
    - name: nomad
    - onchanges:
      - nomad-build

nomad-install-service:
   file.copy:
    - name: /etc/systemd/system/nomad.service
    - source: /tmp/nomad-v{{ nomad.version }}/go/src/github.com/hashicorp/nomad/dist/systemd/nomad.service
    - onchanges:
      - nomad-build
   module.run:
    - name: service.systemctl_reload
    - onchanges:
      - file: nomad-install-service

nomad-cleanup-build:
  file.absent:
    - name: '/tmp/nomad-v{{ nomad.version }}'
    - require:
      - nomad-install-binary
      - nomad-install-service
    - onchanges:
      - nomad-build

{% else %} 
nomad-install-binary:
  archive.extracted:
    - name: {{ nomad.bin_dir }}
    - source: https://releases.hashicorp.com/nomad/{{ nomad.version }}/nomad_{{ nomad.version }}_{{ grains['kernel'] | lower }}_{{ nomad.arch }}.zip
    - source_hash: https://releases.hashicorp.com/nomad/{{ nomad.version }}/nomad_{{ nomad.version }}_SHA256SUMS
    # If we don't force it here, the mere presence of an older version will prevent an upgrade.
    - overwrite: True 
    # Hashicorp gives a zip with a single binary. Salt doesn't like that.
    - enforce_toplevel: False
    - require:
      - service: nomad-install-binary
    - unless:
      - '{{ nomad.bin_dir }}/nomad -v && {{ nomad.bin_dir }}/nomad -v | grep {{ nomad.version }}'
  file.managed:
    - name: {{ nomad.bin_dir }}/nomad
    - user: root
    - group: root
    - mode: 0755
    - require:
      - archive: nomad-install-binary
  service.dead:
    - name: nomad
    - unless:
      - '{{ nomad.bin_dir }}/nomad -v && {{ nomad.bin_dir }}/nomad -v | grep {{ nomad.version }}'

nomad-install-service:
  file.managed:
    - name: /etc/systemd/system/nomad.service
    - source: https://raw.githubusercontent.com/hashicorp/nomad/v{{ nomad.version }}/dist/systemd/nomad.service
    - source_hash: {{ nomad.service_hash }}
    - onchanges:
      - nomad-install-binary
  module.run:
    - name: service.systemctl_reload
    - onchanges:
      - file: nomad-install-service
{% endif %}



