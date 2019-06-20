# -*- coding: utf-8 -*-
# vim: set softtabstop=2 tabstop=2 shiftwidth=2 expandtab autoindent ft=sls syntax=yaml:

{%- from "nomad/map.jinja" import nomad with context %}

nomad-config:
  file.serialize:
    - name: {{ nomad.config_dir }}/nomad.hcl
    - formatter: json
    - dataset_pillar: nomad:config
    - mode: 640
    - user: root
    - group: root
    {%- if nomad.service != False %}
    - watch_in:
       - service: nomad-service
    {%- endif %}

# Enabling the service here to ensure each state is independent.
nomad-service:
  service.running:
    - name: nomad
    # Restart service if config changes
    - restart: True
    - enable: {{ nomad.service }}
