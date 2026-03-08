---
title: Demo Kurulumu
weight: 3
---

Observer'ı 5 dakikadan kısa sürede tamamen doldurulmuş bir veritabanıyla deneyin. `seed` komutu veritabanını gerçekçi test verileriyle doldurur — kişiler, projeler, destek kayıtları, göç geçmişi, haneler ve daha fazlası — böylece elle veri girmeden her özelliği keşfedebilirsiniz.

## Hızlı başlangıç

```bash
git clone https://github.com/lbrty/observer.git
cd observer
just build
./observer setup
docker compose up -d
./observer migrate up
./observer seed
./observer serve
```

`http://localhost:9000/health` adresini açın — `{"status":"healthy"}` yanıtını görmelisiniz.

Web arayüzünü başlatın:

```bash
cd packages/observer-web
bun install
bun run dev
```

`http://localhost:5173` adresini açın ve `admin@example.com` / `password` ile giriş yapın.

## Seed komutu ne oluşturur

| Ne                    | Ayrıntılar                                                 |
| --------------------- | ---------------------------------------------------------- |
| **Reference data**    | Ülkeler, eyaletler, yerleşim yerleri, ofisler, kategoriler |
| **Users**             | Admin + personel hesapları, hepsinin şifresi `password`    |
| **Projects**          | 2 proje (`--projects` ile yapılandırılabilir)              |
| **People**            | Proje başına 50 kişi (`--people` ile yapılandırılabilir)   |
| **Support records**   | Kişilere bağlı danışma kayıtları                           |
| **Migration records** | Çıkış/varış noktalarıyla hareket geçmişi                   |
| **Households**        | Üyelerle birlikte aile grupları                            |
| **Notes**             | Kişilere eklenmiş vaka notları                             |
| **Pets**              | Etiketli evcil hayvan kayıtları                            |
| **Tags**              | Sınıflandırma için projeye özgü etiketler                  |

### Varsayılan kimlik bilgileri

| Email               | Password   | Role  |
| ------------------- | ---------- | ----- |
| `admin@example.com` | `password` | Admin |

### Özel seed seçenekleri

```bash
# Daha fazla proje ve kişi
./observer seed --projects 5 --people 200

# Tekrarlanabilir veri (aynı seed = aynı çıktı)
./observer seed --seed 42
```

{{% callout type="warning" %}}
Seed komutu veri eklemeden önce **tüm tabloları siler**. Asla üretim veritabanında çalıştırmayın.
{{% /callout %}}

## Adım adım açıklama

### 1. Build

```bash
just build
```

`observer` ikili dosyasını derler.

### 2. Setup

```bash
./observer setup
```

Mantıklı varsayılan değerlerle `.env` dosyası, `keys/` ve `data/uploads/` dizinleri oluşturur ve JWT imzalama için 4096-bit RSA anahtar çifti üretir. `.env` zaten mevcutsa, üzerine yazmadan önce onay ister.

### 3. Servisleri başlatın

```bash
docker compose up -d
```

PostgreSQL ve Redis'i arka planda başlatır.

### 4. Migrate

```bash
./observer migrate up
```

Tüm veritabanı migration'larını uygular.

### 5. Seed

```bash
./observer seed
```

Veritabanını demo verilerle doldurur. Bu adım, boş bir örnek ile paydaşlara gösterebileceğiniz çalışan bir demo arasındaki farkı yaratır.

### 6. Çalıştır

```bash
./observer serve
```

API sunucusunu `http://localhost:9000` adresinde başlatır.

### 7. Frontend (isteğe bağlı)

```bash
cd packages/observer-web
bun install
bun run dev
```

Web arayüzünü `http://localhost:5173` adresinde başlatır.

## Sıfırlama

Her şeyi silip sıfırdan başlamak için:

```bash
docker compose down -v
```

Ardından 3. adımdan itibaren tekrar çalıştırın.
