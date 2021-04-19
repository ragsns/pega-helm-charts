{{- define "pega.config" -}}
{{- $arg := .mode -}}
# Node type specific configuration for {{ .name }}
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .name }}
  namespace: {{ .root.Release.Namespace }}
data:
{{- if eq $arg "deploy-config" }}
{{- $custom_config := .custom }}

  # Pega deployment prconfig.xml file
  prconfig.xml: |-
{{- if $custom_config.prconfig }}
{{ $custom_config.prconfig | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/prconfig.xml" | indent 6 }}
{{- end }}

  # Pega deployment prlog4j2.xml file
  prlog4j2.xml: |-
{{- if $custom_config.prlog4j2 }}
{{ $custom_config.prlog4j2 | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/prlog4j2.xml" | indent 6 }}
{{- end }}

  # Pega deployment contextXML template file
  context.xml.tmpl: |-
{{- if $custom_config.contextXML }}
{{ $custom_config.contextXML | indent 6 }}
{{ else }}
{{ .root.Files.Get "config/deploy/context.xml.tmpl" | indent 6 }}
{{- end }}
{{- end }}
---
{{- end }}