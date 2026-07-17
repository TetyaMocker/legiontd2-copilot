# Feature Matrix — полное описание логики

## 1. Общая архитектура обучения (как я это вижу)

```
       Game State (JSON)                     Feature Matrix (JSON)                ML Model
  ┌──────────────────────┐    build()    ┌──────────────────────┐    predict()   ┌──────────┐
  │ economy              │ ───────────→  │ economy              │ ───────────→  │ Action:  │
  │ hand[]               │              │ board                │               │   type   │
  │ fieldUnits{}         │              │ available             │               │   msg    │
  │ mercenaries[]        │              │ waves                 │               │   score  │
  │ waves (static)       │              │ coverage              │               └──────────┘
  └──────────────────────┘              │ context               │
                                        └──────────────────────┘
                                                │
                                         (train on logs)
                                                │
                                        ┌────────────────┐
                                        │  SQLite logs    │
                                        │  + ручная       │
                                        │  разметка       │
                                        └────────────────┘
```

### Этапы обучения:

1. **Сбор датасета** — каждый тик (сек) логируем `FeatureMatrix` + какое решение принял игрок/эвристика
2. **Разметка (labeling)** — человек проставляет `ground truth`: "правильная рекомендация = build_units", "надо было строить workers"
3. **Обучение** — нейросеть (простая FFN / градиентный бустинг) учится предсказывать `{ action, priority, message }` по `FeatureMatrix`
4. **Инференс** — в рантайме матрица подаётся в модель → на выходе top-K рекомендаций

---

## 2. Полная схема Feature Matrix

### 2.1 economy — экономика

| Поле        | Тип     | Откуда                        | Пример       |
|-------------|---------|-------------------------------|--------------|
| gold        | int     | `refreshGold`                 | 253          |
| mythium     | int     | `refreshMythium`              | 20           |
| supply      | int     | `refreshSupply`               | 0            |
| supplyCap   | int     | `refreshSupplyCap`            | 256          |
| income      | int     | `refreshGoldRemaining`        | 3            |
| wave        | int     | `refreshWaveNumber`           | 1            |
| phase       | string  | вычисляется из enemiesRemaining | "building" / "fighting" |
| timer       | int     | `refreshWaveTime`             | 71           |
| kingHp      | float   | `refreshLeftKingMaxHp`        | 3000         |
| enemyKingHp | float   | `refreshRightKingMaxHp`       | 3000         |

**Правила:**
- `phase = "fighting"` если `enemiesRemainingWest > 0 || enemiesRemainingEast > 0`, иначе `"building"`
- `timer` — сколько секунд до конца волны (0 если бой уже идёт)

---

### 2.2 board — что уже на столе

| Поле        | Тип          | Описание                              |
|-------------|--------------|---------------------------------------|
| totalUnits  | int          | количество юнитов на поле              |
| totalHp     | float64      | сумма HP всех юнитов на поле           |
| byDamage    | DamageCount  | разбивка по типам урона               |
| byRole      | RoleCount    | разбивка по ролям                     |
| units       | BoardUnit[]  | детальный список юнитов на столе       |

**DamageCount:**
```
{ normal: 2, pierce: 3, magic: 0 }
```

**RoleCount:**
```
{ tank: 1, dps: 3, balanced: 1 }
```

**BoardUnit:**
```json
{
  "name": "Peewee",
  "role": "role_tank",
  "damageType": "Pierce",
  "hp": 3600
}
```

**Правила:**
- Юниты на поле берутся из `state.fieldUnits{}`
- Тип урона юнита — из `FighterAttack` / `MercAttack` (статичные map в `internal/unitdata`)
- Если у юнита нет записи в map — тип урона = `"unknown"`
- Роль юнита — из `action.role` (присылает игра: `role_tank`, `role_dps`, `role_balanced`)
- Один юнит на поле = +1 к его damageType и role

---

### 2.3 available — что доступно для постройки

**fighters[]** и **mercs[]** — массивы:

```json
{
  "name": "Peewee",
  "goldCost": 25,
  "mythCost": 0,
  "supplyCost": 1,
  "role": "role_tank",
  "damageType": "Pierce",
  "stacks": 1,
  "affordable": true
}
```

| Поле        | Тип     | Описание                                  |
|-------------|---------|--------------------------------------------|
| goldCost    | int     | из `action.subheader` (Gold.png), или static map, или 0 |
| mythCost    | int     | из `action.subheader` (Mythium.png), или static map, или 0 |
| supplyCost  | int     | из `action.subheader` (Supply.png), или static map, или 0 |
| stacks      | int     | сколько раз можно купить (из `refreshActionStock`) |
| affordable  | bool    | хватает ли ресурсов для покупки            |

