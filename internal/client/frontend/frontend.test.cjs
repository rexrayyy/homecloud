const {test}=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');
const path=require('node:path');
const html=fs.readFileSync(path.join(__dirname,'index.html'),'utf8');
const script=fs.readFileSync(path.join(__dirname,'app.js'),'utf8');
function storageMock(){const s={};return {getItem:k=>k in s?s[k]:null,setItem:(k,v)=>{s[k]=String(v)},removeItem:k=>{delete s[k]},clear:()=>{for(const k of Object.keys(s))delete s[k]}};}
function harness(forceLang){
  const elements={};
  function element(){return {value:'',textContent:'',disabled:false,checked:false,hidden:false,dataset:{},classList:{toggle(){}},addEventListener(){},setAttribute(){},getAttribute(){return null},append(){},replaceChildren(){},querySelector(){return null;},getBoundingClientRect(){return {width:800,height:280};}};}
  for(const m of html.matchAll(/id="([^"]+)"/g))elements[m[1]]=element();
  elements.port.value='22';elements.serverPort.value='8081';elements.deployPort.value='8081';
  const store=storageMock();
  const context=vm.createContext({
    document:{documentElement:{lang:'en'},getElementById:id=>{assert.ok(elements[id],`missing DOM id ${id}`);return elements[id];},querySelectorAll:()=>[],createElement:element},
    window:{addEventListener(){},localStorage:store},
    localStorage:store,
    confirm:()=>true,setTimeout,clearTimeout,console
  });
  vm.runInContext(script,context);
  const h={elements,store,run:code=>vm.runInContext(code,context)};
  if(forceLang)h.run(`setLanguage(${JSON.stringify(forceLang)},false)`);
  return h;
}
test('all referenced elements exist; static preview boots',()=>{const h=harness('en');assert.match(h.elements['runtime-state'].textContent,/Static preview/);assert.equal(h.elements.connect.disabled,true);});
test('host probe never sends password; verified connection required',()=>{const h=harness('en');h.elements.host.value='server.local';h.elements.user.value='alice';h.elements.password.value='secret';assert.equal(h.run('config(false,true).password'),'');assert.throws(()=>h.run('config()'),/verify/i);h.run("state.fingerprint='SHA256:abc'");h.elements.trust.checked=true;assert.equal(h.run('config().host'),'server.local');h.run('resetTrust()');assert.equal(h.elements.trust.checked,false);assert.equal(h.elements.connect.disabled,true);});
test('business validation rejects invalid host, port and traversal',()=>{const h=harness('en');h.elements.host.value='https://server';assert.throws(()=>h.run('config(false,true)'),/hostname/);h.elements.host.value='server';h.elements.port.value='0';assert.throws(()=>h.run('config(false,true)'),/SSH port/);assert.throws(()=>h.run("safePath('../secret')"),/traversal/);assert.throws(()=>h.run("safePath('foo\\\\bar')"),/backslashes/);assert.equal(h.run("safePath('/music/song.mp3')"),'/music/song.mp3');});
test('errors redact bearer and token-like hex; unavailable operations do not replay',async()=>{const h=harness('en');assert.equal(h.run("normalize('Bearer abc123')"),'Bearer [redacted]');assert.equal(h.run(`normalize('${'a'.repeat(64)}')`),'[redacted]');await assert.rejects(h.run("call('Request','/api/v1/stats','')"),/unavailable/);assert.equal(h.run('state.queue[0]'),'Request');assert.throws(()=>h.run("requireConnection('UploadFile')"),/not been executed/);});
test('terminal output is bounded',()=>{const h=harness('en');h.run("terminalOutput('hello')");assert.equal(h.run('terminalTail'),'hello');h.run("terminalOutput('x'.repeat(160000))");assert.equal(h.run('terminalTail.length'),150000);});
test('native method literals exist on client.App',()=>{const backend=['backend_windows.go','deploy_windows.go','devices_windows.go','files_sftp_windows.go','rdp_windows.go'].map(f=>fs.readFileSync(path.join(__dirname,'../'+f),'utf8')).join('\n');for(const m of script.matchAll(/call\('([A-Z]\w+)'/g))assert.match(backend,new RegExp(`func \\(a \\*App\\) ${m[1]}\\(`));for(const name of ['Connect','Deploy'])assert.match(backend,new RegExp(`func \\(a \\*App\\) ${name}\\(`));});
test('language defaults to Russian, switches to English and persists',()=>{
  const h=harness();
  assert.equal(h.run('language'),'ru');
  assert.match(h.elements['runtime-state'].textContent,/Статический просмотр/);
  assert.equal(h.elements['language'].value,'ru');
  assert.equal(h.run("String(t('Dashboard'))"),'Обзор');
  h.run("setLanguage('en')");
  assert.equal(h.run('language'),'en');
  assert.match(h.elements['runtime-state'].textContent,/Static preview/);
  assert.equal(h.elements['language'].value,'en');
  assert.equal(h.run("String(t('Dashboard'))"),'Dashboard');
  assert.equal(h.store.getItem('homecloud.language'),'en');
  h.run("setLanguage('ru')");
  assert.equal(h.run('language'),'ru');
  assert.match(h.elements['runtime-state'].textContent,/Статический просмотр/);
  assert.equal(h.store.getItem('homecloud.language'),'ru');
});
test('saved language is restored on boot',()=>{
  const elements={};
  function element(){return {value:'',textContent:'',disabled:false,checked:false,hidden:false,dataset:{},classList:{toggle(){}},addEventListener(){},setAttribute(){},getAttribute(){return null},append(){},replaceChildren(){},querySelector(){return null;},getBoundingClientRect(){return {width:800,height:280};}};}
  for(const m of html.matchAll(/id="([^"]+)"/g))elements[m[1]]=element();
  elements.port.value='22';elements.serverPort.value='8081';elements.deployPort.value='8081';
  const store=storageMock();store.setItem('homecloud.language','en');
  const context=vm.createContext({document:{documentElement:{lang:'ru'},getElementById:id=>elements[id],querySelectorAll:()=>[],createElement:element},window:{addEventListener(){},localStorage:store},localStorage:store,confirm:()=>true,setTimeout,clearTimeout,console});
  vm.runInContext(script,context);
  assert.equal(vm.runInContext('language',context),'en');
  assert.match(elements['runtime-state'].textContent,/Static preview/);
});