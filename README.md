# reftype

Sample to show how to query and show manifest of using OCI reference types 


## List reference type manifests 

```
➜ go run . ls localhost:5000/reftype-test:v1 | jq
{
  "mediaType": "application/vnd.oci.artifact.manifest.v1+json",
  "digest": "sha256:28ffe613ed468b4aa3b60061f529ef676144d47bfc7067bd465d58e370f07718",
  "size": 491
}
```



## View the new artifact manifest 

```
➜ go run . manifest localhost:5000/reftype-test@sha256:28ffe613ed468b4aa3b60061f529ef676144d47bfc7067bd465d58e370f07718
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.artifact.manifest.v1+json",
  "blobs": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar",
      "digest": "sha256:db020b5530356d43394688e483afd350acc008b8d4bdb37e9906105f8b367fbb",
      "size": 12422
    }
  ],
  "reference": {
    "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
    "digest": "sha256:65b3a80ebe7471beecbc090c5b2cdd0aafeaefa0715f8f12e40dc918a3a70e32",
    "size": 528
  }
}
```