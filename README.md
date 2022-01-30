# dami-operator

> **Disclaimer:** This project does not assert to be fully correct. It just built for *educational purposes*.


dami-operator is a Kubernetes operator for [dami].

## Install and Run

Before running the operator, please install CRDs as follows:

```bash
$ make install
```

Once you installed CRDs, please make sure that you have up and running [dami] instance. Then, just run:

```bash
$ make run ENABLE_WEBHOOKS=false
```

## DamiDefinition CRD

You can find one example DamiDefinition under [samplk delee](./config/samples/damidefinition.yaml) folder.

```yaml
apiVersion: damigroup.dami.io/v1alpha1
kind: DamiDefinition
metadata:
  name: damidefinition-sample-01
spec:
  resp: "hello response YAML"
```

The `resp` field configures `resp` configuration on the [dami].

Based on the above DamiDefinition, any incoming requests to `/api` endpoint of the [dami] returns `"hello response YAML"`.

```bash
$ curl -isS http://localhost:8001/api
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Content-Type,access-control-allow-origin, access-control-allow-headers
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Sun, 30 Jan 2022 10:38:55 GMT
Content-Length: 35

{"document":"hello response YAML"}
```

[dami]: https://github.com/buraksekili/dami