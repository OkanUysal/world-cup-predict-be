package swagger

import (
	"fmt"
	"net/http"
	"strings"
)

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = strings.Split(proto, ",")[0]
		scheme = strings.TrimSpace(scheme)
	}

	host := r.Host
	if fwd := r.Header.Get("X-Forwarded-Host"); fwd != "" {
		host = strings.Split(fwd, ",")[0]
		host = strings.TrimSpace(host)
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

func buildSpec(baseURL string) map[string]any {
	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "World Cup Predict API",
			"description": "Dünya Kupası tahmin uygulaması backend API'si. Channel bazlı yarışma, JWT auth ve admin event yönetimi.",
			"version":     "1.0.0",
		},
		"servers": []map[string]string{
			{"url": baseURL, "description": "Current host"},
		},
		"tags": []map[string]string{
			{"name": "Health", "description": "Sağlık kontrolü"},
			{"name": "Auth", "description": "Kayıt ve giriş"},
			{"name": "User", "description": "Kullanıcı işlemleri"},
			{"name": "Events", "description": "Event listeleme ve tahminler"},
			{"name": "Admin Channels", "description": "Admin channel yönetimi"},
			{"name": "Admin Events", "description": "Admin event yönetimi"},
		},
		"paths": buildPaths(),
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
					"description":  "Login veya register sonrası dönen access_token",
				},
			},
			"schemas": buildSchemas(),
		},
	}
}

