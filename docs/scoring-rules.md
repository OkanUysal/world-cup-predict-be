# Puanlama Kuralları

Toplam puan, channel içindeki tüm eventlerden kazanılan `points_awarded` değerlerinin toplanmasıyla hesaplanır.

---

## 1. Maç skoru (`match_score`)

Grup maçları ve eleme turları (son 32, son 16, çeyrek final, yarı final, final, 3.’lük maçı) için geçerlidir.

### Tahmin formatı

```json
{ "home_score": 2, "away_score": 1 }
```

### Eleme maçlarında skor tanımı

- **Uzatmalar dahildir** — admin sonucu ve kullanıcı tahmini, maç bittikten sonraki skoru (90 dk + uzatma) yansıtır.
- **Penaltı atışları dahil değildir** — penaltılarla belirlenen kazanan, skora yansımaz.

**Örnek:** Normal süre 1-1, uzatma 0-0, penaltılarla A takımı kazanır.

- Resmi skor (tahmin/sonuç): **1-1** (beraberlik)
- Penaltı kazananı skora yazılmaz

### Puan tablosu

| Durum | Puan | Açıklama |
|-------|------|----------|
| Tam skor doğru | **3** | Ev ve deplasman golü birebir aynı |
| Sonuç doğru, skor yanlış | **1** | Kazanan aynı veya ikisi de beraberlik tahmin etmiş |
| Yanlış | **0** | Kazanan/beraberlik tutmuyor |

**Sonuç doğru** = ev sahibi kazanır / deplasman kazanır / beraberlik aynı.

### Örnekler

| Tahmin | Sonuç | Puan | Neden |
|--------|-------|------|-------|
| 2-1 | 2-1 | 3 | Tam skor |
| 3-0 | 2-1 | 1 | Ev sahibi kazandı (ikisi de) |
| 1-1 | 2-2 | 1 | Beraberlik |
| 0-1 | 2-1 | 0 | Kazanan farklı |
| 1-1 | 1-1 (uzatma sonrası, penaltıyla A kazanır) | 3 | Resmi skor 1-1 |

> Tam skor (3 puan) ile sonuç (1 puan) **toplanmaz**; tam skor tutarsa sadece 3 puan verilir.

---

## 2. Şampiyon (`champion`)

Turnuvayı kazanan takım tahmini.

### Tahmin / sonuç formatı

```json
{ "team": "Brezilya" }
```

| Durum | Puan |
|-------|------|
| Doğru takım | **10** |
| Yanlış | **0** |

---

## 3. İkinci (`runner_up`)

Finali kaybeden takım tahmini.

### Tahmin / sonuç formatı

```json
{ "team": "Arjantin" }
```

| Durum | Puan |
|-------|------|
| Doğru takım | **5** |
| Yanlış | **0** |

---

## 4. Üçüncü (`third_place`)

3.’lük maçını kazanan takım tahmini.

### Tahmin / sonuç formatı

```json
{ "team": "Fransa" }
```

| Durum | Puan |
|-------|------|
| Doğru takım | **3** |
| Yanlış | **0** |

---

## Puan hesaplama akışı (admin)

1. Admin maç/event sonucunu girer (`POST /api/v1/admin/events/{id}/result`)
2. Admin puan hesaplamayı tetikler:
   - Tek event: `POST /api/v1/admin/events/{id}/calculate-scores`
   - Toplu: `POST /api/v1/admin/events/calculate-scores`
3. Her tahmin için `points_awarded` güncellenir
4. `user_scores.total_points` yeniden toplanır
5. Event durumu `completed` olur

---

## Özet tablo

| Event type | Doğru tahmin | Puan |
|------------|--------------|------|
| `match_score` | Kazanan / beraberlik | 1 |
| `match_score` | Tam skor (90 dk + uzatma, penaltı hariç) | 3 |
| `champion` | Şampiyon takım | 10 |
| `runner_up` | İkinci takım | 5 |
| `third_place` | Üçüncülük maçı kazananı | 3 |