**Правила affordable:**
- **Файтер:** `stacks > 0 && gold >= goldCost && goldCost > 0 && (supplyCap - supply) >= supplyCost`
- **Мерк:** `stacks > 0 && mythium >= mythCost && mythCost > 0 && (supplyCap - supply) >= supplyCost`
- Если `goldCost == 0` (неизвестно) — affordable = false (не можем определить)
- Если `mythCost == 0` — affordable = false

**Источник стоимости (patcher):**
1. Парсим `action.subheader` — ищем `<img ... Gold.png' /> ЧИСЛО`
2. Если не нашли — проверяем `action.goldCost`, `action.costGold`, `action.cost.gold`
3. Если и там 0 — fallback в статический map FighterCosts/MercCosts (сейчас пустой)

---

### 2.4 waves — анализ волн

```json
{
  "current": 1,
  "upcoming": [
    {
      "number": 2,
      "name": "Wales",
      "armorType": "heavy",
      "attackType": "normal",
      "amount": 12,
      "hasBoss": false,
      "bestDamage": "Magic",
      "multiplier": 1.25
    },
    ...
  ]
}
```

| Поле       | Тип     | Описание                                    |
|------------|---------|----------------------------------------------|
| armorType  | string  | `light`, `medium`, `heavy`, `fortified`      |
| attackType | string  | `normal`, `pierce`, `magic`, `chaos`         |
| amount     | int     | количество врагов (волна + босс)             |
| hasBoss    | bool    | есть ли босс                                 |
| bestDamage | string  | самый эффективный тип урона против этой брони |

**Правила bestDamage:**
```go
DamageMultiplier(atk, armor):
           light  medium  heavy  fortified
Normal     0.80   0.90    1.15   1.15
Pierce     1.20   0.85    1.15   0.80
Magic      1.00   1.25    0.75   1.05
```
Выбирается атака с максимальным множителем. Например:
- light → Pierce (1.20)
- medium → Magic (1.25)
- heavy → Normal (1.15)
- fortified → Normal (1.15)

**Источник данных:** статический map `WavesByNumber` в `internal/wavedata/waves.go` (волны 1-21)

---

### 2.5 coverage — анализ покрытия урона

```json
{
  "boardDamage":     { "normal": 2, "pierce": 3, "magic": 0 },
  "availableDamage": { "normal": 1, "pierce": 2, "magic": 1 },
  "missingTypes":    ["Magic"],
  "recommended":     "Magic",
  "explanation":     "Броня врагов чаще всего уязвима к Magic урону. Построй больше Magic-юнитов"
}
```

**Алгоритм coverage:**

1. **Считаем boardDamage** — сколько юнитов каждого типа урона на поле + в руке (поставленные)
2. **Считаем availableDamage** — сколько юнитов каждого типа доступно для покупки (в руке)
3. **Для каждой волны (upcoming):** смотрим `bestDamage`, считаем сколько волн требуют каждого типа
4. **Если для какого-то типа `bestDamage` нужно в ≥2 волнах, а у нас <2 юнитов этого типа (board + available):**
   - Добавляем в missingTypes
   - Рекомендуем этот тип
5. **Если missingTypes пустой:** `recommended = "balanced"`

**Правила недостаточности:**
- `missingTypes = []` — все типы урона покрыты
- `missingTypes = ["Magic"]` — не хватает магического урона под грядущие волны
- `missingTypes = ["Pierce", "Magic"]` — не хватает обоих

---

### 2.6 context — контекстные флаги

```json
{
  "isFighting": false,
  "isKingLow": false,
  "hasAffordable": true,
  "canBuildWorkers": true
}
```

| Флаг            | Условие                                        |
|-----------------|-------------------------------------------------|
| isFighting      | `phase == "fighting"`                           |
| isKingLow       | `kingHp < 30` (% от макс?)                     |
| hasAffordable   | хотя бы один affordable fighter или merc        |
| canBuildWorkers | `gold >= 50 && phase == "building"`            |

**Note:** сейчас `kingHp` приходит как абсолютное значение (3000). Флаг `isKingLow` считает `kingHp < 30.0` — это баг, нужно сравнивать с maxHp в процентах. @todo

---

## 3. Текущие эвристические правила (heuristic.go)

Правила применяются последовательно, каждое добавляет рекомендацию с приоритетом 0-5.

