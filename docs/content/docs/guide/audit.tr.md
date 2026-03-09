---
title: Denetim Günlüğü
weight: 9
---

Observer, vaka verileri üzerindeki her oluşturma, güncelleme ve silme işlemi için bir denetim günlüğü kaydeder. Denetim günlükleri yalnızca ekleme amaçlıdır — düzenlenemez veya silinemez.

## Neler Günlüğe Kaydedilir

| Kategori          | Kaydedilen işlemler                                    |
| ----------------- | ------------------------------------------------------ |
| Kişiler           | Kişi kayıtlarını oluşturma, güncelleme, silme          |
| Destek kayıtları  | Danışmanlıkları oluşturma, güncelleme, silme           |
| Göç kayıtları     | Hareket kayıtlarını oluşturma, güncelleme, silme       |
| Hanehalkları      | Hane ve üye kayıtlarını oluşturma, güncelleme, silme   |
| Notlar            | Vaka notlarını oluşturma, güncelleme, silme            |
| Belgeler          | Yükleme, meta veri güncelleme, belge silme             |
| Evcil hayvanlar   | Evcil hayvan kayıtlarını oluşturma, güncelleme, silme  |
| İzinler           | Proje izinlerini atama, güncelleme, iptal etme         |

Kimlik doğrulama olayları (giriş, çıkış, token yenileme) proje denetim günlüğünde yer almaz — sunucu günlüklerinde görünürler.

## Denetim Günlüğünü Görüntüleme

Denetim günlüğüne yalnızca proje yöneticileri ve sahipleri erişebilir.

```http
GET /projects/:project_id/audit?page=1&per_page=50
```

| Parametre  | Açıklama                                        |
| ---------- | ----------------------------------------------- |
| `page`     | Sayfa numarası (varsayılan 1)                   |
| `per_page` | Sayfa başına sonuç (varsayılan 50)              |
| `actor_id` | Değişikliği yapan kullanıcıya göre filtrele     |
| `start`    | Tarihe göre filtrele (YYYY-AA-GG)               |
| `end`      | Tarihe göre filtrele (YYYY-AA-GG)               |

### Yanıt formatı

Her denetim kaydı şunları içerir:

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

| Alan          | Açıklama                                                   |
| ------------- | ---------------------------------------------------------- |
| `actor_id`    | İşlemi gerçekleştiren kullanıcı                            |
| `actor_ip`    | İsteğin IP adresi                                          |
| `action`      | `create`, `update` veya `delete`                           |
| `entity_type` | Kayıt türü (`person`, `support_record`, `note` vb.)        |
| `entity_id`   | Etkilenen kaydın ULID'i                                    |

## Web Arayüzünde Denetim

1. Bir proje açın
2. Gezinti bölümünde **Denetim Günlüğü** seçeneğine tıklayın
3. Tarih aralığına veya kullanıcıya göre filtreleyin

{{% callout type="note" %}}
Denetim günlüğü erişimi `manager` veya `owner` proje rolü gerektirir. Danışmanlar ve görüntüleyiciler denetim günlüğünü göremez.
{{% /callout %}}
