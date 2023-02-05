---
title: Encryption keys 🔐
---

## Encryption keys 🔐

RSA keys are used to

1. To encrypt certain personal information,
2. Issue and validate JWT access and refresh tokens,
3. AES secrets we generate to encrypt uploaded documents.

## Generating a key

It is possible to use `openssl` to generate keys

```sh
openssl genrsa -out key.pem 2048
```

You can also use built-in CLI generator

```sh
python -m observer keys generate --size 2048 -o key.pem 
```
