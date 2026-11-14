const path = process.argv[2];
const list = await (await fetch('http://127.0.0.1:9223/json')).json();
const page = list.find(t => t.type === 'page' && t.webSocketDebuggerUrl) || list[0];
console.log('PAGE=' + (page?.title || '') + ' URL=' + (page?.url || ''));
const ws = new WebSocket(page.webSocketDebuggerUrl);
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let seq = 0; const pending = new Map();
ws.onmessage = e => { const m = JSON.parse(e.data); if (m.id && pending.has(m.id)) { pending.get(m.id)(m); pending.delete(m.id); } };
const call = (method, params) => new Promise(res => { const id = ++seq; pending.set(id, res); ws.send(JSON.stringify({ id, method, params })); });
async function evalJs(expression, awaitPromise = true) {
  const r = await call('Runtime.evaluate', { expression, awaitPromise, returnByValue: true });
  if (r.result?.exceptionDetails) return { error: r.result.exceptionDetails.text, detail: r.result.exceptionDetails.exception?.description };
  return r.result?.result?.value;
}
console.log('BINDINGS=' + JSON.stringify(await evalJs('Object.keys(window.go?.main?.App || {}).length', false)));
const models = await evalJs('window.go.main.App.GetLLMModels().then(m => ({count:m.length, lfm:m.filter(x=>x.name.includes("LFM")).map(x=>({name:x.name,sizeMB:Math.round(x.sizeMB)}))}))');
console.log('LLM_MODELS=' + JSON.stringify(models));
const res = await evalJs(`window.go.main.App.LoadModel(${JSON.stringify(path)})`);
console.log('LOAD_RESULT=' + JSON.stringify(res ?? 'resolved-null'));
ws.close();
