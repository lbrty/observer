---
title: Raporlor
weight: 7
---

Observer ukrainalyk ökümdük emös ujumdardyn donorduk rapor mildettemelerine ylajyq iştelip çyqqan 39 strukturalyq rapor türün caratat. Bardyq raporlor bir proektke cardalat cana küntizme aralyğy boyunça süzülöt.

## Rapordu Iştetüü

Bardyq Raporlor tömöndögüdöj cetkiliktüü:

```http
GET /projects/:project_id/reports/:report_type?start=2024-01-01&end=2024-12-31
```

`:report_type` orniña tömdögü tablitsalardagy identifikatorlordun birini qojuñuz.

## Sanoo Türlörü

Ar bir Rapor üç sanoo ädisinin birini qoldonot:

| Tür          | Maaany                                                                                 |
| ------------ | -------------------------------------------------------------------------------------- |
| **Okuja**    | `support_records` satyrlary — üç konsultasija alğan bir adam 3 dep eseptelet           |
| **Adamdar**  | Ar türdüü şahsiyattar — bir adam qançalağan cazuuğa iye bolbo birdej bir ret eseptelet |
| **Bölüktör** | Ar türdüü üj-büluö bölüktörü (bir üj-büluö = bir bölük)                                |

Küntizme süzüüsü konsultasijalar üçün `support_records.provided_at` cana registrasija raporloru üçün `people.registered_at` qoldonot — `created_at` emes.

## Reapor Türlörü

### 1-topt — Calpy Konsultasija Sany

| #   | Rapor                                        | Tür   |
| --- | -------------------------------------------- | ----- |
| 1   | Bardyq türdögü konsultacijalardyn calpy sany | Okuja |
| 2   | Bardyq uquqtuq konsultasijalar               | Okuja |
| 3   | Bardyq sotsialdyq konsultasijalar            | Okuja |

### 2-topt — Jinsi Bölünüşü

| #   | Rapor                                  | Tür     |
| --- | -------------------------------------- | ------- |
| 12  | Maalymatta kattalgan erkekter          | Adamdar |
| 13  | Maalymatta kattalgan ayaldar           | Adamdar |
| 14  | Uquq konsultasija alğan ayaldar        | Adamdar |
| 15  | Sotsialdyq konsultasija alğan ayaldar  | Adamdar |
| 16  | Uquq konsultasija alğan erkekter       | Adamdar |
| 17  | Sotsialdyq konsultasija alğan erkekter | Adamdar |

14–17 raporlor ar türdüü adamdardy sanajt, konsultasija okujalaryn emes.

### 3-topt — Geografijalyk / IDP Statusu

| #    | Rapor                                                          | Tür     |
| ---- | -------------------------------------------------------------- | ------- |
| 4    | calpy kattalğan adamdar                                        | Adamdar |
| 5–6  | Konflikt aymağyndan kattalğan adamdar                          | Adamdar |
| 7–10 | Konflikt aymağyndan uquq/sotsialdyq konsultasija alğan adamdar | Adamdar |
| 11   | IDP emes kattalğandar                                          | Adamdar |

IDP statusu avtomattyk türdö anyktalat: `origin_place_id → places → states.conflict_zone`. 5–6 cana 7–10 raporloru konflikt aymağynyn belgisi boyunça parameterlenet.

### 4-topt — Açuura Kategorijasy Boyunça Bölünüş

| #   | Rapor                                                      | Tür     |
| --- | ---------------------------------------------------------- | ------- |
| 18  | Kattalğan adamdar — açuura kategorijasy boyunça            | Adamdar |
| 19  | Sotsialdyq konsultasija alğan adamdar — kategorija boyunça | Adamdar |
| 20  | Uquq konsultasija alğan adamdar — kategorija boyunça       | Adamdar |

Kategorija tayyndalbağan adamdar "kategorijasyz" tobiñe kiret.

### 5-topt — Azyrky Turğan Aymağy

| #   | Rapor                                                 | Tür     |
| --- | ----------------------------------------------------- | ------- |
| 21  | Kattalğan adamdar — azyrky aymaq boyunça              | Adamdar |
| 22  | Uquq konsultasija alğan adamdar — aymaq boyunça       | Adamdar |
| 23  | Sotsialdyq konsultasija alğan adamdar — aymaq boyunça | Adamdar |

### 6-topt — Qoldoo Sferas Boyunça Bölünüş

