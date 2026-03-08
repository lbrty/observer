---
title: CLI Referansı
weight: 4
---

Observer; sunucu yönetimi, veritabanı geçişleri, anahtar üretimi, kullanıcı yönetimi ve geliştirme araçları için bir CLI sağlar.

## Kurulum

```bash
# Kaynaktan derleme
just build

# Veya doğrudan kurulum
go install github.com/lbrty/observer/cmd/observer@latest
```

## Komutlar

### serve

HTTP sunucusunu başlatır.

```bash
observer serve [flags]
```

| Flag     | Type   | Default     | Açıklama                                                     |
| -------- | ------ | ----------- | ------------------------------------------------------------ |
| `--host` | string | `localhost` | Sunucu adresi (`SERVER_HOST` env değişkenini geçersiz kılar) |
| `--port` | int    | `9000`      | Sunucu portu (`SERVER_PORT` env değişkenini geçersiz kılar)  |

**Örnekler:**

```bash
# Varsayılan ayarlarla başlat (localhost:9000)
observer serve

# Özel adres ve port
observer serve --host 0.0.0.0 --port 8080

# Ortam değişkenleri ile yapılandırma
DATABASE_DSN="postgres://..." REDIS_URL="redis://..." observer serve
```

Üretim derlemelerinde, gömülü geçişler başlangıçta otomatik olarak uygulanır. Sunucu, SIGINT/SIGTERM sinyallerinde 30 saniyelik bir zaman aşımı ile düzgün bir şekilde kapanır.

---

### migrate

Veritabanı geçiş yönetimi.

#### migrate up

Bekleyen tüm geçişleri uygular.

```bash
observer migrate up [flags]
```

| Flag     | Type   | Default      | Açıklama                   |
| -------- | ------ | ------------ | -------------------------- |
| `--path` | string | `migrations` | Geçiş dosyaları dizin yolu |

```bash
# Bekleyen tüm geçişleri uygula
observer migrate up

# Özel geçiş dizini kullan
observer migrate up --path ./db/migrations
```

#### migrate create

Yeni bir ileri-yönlü geçiş dosyası oluşturur.

```bash
observer migrate create [name] [flags]
```

| Flag     | Type   | Default      | Açıklama                   |
| -------- | ------ | ------------ | -------------------------- |
| `--path` | string | `migrations` | Geçiş dosyaları dizin yolu |
| `--seq`  | uint   | auto         | Açık sıra numarası         |

```bash
# Geçiş oluştur (otomatik numaralandırma)
observer migrate create add_audit_log

# Açık sıra numarası ile oluştur
observer migrate create add_audit_log --seq 25
```

#### migrate version

Mevcut geçiş sürümünü gösterir.

```bash
observer migrate version [flags]
```

| Flag     | Type   | Default      | Açıklama                   |
| -------- | ------ | ------------ | -------------------------- |
| `--path` | string | `migrations` | Geçiş dosyaları dizin yolu |

---

### keygen

JWT imzalama için RSA anahtar çifti üretir.

```bash
observer keygen [flags]
```

| Flag       | Type   | Default | Açıklama                            |
| ---------- | ------ | ------- | ----------------------------------- |
| `--bits`   | int    | `4096`  | RSA anahtar boyutu (minimum 4096)   |
| `--output` | string | `.`     | Anahtar dosyaları için çıktı dizini |

**Örnekler:**

```bash
# Mevcut dizinde anahtar üret
observer keygen

# keys/ dizininde 8192-bit anahtar üret
observer keygen --bits 8192 --output keys
```

Çıktı dosyaları: `private_key.pem` (0600) ve `public_key.pem` (0644).

---

### create-admin

Platform yönetici hesabı oluşturur.

```bash
observer create-admin [flags]
```

| Flag           | Type   | Zorunlu | Açıklama                           |
| -------------- | ------ | ------- | ---------------------------------- |
| `--email`      | string | evet    | Yönetici e-posta adresi            |
| `--password`   | string | evet    | Yönetici parolası (min 8 karakter) |
| `--first-name` | string | hayır   | Ad                                 |
| `--last-name`  | string | hayır   | Soyad                              |
| `--phone`      | string | hayır   | Telefon numarası                   |

**Örnekler:**

