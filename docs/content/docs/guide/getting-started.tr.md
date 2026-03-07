---
title: Başlangıç
weight: 2
---

## 5 dakikada çalışır halde görün

Bir sunucuya, barındırma sağlayıcısına veya BT departmanına ihtiyacınız yok. Docker kurulu bir dizüstü bilgisayarınız varsa, Observer'ı hemen şimdi çalışır halde görebilirsiniz.

```bash
git clone https://github.com/lbrty/observer.git
cd observer
cp .env.example .env
just generate-keys
just docker-up
just run
```

Tarayıcınızda `http://localhost:9000/health` adresini açın. `"status":"healthy"` görüyorsanız, backend çalışıyor demektir.

Ardından web arayüzünü başlatın:

```bash
just web-dev
```

`http://localhost:5173` adresini açın — Observer'a bakıyorsunuz.

## Ne başlattınız

- API'yi sunan bir **backend** — kimlik doğrulama, veri depolama ve raporları yönetir
- Kişiler, hanehalkları, destek kayıtları, göç geçmişi, belgeler ve evcil hayvanlar için tabloları olan bir **veritabanı**
- Proje yönetimi, rol tabanlı erişim ve yerleşik raporlama içeren bir **web arayüzü**
- **Otomatik giriş güvenliği** — her oturum yenilemesinde token'lar döndürülür

Tüm bunlar tek bir makinede çalışır. Üretim ortamında, herhangi bir sunucuya kopyalayabileceğiniz tek bir dosyaya derlenir.

## Gerçek dağıtıma hazır mısınız?

"Deniyorum" aşamasından "ekibim bunu her gün kullanıyor" aşamasına geçmek için şunlara ihtiyacınız var:

| Ne | Neden |
| --- | --- |
| Bir sunucu (VPS veya yerinde) | Observer kendi sunucunuzda barındırılır — verileriniz altyapınızdan asla çıkmaz |
| PostgreSQL | Observer'ın ihtiyaç duyduğu tek harici hizmet |
| Yaklaşık 30 dakika | Alan adınız yönlendirilmiş bir sunucuda `docker compose up` çalıştırın |

Abonelik yok. Kullanıcı başına ücret yok. Satıcı bağımlılığı yok. Verilerin ve dağıtımın sahibi sizsiniz.

Adım adım üretim kurulumu için [Dağıtım](/docs/guide/deployment/) sayfasına bakın.

## Geliştiriciler için: yerel kurulum

Observer'ın kendisi üzerinde çalışmak istiyorsanız, şu araçların kurulu olması gerekir:

| Araç | Sürüm | Kurulum |
| --- | --- | --- |
| Go | 1.25.* | https://go.dev/dl/ |
| Bun | latest | https://bun.sh/ |
| Docker + Compose | latest | https://docs.docker.com/get-docker/ |
| Just | latest | https://github.com/casey/just#installation |

### 1. Klonlayın ve bağımlılıkları yükleyin

```bash
git clone https://github.com/lbrty/observer.git
cd observer
go mod download
bun install
```

### 2. Ortamı yapılandırın

```bash
cp .env.example .env
```

Varsayılan değerler, sağlanan `docker-compose.yml` ile kutudan çıktığı gibi çalışır. Düzenleme gerekmez.

### 3. İmzalama anahtarlarını oluşturun

```bash
just generate-keys
```

Bu, Observer'ın giriş token'larını imzalamak için kullandığı bir anahtar çifti oluşturur. `keys/` dizini gitignore'dadır — her geliştirici kendi anahtarlarını oluşturur.

### 4. Her şeyi başlatın

```bash
just docker-up    # PostgreSQL ve Redis'i başlatır
just run          # backend'i :9000 portunda başlatır (migration'ları otomatik çalıştırır)
just web-dev      # ön yüzü :5173 portunda başlatır
```

## Bir şey çalışmıyor mu?

**Port 5432 zaten kullanımda** — Muhtemelen yerel bir PostgreSQL çalışıyor. Durdurun veya `docker-compose.yml` dosyasında portu değiştirin.

**Anahtar yolları için "no such file or directory"** — Önce `just generate-keys` komutunu çalıştırmanız gerekiyor.

**Migration "connection refused" hatası veriyor** — Veritabanı konteyneri henüz hazır olmayabilir. `just docker-up` komutundan sonra birkaç saniye bekleyin ve tekrar deneyin.
