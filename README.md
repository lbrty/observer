[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/banner2-direct.svg)](https://stand-with-ukraine.pp.ua)
# 🎩 Observer

This project seeks to bridge the gap in tools and systems to support and help internally displaced persons
or in general people in need any kind of support humanitarian, legal or seeking medical help.

Made in 🇩🇪 with love from 🇰🇬 & 🇺🇦 for people in need.

## 🌈 Session handling

Motivation behind this is to avoid database queries to verify sessions and keep it as stateless as possible.
To manage sessions we use JWT access tokens and refresh tokens where access token is a short living
token, while refresh tokens will last a lot longer they are used to obtain new access tokens.

## 🍄 Encryption

**NOTE:**

Application depends on set of private keys which then used to encrypt/decrypt generated private keys.

### ❄️ Entity level encryption

For some entities we use encryption for sensitive personal information each of entities
contain a field `encryption_key` which has the following format `system_key_hash:base64_key_contents`
where `system_key_hash` is a hash of a key which is used to encrypt generated private key `base64_key_content`.

### ⚡️ Auth & security

To hash and verify user's passwords we use `passlib` with `bcrypt` and for MFA (TOTP) we use `totp` packages.
