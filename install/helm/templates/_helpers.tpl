{{- define "ImagePullSecret"}}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.global.registry.url (printf "%s:%s" .Values.global.registry.user .Values.global.registry.password | b64enc) | b64enc }}
{{- end}}