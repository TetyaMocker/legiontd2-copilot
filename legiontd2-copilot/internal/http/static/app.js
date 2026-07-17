async function fetchState() {
  try {
    const resp = await fetch('/api/state');
    if (!resp.ok) throw new Error(resp.statusText);
    return await resp.json();
  } catch {
    return null;
  }
}

function buildUnitCard(u, gold, supply, supplyCap, mythium) {
  const div = document.createElement('div');
  div.className = 'unit-card';
  const canAfford = gold >= (u.costGold || 0)
    && (supplyCap - supply) >= (u.costSupply || 0)
    && u.stacks > 0;
  if (!canAfford) div.className += ' dim';

  const nameSpan = document.createElement('span');
  nameSpan.className = 'u-name';
  nameSpan.textContent = u.name || '?';

  const btn = document.createElement('button');
  btn.className = 'action-btn';
  btn.textContent = 'Ставь';
  btn.disabled = !canAfford;

  const info = document.createElement('span');
  info.className = 'u-info';
  const parts = [];
  if (u.costGold > 0) parts.push(u.costGold + 'g');
  if (u.costMythium > 0) parts.push(u.costMythium + 'm');
  if (u.costSupply > 0) parts.push(u.costSupply + 's');
  if (u.stacks > 0) parts.push('x' + u.stacks);
  info.textContent = parts.join(' ') || '';

  div.appendChild(nameSpan);
  div.appendChild(info);
  div.appendChild(btn);
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
}

async function tick() {
  var data = await fetchState();
  if (data) updateUI(data);
}

setInterval(tick, 1000);
tick();
