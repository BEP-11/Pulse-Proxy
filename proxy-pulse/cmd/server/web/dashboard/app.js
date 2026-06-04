const metrics = document.getElementById('metrics');
const keys = ['proxy_pulse_active_connections_total', 'go_goroutines'];
keys.forEach(k => {
  const div = document.createElement('div');
  div.style.padding='1rem';div.style.borderRadius='8px';div.style.background='#f4f4f5';
  div.id = k; metrics.appendChild(div);
});

setInterval(async () => {
  const res = await fetch('/metrics');
  const text = await res.text();
  keys.forEach(k => {
    const match = text.match(new RegExp(`${k}[\\s]+([0-9.]+)`));
    const el = document.getElementById(k);
    if(match) el.innerHTML = `<strong>${k.replace(/_/g,' ')}</strong><br><span style="font-size:1.5rem;color:#2563eb">${match[1]}</span>`;
  });
}, 2000);
