# ADR-010: ~~OCR как fallback, Overwolf как upgrade~~ — SUPERSEDED

**Статус:** Superseded (2026-07-17)  
**Причина:** OCR и Overwolf удалены из архитектуры. Единственный источник live-данных — HudApi через `patches/copilot-patcher.js`.  
**История:** Ранее предполагался многоуровневый подход к получению данных (Overwolf > OCR > CV). После анализа стало ясно, что HudApi даёт все необходимые данные без OCR/CV.
