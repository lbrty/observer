---
title: Dağıtım
weight: 4
---

Bu rehber, Observer'ı ekibinizin kullanabileceği bir sunucuya kurma sürecinde size yol gösterir. Derin teknik bilgiye ihtiyacınız yok — bir sunucuya SSH ile bağlanıp birkaç komut çalıştırabiliyorsanız, bunu yapabilirsiniz.

## Docker (önerilen)

Bu en basit yoldur. Observer, web arayüzü zaten dahil edilmiş tek bir Docker imajı olarak sunulur — ön yüz tarafında kurulacak veya yapılandırılacak ekstra bir şey yoktur.

### Gerekenler

- Docker ve Docker Compose kurulu bir sunucu
- Sunucunuza yönlendirilmiş bir alan adı (HTTPS için)

### Adım 1: İmzalama anahtarlarını oluşturun

Observer, giriş token'larını imzalamak için RSA anahtarları kullanır. Bunları oluşturmak için sunucunuzda şu komutları çalıştırın:

```bash
mkdir -p keys
openssl genrsa -out keys/jwt_rsa 4096
openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
```

Bu anahtarları güvende tutun. Kaybederseniz, herkesin yeniden giriş yapması gerekir.

### Adım 2: Ortamınızı yapılandırın

Örnek ortam dosyasını kopyalayın ve kurulumunuza göre düzenleyin:

```bash
cp .env.example .env
```

En önemli değişkenler:

| Değişken               | Ne işe yarar                                      | Varsayılan                    |
| ---------------------- | ------------------------------------------------- | -------------------------- |
| `DATABASE_DSN`         | Observer'ın PostgreSQL'e nasıl bağlanacağı               | _(ayarlanmalı)_            |
| `REDIS_URL`            | Observer'ın Redis'e nasıl bağlanacağı                    | `redis://localhost:6379/0` |
| `JWT_PRIVATE_KEY_PATH` | Adım 1'deki özel anahtarın konumu         | `keys/jwt_rsa`             |
| `JWT_PUBLIC_KEY_PATH`  | Adım 1'deki genel anahtarın konumu          | `keys/jwt_rsa.pub`         |
| `CORS_ORIGINS`         | Alan adınız (ör. `https://observer.yourorg.org`) | `http://localhost:5173`    |
| `COOKIE_SECURE`        | HTTPS kullanırken `true` olarak ayarlayın (ayarlamalısınız)       | `true`                     |
| `SERVER_HOST`          | Dinlenecek adres                        | `localhost`                |
| `SERVER_PORT`          | Dinlenecek port                           | `9000`                     |

Tam liste için [Ortam Değişkenleri](/docs/developers/reference/variables/) sayfasına bakın.

### Adım 3: Observer'ı başlatın

```bash
docker compose up -d
```

Bu, PostgreSQL, Redis ve Observer'ı başlatır. Veritabanı şeması ilk başlatmada otomatik olarak oluşturulur — manuel migration adımı gerekmez.

### Adım 4: Çalıştığını doğrulayın

```bash
curl http://localhost:9000/health
```

Şunu görmelisiniz:

```json
{ "status": "healthy", "database": "connected", "timestamp": "..." }
```

Bunu görüyorsanız, Observer hazırdır. Web arayüzüne erişmek için alan adınızı tarayıcıda açın.

## Docker olmadan (VPS / bare metal)

Observer'ı doğrudan çalıştırmayı tercih ederseniz, binary'yi derleyin:

```bash
CGO_ENABLED=0 go build -tags production -ldflags="-s -w" -o observer ./cmd/observer
```

`-tags production` bayrağı, web arayüzünü binary'ye gömer. Her yere kopyalayabileceğiniz tek bir dosya elde edersiniz.

Çalıştırın:

```bash
./observer serve --host 0.0.0.0
```

PostgreSQL ve Redis'in ayrıca çalışıyor olması gerekir. `DATABASE_DSN` ve `REDIS_URL` değişkenlerini onlara yönlendirin.

## HTTPS kurulumu

Observer'ı her zaman HTTPS'yi yöneten bir reverse proxy arkasında çalıştırmalısınız. Bu, giriş kimlik bilgilerinin ve kişisel verilerin aktarım sırasında şifrelenmesini sağlar.

[Caddy](https://caddyserver.com/) en kolay seçenektir — sertifikaları otomatik olarak yönetir:

```
observer.yourorg.org {
    reverse_proxy localhost:9000
}
```

Nginx veya başka bir proxy kullanıyorsanız, şunları ayarladığınızdan emin olun:

- Ortamınızda `COOKIE_SECURE=true`
- `CORS_ORIGINS` değerini gerçek alan adınıza ayarlayın (ör. `https://observer.yourorg.org`)
