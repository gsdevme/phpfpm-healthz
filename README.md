# phpfpm-healthz

Why? Running a health check on php-fpm when deployed as a standalone container without nginx can be abit of a painful experience. Options include using a bash script implementation like this https://github.com/renatomefi/php-fpm-healthcheck 

## Kubernetes

```yaml
# Check its listening on 127.0.0.1:9000 (running)
livenessProbe:
    exec:
      command:
        - echo
        - -n
        - >
        - /dev/tcp/127.0.0.1/9000
    failureThreshold: 1
# Check we are ready to receive traffic (ready)
readinessProbe:
    exec:
      command:
        - phpfpm-healthz
        - --uri=/api/my-custom-endpoint/healthz
    initialDelaySeconds: 30
    periodSeconds: 5
    timeoutSeconds: 5
    failureThreshold: 1
```

## Running it

```bash
$: phpfpm-healthz -h                                                           
Usage:
  php-fpm healthz fastcgi checker [flags]

Flags:
      --file string   The path to the script filename (default "/app/public/index.php")
  -h, --help          help for php-fpm
      --uri string    The Request URI that you want to hit (default "/healthz")
```

It will print to stdout also in the even of success or failure

```bash
./phpfpm-healthz --uri=/api/v2/healthz
Endpoint /app/public/index.php/api/v2/healthz
Success, status code: 200
```

```bash
./phpfpm-healthz --uri=/not_a_real_endpoint
Endpoint /app/public/index.php/not_a_real_endpoint
Non zero status code returned 404
```