| # | Условие                                                        | Action          | Приоритет | Сообщение (пример)                                   |
|---|----------------------------------------------------------------|-----------------|-----------|------------------------------------------------------|
| 1 | `phase == building && wave <= 21 && есть counter-юниты`        | `counter_wave`  | 4         | "Волна 2 (Wales) — броня heavy, эффективен Magic. Есть в руке: ..." |
| 2 | `phase == building && wave <= 21 && нет counter-юнитов`        | `wave_info`     | 2         | "Волна 2 (Wales) — броня heavy, эффективен Magic. Нет контр-юнитов в руке" |
| 3 | `wave.armor == light && best == pierce`                        | `tip`           | 1         | "Light броня — Pierce урон на 20% эффективнее"       |
| 4 | `wave.armor == fortified && best == normal`                    | `tip`           | 1         | "Fortified броня — Normal урон наиболее эффективен"   |
| 5 | `timer > 0 && timer <= 30 && есть affordable`                  | `build_units`   | 5         | "До волны 11s — поставь: Peewee, Ranger"             |
| 6 | `timer > 10 && canWork && hand пуст && supply < 5`             | `build_workers` | 4         | "Нет юнитов для постановки — создай рабочих"          |
| 7 | `timer > 10 && canWork && gold >= 100 && wave <= 5`           | `build_workers` | 3         | "Ранняя игра — создай рабочих (есть 253 золота)"      |
| 8 | `timer > 0 && timer <= 10 && mythium > 60`                     | `send_mercs`    | 5         | "Волна вот-вот начнётся — отправь наёмников"         |
| 9 | `mythium > 120 && phase == building`                           | `save_mythium`  | 3         | "Копи мифиум (120) — отправь наёмников на волне противника" |
| 10| `kingHp < 30 && mythium > 40 && phase == building`            | `upgrade_king`  | 5         | "HP короля 20% — улучши короля для защиты"           |
| 11| `phase == building && hand > 0 && affordable == 0 && gold < 50 && mythium < 60` | `save_gold` | 2 | "Мало ресурсов — жди дохода" |
| 12| Если нет других рекоммендаций: `phase == building`             | `hold`          | 1         | "Стройся — ситуация стабильная"                      |
| 13| Если нет других рекоммендаций: `phase == fighting`            | `fight`         | 1         | "Идёт бой — наблюдает"                               |
| 14| Всегда                                                         | `save_mythium`  | 0         | "Mythium: 20 | Доход: 3 | Золото: 253 | Снабжение: 0/256" (инфо-строка) |

---

## 4. Потенциальные улучшения для ML

### 4.1 Какие признаки добавить

| Признак | Почему важен |
|---------|-------------|
| `board.byRole` — баланс ролей | Нет танков → доска падает, нет дамагеров → не убивают волну |
| `board.averageHp` — среднее HP | Показывает, насколько доска "живучая" |
| `coverage.missingTypes` | Прямой сигнал — докупать Magic/Pierce/Normal |
| `waves[].attackType` — чем атакует враг | Нужно знать, переживёт ли доска урон врага |
| `income` — темп экономики | Высокий доход → можно агрессивнее тратить мифиум |
| `kingHpPercent` | Абсолютное значение HP бесполезно, нужно в % от максимума |
| `enemyIncome` (если доступно) | Сравнение темпа с противником |

### 4.2 Формат размеченных данных для обучения

```
features: FeatureMatrix
label: {
  action: "build_units" | "build_workers" | "send_mercs" | "upgrade_king" 
        | "save_mythium" | "save_gold" | "counter_wave" | "hold",
  priority: 0-5,
  targetUnit: "Peewee" | "" (опционально, какой юнит строить)
  comment: "надо было ставить Magic, не хватает"
}
```

### 4.3 Предлагаемый pipeline

1. **Логирование:** каждый тик (1 сек) пишем в SQLite `(match_id, tick, FeatureMatrix, heuristic_action, player_action)`
2. **Разметка:** человек открывает лог матча, смотрит каждый тик, исправляет recommendation если эвристика ошиблась
3. **Обучение:** X = FeatureMatrix, y = { action, priority }
4. **Валидация:** на отложенной выборке (20% матчей) считаем accuracy поверх эвристики
5. **Замена:** когда ML accuracy > heuristic accuracy на 5% — переключаемся

---

## 5. Недостатки текущей реализации (что нужно доработать)

1. **kingHp** — приходит как абсолютное значение (3000). Нужно получать maxHp для расчёта процента
2. **supplyCost** — не найдено в `action.subheader`, возможно лежит в другом свойстве. Нужен дамп merc action
3. **unitAttackType** — много юнитов без записи в `FighterAttack`/`MercAttack` (Masked Spirit, Sacred Steed, Elite Archer, Mr Brewpot и др.)
4. **affordable logic** — если `goldCost == 0` (неизвестен), affordable = false. Юнит не рекомендуется даже если дешёвый
5. **Нет данных о вражеском столе** — не знаем, что строит противник
6. **coverage** не учитывает уже построенных юнитов (fieldUnits) — только hand + mercs в руке

---

## 6. Вопросы, которые нужно решить

1. Как получать `maxKingHp`? (нужно найти событие в engine.on)
2. Как получать supply юнитов? (subheader не содержит Supply.png)
3. Как помечать unknown юниты — `damageType: "unknown"` или исключать из подсчёта?
4. Как часто переобучать модель? (после каждых N матчей?)
5. Какой формат разметки удобен человеку? (CSV? JSON? GUI?)
