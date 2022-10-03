# 🎩 Observer

## Encryption

**NOTE:**

Application depends on set of private keys which then used to encrypt/decrypt generated private keys.

### Entity level encryption

For some entities we use encryption for sensitive personal information each of entities
contain a field `encryption_key` which has the following format `system_key_hash:base64_key_contents`
where `system_key_hash` is a hash of a key which is used to encrypt generated private key `base64_key_content`.
