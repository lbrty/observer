---
title: Baştoo
weight: 2
---

## 5 minutada iştep körüñüz

Server, hosting provayder ce IT bölümü kerek emes. Eger laptopuñuzda Docker ornotulğan bolso, Observer azyr ele iştep catqanyn körsöñüz bolot.

```bash
git clone https://github.com/lbrty/observer.git
cd observer
cp .env.example .env
just generate-keys
just docker-up
just run
```

Brauzerden `http://localhost:9000/health` açyñyz. Eger `"status":"healthy"` körsöñüz, backend iştep catat.

Andan kiyin web interfejsti baştañyz:

```bash
just web-dev
```

`http://localhost:5173` açyñyz — Observer aldyñyzda.

## Emne baştadyñyz

- **Backend** API qyzmat körsötöt — autentifikasija, maalymat saqtoo cana esepdemeler
- **Maalymat bazasy** adamdar, üj çarbalar, qoldoo cazuulary, migrasija taryhy, dokumentter cana üj cajandyqtar üçün tablitsalar menen
- **Web interfejs** proekt başqaruu, rolğo tajanğan cetüü cana kirgizilgen esepdemeler menen
- **Avtomattyq login qoopsuzduğu** — tokender ar bir sessija cañylanuuda ajlanat

Munun bardyğy bir maşinada iştejt. Produksionda bir fajlğa kompilasijalanat, any qajsy bolbosun serverge köçürö alasyz.

## Çynığy ornottooğo dajynsyzby?

"Synap körüp catam"dan "komandam kün sajyn qoldonot"ko ötüü üçün sizge kerek:

| Emne                        | Emne üçün                                                                         |
| --------------------------- | --------------------------------------------------------------------------------- |
| Server (VPS ce cergiliktüü) | Observer öz serverinde iştejt — maalymatyñyz infrastukturañyzdan eç qaçan çyqpajt |
| PostgreSQL                  | Observer kerek qylğan calğyz tyşqy qyzmat                                         |
| 30 minutadaj                | Domeniñiz qaratylğan serverde `docker compose up` iştetiñiz                       |

Catyluu coq. Ar bir qoldonuuçuğa baasy coq. Vendor bajlama coq. Maalymat cana ornottoo sizdin.

Qadam-qadam produksija ornotuu üçün [Ornottoo](/docs/guide/deployment/) qarañyz.

## Iştep çyğuuçular üçün: cergeliktüü ornottoo

Eger Observerdin özünde iştegisi kelse, bul quraldardy ornotuñuz:

| Qural            | Versija | Ornottoo                                   |
| ---------------- | ------- | ------------------------------------------ |
| Go               | 1.25.\* | https://go.dev/dl/                         |
| Bun              | aqyrqy  | https://bun.sh/                            |
| Docker + Compose | aqyrqy  | https://docs.docker.com/get-docker/        |
| Just             | aqyrqy  | https://github.com/casey/just#installation |

### 1. Klondoo cana çöjrönü ornotuu

```bash
git clone https://github.com/lbrty/observer.git
cd observer
go mod download
bun install
```

### 2. Çöjrönü tuuraloo

```bash
cp .env.example .env
```

Defolt maaniler berilgen `docker-compose.yml` menen dajyn iştejt. Özgörtüü kerek emes.

### 3. Qol qojuu açqyçtaryn tüzüü

```bash
just generate-keys
```

Bul Observerdin login tokenderine qol qojuu üçün açqyç cuptaryn tüzöt. `keys/` papkasy gitignored — ar bir iştep çyğuuçu özünün açqyçtaryn tüzöt.

### 4. Bardyğyn baştoo

```bash
just docker-up    # PostgreSQL cana Redis baştajt
just run          # backend :9000 portunda baştajt (migrasijalar avtomattyq işlejt)
just web-dev      # frontend :5173 portunda baştajt
```

## Bir nerse iştebej cataby?

**Port 5432 qoldonuuda** — Cergeliktüü PostgreSQL iştep catqan boluşu mümkün. Any toqtotuñuz, ce `docker-compose.yml` içindegi porttu özgörtüñüz.

**Açqyç coldoru üçün "no such file or directory"** — Adegende `just generate-keys` iştetiñiz.

**Migrasija "connection refused" qatasy menen** — Maalymat bazasy kontejneri dajyn emes boluşu mümkün. `just docker-up` kejin bir neçe sekunda kütüp, qajra sunañyz.
