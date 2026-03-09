---
title: İzinler
weight: 6
---

Observer iki seviyeli bir izin modeli kullanır. Her kullanıcının yönetici işlevlerine erişimi kontrol eden bir **platform rolü** vardır. Her proje içinde, kullanıcıların dava verileriyle neler yapabileceğini kontrol eden bir **proje rolü** de vardır. Bunun üzerine, üç hassasiyet bayrağı yanıtlarda hangi alanların görüneceğini kısıtlar.

## Platform Rolleri

Platform rolleri kullanıcı hesabı oluşturulduğunda ayarlanır ve yalnızca bir yönetici tarafından değiştirilebilir.

| Rol          | Ne yapabilirler                                                           |
| ------------ | ------------------------------------------------------------------------- |
| `admin`      | Her şeye tam erişim — kullanıcılar, projeler, referans verileri, tüm davalar |
| `staff`      | Referans verileri ve projeleri yönetme; diğer kullanıcıları yönetemez    |
| `consultant` | Atandıkları projelerde çalışma; yönetici paneline erişim yok             |
| `guest`      | Atandıkları projelerde salt okunur erişim                                |

Platform yöneticileri tüm projelerde otomatik olarak **proje sahibi** olarak hareket eder — proje düzeyindeki izin kontrollerini tamamen atlarlar.

## Proje Rolleri

Proje rolleri belirli bir proje içinde geçerlidir. Bir kullanıcının projeye erişebilmesi için açıkça atanması gerekir.

| Rol          | Sıra | Ne yapabilirler                                        |
| ------------ | ---- | ------------------------------------------------------ |
| `owner`      | 4    | Projeyi silme dahil her şey                            |
| `manager`    | 3    | Tüm dava veri işlemleri + ekip üyelerini yönetme       |
| `consultant` | 2    | Dava oluşturma, okuma, güncelleme ve veri dışa aktarma |
| `viewer`     | 1    | Dava verilerine salt okunur erişim                     |

### Eylem Gereksinimleri

| Eylem              | Minimum proje rolü |
| ------------------ | ------------------ |
| Veri okuma         | `viewer`           |
| Kayıt oluşturma    | `consultant`       |
| Kayıt güncelleme   | `consultant`       |
| Kayıt silme        | `manager`          |
| Üye yönetimi       | `manager`          |
| Veri dışa aktarma  | `consultant`       |

## Hassasiyet Bayrakları

Her proje izninin bir yönetici veya sahip tarafından kullanıcı başına, proje başına ayarlanan üç bağımsız boolean bayrağı vardır.

| Bayrak                | Neyi kontrol eder                                                           |
| --------------------- | --------------------------------------------------------------------------- |
| `can_view_contact`    | Kişi kayıtlarındaki telefon numaraları ve e-posta adresleri                 |
| `can_view_personal`   | Tam ad, doğum tarihi, ulusal kimlik (`external_id`), onay bilgisi           |
| `can_view_documents`  | Belge dosyası erişimi ve belge meta verileri                                |

Bir bayrak kapalı olduğunda, ilgili alanlar API yanıtlarından ve CSV dışa aktarmalarından çıkarılır — veriler veritabanında kalır ancak o kullanıcıya gönderilmez. `can_view_personal: false` olan bir danışman kayıt oluşturabilse de dışa aktarma yoluyla ulusal kimliklere veya doğum tarihlerine erişemez.

## İzin Atama

Yalnızca platform yöneticileri ve proje sahipleri/yöneticileri proje izinleri atayabilir.

### Kullanıcıyı projeye ekle

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

### Mevcut izni güncelle

```http
PUT /admin/projects/:project_id/permissions/:permission_id
```

### Kullanıcıyı projeden kaldır

```http
DELETE /admin/projects/:project_id/permissions/:permission_id
```

## Yaygın Ayarlar

### Saha çalışanı (consultant, kısıtlı kişisel veri)

```json
{
  "role": "consultant",
  "can_view_contact": true,
  "can_view_personal": false,
  "can_view_documents": false
}
```

Koordinasyon için telefon numaraları görünür; ulusal kimlik veya doğum tarihlerine erişim yok.

### Süpervizör (manager, tam erişim)

```json
{
  "role": "manager",
  "can_view_contact": true,
  "can_view_personal": true,
  "can_view_documents": true
}
```

### Dış denetçi (viewer, hassas veri yok)

```json
{
  "role": "viewer",
  "can_view_contact": false,
  "can_view_personal": false,
  "can_view_documents": false
}
```
