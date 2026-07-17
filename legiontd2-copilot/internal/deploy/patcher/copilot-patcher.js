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

    // Dump all enumerable properties of a raw action object (first time only, for debugging)
    var _rawDumped = false;
    function dumpRawAction(action, label) {
        if (_rawDumped || !action) return;
        _rawDumped = true;
        var result = {};
        try {
            for (var key in action) {
                try {
                    var val = action[key];
                    if (typeof val !== 'function' && typeof val !== 'undefined' && val !== null) {
                        if (typeof val === 'object' && !Array.isArray(val) && val.toString === Object.prototype.toString) {
                            var sub = {};
                            for (var k2 in val) {
                                try { sub[k2] = val[k2]; } catch(e) {}
                            }
                            result[key] = sub;
                        } else if (Array.isArray(val)) {
                            result[key] = '[array length=' + val.length + ']';
                        } else {
                            result[key] = val;
                        }
                    }
                } catch(e) {}
            }
        } catch(e) {}
        gameState._actionSample = result;
    }

    // Parse unit info from action header HTML
    function parseAction(action) {
        if (!action) return null;

        // Extract name — try every possible property the game might expose
        var name = action.name || action.title || action.label || action.tooltip
                 || action.description || action.unitName || action.fighterName
                 || (action.display && (action.display.name || action.display.title))
                 || '';
        if (!name && action.header) {
            name = action.header;
            if (action.header.indexOf('<') >= 0) {
                var m = action.header.match(/>([^<]+)</);
                if (m) name = m[1].trim();
                else name = action.header.replace(/<[^>]*>/g, '').trim();
            } else {
                name = action.header.trim();
            }
        }

        // Fallback: extract real unit name from icon filename
        // Game sends full HTML tooltip as name — detect by HTML tags or button label brackets
        if (!name || name.indexOf('<') >= 0 || name.indexOf('[') === 0) {
            var iconPath = action.image || action.icon || '';
            if (iconPath) {
                var iconName = iconPath.replace(/.*[\/\\]/, '').replace(/\.[^.]+$/, '');
                if (iconName) {
                    name = iconName.replace(/([a-z])([A-Z])/g, '$1 $2');
                }
            }
        }

        // Extract costs from the action object (Coherent GT exposes these as direct properties)
        function getCost(src, fallback) {
            if (typeof src === 'number') return src;
            if (src != null && typeof src.gold === 'number') return src.gold;
            return fallback;
        }
        var costGold   = getCost(action.goldCost,   getCost(action.costGold,   0));
        var costMythium= getCost(action.mythiumCost,getCost(action.costMythium,0));
        var costSupply = getCost(action.supplyCost, getCost(action.costSupply, 0));

        // Also check a unified cost object
        if (action.cost) {
            if (!costGold)    costGold    = action.cost.gold    || 0;
            if (!costMythium) costMythium = action.cost.mythium || 0;
            if (!costSupply)  costSupply  = action.cost.supply  || 0;
        }

        // Extract costs from subheader HTML, e.g.:
        // "<img class='tooltip-icon' src='hud/img/icons/Gold.png' /> 25"
        // "<img class='tooltip-icon' src='hud/img/icons/Mythium.png' /> 5"
        // "<img class='tooltip-icon' src='hud/img/icons/Supply.png' /> 2"
        if (action.subheader && typeof action.subheader === 'string') {
            if (!costGold) {
                var g = action.subheader.match(/Gold\.png[^0-9]*(\d+)/i);
                if (g) costGold = parseInt(g[1], 10);
            }
            if (!costMythium) {
                var m = action.subheader.match(/Mythium\.png[^0-9]*(\d+)/i);
                if (m) costMythium = parseInt(m[1], 10);
            }
            if (!costSupply) {
                var s = action.subheader.match(/Supply\.png[^0-9]*(\d+)/i);
                if (s) costSupply = parseInt(s[1], 10);
            }
        }

        return {
            actionId: action.actionId,
            name: name,
            icon: action.image || '',
            costGold: costGold,
            costMythium: costMythium,
            costSupply: costSupply,
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
            sendState();
        });

        engine.on('refreshSupply', function (value) {
            gameState.supply = value;
            sendState();
        });

        engine.on('refreshSupplyCap', function (value) {
            gameState.supplyCap = value;
            sendState();
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
            if (actions.length > 0) dumpRawAction(actions[0], 'hand');
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
            if (actions.length > 0 && !gameState._mercActionSample) {
                var s = {};
                for (var k in actions[0]) { try { if (typeof actions[0][k] !== 'function' && typeof actions[0][k] !== 'undefined' && actions[0][k] !== null) { s[k] = actions[0][k]; } } catch(e) {} }
                gameState._mercActionSample = s;
            }
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

        // ===== Scoreboard info (all players' economy + field units) =====
        engine.on('refreshScoreboardInfo', function (scoreboardInfo, spectator) {
            if (!scoreboardInfo || !Array.isArray(scoreboardInfo)) return;
            gameState.scoreboardInfo = scoreboardInfo;
            sendState();
        });

        // ===== Team gold =====
        engine.on('refreshTeamGold', function (leftGold, rightGold) {
            gameState.teamGoldLeft = leftGold;
            gameState.teamGoldRight = rightGold;
            sendState();
        });

        // ===== King upgrades =====
        engine.on('refreshKingUpgrades', function (leftUpgrades, rightUpgrades) {
            gameState.kingUpgradesLeft = leftUpgrades;
            gameState.kingUpgradesRight = rightUpgrades;
            sendState();
        });

        // ===== Moneylender =====
        engine.on('refreshMoneylender', function (gold, cost, enabled) {
            gameState.moneylender = { gold: gold, cost: cost, enabled: enabled };
        });

        // ===== Reset on game end =====
        engine.on('refreshIsInGame', function (inGame) {
            if (!inGame) {
                _rawDumped = false;
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
                    selectedUnit: null,
                    scoreboardInfo: [],
                    teamGoldLeft: 0, teamGoldRight: 0,
                    kingUpgradesLeft: null, kingUpgradesRight: null,
                    moneylender: null
                };
                sendState();
            }
        });

        connect();
    }

    init();
})();
