---
title: Raporlar
weight: 7
---

Observer, gerçek Ukraynalı STK bağışçı raporlama yükümlülüklerine uyacak şekilde tasarlanmış 39 yapılandırılmış rapor türü üretir. Tüm raporlar tek bir projeye kapsamlıdır ve bir tarih aralığına göre filtrelenir.

## Rapor Çalıştırma

Tüm raporlara şu adresten erişilebilir:

```http
GET /projects/:project_id/reports/:report_type?start=2024-01-01&end=2024-12-31
```

`:report_type` yerine aşağıdaki tablolardaki tanımlayıcılardan birini kullanın.

## Sayım Türleri

Her rapor üç sayım yönteminden birini kullanır:

| Tür        | Anlam                                                                       |
| ---------- | --------------------------------------------------------------------------- |
| **Olaylar** | `support_records` satırları — üç danışmanlık alan bir kişi 3 olarak sayılır |
| **Kişiler** | Benzersiz bireyler — aynı kişi kaç kayda sahip olursa olsun bir kez sayılır |
| **Birimler** | Benzersiz aile birimleri (bir hane = bir birim)                            |

Tarih filtrelemesi danışmanlıklar için `support_records.provided_at`, kayıt raporları için `people.registered_at` kullanır — `created_at` değil.

## Rapor Grupları

### Grup 1 — Genel Danışmanlık Sayıları

| # | Rapor | Tür |
|---|-------|-----|
| 1 | Tüm türlerdeki toplam danışmanlıklar | Olaylar |
| 2 | Toplam hukuki danışmanlıklar | Olaylar |
| 3 | Toplam sosyal danışmanlıklar | Olaylar |

### Grup 2 — Cinsiyet Dağılımı

| # | Rapor | Tür |
|---|-------|-----|
| 12 | Dönemde kayıtlı erkekler | Kişiler |
| 13 | Dönemde kayıtlı kadınlar | Kişiler |
| 14 | Hukuki danışmanlık alan kadınlar | Kişiler |
| 15 | Sosyal danışmanlık alan kadınlar | Kişiler |
| 16 | Hukuki danışmanlık alan erkekler | Kişiler |
| 17 | Sosyal danışmanlık alan erkekler | Kişiler |

Raporlar 14–17, danışmanlık olaylarını değil, benzersiz kişileri sayar.

### Grup 3 — Coğrafi / YDB Durumu

| # | Rapor | Tür |
|---|-------|-----|
| 4 | Toplam kayıtlı kişiler | Kişiler |
| 5–6 | Çatışma bölgesinden kayıtlı kişiler | Kişiler |
| 7–10 | Çatışma bölgesinden hukuki/sosyal danışmanlık alan kişiler | Kişiler |
| 11 | YDB olmayan kayıtlı kişiler | Kişiler |

YDB durumu otomatik olarak türetilir: `origin_place_id → places → states.conflict_zone`. Raporlar 5–6 ve 7–10, çatışma bölgesi etiketine göre parametrelendirilmiştir.

### Grup 4 — Kırılganlık Kategorisi Dağılımı

| # | Rapor | Tür |
|---|-------|-----|
| 18 | Kayıtlı kişiler — kırılganlık kategorisine göre | Kişiler |
| 19 | Sosyal danışmanlık alan kişiler — kategoriye göre | Kişiler |
| 20 | Hukuki danışmanlık alan kişiler — kategoriye göre | Kişiler |

Atanmış kategorisi olmayan kişiler "kategorisiz" grubuna girer.

### Grup 5 — Mevcut Kalış Bölgesi

| # | Rapor | Tür |
|---|-------|-----|
| 21 | Kayıtlı kişiler — mevcut bölgeye göre | Kişiler |
| 22 | Hukuki danışmanlık alan kişiler — bölgeye göre | Kişiler |
| 23 | Sosyal danışmanlık alan kişiler — bölgeye göre | Kişiler |

### Grup 6 — Destek Alanı Dağılımı

