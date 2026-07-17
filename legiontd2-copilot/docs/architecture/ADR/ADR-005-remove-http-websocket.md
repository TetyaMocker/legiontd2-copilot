# ADR-005: Удаление HTTP/WebSocket из архитектуры

**Статус:** Active  
**Дата:** 2026-07-17  
**Контекст:** Проект — Desktop-приложение, но текущая реализация содержит полноценный HTTP-сервер с REST API, WebSocket-хаб и браузерный UI (SPA). Это архитектурное рассогласование с ТЗ (раздел 1: "desktop-приложение под Windows").  
**Решение:** Удалить встроенный HTTP-сервер, REST-эндпоинты, WebSocket hub и браузерный UI. Заменить их на:
  - gRPC для Go↔Python IPC
  - Wails native window для UI
  - Прямые вызовы в Go для advisor↔state  
**Последствия:**  
  - Positive: чистая desktop-архитектура; устранение поверхности атаки (HTTP); устранение WebSocket complexity.  
  - Negative: потеря возможности открыть UI в любом браузере на локальной сети. (Не является требованием для личного Desktop-приложения.)  
**Миграция:**  
  - Phase 1: удалить `internal/http/server.go` (endpoints и embed static)  
  - Phase 2: удалить `internal/ws/hub.go` (WebSocket)  
  - Phase 3: перевести go-core на gRPC client  
**Связанные ADR:** ADR-004 (gRPC), ADR-006 (Wails)
