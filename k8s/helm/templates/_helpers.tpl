{{- define "k8s-ces-setup.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "k8s-ces-setup-finisher.name" -}}
{{- "k8s-ces-setup-finisher"}}
{{- end }}



{{/* All-in-one labels */}}
{{- define "k8s-ces-setup.labels" -}}
app: ces
helm.sh/chart:  {{- printf " %s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{ include "k8s-ces-setup.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Selector labels */}}
{{- define "k8s-ces-setup.selectorLabels" -}}
app.kubernetes.io/name: {{ include "k8s-ces-setup.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "k8s-ces-setup-finisher.labels" -}}
app: ces
helm.sh/chart:  {{- printf " %s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/name: {{ include "k8s-ces-setup-finisher.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}


{{/*
Creates the docker config json string used as a docker secret.
*/}}
{{- define "docker_config_json" }}
  {{- $url := index . 0 }}
  {{- $username := index . 1 }}
  {{- $password := index . 2 | b64dec }}
  {"auths":{"{{ $url }}":{"username":"{{ $username }}","password":{{ $password | toJson }},"email":"test@mtest.de","auth":"{{ print $username ":" $password | b64enc}}"}}}
{{- end }}

{{- define "helm_config_json" }}
  {{- $host := index . 0 }}
  {{- $username := index . 1 }}
  {{- $password := index . 2 | b64dec }}
{{/*  Helm auth does not work with protocols in config file. Remove them to be sure.*/}}
  {"auths": {"{{ $host | replace "oci://" "" | replace "http://" "" | replace "https://" "" }}": {"auth": "{{ print $username ":" $password | b64enc}}"}}}
{{- end }}


{{- define "printCloudoguLogo" }}
{{- printf "\n" }}
...
                    ./////,
                ./////==//////*
               ////.  ___   ////.
        ,**,. ////  ,////A,  */// ,**,.
   ,/////////////*  */////*  *////////////A
  ////'        \VA.   '|'   .///'       '///*
 *///  .*///*,         |         .*//*,   ///*
 (///  (//////)**--_./////_----*//////)   ///)
  V///   '°°°°      (/////)      °°°°'   ////
   V/////(////////\. '°°°' ./////////(///(/'
      'V/(/////////////////////////////V'
{{- printf "\n" }}
{{- end }}
