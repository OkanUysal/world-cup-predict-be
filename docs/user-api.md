# User API — DTO & Response Dokümantasyonu

Base URL: `/api/v1`

Tüm tarih/saat alanları **UTC (GMT)** formatındadır: `2006-01-02T15:04:05Z07:00`

---

## Kimlik Doğrulama

Korumalı endpoint'ler için header:

```
Authorization: Bearer <access_token>
```

Login veya register sonrası dönen `access_token` kullanılır. Token süresi **24 saat**tir.

---

## Ortak Tipler

### ErrorResponse

Hata durumlarında dönen standart yapı:

```json
{
  "error": "hata mesajı"
}
```

### UserProfile

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `id` | string (uuid) | Evet | Kullanıcı ID |
| `name` | string | Evet | Giriş kullanıcı adı (değişmez) |
| `nickname` | string | Hayır | Görünen isim (boş olabilir) |
| `display_name` | string | Evet | `nickname` doluysa nickname, değilse `name` |
| `role` | string | Evet | `user` veya `admin` |
| `channel_id` | string (uuid) | Hayır | Kullanıcının channel ID'si |
| `channel` | ChannelSummary | Hayır | Kanal bilgisi (sadece `/me`) |
| `total_points` | integer | Hayır | Toplam puan (sadece `/me`) |

### ChannelSummary

| Alan | Tip | Açıklama |
|------|-----|----------|
| `id` | string (uuid) | Channel ID |
| `code` | string | Channel kodu |
| `name` | string | Kanal açıklaması (admin'in girdiği isim) |

### Event

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `id` | string (uuid) | Evet | Event ID |
| `type` | string | Evet | `match_score`, `champion`, `runner_up`, `third_place` |
| `title` | string | Evet | Event başlığı |
| `metadata` | object | Evet | Event tipine göre değişir (aşağıda) |
| `deadline` | string (datetime) | Evet | Tahmin deadline (UTC) |
| `status` | string | Evet | `open`, `locked`, `completed` |
| `result` | object | Hayır | Admin sonucu (sadece girildiyse) |
| `created_at` | string (datetime) | Evet | Oluşturulma zamanı (UTC) |

**Event status açıklaması**

| Status | Anlam |
|--------|-------|
| `open` | Deadline öncesi; tahmin girilebilir/güncellenebilir |
| `locked` | Deadline geçti, sonuç yok; channel tercihleri görünür |
| `completed` | Sonuç girildi ve puan hesaplandı |

### Prediction

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `id` | string (uuid) | Evet | Tahmin ID |
| `event_id` | string (uuid) | Evet | Event ID |
| `user_id` | string (uuid) | Evet | Kullanıcı ID |
| `user_name` | string | Hayır | Görünen isim: nickname veya name |
| `choice` | object | Evet | Tahmin içeriği (event tipine göre) |
| `points_awarded` | integer | Evet | Kazanılan puan |
| `created_at` | string (datetime) | Evet | Oluşturulma (UTC) |
| `updated_at` | string (datetime) | Evet | Son güncelleme (UTC) |

### UserScore

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `user_id` | string (uuid) | Evet | Kullanıcı ID |
| `user_name` | string | Evet | Görünen isim: nickname veya name |
| `channel_id` | string (uuid) | Evet | Channel ID |
| `total_points` | integer | Evet | Toplam puan |
| `updated_at` | string (datetime) | Evet | Son güncelleme (UTC) |

---

## Event Metadata & Choice Formatları

### `match_score` — Maç skoru

**metadata** (event detayında):

```json
{
  "home_team": "Meksika",
  "away_team": "Güney Afrika",
  "group": "A",
  "match_no": 1,
  "venue": "Mexico City",
  "kickoff_gmt": "2026-06-11T19:00:00Z"
}
```

**choice** (tahmin gönderirken):

```json
{
  "home_score": 2,
  "away_score": 1
}
```

### `champion` / `runner_up` / `third_place` — Tek takım seçimi

**metadata**:

```json
{
  "teams": ["ABD", "Almanya", "Arjantin", "..."]
}
```

**choice**:

```json
{
  "team": "Brezilya"
}
```

---

## Auth Endpoint'leri (Public)

### POST `/auth/register`

Yeni kullanıcı kaydı. Admin tarafından oluşturulmuş bir channel koduna kayıt olunur.

**Request DTO**

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `name` | string | Evet | Kullanıcı adı (channel içinde unique) |
| `password` | string | Evet | Min. 6 karakter |
| `channel_code` | string | Evet | Channel kodu (örn. `ABC123`) |

```json
{
  "name": "ahmet",
  "password": "secret123",
  "channel_code": "ABC123"
}
```

**Response `201 Created` — AuthResponse**

| Alan | Tip | Açıklama |
|------|-----|----------|
| `access_token` | string | JWT access token |
| `expires_at` | string | Token bitiş zamanı (UTC) |
| `user` | UserProfile | Kullanıcı bilgisi |

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2026-06-07T20:00:00Z",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "ahmet",
    "role": "user",
    "channel_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
  }
}
```

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid request body` / validasyon mesajı |
| 404 | `channel not found` |
| 409 | `user already exists in this channel` |

