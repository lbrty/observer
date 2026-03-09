---
title: Uquktar
weight: 6
---

Observer eki deñgeeli uquk modelin qoldonot. Ar bir qoldonuuçunun admin funksijalarğa cetüünü basqarğan **platforma rolu** bar. Ar bir proekttin içinde qoldonuuçulardyn öz maalymat bazasyn basqarğan **proekt rolu** da bar. Bundan tyşqary, üç qaptamçy belgi cafapta qaysyl ölçömdör çyğarylğanyn çektejt.

## Platforma Roldaru

Platforma roldaru qoldonuuçu akkauntu dörötkön kezde ornotuluşat cana tek admin taraby özgörtülö alaşat.

| Rol          | Emne isteşe bolot                                                          |
| ------------ | -------------------------------------------------------------------------- |
| `admin`      | Bardyq nersege toluk cetüü — qoldonuuçular, proektter, maalymattyb bardyğy |
| `staff`      | Maalymatty cana proektterdi basqaruu; basqa qoldonuuçulardy basqara albajt |
| `consultant` | Bölüngön proektterde iştöö; admin paneline cetüü coq                       |
| `guest`      | Bölüngön proektterde oqup qana cetüü                                       |

Platforma adminderi bardyq proektterde avtomattyk türdö **proekt eesi** bolup, proekt deñgeelindegi uquk tekşerüüdön ttü.

## Proekt Roldory

Proekt roldaru belgili bir proektin içinde işteşet. Qoldonuuçu proektke cetüü üçün anğa açyk bölünüşü kerek.

| Rol          | Ranq | Emne isteşe bolot                                          |
| ------------ | ---- | ---------------------------------------------------------- |
| `owner`      | 4    | Proektti öçürüü qoşo bardyğy                               |
| `manager`    | 3    | Bardyq mazmunat operasijalary + komanda müçölörün basqaruu |
| `consultant` | 2    | Proektterde iş alyp barat, maalymat, qoşuu cana cañyrtuu   |
| `viewer`     | 1    | Mazmunatqa oqup qana cetüü                                 |

### Amaldar Talaptary

| Amal               | Minimal proekt rolu |
| ------------------ | ------------------- |
| Maalymat oquu      | `viewer`            |
| Cazuu tüzüü        | `consultant`        |
| Cazuu cañyrtuu     | `consultant`        |
| Cazuu öçürüü       | `manager`           |
| Müçölördü basqaruu | `manager`           |
| Maalymat çyğaruu   | `consultant`        |

## Qaptamçy Belgileri

Ar bir proekt uquğunda üç özümçöl bool belgi bar. Bular manager ce eeçi taraby ar bir qoldonuuçu, ar bir proekt üçün ornotuluşat.

| Belgi                | Emne basqarat                                                            |
| -------------------- | ------------------------------------------------------------------------ |
| `can_view_contact`   | Adam cazuularyndağy telefon nomerleri cana email adresteri               |
| `can_view_personal`  | Toluk at, tuulğan künü, uluttuk ID (`external_id`), macburlyq maaalymaty |
| `can_view_documents` | Dokumentterge uruk saaty bar                                             |

Belgi öçük bolso, tiyiştüü ölçömdör API coobunan cana CSV eksportton çetke cyğarylat — maalymat maalymat bazasynda qalat, biroq ağa qoldonuuçuğa cönötülböjt. `can_view_personal: false` belgilengen konsultant cazuulardy tüzö alsa da, eksport arqyluu uluttuk ID ce tuulğan künün cb. maalymatyn ala albajt.

## Uquqtar Bölüştürüü

Proekt uquqtaryn tek platforma adminderi cana proekt eeleri/managerleri bölüştürö alat.

### Qoldonuuçunu proektke qoşuu

```http
POST /admin/projects/:project_id/permissions
```

```json
{
  "user_id": "01J...",
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

### Bar uquqtu cañyrtuv

```http
PUT /admin/projects/:project_id/permissions/:permission_id
```

### Qoldonuuçunu proektden çyğaruu

```http
DELETE /admin/projects/:project_id/permissions/:permission_id
```

## Tarqatylğan Ornotuular

### Tala işçisi (consultant, çektelgen ceke maalymat)

```json
{
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

Koordinasija üçün telefon nomerleri açyk; uluttuk ID ce tuulğan küngö cetüü coq.

### Nazarçy (manager, toluk cetüü)

```json
{
  "role": "manager",
  "can_view_contact": true,
  "can_view_personal": true,
  "can_view_documents": true
}
```

### Tyşky audiitor (viewer, qaptamçy maalymat coq)

```json
{
  "role": "viewer",
  "can_view_contact": false,
  "can_view_personal": false,
  "can_view_documents": false
}
```