| # | Rapor | Tür |
|---|-------|-----|
| 24 | Hukuki danışmanlık sayısı — alana göre | Olaylar |
| 25 | Hukuki danışmanlık alan kişiler — alana göre | Kişiler |
| 29 | Sosyal danışmanlık sayısı — alana göre | Olaylar |
| 30 | Sosyal danışmanlık alan kişiler — alana göre | Kişiler |

**Destek alanları:**

| Değer | Açıklama |
|-------|----------|
| `housing_assistance` | Konut hakları, tahliye, sosyal konut |
| `document_recovery` | Pasaportlar, doğum belgeleri, mülk belgeleri |
| `social_benefits` | YDB kaydı, sosyal ödemeler |
| `property_rights` | İşgal altındaki topraklarda bırakılan mülk |
| `employment_rights` | İş hukuku, işten çıkarma, iş bulma |
| `family_law` | Boşanma, velayet, nafaka |
| `healthcare_access` | Sağlık sigortası, engellilik belgelendirmesi |
| `education_access` | Okul kaydı, eğitim hakları |
| `financial_aid` | Acil mali yardım |
| `psychological_support` | Ruh sağlığı yönlendirmeleri, danışmanlık |
| `other` | Listede olmayan veya kesişen konular |

### Grup 7 — Ofis Dağılımı

| # | Rapor | Tür |
|---|-------|-----|
| 28 | Hukuki danışmanlık sayısı — ofise göre | Olaylar |
| 32 | Sosyal danışmanlık sayısı — ofise göre | Olaylar |
| 33 | Toplam danışmanlık sayısı — ofise göre | Olaylar |

### Grup 8 — Yaş Grubu Dağılımı

| # | Rapor | Tür |
|---|-------|-----|
| 26 | Hukuki danışmanlık sayısı — yaş grubuna göre | Olaylar |
| 27 | Hukuki danışmanlık alan kişiler — yaş grubuna göre | Kişiler |
| 31a | Sosyal danışmanlık sayısı — yaş grubuna göre | Olaylar |
| 31b | Sosyal danışmanlık alan kişiler — yaş grubuna göre | Kişiler |
| 34 | Toplam danışmanlık sayısı — yaş grubuna göre | Olaylar |

**Yaş grupları:**

| Değer | Yaş aralığı |
|-------|-------------|
| `infant` | 0–1 |
| `toddler` | 1–3 |
| `pre_school` | 3–6 |
| `middle_childhood` | 6–12 |
| `young_teen` | 12–14 |
| `teenager` | 14–18 |
| `young_adult` | 18–25 |
| `early_adult` | 25–35 |
| `middle_aged_adult` | 35–55 |
| `old_adult` | 55+ |

`birth_date` ayarlandığında ve `age_group` null olduğunda, uygulama grubu otomatik olarak hesaplar.

### Grup 9 — Etiket Arama

| # | Rapor | Tür |
|---|-------|-----|
| 35 | Belirli etiketlere sahip kişiler için destek kayıtları | Olaylar |
| 36 | Belirli etiketlerle kayıtlı kişiler | Kişiler |

Bir veya daha fazla etiket ID'sini parametre olarak geçin. Geçici bağışçı sorguları için kullanışlıdır.

### Grup 10 — Aile Birimleri

| # | Rapor | Tür |
|---|-------|-----|
| 37 | Hukuki danışmanlık alan kişiler ve aile üyeleri | Kişiler + Birimler |
| 38 | Sosyal danışmanlık alan kişiler ve aile üyeleri | Kişiler + Birimler |
| 39 | Dönemde kayıtlı kişiler ve aile üyeleri | Kişiler + Birimler |

Bir aile birimi, bir hane kaydıdır. Sayımlar hem toplam bireyler hem de benzersiz hane birimleri olarak döndürülür.

## Evcil Hayvan Raporları

Hayvanlarla ilgili raporlar ayrıca şu adreste mevcuttur:

```http
GET /projects/:project_id/pet-reports/:report_type
```

Evcil hayvan raporları tür dağılımını, aşılama durumunu, kısırlaştırma durumunu ve hayvan-insan oranlarını kapsar.
