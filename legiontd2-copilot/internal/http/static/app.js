async function fetchState() {
  try {
    const resp = await fetch('/api/state');
    if (!resp.ok) throw new Error(resp.statusText);
    return await resp.json();
  } catch {
    return null;
  }
}

const HOTKEY_LABELS = ['Q','W','E','R','T','Y','U','I'];

function buildUnitCard(u, gold, supply, supplyCap, mythium) {
  const div = document.createElement('div');
  div.className = 'unit-card';

  const needGold = (u.costGold || 0) > gold;
  const needSupply = (u.costSupply || 0) > (supplyCap - supply);
  const noStacks = !u.stacks || u.stacks < 1;
  const blocked = needGold || needSupply || noStacks;

  if (blocked) div.className += ' dim';

  const iconImg = document.createElement('img');
  iconImg.className = 'u-icon';
  iconImg.src = '/icons/' + (u.icon || '').replace(/^Icons\//, '');
  iconImg.alt = u.name || '';
  iconImg.loading = 'lazy';

  const hotkey = document.createElement('span');
  hotkey.className = 'u-hotkey';
  hotkey.textContent = HOTKEY_LABELS[u.actionId] || '?';

  const infoWrap = document.createElement('div');
  infoWrap.className = 'u-info-wrap';

  const nameSpan = document.createElement('div');
  nameSpan.className = 'u-name';
  nameSpan.textContent = u.name || '?';

  const meta = document.createElement('div');
  meta.className = 'u-meta';

  const parts = [];
  if (u.costGold > 0) parts.push(u.costGold + 'g');
  if (u.costMythium > 0) parts.push(u.costMythium + 'm');
  if (u.costSupply > 0) parts.push(u.costSupply + 's');
  if (u.stacks > 0) parts.push('x' + u.stacks);
  meta.textContent = parts.join(' | ') || '';

  if (blocked) {
    const reason = document.createElement('div');
    reason.className = 'u-reason';
    if (noStacks) reason.textContent = 'нет в руке';
    else if (needGold) reason.textContent = 'не хватает ' + (u.costGold - gold) + 'g';
    else if (needSupply) reason.textContent = 'не хватает ' + (u.costSupply - (supplyCap - supply)) + ' снаряжения';
    infoWrap.appendChild(nameSpan);
    infoWrap.appendChild(meta);
    infoWrap.appendChild(reason);
  } else {
    infoWrap.appendChild(nameSpan);
    infoWrap.appendChild(meta);
  }

  div.appendChild(iconImg);
  div.appendChild(hotkey);
  div.appendChild(infoWrap);
  return div;
}

function updateUI(data) {
  const state = data.state || data;

  document.getElementById('mythium').textContent = state.mythium ?? '—';
  document.getElementById('gold').textContent = state.gold ?? '—';
  document.getElementById('income').textContent = state.income ?? '—';
  document.getElementById('wave').textContent = state.wave ?? '—';
  document.getElementById('phase').textContent = state.phase === 'fighting' ? 'БОЙ' : 'Стройка';

  const timer = state.waveTimer;
  if (timer != null) {
    const m = Math.floor(timer / 60);
    const s = timer % 60;
    document.getElementById('timer').textContent = m + ':' + (s < 10 ? '0' : '') + s;
  } else {
    document.getElementById('timer').textContent = '—';
  }

  document.getElementById('supply').textContent = (state.supply != null ? state.supply : '—') + '/' + (state.supplyCap ?? '—');

  const statusEl = document.getElementById('status');
  const connected = state.mythium > 0 || state.gold > 0 || (state.hand && state.hand.length > 0);
  if (connected) {
    statusEl.innerHTML = '● Статус: подключено';
    statusEl.style.color = '#3fb950';
  } else {
    statusEl.innerHTML = '● Статус: ожидание игры';
    statusEl.style.color = '#8b949e';
  }

  const handEl = document.getElementById('handUnits');
  handEl.innerHTML = '';
  if (state.hand && state.hand.length > 0) {
    state.hand.forEach(function(u) {
      handEl.appendChild(buildUnitCard(u, state.gold, state.supply, state.supplyCap, state.mythium));
    });
  } else {
    handEl.innerHTML = '<div class="dim" style="padding:8px">Нет данных</div>';
  }

  const mercEl = document.getElementById('mercs');
  mercEl.innerHTML = '';
  if (state.mercenaries && state.mercenaries.length > 0) {
    state.mercenaries.forEach(function(m) {
      mercEl.appendChild(buildUnitCard(m, state.gold, state.supply, state.supplyCap, state.mythium));
    });
  } else {
    mercEl.innerHTML = '<div class="dim" style="padding:8px">Нет данных</div>';
  }

  const recsEl = document.getElementById('recommendations');
  recsEl.innerHTML = '';
  var recs = data.recommendations || [];
  if (recs.length > 0) {
    recs.forEach(function(r) {
      var li = document.createElement('li');
      var badge = document.createElement('span');
      badge.className = 'badge badge-' + (r.action || 'info');
      badge.textContent = r.action || r.kind || '?';
      li.appendChild(badge);
      li.appendChild(document.createTextNode(r.message || r.explanation || ''));
      recsEl.appendChild(li);
    });
  } else {
    recsEl.innerHTML = '<li class="dim">Нет данных</li>';
  }

  // Debug: raw JSON dump
  var debugEl = document.getElementById('rawDebug');
  if (debugEl) {
    var handNames = (state.hand || []).map(function(u) { return JSON.stringify({id: u.actionId, name: u.name}); }).join('\n');
    var mercNames = (state.mercenaries || []).map(function(u) { return JSON.stringify({id: u.actionId, name: u.name}); }).join('\n');
    var info = '';
    if (state._actionSample) {
      info += '--- Raw action sample (first hand unit) ---\n' + JSON.stringify(state._actionSample, null, 2) + '\n\n';
    }
    if (state._mercActionSample) {
      info += '--- Raw merc action sample (first merc) ---\n' + JSON.stringify(state._mercActionSample, null, 2) + '\n\n';
    }
    if (data.matrix) {
      info += '--- Feature Matrix (AI) ---\n' + JSON.stringify(data.matrix, null, 2) + '\n\n';
    }
    info += '--- Hand unit names ---\n' + (handNames || '(empty)');
    info += '\n\n--- Merc names ---\n' + (mercNames || '(empty)');
    info += '\n\n--- Full state ---\n' + JSON.stringify(state, null, 2);
    debugEl.textContent = info;

    var copyBtn = document.getElementById('copyDebug');
    if (copyBtn) {
      copyBtn.onclick = function() {
        navigator.clipboard.writeText(info).then(function() {
          copyBtn.textContent = 'Скопировано!';
          setTimeout(function() { copyBtn.textContent = 'Копировать'; }, 2000);
        });
      };
    }
  }
}

async function tick() {
  var data = await fetchState();
  if (data) updateUI(data);
}

setInterval(tick, 1000);
tick();
