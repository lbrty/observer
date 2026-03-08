---
title: Demo ornotuu
weight: 3
---

Observer'du toluk malımat bazası menen 5 mınötkö chetinçe sınap körüñüz. `seed` buyruğu bazağa realduv test malımattar kiyiret — adamdar, proyektter, koldoo jazuuları, köçüü tarıhı, üy bölüktörü jana başkalar — baar mümkünçülüktördü koldon kiyirgenden malımat terbebey izildey alasız.

## Tez baştoo

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

`http://localhost:9000/health` açıñız — `{"status":"healthy"}` körünüşü kerek.

Veb-interfeysti iştötiñiz:

```bash
cd packages/observer-web
bun install
bun run dev
```

`http://localhost:5173` açıñız jana `admin@example.com` / `password` menen kiriñiz.

## Seed buyruğu emne tuzdayt

| Emne                   | Çoo malımat                                          |
| ---------------------- | ---------------------------------------------------- |
| **Spravka malımatı**   | Ölkölör, oblusttar, cerler, ofister, kategoriyalar   |
| **Koldonuuçular**      | Admin + kızmatker akkaunttar, baarının sırsozu `password` |
| **Proyektter**         | 2 proyekt (`--projects` menen öztürgö bolot)         |
| **Adamdar**            | Ar proyektke 50 (`--people` menen öztürgö bolot)     |
| **Koldoo jazuuları**   | Adamdarga baylanışkan konsultatsiya jazuuları         |
| **Köçüü jazuuları**    | Çıkkaan/baruuçu cerler menen kozğoluu tarıhı         |
| **Üy bölüktörü**       | Müçölörü menen üy-bülö toptoşturuuları               |
| **Beleşmeler**         | Adamdarga tirkölgön iş beleşmeleri                   |
| **Üy aybandarı**       | Tegder menen ayban jazuuları                          |
| **Tegder**             | Proyekt çegindegi kategoriyalöö tegderi               |

### Baştapkı kiriş malımattar

| Email               | Sırsöz     | Rol   |
| ------------------- | ---------- | ----- |
| `admin@example.com` | `password` | Admin |

### Seed parametrleri

```bash
# Köbüröök proyekt jana adam
./observer seed --projects 5 --people 200

# Kayta çığaruuçu malımat (birdey seed = birdey natıyca)
./observer seed --seed 42
```

{{% callout type="warning" %}}
Seed buyruğu malımat kiyirgenge çeyin **baarık tabliçalardı tazalayt**. Anı eç kaçan produktion bazağa iştötpöñüz.
{{% /callout %}}

## Kadamdap tüşündürmö

### 1. Kuruştuuruu

```bash
just build
```

`observer` binarnik faylın kompiliyatsiyalayt.

### 2. Ornotuu

```bash
./observer setup
```

Ilayıktuu baştapkı maan menen `.env` faylın, `keys/` jana `data/uploads/` papkaların tuzdayt, jana JWT kol koyuu üçün 4096-bit RSA açkıç cubun çığarat. Eger `.env` bar bolso, üstünön cazuudan murda surayt.

### 3. Servisterdi iştötüü

```bash
docker compose up -d
```

PostgreSQL jana Redis'ti fonduk recimde iştötöt.

### 4. Migratsiya

```bash
./observer migrate up
```

Baarık baza migratsiyalardı koldonot.

### 5. Seed

```bash
./observer seed
```

Bazanı demo malımattar menen toltuarat. Bul kadamdın ayırması — boş instansiya menen kızıkdar tüzgö körsötö turgan iştegen demo ortosundağı ayırma.

### 6. İştötüü

```bash
./observer serve
```

API serverin `http://localhost:9000` dareğinde iştötöt.

### 7. Frontend (kaalaganıñızça)

```bash
cd packages/observer-web
bun install
bun run dev
```

Veb-interfeysti `http://localhost:5173` dareğinde iştötöt.

## Tazalap baştoo

Baarın öçürüp kaytadan baştoo üçün:

```bash
docker compose down -v
```

Annan 3-kadamdan baştañız.
