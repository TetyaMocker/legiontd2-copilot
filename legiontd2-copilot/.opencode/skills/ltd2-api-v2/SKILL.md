---
name: ltd2-api-v2
description: Справочная информация по официальному Legion TD 2 API v2 — эндпоинты, лимиты, формат данных, известные особенности парсинга. Загружай при работе с internal/api, internal/dataset, cmd/dataset, или любым кодом, обращающимся к apiv2.legiontd2.com.
license: internal
compatibility: opencode
metadata:
  audience: go-core, ml-dataset
---

## Базовые факты

- Base URL: `https://apiv2.legiontd2.com`
- Авторизация: заголовок `x-api-key`, ключ выпускается на `developer.legiontd2.com`.
- Rate limit: 5 req/s, 100 burst, 1000 запросов/день. Не реализовывать retry/backoff
  внутри низкоуровневого клиента — превышение лимита должно быть видно вызывающему
  коду (батч-джобу сбора датасета), а не скрыто автоматическими повторами.
- API отдаёт матчи **не старше одного года**. Более старые — отдельный архив
  `games-archive/2018-2021/games.bson.zip.00N` (6 частей), восстанавливается через
  `mongorestore`, не через REST.
- **API v2 НЕ даёт live-состояние текущего матча.** Только историю завершённых игр и
  версионированные справочники. Любой код, ожидающий получить от этого API "что
  происходит в матче прямо сейчас" — ошибка проектирования.

## Ключевая находка: buildPerWave

При запросе с `includeDetails=true` (`/games`, `/games/byId/{id}`,
`/players/matchHistory/{id}`) в `playersData` приходят массивы **по волнам**:

- `buildPerWave: [][]string` — расстановка юнитов формата `"unit_id:x|y"` на каждой волне;
- `leaksPerWave` — какие юниты противника прошли (сигнал "волна не удержана");
- `netWorthPerWave`, `workersPerWave`, `incomePerWave` — экономика по волнам;
- `mercenariesSentPerWave` / `mercenariesReceivedPerWave`;
- `kingUpgradesPerWave` / `opponentKingUpgradesPerWave`.

Это единственный способ получить размеченные примеры "расстановка → результат волны"
**без всякого CV**, на реальных матчах. Приоритет использования этих данных для
обучения выше, чем ручная разметка скриншотов — не начинай с CV-детекции юнитов,
пока не исчерпана эта возможность.

## Особенности формата, ломающие наивный парсинг

- Многие числовые поля в ответах API — **JSON-строки**, а не числа
  (`"mythiumCost": "7"`, не `"mythiumCost": 7`). Наивный `json.Unmarshal` в
  `int`/`float64` упадёт. См. `numString` в `internal/api/client.go` как пример
  правильного подхода (кастомный `UnmarshalJSON`).
- Строки в `mercenariesSentPerWave` иногда содержат экранированную лишнюю кавычку
  (`"Snail\""` вместо `"Snail"`) — это особенность самого API, не баг парсера.
  Не пытаться "исправить" данные при парсинге — сохранять как есть, обрабатывать
  на этапе очистки датасета, задокументировав это явно.
- `buildPerWave` entry формат: `unit_id:x|y`, где x/y — float с точкой (не запятой).
  Парсить через `SplitN` по `:`, затем по `|`, см. `ParseBuildEntry` как reference.

## Основные эндпоинты (по итогам разбора Swagger)

| Метод | Путь | Назначение |
|---|---|---|
| GET | `/players/byId/{id}`, `/players/byName/{name}` | Профиль игрока |
| GET | `/players/matchHistory/{id}` | История матчей игрока |
| GET | `/players/stats/{id}`, `/players/stats` | Статистика (винрейт и т.п.) |
| GET | `/units/byVersion/{version}` | Справочник юнитов на патч (`limit=300` = все) |
| GET | `/games` | Фильтр по матчам (version, queueType, дате, elo) |
| GET | `/games/byId/{id}` | Один матч |
| GET | `/info/legions`, `/info/waves`, `/info/spells`, `/info/abilities`, `/info/research` | Статичные справочники |

Полную схему см. в `internal/api/client.go` — она уже написана по реальной Swagger-выгрузке,
не по догадкам. Расхождения между кодом и этим списком означают, что либо API поменялся,
либо клиент устарел — не доверяй слепо ни тому, ни другому, сверяйся с `/openapi.yaml`.
