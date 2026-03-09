---
title: Audit Curnaly
weight: 9
---

Observer iş maalymatdaryndağy ar bir tüzüü, cañyrtuu cana coq qyluu operatsijasy üçün audit curnalynyn cazuusun saqtajt. Audit curnaldary qoşuu ğana — alardy özgörtüügö ce coq qyluuğa bolbojt.

## Emne cazylat

| Kategorija         | Cazylğan operatsijalar                              |
| ------------------ | --------------------------------------------------- |
| Adamdar            | Adam cazuularyn tüzüü, cañyrtuu, coq qyluu         |
| Qoldoo cazuulary   | Konsultatsijalardy tüzüü, cañyrtuu, coq qyluu      |
| Migratsija cazuulary | Qyjmyl cazuularyn tüzüü, cañyrtuu, coq qyluu     |
| Üj-bülölör         | Üj-bülö cana müçö cazularyn tüzüü, cañyrtuu, coq qyluu |
| Eskertmeler        | Iş eskertmelerini tüzüü, cañyrtuu, coq qyluu       |
| Dokumentter        | Cüktöö, metadatalardy cañyrtuu, dokumentterdi coq qyluu |
| Üj canybarları     | Üj canybarlarynyn cazularyn tüzüü, cañyrtuu, coq qyluu |
| Uruqsattar         | Dolboor uruqsattaryn berüü, cañyrtuu, qajtaryp aluu |

Autentifikasija okujalary (kirüü, çyğuu, tokendi cañyrtuu) dolboordun audit curnalynda coq — alar server curnaldarynda pajda bolot.

## Audit curnaldy körüü

Audit curnalyna calgyz dolboor menecerleri cana eeleri kire alat.

```http
GET /projects/:project_id/audit?page=1&per_page=50
```

| Parametr   | Süröttömö                              |
| ---------- | -------------------------------------- |
| `page`     | Baraqça nomeri (demejki 1)             |
| `per_page` | Baraqçadağy natyjcalar (demejki 50)    |
| `actor_id` | Özgörtüü casağan qoldonuuçu bojunça çypqaloo |
| `start`    | Kün bojunça çypqaloo (CCCC-AA-KK)      |
| `end`      | Kün bojunça çypqaloo (CCCC-AA-KK)      |

### Coop formaty

Ar bir audit cazuusu tömönkülördü qamtyjt:

```json
{
  "id": "01J...",
  "project_id": "01J...",
  "actor_id": "01J...",
  "actor_ip": "192.168.1.1",
  "action": "create",
  "entity_type": "person",
  "entity_id": "01J...",
  "created_at": "2024-06-15T10:30:00Z"
}
```

| Tala          | Süröttömö                                                |
| ------------- | -------------------------------------------------------- |
| `actor_id`    | Araketti casağan qoldonuuçu                              |
| `actor_ip`    | Suroo-talaptyn IP daregi                                 |
| `action`      | `create`, `update` ce `delete`                           |
| `entity_type` | Cazuu türü (`person`, `support_record`, `note`, c.b.)   |
| `entity_id`   | Taasir etken cazuunun ULID'i                             |

## Web-interfejste audit

1. Dolboordu açyñyz
2. Navigasijadağy **Audit Curnaly** degendi basyñyz
3. Kün aralyğy ce qoldonuuçu bojunça çypqalañyz

{{% callout type="note" %}}
Audit curnalyna kirüü `manager` ce `owner` dolboor rolun talap qylat. Konsultanttar cana körüüçülör audit curnalın körö alışpajt.
{{% /callout %}}
