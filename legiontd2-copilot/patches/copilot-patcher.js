// LT2 Copilot Patcher
// Intercepts game events via Coherent GT engine.on() and forwards them via WebSocket.
// Place this file in: Legion TD 2_Data/uiresources/AeonGT/hud/js/
// Add to gateway.html: <script src="hud/js/copilot-patcher.js"></script>

(function () {
    'use strict';

    var COPILOT_WS_URL = 'ws://localhost:8080/ws';

    var ws = null;
    var reconnectTimer = null;
    var connected = false;

    // Cached game state
    var gameState = {
        mythium: 0,
        gold: 0,
        supply: 0,
        supplyCap: 0,
        income: 0,
        mythiumGatherRate: 0,
        estimatedMythium: 0,
        wave: 0,
        waveTimer: 0,
        kingHp: 0,
        enemyKingHp: 0,
        enemiesRemainingWest: 0,
        enemiesRemainingEast: 0,
        timeElapsed: 0,
        phase: 'unknown'
    };

    function connect() {
        if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) {
            return;
        }
        try {
            ws = new WebSocket(COPILOT_WS_URL);
            ws.onopen = function () {
                connected = true;
                sendState();
            };
            ws.onclose = function () {
                connected = false;
                scheduleReconnect();
            };
            ws.onerror = function () {
                connected = false;
            };
        } catch (e) {
            scheduleReconnect();
        }
    }

    function scheduleReconnect() {
        if (reconnectTimer) return;
        reconnectTimer = setTimeout(function () {
            reconnectTimer = null;
            connect();
        }, 3000);
    }

    function sendState() {
        if (!ws || ws.readyState !== WebSocket.OPEN) return;
        try {
            ws.send(JSON.stringify(gameState));
        } catch (e) { /* ignore */ }
    }

    // Register Coherent GT event listeners
    // These run alongside the game's existing handlers (multiple handlers per event are supported)
    function init() {
        if (typeof engine === 'undefined') {
            // Not running inside Coherent GT; retry later
            setTimeout(init, 1000);
            return;
        }

        // Economy
        engine.on('refreshMythium', function (value) {
            gameState.mythium = value;
            sendState();
        });

        engine.on('refreshGold', function (value) {
            gameState.gold = value;
            sendState();
        });

        engine.on('refreshSupply', function (value) {
            gameState.supply = value;
        });

        engine.on('refreshSupplyCap', function (value) {
            gameState.supplyCap = value;
        });

        engine.on('refreshGoldRemaining', function (goldNextWave, goldRemaining, income) {
            gameState.income = income || 0;
        });

        engine.on('refreshMythiumGatherRate', function (player, rate) {
            gameState.mythiumGatherRate = rate || 0;
        });

        engine.on('refreshEstimatedMythium', function (value) {
            gameState.estimatedMythium = value;
        });

        // Wave
        engine.on('refreshWaveNumber', function (waveNumber) {
            gameState.wave = waveNumber;
            sendState();
        });

        engine.on('refreshWaveTime', function (value) {
            gameState.waveTimer = value;
        });

        // King HP
        engine.on('refreshLeftKingMaxHp', function (hp) {
            gameState.kingHp = hp;
            sendState();
        });

        engine.on('refreshRightKingMaxHp', function (hp) {
            gameState.enemyKingHp = hp;
        });

        // Enemies
        engine.on('refreshWestEnemiesRemaining', function (value) {
            gameState.enemiesRemainingWest = value;
            gameState.phase = value > 0 ? 'fighting' : 'building';
            sendState();
        });

        engine.on('refreshEastEnemiesRemaining', function (value) {
            gameState.enemiesRemainingEast = value;
            gameState.phase = value > 0 ? 'fighting' : 'building';
        });

        // Time
        engine.on('refreshTimeElapsed', function (seconds) {
            gameState.timeElapsed = seconds;
        });

        // Game state
        engine.on('refreshIsInGame', function (inGame) {
            if (!inGame) {
                // Reset state on game exit
                gameState = {
                    mythium: 0, gold: 0, supply: 0, supplyCap: 0,
                    income: 0, mythiumGatherRate: 0, estimatedMythium: 0,
                    wave: 0, waveTimer: 0, kingHp: 0, enemyKingHp: 0,
                    enemiesRemainingWest: 0, enemiesRemainingEast: 0,
                    timeElapsed: 0, phase: 'unknown'
                };
            }
        });

        // Connect WebSocket
        connect();
    }

    init();
})();
