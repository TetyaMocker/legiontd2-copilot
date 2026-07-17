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

// Try WebSocket for real-time updates, fall back to polling
let ws = null;
function connectWS() {
  try {
    ws = new WebSocket('ws://' + location.host + '/ws');
    ws.onmessage = function(e) {
      try {
        const data = JSON.parse(e.data);
        if (data.type === 'recommendation') {
          const state = awaitFetchState();
          state.then(s => { if (s) updateUI(s); });
        }
      } catch {}
    };
    ws.onclose = function() { setTimeout(connectWS, 3000); };
  } catch {
    setTimeout(connectWS, 5000);
  }
}

async function awaitFetchState() {
  return await fetchState();
}

// Fallback polling
async function tick() {
  const data = await fetchState();
  if (data) updateUI(data);
}

connectWS();
setInterval(tick, 1000);
tick();