```bash
# Zorunlu alanlarla yönetici oluştur
observer create-admin --email admin@example.com --password "s3cure-p4ss"

# İsteğe bağlı profil alanları ile
observer create-admin \
  --email admin@example.com \
  --password "s3cure-p4ss" \
  --first-name Admin \
  --last-name User \
  --phone "+1234567890"
```

Veritabanına bağlanır, parolayı Argon2id ile hashler ve kullanıcıyı admin rolüyle ekler (doğrulanmış + aktif). Yinelenen e-posta ve telefon numaralarını reddeder.

---

### seed

Geliştirme için veritabanını gerçekçi test verileriyle doldurur.

```bash
observer seed [flags]
```

| Flag         | Type  | Default | Açıklama                      |
| ------------ | ----- | ------- | ----------------------------- |
| `--people`   | int   | `50`    | Proje başına kişi sayısı      |
| `--projects` | int   | `2`     | Proje sayısı                  |
| `--seed`     | int64 | `0`     | Rastgele tohum (0 = rastgele) |

**Örnekler:**

```bash
# Varsayılan değerlerle doldur (2 proje, her biri 50 kişi)
observer seed

# Özel sayılar
observer seed --projects 5 --people 200

# Tekrarlanabilir tohum
observer seed --seed 42
```

{{% callout type="warning" %}}
Bu komut veri eklemeden önce TÜM tabloları siler. Üretim veritabanında çalıştırmayın.
{{% /callout %}}

Referans verileri (ülkeler, bölgeler, yerler, ofisler, kategoriler), bilinen parolalarla (`password`) kullanıcılar, izinli projeler oluşturur ve kişileri destek kayıtları, göç kayıtları, notlar, evcil hayvanlar ve hanehalkları ile doldurur.

---

### setup

İlk kurulum işlemini etkileşimli olarak çalıştırır.

```bash
observer setup
```

Bu komut:

1. Makul varsayılanlarla bir `.env` dosyası oluşturur (üzerine yazmadan önce onay ister)
2. Gerekli dizinleri oluşturur (`keys/`, `data/uploads/`)
3. JWT imzalama için 4096-bit RSA anahtar çifti üretir
4. Sonraki adımlar için talimatları yazdırır

**Örnek çıktı:**

```
Created .env with default configuration.
Created directory: keys
Created directory: data/uploads
Generating 4096-bit RSA key pair...
Private key written to: keys/private_key.pem
Public key written to: keys/public_key.pem

Setup complete!

Next steps:
  1. Start Postgres and Redis:
     docker compose up -d

  2. Run database migrations:
     observer migrate up

  3. Create an admin user:
     observer create-admin --email admin@example.com --password "your-password"

  4. Start the server:
     observer serve
```

---

## Yaygın İş Akışları

### İlk kurulum

```bash
# 1. Etkileşimli kurulumu çalıştır (.env, anahtarlar, dizinler oluşturur)
observer setup

# 2. Postgres ve Redis'i başlat
docker compose up -d

# 3. Geçişleri uygula
observer migrate up

# 4. Yönetici hesabı oluştur
observer create-admin --email admin@example.com --password "your-password"

# 5. Sunucuyu başlat
observer serve
```

Örnek verilerle tam demo kurulumu için [demo kurulum kılavuzuna](/docs/guide/demo-setup/) bakın.

### Yeni geçiş ekleme

```bash
# Geçiş dosyasını oluştur
observer migrate create add_audit_log

# Oluşturulan SQL dosyasını düzenle
$EDITOR migrations/000025_add_audit_log.up.sql

# Uygula
observer migrate up
```

### Geliştirme verilerini doldurma

```bash
# Önce geçişlerin uygulandığından emin ol
observer migrate up

# Varsayılan değerlerle doldur
observer seed

# Giriş bilgileri: admin@example.com / password
```

### Yeni JWT anahtarları üretme

```bash
# Yeni anahtarlar üret
observer keygen --output keys

# Yollar varsayılanlardan farklıysa .env dosyasını güncelle
JWT_PRIVATE_KEY_PATH=keys/private_key.pem
JWT_PUBLIC_KEY_PATH=keys/public_key.pem

# Sunucuyu yeniden başlat — mevcut tokenlar geçersiz olacaktır
observer serve
```
