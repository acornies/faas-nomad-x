# -*- coding: utf-8 -*-
# vim: set softtabstop=2 tabstop=2 shiftwidth=2 expandtab autoindent ft=sls syntax=yaml:

{% from "nomad/map.jinja" import nomad with context %}

nomad-remove-service:
  service.dead:
    - name: nomad
  file.absent:
    - name: /etc/systemd/system/nomad.service
    - require:
      - service: nomad-remove-service
  cmd.run:
    - name: 'systemctl daemon-reload'
    - onchanges:
      - file: nomad-remove-service

nomad-remove-binary:
   file.absent:
    - name: {{ nomad.bin_dir }}/nomad
    - require:
      - nomad-remove-service

nomad-remove-config:
    file.absent:
      - name: {{ nomad.config_dir }}
      - require:
        - nomad-remove-binary

nomad-remove-datadir:
    file.absent:
      - name: {{ nomad.config.data_dir }}
      - require:
        - nomad-remove-config
      - onlyif: 
        - 'command -v {{ nomad.bin_dir }}/nomad'