func buildPaths() map[string]any {
	return map[string]any{
		"/health": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Health"},
				"summary":     "Health check",
				"operationId": "healthCheck",
				"responses": map[string]any{
					"200": jsonResponse("OK", ref("HealthResponse")),
				},
			},
		},
		"/api/v1/auth/register": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Auth"},
				"summary":     "Kullanıcı kaydı",
				"description": "Mevcut bir channel koduna kayıt olur.",
				"operationId": "register",
				"requestBody": jsonBody(ref("RegisterRequest"), "Kayıt bilgileri"),
				"responses": map[string]any{
					"201": jsonResponse("Kayıt başarılı", ref("AuthResponse")),
					"400": errorResponse("Geçersiz istek"),
					"404": errorResponse("Channel bulunamadı"),
					"409": errorResponse("Kullanıcı zaten mevcut"),
				},
			},
		},
		"/api/v1/auth/login": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Auth"},
				"summary":     "Giriş",
				"description": "Kullanıcı channel_code ile giriş yapar. Admin kullanıcılar channel_code olmadan giriş yapabilir.",
				"operationId": "login",
				"requestBody": jsonBody(ref("LoginRequest"), "Giriş bilgileri"),
				"responses": map[string]any{
					"200": jsonResponse("Giriş başarılı", ref("AuthResponse")),
					"400": errorResponse("Geçersiz istek"),
					"401": errorResponse("Geçersiz kimlik bilgileri"),
					"404": errorResponse("Channel bulunamadı"),
				},
			},
		},
		"/api/v1/me": map[string]any{
			"get": map[string]any{
				"tags":        []string{"User"},
				"summary":     "Profil bilgisi",
				"operationId": "getMe",
				"security":    bearer(),
				"responses": map[string]any{
					"200": jsonResponse("Profil", ref("UserProfile")),
					"401": errorResponse("Yetkisiz"),
				},
			},
		},
		"/api/v1/leaderboard": map[string]any{
			"get": map[string]any{
				"tags":        []string{"User"},
				"summary":     "Channel sıralaması",
				"description": "Kullanıcının channel'ındaki puan sıralaması.",
				"operationId": "getLeaderboard",
				"security":    bearer(),
				"responses": map[string]any{
					"200": jsonResponse("Sıralama listesi", arrayRef("UserScore")),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Channel üyeliği gerekli"),
				},
			},
		},
		"/api/v1/events": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Events"},
				"summary":     "Event listesi",
				"description": "Deadline'a göre sıralı event listesi. pending = deadline geçmiş, sonuç bekleyen.",
				"operationId": "listEvents",
				"security":    bearer(),
				"parameters": []map[string]any{
					{
						"name":        "status",
						"in":          "query",
						"required":    false,
						"schema":      map[string]any{"type": "string", "enum": []string{"open", "locked", "pending", "completed"}, "default": "open"},
						"description": "Event durumu filtresi",
					},
				},
				"responses": map[string]any{
					"200": jsonResponse("Event listesi", arrayRef("EventWithPrediction")),
					"400": errorResponse("Geçersiz filtre"),
					"401": errorResponse("Yetkisiz"),
				},
			},
		},
		"/api/v1/events/{id}": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Events"},
				"summary":     "Event detayı",
				"operationId": "getEvent",
				"security":    bearer(),
				"parameters":  pathID("Event UUID"),
				"responses": map[string]any{
					"200": jsonResponse("Event detayı", ref("EventDetailResponse")),
					"401": errorResponse("Yetkisiz"),
					"404": errorResponse("Event bulunamadı"),
				},
			},
		},
		"/api/v1/events/{id}/prediction": map[string]any{
			"put": map[string]any{
				"tags":        []string{"Events"},
				"summary":     "Tahmin oluştur/güncelle",
				"description": "Sadece açık eventlerde ve deadline öncesinde.",
				"operationId": "upsertPrediction",
				"security":    bearer(),
				"parameters":  pathID("Event UUID"),
				"requestBody": jsonBody(ref("PredictionRequest"), "Tahmin"),
				"responses": map[string]any{
					"200": jsonResponse("Tahmin kaydedildi", ref("Prediction")),
					"400": errorResponse("Geçersiz istek"),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Event kapalı"),
					"404": errorResponse("Event bulunamadı"),
				},
			},
		},
		"/api/v1/events/{id}/predictions": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Events"},
				"summary":     "Channel tahminleri",
				"description": "Deadline geçtikten sonra aynı channel'daki kullanıcıların tercihleri görünür.",
				"operationId": "listEventPredictions",
				"security":    bearer(),
				"parameters":  pathID("Event UUID"),
				"responses": map[string]any{
					"200": jsonResponse("Tahmin listesi", arrayRef("Prediction")),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Tahminler henüz görünür değil"),
					"404": errorResponse("Event bulunamadı"),
				},
			},
		},
		"/api/v1/admin/channels": map[string]any{
			"get": map[string]any{
				"tags":        []string{"Admin Channels"},
				"summary":     "Channel listesi",
				"operationId": "adminListChannels",
				"security":    bearer(),
				"responses": map[string]any{
					"200": jsonResponse("Channel listesi", arrayRef("Channel")),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Admin gerekli"),
				},
			},
			"post": map[string]any{
				"tags":        []string{"Admin Channels"},
				"summary":     "Channel oluştur",
				"operationId": "adminCreateChannel",
				"security":    bearer(),
				"requestBody": jsonBody(ref("CreateChannelRequest"), "Channel bilgileri"),
				"responses": map[string]any{
					"201": jsonResponse("Channel oluşturuldu", ref("Channel")),
					"400": errorResponse("Geçersiz istek"),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Admin gerekli"),
					"409": errorResponse("Channel kodu zaten mevcut"),
				},
			},
		},
		"/api/v1/admin/events": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Admin Events"},
				"summary":     "Event oluştur",
				"operationId": "adminCreateEvent",
				"security":    bearer(),
				"requestBody": jsonBody(ref("CreateEventRequest"), "Event bilgileri"),
				"responses": map[string]any{
					"201": jsonResponse("Event oluşturuldu", ref("Event")),
					"400": errorResponse("Geçersiz istek"),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Admin gerekli"),
				},
			},
		},
		"/api/v1/admin/events/{id}": map[string]any{
			"patch": map[string]any{
				"tags":        []string{"Admin Events"},
				"summary":     "Event güncelle",
				"description": "Sadece açık (open) eventler güncellenebilir.",
				"operationId": "adminUpdateEvent",
				"security":    bearer(),
				"parameters":  pathID("Event UUID"),
				"requestBody": jsonBody(ref("UpdateEventRequest"), "Güncellenecek alanlar"),
				"responses": map[string]any{
					"200": jsonResponse("Event güncellendi", ref("Event")),
					"400": errorResponse("Geçersiz istek"),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Admin gerekli"),
					"404": errorResponse("Event bulunamadı"),
				},
			},
		},
		"/api/v1/admin/events/{id}/result": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Admin Events"},
				"summary":     "Event sonucu gir",
				"operationId": "adminSetEventResult",
				"security":    bearer(),
				"parameters":  pathID("Event UUID"),
				"requestBody": jsonBody(ref("SetResultRequest"), "Sonuç"),
				"responses": map[string]any{
					"200": jsonResponse("Sonuç kaydedildi", ref("Event")),
					"400": errorResponse("Geçersiz istek"),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Admin gerekli"),
					"404": errorResponse("Event bulunamadı"),
				},
			},
		},
		"/api/v1/admin/events/{id}/calculate-scores": map[string]any{
			"post": map[string]any{
				"tags":        []string{"Admin Events"},
				"summary":     "Puan hesapla",
				"description": "Event sonucuna göre tahmin puanlarını hesaplar ve event'i completed yapar.",
				"operationId": "adminCalculateScores",
				"security":    bearer(),
				"parameters":  pathID("Event UUID"),
				"responses": map[string]any{
					"200": jsonResponse("Puanlar hesaplandı", ref("Event")),
					"400": errorResponse("Sonuç girilmemiş"),
					"401": errorResponse("Yetkisiz"),
					"403": errorResponse("Admin gerekli"),
					"404": errorResponse("Event bulunamadı"),
					"409": errorResponse("Puanlar zaten hesaplanmış"),
				},
			},
		},
	}
}

