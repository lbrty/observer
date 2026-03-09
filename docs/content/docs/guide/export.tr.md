---
title: Dışa Aktarma
weight: 8
---

Observer, vaka verilerini ve raporları CSV dosyaları olarak dışa aktarabilir. Dışa aktarma için en az `consultant` proje rolü gereklidir.

## Vaka Verilerini Dışa Aktarma

Bir proje için tüm kişileri ve ilgili kayıtları dışa aktarın:

```http
GET /projects/:project_id/export/people?format=csv
```

İsteğe bağlı sorgu parametreleri:

| Parametre     | Açıklama                                         |
| ------------- | ------------------------------------------------ |
| `start`       | Kayıt tarihine göre filtrele (YYYY-AA-GG)        |
| `end`         | Kayıt tarihine göre filtrele (YYYY-AA-GG)        |
| `category_id` | Kırılganlık kategorisine göre filtrele           |
| `tag_id`      | Etikete göre filtrele                            |

Yanıt akış olarak iletilir — büyük veri kümeleri bağlantının zaman aşımına uğramaması için aşamalı olarak gönderilir.

## Destek Kayıtlarını Dışa Aktarma

Bir proje için danışmanlık kayıtlarını dışa aktarın:

```http
GET /projects/:project_id/export/support-records?format=csv&start=2024-01-01&end=2024-12-31
```

| Parametre | Açıklama                                          |
| --------- | ------------------------------------------------- |
| `start`   | `provided_at` tarihine göre filtrele (YYYY-AA-GG) |
| `end`     | `provided_at` tarihine göre filtrele (YYYY-AA-GG) |
| `type`    | `legal` veya `social`                             |
| `sphere`  | Destek alanı (örn. `housing_assistance`)          |

## Göç Kayıtlarını Dışa Aktarma

Hareket/yerinden edilme geçmişini dışa aktarın:

```http
GET /projects/:project_id/export/migration-records?format=csv
```

## CSV Formatı

Tüm dışa aktarmalar başlık satırlı UTF-8 CSV kullanır. Alanlar, API yanıtlarıyla aynı yapıyı takip eder. Hassas alanlar (`contact`, `personal`, `documents`), proje izin bayraklarınıza göre dahil edilir veya çıkarılır — API okumalarına uygulanan aynı kurallar dışa aktarmalara da uygulanır.

## Web Arayüzünde İndirme

1. Bir proje açın
2. Gezinti bölümünde **Dışa Aktarma** seçeneğine tıklayın
3. Kayıt türünü ve tarih aralığını seçin
4. **CSV İndir** seçeneğine tıklayın

Dosya doğrudan tarayıcınıza indirilir.

{{% callout type="note" %}}
Dışa aktarmalar projenizle sınırlıdır. Atanmadığınız bir projeden veri dışa aktaramazsınız ve hassasiyet bayraklarınız uygulanır — `can_view_personal` kapalıysa, dışa aktarılan CSV ulusal kimlik numaralarını veya doğum tarihlerini içermeyecektir.
{{% /callout %}}
