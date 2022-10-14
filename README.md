# waf-builder

waf-builder provides a command to build nginx/openresty along with modsecurity, HTTP/2, and any other desired nginx modules seamlessly.

waf-builder is a fork of [nginx-build](https://github.com/cubicdaiya/nginx-build). It adds a few features, including options to easily build http2 and ModSecurity modules into the final product. It also fixes a few issues observed with `nginx-build`, and gets rid of third party system requirements like `git` without sacrificing functionality.

*Note that this documentation may be incomplete in some areas*

## Build Support

 * [nginx](https://nginx.org/)
 * [OpenResty](https://openresty.org/)
 * [ModSecurity](https://github.com/SpiderLabs/ModSecurity-nginx)

## Custom Configuration

waf-builder provides a mechanism for customizing configuration for building nginx.

### Configuration for building nginx

#### About `--add-module` and `--add-dynamic-module`

`waf-builder` allows to use `--add-module`.

```bash
$ waf-builder \
-d work \
--add-module=/path/to/ngx_http_hello_world
```

Also, `waf-builder` allows to use `--add-dynamic-module`.

```bash
$ waf-builder \
-d work \
--add-dynamic-module=/path/to/ngx_http_hello_world
```

<details>
<summary>Embedding static libraries</summary>

### Embedding zlib statically

Give `-zlib` to `waf-builder`.

```bash
$ waf-builder -d work -zlib
```

`-zlibversion` is an option to set a version of zlib.

### Embedding PCRE statically

Give `-pcre` to `waf-builder`.

```bash
$ waf-builder -d work -pcre
```

`-pcreversion` is an option to set a version of PCRE.

### Embedding OpenSSL statically

Give `-openssl` to `waf-builder`.

```bash
$ waf-builder -d work -openssl
```

`-opensslversion` is an option to set a version of OpenSSL.

*Note: this likely shouldn't be done if the target deployment is using OpenSSL in FIPS mode, nginx should rely on the shared library in that scenario.*

### Embedding LibreSSL statically

Give `-libressl` to `waf-builder`.

```bash
$ waf-builder -d work -libressl
```

`-libresslversion` is an option to set a version of LibreSSL.

</details>

<details>Embedding nginx modules</details>
<summary>

### Embedding 3rd-party modules

`waf-builder` provides a mechanism for embedding 3rd-party modules.
Prepare a json file below.

```json
[
  {
    "name": "ngx_http_hello_world",
    "form": "git",
    "url": "https://github.com/cubicdaiya/ngx_http_hello_world",
    "dynamic": false
  }
]
```

(Change dynamic accordingly)

Give this file to `waf-builder` with `-m`.

```bash
$ waf-builder -d work -m modules.json.example
```

waf-builder will use a built-in git client (a new feature of this fork in contrast to `nginx-build`) to fetch the module. if it is a valid nginx module, it will build nginx/openresty with it included.

#### Provision for 3rd-party module

There are some 3rd-party modules expected provision. `waf-builder` provides the options such as `shprov` and `shprovdir` for this problem.
There is the example configuration below.

```json
[
  {
    "name": "njs/nginx",
    "form": "hg",
    "url": "https://hg.nginx.org/njs",
    "shprov": "./configure && make",
    "shprovdir": ".."
  }
]
```

</details>

<details>
<summary>Patches and Idempotent builds</summary>

## Applying patch before building nginx

`waf-builder` provides the options such as `-patch` and `-patch-opt` for applying patch to nginx.

```console
waf-builder \
 -d work \
 -patch something.patch \
 -patch-opt "-p1"
```

## Idempotent build

`waf-builder` supports a certain level of idempotent build of nginx.
If you want to ensure a build of nginx idempotent and do not want to build nginx as same as already installed nginx,
give `-idempotent` to `waf-builder`.

```bash
$ waf-builder -d work -idempotent
```

`-idempotent` ensures an idempotent by checking the software versions below.

* nginx
* PCRE
* zlib
* OpenSSL

On the other hand, `-idempotent` does not cover versions of 3rd party modules and dynamic linked libraries.

</details>

## Build OpenResty

`waf-builder` supports to build [OpenResty](https://openresty.org/).

```bash
$ waf-builder -d work -openresty -pcre -openssl
```

If you don't install PCRE and OpenSSL on your system, it is required to add the option `-pcre` and `-openssl`.


And there is the limitation for the support of OpenResty.
`waf-builder` does not allow to use OpenResty's unique configure options directly.
If you want to use OpenResty's unique configure option, [Configuration for building nginx](#configuration-for-building-nginx) is helpful.
