async function fetchState() {
  try {
    const resp = await fetch('/api/state');
    if (!resp.ok) throw new Error(resp.statusText);
    return await resp.json();
  } catch {
    return null;
  }
}

function updateUI(state) {
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

  document.getElementById('kingHp').textContent = state.kingHp != null ? state.kingHp + '%' : '—';

  const conf = state.confidence;
  const statusEl = document.getElementById('status');
  if (conf == null) {
    statusEl.textContent = 'Нет данных';
    statusEl.className = 'status err';
    document.getElementById('confidence').textContent = '—';
  } else {
    const pct = Math.round(conf * 100);
    document.getElementById('confidence').textContent = pct + '%';
    if (conf >= 0.7) {
      statusEl.textContent = 'Распознавание стабильно';
      statusEl.className = 'status ok';
    } else if (conf >= 0.3) {
      statusEl.textContent = 'Неуверенное распознавание';
      statusEl.className = 'status warn';
    } else {
      statusEl.textContent = 'Распознавание недоступно';
      statusEl.className = 'status err';
    }
  }

  const recsEl = document.getElementById('recommendations');
  recsEl.innerHTML = '';
  if (state.recommendations && state.recommendations.length > 0) {
    state.recommendations.forEach(r => {
      const li = document.createElement('li');
      const badge = document.createElement('span');
      badge.className = 'badge badge-' + r.kind;
      badge.textContent = r.kind;
      li.appendChild(badge);
      li.appendChild(document.createTextNode(r.explanation));
      recsEl.appendChild(li);
    });
  } else {
    recsEl.innerHTML = '<li style="color:#8b949e">Нет рекомендаций</li>';
  }
}

async function tick() {
  const state = await fetchState();
  if (state) updateUI(state);
}

setInterval(tick, 1000);
tick();