---

### POST `/auth/login`

Kullanıcı girişi.

**Request DTO**

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `name` | string | Evet | Kullanıcı adı |
| `password` | string | Evet | Şifre |
| `channel_code` | string | Hayır | Channel kodu (admin için opsiyonel) |

```json
{
  "name": "ahmet",
  "password": "secret123",
  "channel_code": "ABC123"
}
```

**Response `200 OK` — AuthResponse**

Register ile aynı yapı.

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid request body` |
| 401 | `invalid credentials` |
| 404 | `channel not found` |

---

## Kullanıcı Endpoint'leri (JWT gerekli)

### GET `/me`

Giriş yapmış kullanıcının profil ve puan bilgisi.

**Request:** Body yok.

**Response `200 OK` — UserProfile**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "ahmet",
  "nickname": "AhmetGK",
  "display_name": "AhmetGK",
  "role": "user",
  "channel_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "channel": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "code": "ABC123",
    "name": "Okul Ligi"
  },
  "total_points": 15
}
```

**Hatalar**

| HTTP | error |
|------|-------|
| 401 | `unauthorized` / `invalid or expired token` |

---

### PATCH `/me/nickname`

Kullanıcının görünen ismini (nickname) günceller. Boş string gönderilirse nickname temizlenir.

**Request DTO**

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `nickname` | string | Evet | 2–64 karakter; `""` ile temizlenir |

```json
{
  "nickname": "AhmetGK"
}
```

**Response `200 OK` — UserProfile** (güncellenmiş profil)

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid request body` / validasyon mesajı |
| 401 | `unauthorized` |

---

### GET `/leaderboard`

Kullanıcının channel'ındaki puan sıralaması. `user_name` alanı nickname varsa nickname, yoksa giriş adını (`name`) gösterir.

**Request:** Body yok.

**Response `200 OK` — UserScore[]**

```json
[
  {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_name": "ahmet",
    "channel_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "total_points": 15,
    "updated_at": "2026-06-11T22:00:00Z"
  },
  {
    "user_id": "660e8400-e29b-41d4-a716-446655440001",
    "user_name": "mehmet",
    "channel_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "total_points": 12,
    "updated_at": "2026-06-11T22:00:00Z"
  }
]
```

Sıralama: `total_points` azalan, eşitlikte `user_name` artan.

**Hatalar**

| HTTP | error |
|------|-------|
| 401 | `unauthorized` |
| 403 | `channel membership required` |

---

### GET `/events`

Event listesi. Deadline'a göre artan sıralı.

**Query parametreleri**

| Param | Tip | Varsayılan | Açıklama |
|-------|-----|------------|----------|
| `status` | string | `open` | `open`, `locked`, `pending`, `completed` |

**Status filtreleri**

| Filtre | Açıklama |
|--------|----------|
| `open` | Tahmin girilebilir eventler |
| `locked` | Deadline geçmiş, sonuç bekleyen |
| `pending` | `locked` ile aynı (deadline geçmiş, sonuç yok) |
| `completed` | Sonuçlanmış eventler |

**Response `200 OK` — EventWithPrediction[]**

```json
[
  {
    "event": {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "type": "match_score",
      "title": "WC2026: Grup A: Meksika vs Güney Afrika",
      "metadata": {
        "home_team": "Meksika",
        "away_team": "Güney Afrika",
        "group": "A",
        "match_no": 1,
        "venue": "Mexico City",
        "kickoff_gmt": "2026-06-11T19:00:00Z"
      },
      "deadline": "2026-06-11T18:00:00Z",
      "status": "open",
      "created_at": "2026-06-06T20:00:00Z"
    },
    "my_prediction": {
      "id": "...",
      "event_id": "...",
      "user_id": "...",
      "choice": { "home_score": 2, "away_score": 1 },
      "points_awarded": 0,
      "created_at": "2026-06-06T21:00:00Z",
      "updated_at": "2026-06-06T21:00:00Z"
    }
  }
]
```

`my_prediction` alanı tahmin yoksa response'ta bulunmaz (`omitempty`).

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid status filter` |
| 401 | `unauthorized` |

