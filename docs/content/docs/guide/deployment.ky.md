---
title: Ornotuu
weight: 4
---

Bul qoldonmo Observerdi komandañyz qoldono ala turğan serverge ornottoonu cetektejt. Tereñ tehnikalyq bilim kerek emes — eger serverge SSH qylyp, biraz komanda iştete alsañyz, munu qyla alasyz.

## Docker (sunuştalat)

Bul eñ ceñil col. Observer web interfejs menen dajyn bir Docker image katary kelet - frontend cağynan qoşumça ornotuu ce konfigurasija kerek emes.

### Emne kerek

- Docker cana Docker Compose ornotulğan server
- Serveriñizge qaratylğan domen aty (HTTPS üçün)

### 1-qadam: Açqyçtaryn tüzüü

Observer login tokenderine qol qojuu üçün RSA açqyçtaryn qoldonot. Bulardy tüzüü üçün serveriñizde bul komandalardy işletiñiz:

```sh
mkdir -p keys
openssl genrsa -out keys/jwt_rsa 4096
openssl rsa -in keys/jwt_rsa -pubout -out keys/jwt_rsa.pub
```

Bul açqyçtardy qoopsuz saqtañyz. Eger coğoltsoñuz, bardyğy qajra login boluşu kerek.

### 2-qadam: Çöjrönü tuuraloo

Çöjrö fajl ülgüsün köçürüp, öz ornottooñuzğa tuuralañyz:

```sh
cp .env.example .env
```

Eñ maanilüü özgörmölör:

| Özgörmö                | Emne qylat                                              | Default                    |
| ---------------------- | ------------------------------------------------------- | -------------------------- |
| `DATABASE_DSN`         | Observer PostgreSQLğa qandaj bajlanat                   | _(ornotuluşu ylazym)_      |
| `REDIS_URL`            | Observer Rediske qandaj bajlanat                        | `redis://localhost:6379/0` |
| `JWT_PRIVATE_KEY_PATH` | 1-qadamdağy ceke açqyçtyn colu                          | `keys/jwt_rsa`             |
| `JWT_PUBLIC_KEY_PATH`  | 1-qadamdağy açyq açqyçtyn colu                          | `keys/jwt_rsa.pub`         |
| `CORS_ORIGINS`         | Siziñ domeniñiz (misaly `https://observer.yourorg.org`) | `http://localhost:5173`    |
| `COOKIE_SECURE`        | HTTPS qoldonğonda `true` qojuñuz (kerek)                | `true`                     |
| `SERVER_HOST`          | Qajsy darekte tuñdoo                                    | `localhost`                |
| `SERVER_PORT`          | Qajsy porttu tuñdoo                                     | `9000`                     |

Toluq tizme üçün [Çöjrö Özgörmölörü](/docs/developers/reference/variables/) qarañyz.

### 3-qadam: Observerdi baştoo

```sh
docker compose up -d
```

Bul PostgreSQL, Redis cana Observerdi baştajt. Maalymat baza shemasy birinçi işletüüdö avtomattyq türdö tüzülöt — qol menen migrasija qadamy kerek emes.

### 4-qadam: Iştep catqanyn tekşerüü

```sh
curl http://localhost:9000/health
```

Munu körüşüñüz kerek:

```json
{ "status": "healthy", "database": "connected", "timestamp": "..." }
```

Eger munu körsöñüz, Observer dajyn. Web interfejske cetüü üçün domeniñizdi brauzerde açyñyz.

## Dockersuz (VPS / bare metal)

Eger Observerdi tüzdön-tüz iştetkiñiz kelse, binary tüzüñüz:

```sh
CGO_ENABLED=0 go build -tags production -ldflags="-s -w" -o observer ./cmd/observer
```

`-tags production` flagy web interfejsti binary içine kirgizet. Qajsy bolbosun cerge köçürö ala turğan bir fajl alasyz.

Işletiñiz:

```sh
./observer serve --host 0.0.0.0
```

PostgreSQL cana Redis özünçö iştep turuşu kerek. `DATABASE_DSN` cana `REDIS_URL` alarğa qaratyñyz.

## HTTPS ornottoo

Observerdi HTTPS işletken revers proxy artynda dajyma işletüü kerek. Bul login maalymattaryn cana ceke maalymattardy tranzitte şifrlejt.

[Caddy](https://caddyserver.com/) eñ ceñil variant — sertifikattardy avtomattyq işlejt:

```
observer.yourorg.org {
    reverse_proxy localhost:9000
}
```

Eger Nginx ce başqa proxy qoldonsoñuz, bulardy ornotuñuz:

- Çöjröñüzdö `COOKIE_SECURE=true`
- `CORS_ORIGINS` çynyqy domeniñizge (misaly `https://observer.yourorg.org`)
