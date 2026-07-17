// LT2 Copilot Patcher v2
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
        phase: 'unknown',
        // Hand (available buildings/units)
        hand: [],
        // Units currently on the field (tracked by unit ID)
        fieldUnits: {},
        // Queue of units being built
        buildQueue: [],
        // Purchasing queue
        purchaseQueue: []
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

    // Parse unit info from action header HTML
    function parseAction(action) {
        if (!action) return null;
        var name = '';
        var m = action.header && action.header.match(/>(.*?)</);
        if (m) name = m[1].trim();
        return {
            actionId: action.actionId,
            name: name,
            icon: action.image || '',
            costGold: 0,
            costMythium: 0,
            costSupply: 0,
            stacks: action.stacks || 1,
            role: action.role || ''
        };
    }

    // Register Coherent GT event listeners
    function init() {
        if (typeof engine === 'undefined') {
            setTimeout(init, 1000);
            return;
        }

        // ===== Economy =====
        engine.on('refreshMythium', function (value) {
            gameState.mythium = value;
            sendState();
        });

        engine.on('refreshGold', function (value) {
            gameState.gold = value;
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

        // ===== Wave =====
        engine.on('refreshWaveNumber', function (waveNumber) {
            gameState.wave = waveNumber;
            sendState();
        });

        engine.on('refreshWaveTime', function (value) {
            gameState.waveTimer = value;
        });

        // ===== King HP =====
        engine.on('refreshLeftKingMaxHp', function (hp) {
            gameState.kingHp = hp;
        });

        engine.on('refreshRightKingMaxHp', function (hp) {
            gameState.enemyKingHp = hp;
        });

        // ===== Enemies & Phase =====
        engine.on('refreshWestEnemiesRemaining', function (value) {
            gameState.enemiesRemainingWest = value;
            gameState.phase = value > 0 ? 'fighting' : (gameState.enemiesRemainingEast > 0 ? 'fighting' : 'building');
            if (gameState.phase === 'fighting') sendState();
        });

        engine.on('refreshEastEnemiesRemaining', function (value) {
            gameState.enemiesRemainingEast = value;
            gameState.phase = value > 0 ? 'fighting' : (gameState.enemiesRemainingWest > 0 ? 'fighting' : 'building');
        });

        // ===== Time =====
        engine.on('refreshTimeElapsed', function (seconds) {
            gameState.timeElapsed = seconds;
        });

        // ===== Hand: Available buildings (Dashboard) =====
        engine.on('refreshDashboardActions', function (actions) {
            if (!actions || !Array.isArray(actions)) return;
            gameState.hand = actions.map(parseAction).filter(Boolean);
            sendState();
        });

        // ===== Hand: Action stock (how many available) =====
        engine.on('refreshActionStock', function (actionId, value) {
            for (var i = 0; i < gameState.hand.length; i++) {
                if (gameState.hand[i].actionId === actionId) {
                    gameState.hand[i].stacks = value;
                    break;
                }
            }
            sendState();
        });

        // ===== Build queue =====
        engine.on('refreshActionQueue', function (actionId, value) {
            // value = number in queue
            var found = false;
            for (var i = 0; i < gameState.buildQueue.length; i++) {
                if (gameState.buildQueue[i].actionId === actionId) {
                    if (value <= 0) {
                        gameState.buildQueue.splice(i, 1);
                    } else {
                        gameState.buildQueue[i].count = value;
                    }
                    found = true;
                    break;
                }
            }
            if (!found && value > 0) {
                gameState.buildQueue.push({ actionId: actionId, count: value });
            }
            sendState();
        });

        // ===== Purchasing queue =====
        engine.on('refreshActionPurchasingQueue', function (actionId, value) {
            var found = false;
            for (var i = 0; i < gameState.purchaseQueue.length; i++) {
                if (gameState.purchaseQueue[i].actionId === actionId) {
                    if (value <= 0) {
                        gameState.purchaseQueue.splice(i, 1);
                    } else {
                        gameState.purchaseQueue[i].count = value;
                    }
                    found = true;
                    break;
                }
            }
            if (!found && value > 0) {
                gameState.purchaseQueue.push({ actionId: actionId, count: value });
            }
        });

        // ===== Field units: Track HP updates =====
        // Unit IDs 1+ correspond to units on the field
        engine.on('refreshUnitHp', function (unitId, amount) {
            if (unitId < 1) return;
            if (amount <= 0) {
                // Unit died or removed
                if (gameState.fieldUnits[unitId]) {
                    delete gameState.fieldUnits[unitId];
                    sendState();
                }
            } else {
                if (!gameState.fieldUnits[unitId]) {
                    gameState.fieldUnits[unitId] = { unitId: unitId, hp: amount, firstSeen: Date.now() };
                } else {
                    gameState.fieldUnits[unitId].hp = amount;
                }
            }
        });

        // ===== Target properties (selected unit details) =====
        engine.on('setTargetProperties', function (properties) {
            if (!properties) return;
            // Try to match selected unit to fieldUnits by refreshing info
            if (properties.fighterName && properties.currentHp != null) {
                // A unit was clicked — store its details
                gameState.selectedUnit = {
                    name: properties.fighterName,
                    hp: properties.currentHp,
                    maxHp: properties.maxHp,
                    title: properties.fighterTitle || ''
                };
                sendState();
            }
        });

        // ===== Mercenaries (Windshield) =====
        engine.on('refreshWindshieldActions', function (actions) {
            if (!actions || !Array.isArray(actions)) return;
            gameState.mercenaries = actions.map(parseAction).filter(Boolean);
            sendState();
        });

        // ===== Town actions (Glovebox) =====
        engine.on('refreshGloveboxActions', function (actions) {
            if (!actions || !Array.isArray(actions)) return;
            gameState.townActions = actions.map(parseAction).filter(Boolean);
            sendState();
        });

        // ===== Player scores (all players' resources) =====
        engine.on('refreshPlayerScores', function (playerScores) {
            if (!playerScores || !Array.isArray(playerScores)) return;
            gameState.playerScores = playerScores;
        });

        // ===== Reset on game end =====
        engine.on('refreshIsInGame', function (inGame) {
            if (!inGame) {
                gameState = {
                    mythium: 0, gold: 0, supply: 0, supplyCap: 0,
                    income: 0, mythiumGatherRate: 0, estimatedMythium: 0,
                    wave: 0, waveTimer: 0, kingHp: 0, enemyKingHp: 0,
                    enemiesRemainingWest: 0, enemiesRemainingEast: 0,
                    timeElapsed: 0, phase: 'unknown',
                    hand: [],
                    fieldUnits: {},
                    buildQueue: [],
                    purchaseQueue: [],
                    mercenaries: [],
                    townActions: [],
                    playerScores: [],
                    selectedUnit: null
                };
                sendState();
            }
        });

        connect();
    }

    init();
})();
