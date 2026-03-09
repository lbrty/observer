---
title: Eksport
weight: 8
---

Observer iş maalymattaryn cana otçöttardy CSV fajldary qatary eksporttoj alat. Eksport üçün keminde `consultant` dolboor rolu talap qylynat.

## Iş maalymattaryn eksporttoo

Dolboor bojunça bardyq adamdardy cana alardyn tieşelüü cazuularyn eksporttoo:

```http
GET /projects/:project_id/export/people?format=csv
```

Qoşumça suroo parametrleri:

| Parametr      | Süröttömö                                        |
| ------------- | ------------------------------------------------ |
| `start`       | Qattalğan kün bojunça çypqaloo (CCCC-AA-KK)      |
| `end`         | Qattalğan kün bojunça çypqaloo (CCCC-AA-KK)      |
| `category_id` | Ajaluuluq kategorijasy bojunça çypqaloo             |
| `tag_id`      | Teg bojunça çypqaloo                              |

Coop ağym qatary cönötülöt — çoñ maalymat toptomdoru progressivdüü cönötülöt, oşonduktan bajlanyş üzülböjt.

## Qoldoo cazuularyn eksporttoo

Dolboor bojunça konsultasija cazuularyn eksporttoo:

```http
GET /projects/:project_id/export/support-records?format=csv&start=2024-01-01&end=2024-12-31
```

| Parametr  | Süröttömö                                         |
| --------- | ------------------------------------------------- |
| `start`   | `provided_at` künü bojunça çypqaloo (CCCC-AA-KK)  |
| `end`     | `provided_at` künü bojunça çypqaloo (CCCC-AA-KK)  |
| `type`    | `legal` ce `social`                               |
| `sphere`  | Qoldoo çöjrösü (mis. `housing_assistance`)        |

## Migratsija cazuularyn eksporttoo

Qyjmyl/cer qotoruu taryhyn eksporttoo:

```http
GET /projects/:project_id/export/migration-records?format=csv
```

## CSV formaty

Bardyq eksporttor baş qatar menen UTF-8 CSV qoldonot. Talalar API cooptoru menen birdej strukturağa ee. Sezgiç talalar (`contact`, `personal`, `documents`) dolboor uruqsat celekterine carasa kiret ce çyğarylat — API oquusuna qoldonulğan ereceler eksportko dağy qoldonulat.

## Web-interfejste cüktöp aluu

1. Dolboordu açyñyz
2. Navigasijadağy **Eksport** degendi basyñyz
3. Cazuu türün cana kün aralyğyn tandañyz
4. **CSV cüktöp aluu** degendi basyñyz

Fajl tüzdön-tüz brauzeriñizge cüktölöt.

{{% callout type="note" %}}
Eksporttor dolbooruñuzğa çektelgen. Siz dajyndalbağan dolboordon maalymattardy eksporttoj albajsyz cana sezgiçtik celekteriñiz qoldonulat — eger `can_view_personal` öçürülgön bolso, eksporttalğan CSV uluttuq ID ce tuulğan kündördü qamtybajt.
{{% /callout %}}