func buildSchemas() map[string]any {
	return map[string]any{
		"ErrorResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"error": map[string]any{"type": "string", "example": "error message"},
			},
		},
		"HealthResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"status": map[string]any{"type": "string", "example": "ok"},
			},
		},
		"RegisterRequest": map[string]any{
			"type":     "object",
			"required": []string{"name", "password", "channel_code"},
			"properties": map[string]any{
				"name":         map[string]any{"type": "string", "example": "ahmet"},
				"password":     map[string]any{"type": "string", "format": "password", "minLength": 6, "example": "secret123"},
				"channel_code": map[string]any{"type": "string", "example": "ABC123"},
			},
		},
		"LoginRequest": map[string]any{
			"type":     "object",
			"required": []string{"name", "password"},
			"properties": map[string]any{
				"name":         map[string]any{"type": "string", "example": "ahmet"},
				"password":     map[string]any{"type": "string", "format": "password", "example": "secret123"},
				"channel_code": map[string]any{"type": "string", "example": "ABC123", "description": "Admin için opsiyonel"},
			},
		},
		"AuthResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"access_token": map[string]any{"type": "string"},
				"expires_at":   map[string]any{"type": "string", "format": "date-time"},
				"user":         ref("UserProfile"),
			},
		},
		"UserProfile": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":           map[string]any{"type": "string", "format": "uuid"},
				"name":         map[string]any{"type": "string"},
				"role":         map[string]any{"type": "string", "enum": []string{"user", "admin"}},
				"channel_id":   map[string]any{"type": "string", "format": "uuid", "nullable": true},
				"total_points": map[string]any{"type": "integer", "nullable": true},
			},
		},
		"Channel": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "format": "uuid"},
				"code":       map[string]any{"type": "string", "example": "ABC123"},
				"name":       map[string]any{"type": "string", "example": "Okul Ligi"},
				"created_at": map[string]any{"type": "string", "format": "date-time"},
			},
		},
		"CreateChannelRequest": map[string]any{
			"type":     "object",
			"required": []string{"code", "name"},
			"properties": map[string]any{
				"code": map[string]any{"type": "string", "example": "ABC123"},
				"name": map[string]any{"type": "string", "example": "Okul Ligi"},
			},
		},
		"Event": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":         map[string]any{"type": "string", "format": "uuid"},
				"type":       map[string]any{"type": "string", "enum": []string{"match_score", "champion", "runner_up", "third_place"}},
				"title":      map[string]any{"type": "string", "example": "Brezilya vs Arjantin"},
				"metadata":   ref("EventMetadata"),
				"deadline":   map[string]any{"type": "string", "format": "date-time"},
				"status":     map[string]any{"type": "string", "enum": []string{"open", "locked", "completed"}},
				"result":     map[string]any{"nullable": true, "oneOf": []any{ref("MatchScoreResult"), ref("TeamResult")}},
				"created_at": map[string]any{"type": "string", "format": "date-time"},
			},
		},
		"EventMetadata": map[string]any{
			"type":        "object",
			"description": "match_score: {home_team, away_team} | champion/runner_up/third_place: {teams: [...]}",
			"example":     map[string]any{"home_team": "Brezilya", "away_team": "Arjantin"},
		},
		"CreateEventRequest": map[string]any{
			"type":     "object",
			"required": []string{"type", "title", "deadline"},
			"properties": map[string]any{
				"type":     map[string]any{"type": "string", "enum": []string{"match_score", "champion", "runner_up", "third_place"}},
				"title":    map[string]any{"type": "string", "example": "Brezilya vs Arjantin"},
				"metadata": ref("EventMetadata"),
				"deadline": map[string]any{"type": "string", "format": "date-time", "example": "2026-06-15T18:00:00Z"},
			},
		},
		"UpdateEventRequest": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title":    map[string]any{"type": "string"},
				"metadata": ref("EventMetadata"),
				"deadline": map[string]any{"type": "string", "format": "date-time"},
			},
		},
		"SetResultRequest": map[string]any{
			"type":     "object",
			"required": []string{"result"},
			"properties": map[string]any{
				"result": map[string]any{
					"oneOf": []any{ref("MatchScoreResult"), ref("TeamResult")},
				},
			},
		},
		"MatchScoreResult": map[string]any{
			"type":     "object",
			"required": []string{"home_score", "away_score"},
			"properties": map[string]any{
				"home_score": map[string]any{"type": "integer", "minimum": 0, "example": 2},
				"away_score": map[string]any{"type": "integer", "minimum": 0, "example": 1},
			},
		},
		"MatchScoreChoice": map[string]any{
			"type":     "object",
			"required": []string{"home_score", "away_score"},
			"properties": map[string]any{
				"home_score": map[string]any{"type": "integer", "minimum": 0, "example": 2},
				"away_score": map[string]any{"type": "integer", "minimum": 0, "example": 1},
			},
		},
		"TeamResult": map[string]any{
			"type":        "object",
			"description": "champion, runner_up ve third_place eventleri için tek takım seçimi",
			"required":    []string{"team"},
			"properties": map[string]any{
				"team": map[string]any{"type": "string", "example": "Brezilya"},
			},
		},
		"TeamChoice": map[string]any{
			"type":        "object",
			"description": "champion, runner_up ve third_place eventleri için tek takım tahmini",
			"required":    []string{"team"},
			"properties": map[string]any{
				"team": map[string]any{"type": "string", "example": "Brezilya"},
			},
		},
		"PredictionRequest": map[string]any{
			"type":     "object",
			"required": []string{"choice"},
			"properties": map[string]any{
				"choice": map[string]any{
					"oneOf": []any{ref("MatchScoreChoice"), ref("TeamChoice")},
				},
			},
		},
		"Prediction": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":             map[string]any{"type": "string", "format": "uuid"},
				"event_id":       map[string]any{"type": "string", "format": "uuid"},
				"user_id":        map[string]any{"type": "string", "format": "uuid"},
				"user_name":      map[string]any{"type": "string"},
				"choice":         map[string]any{"oneOf": []any{ref("MatchScoreChoice"), ref("TeamChoice")}},
				"points_awarded": map[string]any{"type": "integer", "example": 3},
				"created_at":     map[string]any{"type": "string", "format": "date-time"},
				"updated_at":     map[string]any{"type": "string", "format": "date-time"},
			},
		},
		"EventWithPrediction": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"event":         ref("Event"),
				"my_prediction": map[string]any{"allOf": []any{ref("Prediction")}, "nullable": true},
			},
		},
		"EventDetailResponse": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"event":         ref("Event"),
				"my_prediction": map[string]any{"allOf": []any{ref("Prediction")}, "nullable": true},
			},
		},
		"UserScore": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"user_id":      map[string]any{"type": "string", "format": "uuid"},
				"user_name":    map[string]any{"type": "string"},
				"channel_id":   map[string]any{"type": "string", "format": "uuid"},
				"total_points": map[string]any{"type": "integer", "example": 15},
				"updated_at":   map[string]any{"type": "string", "format": "date-time"},
			},
		},
	}
}

func ref(name string) map[string]string {
	return map[string]string{"$ref": "#/components/schemas/" + name}
}

func arrayRef(name string) map[string]any {
	return map[string]any{
		"type":  "array",
		"items": ref(name),
	}
}

func bearer() []map[string][]string {
	return []map[string][]string{{"bearerAuth": {}}}
}

func pathID(desc string) []map[string]any {
	return []map[string]any{
		{
			"name":        "id",
			"in":          "path",
			"required":    true,
			"description": desc,
			"schema":      map[string]any{"type": "string", "format": "uuid"},
		},
	}
}

func jsonBody(schema map[string]string, desc string) map[string]any {
	return map[string]any{
		"required": true,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": schema,
			},
		},
		"description": desc,
	}
}

func jsonResponse(desc string, schema any) map[string]any {
	return map[string]any{
		"description": desc,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": schema,
			},
		},
	}
}

func errorResponse(desc string) map[string]any {
	return jsonResponse(desc, ref("ErrorResponse"))
}
