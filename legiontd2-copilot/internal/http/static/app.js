async function fetchState() {
  try {
    const resp = await fetch('/api/state');
    if (!resp.ok) throw new Error(resp.statusText);
    return await resp.json();
  } catch {
    return null;
  }
}

function updateUI(data) {
  const state = data.state || data;

  document.getElementById('mythium').textContent = state.mythium ?? '—';
  document.getElementById('income').textContent = state.income ?? '—';
  document.getElementById('wave').textContent = state.wave ?? '—';

  const timer = state.waveTimer;
  if (timer != null) {
    const m = Math.floor(timer / 60);
    const s = timer % 60;
    document.getElementById('timer').textContent = `${m}:${s.toString().padStart(2, '0')}`;
  } else {
    document.getElementById('timer').textContent = '—';
  }

  document.getElementById('kingHp').textContent = state.kingHp != null ? Math.round(state.kingHp) + '%' : '—';
  document.getElementById('enemyKingHp').textContent = state.enemyKingHp != null ? Math.round(state.enemyKingHp) + '%' : '—';
  document.getElementById('gold').textContent = state.gold ?? '—';
  document.getElementById('supply').textContent = (state.supply != null ? state.supply : '—') + '/' + (state.supplyCap ?? '—');

  // Hand units
  const handEl = document.getElementById('handUnits');
  handEl.innerHTML = '';
  if (state.hand && state.hand.length > 0) {
    state.hand.forEach(u => {
      const div = document.createElement('div');
      div.className = 'unit-card';
      const affordable = state.gold >= (u.costGold || 0) && (state.supplyCap - state.supply) >= (u.costSupply || 0) && u.stacks > 0;
      if (!affordable) div.className += ' dim';
      div.innerHTML = `<span class="u-name">${u.name || '?'}</span><span class="u-stock">x${u.stacks ?? 1}</span>`;
      if (u.costGold > 0) div.innerHTML += `<span class="u-cost"><img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3E%3Ccircle cx='8' cy='8' r='7' fill='%23ffcc00'/%3E%3C/svg%3E" style="width:10px;height:10px">${u.costGold}</span>`;
      handEl.appendChild(div);
    });
  } else {
    handEl.innerHTML = '<div class="dim">Нет данных</div>';
  }

  // Field units
  const fieldEl = document.getElementById('fieldUnits');
  fieldEl.innerHTML = '';
  if (state.fieldUnits) {
    const ids = Object.keys(state.fieldUnits);
    if (ids.length > 0) {
      ids.forEach(id => {
        const u = state.fieldUnits[id];
        const div = document.createElement('div');
        div.className = 'unit-card field';
        div.innerHTML = `<span class="u-name">#${id}</span><span class="u-hp">❤${Math.round(u.hp)}</span>`;
        fieldEl.appendChild(div);
      });
    } else {
      fieldEl.innerHTML = '<div class="dim">Поле пусто</div>';
    }
  }

  // Mercenaries
  const mercEl = document.getElementById('mercs');
  mercEl.innerHTML = '';
  if (state.mercenaries && state.mercenaries.length > 0) {
    state.mercenaries.forEach(m => {
      const div = document.createElement('div');
      div.className = 'unit-card dim';
      const canBuy = state.mythium >= (m.costMythium || 0);
      if (canBuy) div.className = 'unit-card';
      div.innerHTML = `<span class="u-name">${m.name || '?'}</span>`;
      if (m.costMythium > 0) div.innerHTML += `<span class="u-cost">🧪${m.costMythium}</span>`;
      mercEl.appendChild(div);
    });
  } else {
    mercEl.innerHTML = '<div class="dim">Нет данных</div>';
  }

  // Recommendations
  const recsEl = document.getElementById('recommendations');
  recsEl.innerHTML = '';
  const recs = data.recommendations || [];
  if (recs.length > 0) {
    recs.forEach(r => {
      const li = document.createElement('li');
      const badge = document.createElement('span');
      badge.className = 'badge';
      badge.textContent = r.action || r.kind || '';
      li.appendChild(badge);
      li.appendChild(document.createTextNode(r.message || r.explanation || ''));
      recsEl.appendChild(li);
    });
  } else {
    recsEl.innerHTML = '<li style="color:#8b949e">Нет данных</li>';
  }
}

async function tick() {
  const data = await fetchState();
  if (data) updateUI(data);
}

setInterval(tick, 1000);
tick();