---

### GET `/events/{id}`

Tek event detayı ve kullanıcının kendi tahmini.

**Path parametreleri**

| Param | Tip | Açıklama |
|-------|-----|----------|
| `id` | uuid | Event ID |

**Response `200 OK` — EventDetailResponse**

| Alan | Tip | Açıklama |
|------|-----|----------|
| `event` | Event | Event detayı |
| `my_prediction` | Prediction | Kullanıcının tahmini (varsa) |

```json
{
  "event": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "type": "champion",
    "title": "WC2026: Şampiyon",
    "metadata": {
      "teams": ["ABD", "Almanya", "Arjantin", "Brezilya"]
    },
    "deadline": "2026-06-10T19:00:00Z",
    "status": "open",
    "created_at": "2026-06-06T20:00:00Z"
  },
  "my_prediction": null
}
```

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid event id` |
| 401 | `unauthorized` |
| 404 | `event not found` |

---

### PUT `/events/{id}/prediction`

Tahmin oluştur veya güncelle. Sadece `status=open` ve deadline öncesinde.

**Path parametreleri**

| Param | Tip | Açıklama |
|-------|-----|----------|
| `id` | uuid | Event ID |

**Request DTO — PredictionRequest**

| Alan | Tip | Zorunlu | Açıklama |
|------|-----|---------|----------|
| `choice` | object | Evet | Event tipine göre (yukarıda) |

Maç skoru örneği:

```json
{
  "choice": {
    "home_score": 2,
    "away_score": 1
  }
}
```

Şampiyon/ikinci/üçüncü örneği:

```json
{
  "choice": {
    "team": "Brezilya"
  }
}
```

**Response `200 OK` — Prediction**

```json
{
  "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
  "event_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "choice": { "home_score": 2, "away_score": 1 },
  "points_awarded": 0,
  "created_at": "2026-06-06T21:00:00Z",
  "updated_at": "2026-06-06T21:00:00Z"
}
```

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid request body` / `invalid match_score choice: ...` |
| 401 | `unauthorized` |
| 403 | `event is closed for predictions` |
| 404 | `event not found` |

---

### GET `/events/{id}/predictions`

Aynı channel'daki tüm kullanıcıların tahminleri. **Deadline geçtikten sonra** görünür (`locked` veya `completed`).

**Path parametreleri**

| Param | Tip | Açıklama |
|-------|-----|----------|
| `id` | uuid | Event ID |

**Response `200 OK` — Prediction[]**

```json
[
  {
    "id": "...",
    "event_id": "...",
    "user_id": "...",
    "user_name": "ahmet",
    "choice": { "home_score": 2, "away_score": 1 },
    "points_awarded": 3,
    "created_at": "2026-06-06T21:00:00Z",
    "updated_at": "2026-06-06T21:00:00Z"
  },
  {
    "id": "...",
    "event_id": "...",
    "user_id": "...",
    "user_name": "mehmet",
    "choice": { "home_score": 1, "away_score": 1 },
    "points_awarded": 1,
    "created_at": "2026-06-06T21:30:00Z",
    "updated_at": "2026-06-06T21:30:00Z"
  }
]
```

**Hatalar**

| HTTP | error |
|------|-------|
| 400 | `invalid event id` |
| 401 | `unauthorized` |
| 403 | `predictions are not visible until deadline passes` |
| 403 | `channel membership required` |
| 404 | `event not found` |

---

## Puanlama Kuralları

| Event type | Doğru tahmin | Puan |
|------------|--------------|------|
| `match_score` | Kazanan / beraberlik | 1 |
| `match_score` | Tam skor | 3 |
| `champion` | Doğru takım | 10 |
| `runner_up` | Doğru takım | 5 |
| `third_place` | Doğru takım | 3 |

---

## Örnek Akış

```
1. POST /auth/register        → access_token al
2. GET  /events?status=open   → açık eventleri listele
3. PUT  /events/{id}/prediction → tahmin gir
4. GET  /events/{id}          → kendi tahminini kontrol et
   (deadline sonrası)
5. GET  /events/{id}/predictions → channel tercihlerini gör
6. GET  /leaderboard          → sıralamayı gör
7. GET  /me                   → profil + toplam puan
```