| #   | Rapor                                                 | Tür     |
| --- | ----------------------------------------------------- | ------- |
| 24  | Uquq konsultasija sany — sfera boyunça                | Okuja   |
| 25  | Uquq konsultasija alğan adamdar — sfera boyunça       | Adamdar |
| 29  | Sotsialdyq konsultasija sany — sfera boyunça          | Okuja   |
| 30  | Sotsialdyq konsultasija alğan adamdar — sfera boyunça | Adamdar |

**Qoldoo sferalary:**

| Mäán                    | Süröttömö                                             |
| ----------------------- | ----------------------------------------------------- |
| `housing_assistance`    | Turğun üj uquqtary, eviction, sotsialdyq turğun üj    |
| `document_recovery`     | Pasporttor, tuulğandyk jarlyqtar, mülk dokumentteri   |
| `social_benefits`       | IDP kattaluusu, sotsialdyq tölemder                   |
| `property_rights`       | Basip alynğan aymaqtarda qalğan mülk                  |
| `employment_rights`     | Emgek myjzamy, boşotuu, jumušqa turuuğa jardam        |
| `family_law`            | Ajrylyşuu, balalardyn qamqorluğu, alimentter          |
| `healthcare_access`     | Medtsinalyq qamoralğanduuğu, engeldik dokumentaciyasi |
| `education_access`      | Mektepke kirüü, bilim beriü uquqtary                  |
| `financial_aid`         | Taryktuudan bir joluqluk qarcylar cardamy             |
| `psychological_support` | Psihologijalyq çakyryktar, keñeştöö                   |
| `other`                 | Tizimde cok ce aralas taqyryptar                      |

### 7-topt — Biuro Boyunça Bölünüş

| #   | Rapor                                        | Tür   |
| --- | -------------------------------------------- | ----- |
| 28  | Uquq konsultasija sany — biuro boyunça       | Okuja |
| 32  | Sotsialdyq konsultasija sany — biuro boyunça | Okuja |
| 33  | calpy konsultasija sany — biuro boyunça      | Okuja |

### 8-topt — Jaş Tobu Boyunça Bölünüş

| #   | Rapor                                                    | Tür     |
| --- | -------------------------------------------------------- | ------- |
| 26  | Uquq konsultasija sany — jaş tobu boyunça                | Okuja   |
| 27  | Uquq konsultasija alğan adamdar — jaş tobu boyunça       | Adamdar |
| 31a | Sotsialdyq konsultasija sany — jaş tobu boyunça          | Okuja   |
| 31b | Sotsialdyq konsultasija alğan adamdar — jaş tobu boyunça | Adamdar |
| 34  | Calpy konsultasija sany — jaş tobu boyunça               | Okuja   |

**Jaş toptor:**

| Mäán                | Jaş aralyğy |
| ------------------- | ----------- |
| `infant`            | 0–1         |
| `toddler`           | 1–3         |
| `pre_school`        | 3–6         |
| `middle_childhood`  | 6–12        |
| `young_teen`        | 12–14       |
| `teenager`          | 14–18       |
| `young_adult`       | 18–25       |
| `early_adult`       | 25–35       |
| `middle_aged_adult` | 35–55       |
| `old_adult`         | 55+         |

`birth_date` ornotulğanda cana `age_group` null bolğondo, tizme avtomattyk türdö eseptelet.

### 9-topt — Teg Izdöö

| #   | Rapor                                   | Tür     |
| --- | --------------------------------------- | ------- |
| 35  | Tegderi bar adamdardyn qoldoo cazuulary | Okuja   |
| 36  | Tegderi menen kattalğan adamdar         | Adamdar |

Bir ce anağurluq teg ID'lerdi parameter qatary ötkörüñüz. Donordun suuranmalary üçün qoldoniluuču.

### 10-topt — Üj-Büluö Bölüktörü

| #   | Rapor                                                        | Tür                |
| --- | ------------------------------------------------------------ | ------------------ |
| 37  | Uquq konsultasija alğan adamdar cana üj-büluö müçölörü       | Adamdar + Bölüktör |
| 38  | Sotsialdyq konsultasija alğan adamdar cana üj-büluö müçölörü | Adamdar + Bölüktör |
| 39  | Maalymatta kattalğan adamdar cana üj-büluö müçölörü          | Adamdar + Bölüktör |

Üj-büluö bölügü bir cazuusu. Adamdardyn calpy sany qatary sanalat.

## Üj Ahbandar Raporloru

Üj ajbandar menen bajlanyştuu Raporlor bölök jerde cetkiliktüü:

```http
GET /projects/:project_id/pet-reports/:report_type
```

Üj ajbandar raporloru türlör boyunça bölünüştü, eksele statusunu, sterilizasija statusunu cana adamdyn üj ajbandarna qarşy bolgonu sanalat.
