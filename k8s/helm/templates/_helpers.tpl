{{/*
Application labels
*/}}
{{- define "labels" -}}
app: ces
app.kubernetes.io/name: k8s-ces-setup
{{- end }}

{{/*
Creates the docker config json string used as a docker secret.
*/}}
{{- define "docker_config_json" }}
  {{- $url := index . 0 }}
  {{- $username := index . 1 }}
  {{- $passwort := index . 2 }}
  {"auths":{"{{ $url }}":{"username":"{{ $username }}","password":"{{ $passwort }}","email":"test@mtest.de","auth":"{{ printf "%s%s%s" $username ":" $passwort | b64enc}}"}}}
{{- end }}

{{- define "helm_config_json" }}
  {{- $url := index . 0 }}
  {{- $username := index . 1 }}
  {{- $passwort := index . 2 }}
  {"auths": {"{{ $url }}": {"auth": "{{ printf "%s%s%s" $username ":" $passwort | b64enc}}"}}}
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
