'use strict';

// English source keys; owned display text is separate from raw server data.
const translations={
  "Your private workspace": "Ваше личное пространство",
  "Disconnected": "Нет подключения",
  "Token unavailable": "Токен недоступен",
  "Workspace": "Рабочее пространство",
  "Dashboard": "Обзор",
  "Containers": "Контейнеры",
  "Files": "Файлы",
  "Media": "Медиа",
  "Terminal": "Терминал",
  "Connect securely to get started. Your password is never saved.": "Начните с безопасного подключения. Пароль никогда не сохраняется.",
  "Secure connection": "Безопасное подключение",
  "One SSH tunnel. Your services stay private.": "Один SSH-туннель. Ваши сервисы остаются приватными.",
  "Disconnect": "Отключиться",
  "Nearby servers": "Серверы рядом",
  "Scan network": "Сканировать сеть",
  "Looking for HomeCloud advertisements…": "Поиск объявлений HomeCloud…",
  "Discovery is a hint, not proof of identity. Private localhost deployments are not advertised; enter their address below.": "Обнаружение не подтверждает подлинность сервера. Приватные установки localhost не публикуются; введите адрес ниже.",
  "Connect to an existing instance": "Подключиться к существующему экземпляру",
  "Host / IP": "Хост / IP",
  "SSH port": "Порт SSH",
  "HomeCloud port": "Порт HomeCloud",
  "SSH user": "Пользователь SSH",
  "Password": "Пароль",
  "1. Read host fingerprint": "1. Получить отпечаток ключа",
  "No host key inspected": "Ключ сервера ещё не проверен",
  "Compare this fingerprint through a trusted channel on the server, for example:": "Сравните отпечаток через доверенный канал на сервере, например:",
  "I independently verified this exact host fingerprint": "Я независимо проверил именно этот отпечаток ключа сервера",
  "2. Connect securely": "2. Безопасно подключиться",
  "Password cleared after each connection attempt.": "Пароль очищается после каждой попытки подключения.",
  "Install or start HomeCloud": "Установить или запустить HomeCloud",
  "Ubuntu · Intel/AMD x86_64 · no sudo. Uses the host, SSH credentials and verified key above.": "Ubuntu · Intel/AMD x86_64 · без sudo. Используются хост, данные SSH и проверенный ключ выше.",
  "Instance port": "Порт экземпляра",
  "Each port has separate files and settings. The service binds to server localhost. Run again after a server reboot. Deployment does not automatically connect.": "У каждого порта свои файлы и настройки. Сервис слушает localhost сервера. После перезагрузки запустите снова. Установка не подключает автоматически.",
  "Install / start instance": "Установить / запустить экземпляр",
  "No live snapshot yet": "Актуальных данных пока нет",
  "Refresh": "Обновить",
  "CPU utilization": "Загрузка процессора",
  "Memory used": "Занято памяти",
  "Total memory": "Всего памяти",
  "Server snapshot": "Данные сервера",
  "Processes": "Процессы",
  "Top 20 by memory": "20 крупнейших по памяти",
  "Name": "Имя",
  "Memory": "Память",
  "Connect and refresh to load processes.": "Подключитесь и обновите список процессов.",
  "Manage workloads on your server.": "Управляйте нагрузкой на сервере.",
  "Name / image": "Имя / образ",
  "State": "Состояние",
  "Actions": "Действия",
  "Connect and refresh to load containers.": "Подключитесь и обновите список контейнеров.",
  "Stop and restart require confirmation. Docker must be available on the server.": "Остановка и перезапуск требуют подтверждения. На сервере должен быть доступен Docker.",
  "Authenticated WebDAV over your SSH connection.": "WebDAV с аутентификацией через SSH-подключение.",
  "Upload file": "Загрузить файл",
  "Parent folder": "Родительская папка",
  "↑ Up": "↑ Вверх",
  "Directory": "Каталог",
  "Open / refresh": "Открыть / обновить",
  "Size": "Размер",
  "Modified": "Изменён",
  "Action": "Действие",
  "Connect to browse your files.": "Подключитесь для просмотра файлов.",
  "Map a Windows drive": "Подключить диск Windows",
  "Drive letter": "Буква диска",
  "Check mapping": "Проверить подключение диска",
  "Windows WebClient does not support the server's Bearer authentication. No credentials or unauthenticated proxy will be created. File transfers above remain available.": "Windows WebClient не поддерживает Bearer-аутентификацию сервера. Учётные данные и прокси без аутентификации не создаются. Передача файлов выше остаётся доступной.",
  "SSH PTY · live Wails events · plain-text terminal view": "SSH PTY · события Wails · текстовый терминал",
  "Open session": "Открыть сеанс",
  "Close": "Закрыть",
  "Commands run as your SSH user. This safe text view strips ANSI styling; full-screen programs are not supported.": "Команды выполняются от пользователя SSH. Текстовый вывод удаляет ANSI-оформление; полноэкранные программы не поддерживаются.",
  "Terminal output": "Вывод терминала",
  "Command": "Команда",
  "Type a command": "Введите команду",
  "Send ↵": "Отправить ↵",
  "Clear view": "Очистить экран",
  "Activity & errors": "Действия и ошибки",
  "Clear log": "Очистить журнал",
  "Error console": "Консоль ошибок",
  "No errors.": "Ошибок нет.",
  "Unavailable operations": "Недоступные операции",
  "Requests are held as reminders only. Nothing is replayed automatically, especially writes.": "Запросы сохраняются только как напоминания. Ничего не повторяется автоматически, особенно запись.",
  "No pending operations.": "Ожидающих операций нет.",
  "Clear reminders": "Очистить напоминания",
  "HomeCloud · native workspace": "HomeCloud · рабочее пространство",
  "Checking desktop bridge…": "Проверка связи с приложением…",
  "Unknown error": "Неизвестная ошибка",
  "Connect from the SSH tab first. This operation has not been executed.": "Сначала подключитесь на вкладке SSH. Операция не выполнена.",
  "Enter a hostname or IP without a scheme, and an SSH port from 1 to 65535.": "Введите имя хоста или IP без схемы и порт SSH от 1 до 65535.",
  "Instance port must be 1024–65535.": "Порт экземпляра должен быть от 1024 до 65535.",
  "SSH user and password are required.": "Требуются пользователь и пароль SSH.",
  "Independently verify the host fingerprint first.": "Сначала независимо проверьте отпечаток ключа сервера.",
  "Token ready · protected": "Токен готов · защищён",
  "Installing and checking service…": "Установка и проверка сервиса…",
  "Establishing verified SSH connection…": "Установка проверенного SSH-подключения…",
  "Connected securely. Token stays in the native backend.": "Безопасное подключение установлено. Токен остаётся в нативной части приложения.",
  "Read host key": "Получение ключа сервера",
  "Invalid fingerprint returned.": "Получен недопустимый отпечаток.",
  "Host key retrieved but NOT trusted. Verify independently before checking the box.": "Ключ получен, но НЕ является доверенным. Проверьте его независимо, прежде чем установить флажок.",
  "Connect": "Подключение",
  "Deploy": "Установка",
  "Disconnected. Live snapshots may be stale.": "Отключено. Данные могут быть устаревшими.",
  "Scanning mDNS…": "Сканирование mDNS…",
  "Discovery selected. Confirm SSH port and verify the host key before connecting.": "Сервер выбран. Уточните порт SSH и проверьте ключ перед подключением.",
  "No advertisements found. Enter a host manually.": "Объявления не найдены. Введите хост вручную.",
  "Discovery unavailable. Enter a host manually.": "Обнаружение недоступно. Введите хост вручную.",
  "Discovery": "Обнаружение",
  "No processes returned.": "Процессы не получены.",
  "start": "Запустить",
  "stop": "Остановить",
  "restart": "Перезапустить",
  "running": "Работает",
  "exited": "Завершён",
  "created": "Создан",
  "paused": "Приостановлен",
  "restarting": "Перезапускается",
  "removing": "Удаляется",
  "dead": "Не работает",
  "No containers on this server.": "На сервере нет контейнеров.",
  "Use a relative path without traversal or backslashes.": "Используйте относительный путь без переходов .. и обратных косых черт.",
  "Folder": "Папка",
  "Open": "Открыть",
  "Download": "Скачать",
  "Download dialog completed.": "Диалог скачивания завершён.",
  "This folder is empty.": "Папка пуста.",
  "Upload": "Загрузка",
  "Upload files": "Загрузить файлы",
  "Upload a folder": "Загрузить папку",
  "Cancel": "Отмена",
  "Upload a file into this folder? A file with the same name may be overwritten.": "Загрузить файл в эту папку? Файл с таким именем может быть перезаписан.",
  "Upload dialog completed.": "Диалог загрузки завершён.",
  "Drive mapping": "Подключение диска",
  "Choose a drive letter D–Z.": "Выберите букву диска D–Z.",
  "MountDrive · requires a supported authenticated mapping solution": "MountDrive · требуется способ подключения с поддерживаемой аутентификацией",
  "Open terminal": "Открытие терминала",
  "terminal events": "События терминала",
  "Native Wails event bridge unavailable. Terminal cannot display output.": "Мост событий Wails недоступен. Терминал не может показывать вывод.",
  "Close terminal": "Закрытие терминала",
  "Terminal emulator failed to load.": "Не удалось загрузить эмулятор терминала.",
  "Terminal input": "Ввод терминала",
  "This command may destroy data or interrupt services. Execute it?": "Команда может уничтожить данные или прервать работу сервисов. Выполнить?",
  "Interrupt": "Прерывание",
  "Clear terminal view": "Очистка экрана терминала",
  "Clear errors": "Очистка ошибок",
  "Terminal resize": "Изменение размера терминала",
  "Unexpected failure": "Неожиданная ошибка",
  "Static preview · native bridge unavailable": "Статический просмотр · нативный мост недоступен",
  "Open the native app for mDNS discovery; manual fields are available.": "Откройте приложение для обнаружения mDNS; доступен ручной ввод.",
  "Native Wails bridge": "Нативный мост Wails",
  "Connection closed": "Соединение закрыто",
  "Connection status": "Состояние подключения",
  "Language": "Язык",
  "Server response": "Ответ сервера",
  "Snapshot · {0}": "Снимок · {0}",
  "{0} server(s) found. Manual fields below always override discovery.": "Найдено серверов: {0}. Ручной ввод ниже имеет приоритет.",
  "{0} is unavailable in this build. Open the native desktop app or update it.": "{0} недоступен в этой сборке. Откройте или обновите нативное приложение.",
  "Install or start HomeCloud on {0}, instance {1}? This writes the application and starts a background service.": "Установить или запустить HomeCloud на {0}, экземпляр {1}? Будет записано приложение и запущен фоновый сервис.",
  "Container {0}": "Контейнер: {0}",
  "{0} {1}? Running work may be interrupted.": "{0} {1}? Текущая работа может быть прервана.",
  "Container {0} completed.": "Действие с контейнером «{0}» выполнено.",
  "[SSH session opened]": "[Сеанс SSH открыт]",
  "[Session closed]": "[Сеанс закрыт]"
};

// Additional translations for file manager grid view and RDP
Object.assign(translations, {
  "File manager over SSH connection.": "Файловый менеджер через SSH-подключение.",
  "Grid view": "Сетка",
  "Table view": "Таблица",
  "Toggle view": "Переключить вид",
  "Path": "Путь",
  "Remote Desktop": "Удалённый рабочий стол",
  "View the mini-PC screen via SSH (screenshot capture, view-only).": "Просмотр экрана мини-ПК через SSH (захват снимков, только просмотр).",
  "Capture screen": "Сделать снимок",
  "Auto refresh": "Автообновление",
  "Requires screenshot tools (scrot, imagemagick, or xwd) installed on the remote system. Display must be active.": "Требуются инструменты скриншотов (scrot, imagemagick или xwd) на удалённой системе. Дисплей должен быть активен.",
  "Not captured yet.": "Снимок ещё не сделан.",
  "Click 'Capture screen' to view the remote desktop.": "Нажмите «Сделать снимок» для просмотра удалённого рабочего стола.",
  "Screen captured": "Снимок экрана получен",
  "Capturing screen…": "Создание снимка экрана…",
  "Screen capture failed": "Ошибка снимка экрана",
  "Auto refresh started": "Автообновление включено",
  "Auto refresh stopped": "Автообновление выключено",
  "Stop auto refresh": "Остановить автообновление",
  "Screen dimensions: {0}": "Разрешение экрана: {0}",
  "Open folder": "Открыть папку",
  "Download file": "Скачать файл",
  "Folder": "Папка",
  "File": "Файл",
  "Empty folder": "Пустая папка",
  "This folder is empty.": "Эта папка пуста.",
  "Connect to browse your files.": "Подключитесь для просмотра файлов.",
  "Home": "Домой",
  "Home folder": "Домашняя папка",
  "Search this folder": "Поиск в этой папке",
  "Search": "Искать",
  "Searching…": "Поиск…",
  "Search results: {0}": "Найдено: {0}",
  "No matches.": "Совпадений нет.",
  "Enter a search term.": "Введите поисковый запрос.",
  "Search failed": "Ошибка поиска",
  "Downloaded to {0}": "Скачано в {0}",
  "Uploaded {0} file(s).": "Загружено файлов: {0}.",
  "Nothing was uploaded.": "Ничего не загружено.",
  "Upload files into {0}? Existing files with the same name may be overwritten.": "Загрузить файлы в {0}? Существующие файлы с теми же именами могут быть перезаписаны.",
  "Files will be saved to {0}.": "Файлы будут сохранены в {0}.",
  "Uploaded {0} file(s) to {1}.": "Загружено файлов: {0} в {1}.",
  "Connect to the mini-PC with mstsc.exe, or view a screenshot over SSH.": "Подключение к мини-ПК через mstsc.exe или просмотр снимка по SSH.",
  "Connect via RDP": "Подключиться по RDP",
  "Forget credentials": "Забыть учётные данные",
  "RDP host / IP": "Хост / IP RDP",
  "RDP port": "Порт RDP",
  "RDP user": "Пользователь RDP",
  "Width": "Ширина",
  "Height": "Высота",
  "Fullscreen": "Во весь экран",
  "Save credentials in Windows": "Сохранить учётные данные в Windows",
  "Check port": "Проверить порт",
  "Probe xrdp over SSH": "Проверить xrdp по SSH",
  "Not connected yet.": "Ещё не подключено.",
  "Launching Remote Desktop…": "Запуск удалённого рабочего стола…",
  "RDP session launched.": "RDP-сессия запущена.",
  "Checking RDP port…": "Проверка порта RDP…",
  "RDP port reachable.": "Порт RDP доступен.",
  "RDP port unreachable.": "Порт RDP недоступен.",
  "Probing xrdp over SSH…": "Проверка xrdp по SSH…",
  "xrdp ready.": "xrdp готов.",
  "xrdp unavailable.": "xrdp недоступен.",
  "RDP credentials removed.": "Учётные данные RDP удалены.",
  "Enter an RDP hostname or IP without a scheme.": "Введите имя хоста или IP RDP без схемы.",
  "RDP port must be 1–65535.": "Порт RDP должен быть 1–65535.",
  "Saved devices": "Сохранённые устройства",
  "Saved device": "Сохранённое устройство",
  "Load": "Загрузить",
  "Save": "Сохранить",
  "Delete": "Удалить",
  "Device name": "Имя устройства",
  "No saved devices.": "Нет сохранённых устройств.",
  "View-only screenshot over SSH": "Снимок экрана по SSH (только просмотр)",
  "Select a saved device first.": "Сначала выберите сохранённое устройство.",
  "Device loaded.": "Устройство загружено.",
  "Device saved.": "Устройство сохранено.",
  "Remember device": "Сохранение устройства",
  "Saved device restored.": "Сохранённое устройство восстановлено.",
  "Saved devices": "Сохранённые устройства",
  "Device deleted.": "Устройство удалено.",
  "Enter the RDP host and user to save.": "Укажите хост и пользователя RDP для сохранения.",
  "Delete this saved device?": "Удалить это сохранённое устройство?",
  "HomeDir is unavailable in this build.": "HomeDir недоступен в этой сборке.",
});
const languageKey='homecloud.language';
let language='ru';
try{const saved=window.localStorage?.getItem(languageKey);if(saved==='en'||saved==='ru')language=saved;}catch{}
class Message{constructor(render){this.render=render;}toString(){return this.render();}}
function t(key,...args){return new Message(()=>{const template=language==='ru'?(translations[key]||key):key;return String(template).replace(/\{(\d+)\}/g,(_,i)=>String(args[i]??''));});}
function joinText(...parts){return new Message(()=>parts.map(String).join(''));}
class UIError extends Error{constructor(value){super(String(value));this.localized=value;}}
const textBindings=new Map();
function put(element,value){if(value instanceof Message)textBindings.set(element,value);else textBindings.delete(element);element.textContent=String(value??'');}
function initLanguage(){
 for(const el of document.querySelectorAll('[data-i18n]'))put(el,t(el.dataset.i18n));
 $('language').addEventListener('change',e=>setLanguage(e.target.value));
 window.addEventListener('storage',e=>{if(e.key===languageKey)setLanguage(e.newValue,false);});
 setLanguage(language,false);
}
function setLanguage(value,persist=true){
 language=value==='en'?'en':'ru';if(document.documentElement)document.documentElement.lang=language;const selector=$('language');if(selector)selector.value=language;
 if(persist)try{(window.localStorage||localStorage)?.setItem(languageKey,language);}catch{}
 for(const [el,message] of textBindings){if(el.isConnected===false){textBindings.delete(el);continue;}el.textContent=String(message);}
 for(const attr of ['aria-label','placeholder'])for(const el of document.querySelectorAll('[data-i18n-'+attr+']'))el.setAttribute(attr,String(t(el.getAttribute('data-i18n-'+attr))));
}

const $ = id => document.getElementById(id);
const state = {connected:false,terminalOpen:false,fingerprint:'',busy:false,tab:'ssh',errors:[],queue:[],fileView:'grid',rdpAuto:false,rdpTimer:null,currentFiles:[],devices:[]};
const bridge = () => window.go?.client?.App;
function redact(value) { return String(value).replace(/Bearer\s+\S+/gi,'Bearer [redacted]').replace(/\b[a-f0-9]{64}\b/gi,'[redacted]'); }
function normalize(error) { if(error instanceof UIError)return error.localized; let value=error?.message || error; if(typeof value==='object')value=value.error || JSON.stringify(value); return value?redact(value).slice(0,2000):t('Unknown error'); }
function notice(message,error=false) { put($('notice'),message instanceof Message?message:redact(message)); $('notice').classList.toggle('error',error); }
function report(error,label) { const message=normalize(error); notice(joinText(label instanceof Message?label:t(label),': ',message),true); state.errors.push(joinText(new Date().toLocaleTimeString(),' ',label instanceof Message?label:t(label),': ',message)); state.errors=state.errors.slice(-80); put($('errors'),joinText(...state.errors.flatMap((e,i)=>i?['\n',e]:[e]))); put($('error-count'),state.errors.length); }
function queue(label) { if(!state.queue.includes(label)) state.queue.push(label); state.queue=state.queue.slice(-30); $('queue').replaceChildren(...state.queue.map(x=>{const li=document.createElement('li');put(li,x);return li;})); }
async function call(method,...args) { if(typeof bridge()?.[method]!=='function') {queue(method);throw new UIError(t('{0} is unavailable in this build. Open the native desktop app or update it.',method));} return bridge()[method](...args); }
async function request(path,body='') { requireConnection(path); try {const raw=await call('Request',path,body);return typeof raw==='string'?JSON.parse(raw):raw;} catch(e) {if(/unavailable|404|503|not implemented/i.test(normalize(e)))queue(path);throw e;} }
function requireConnection(label) {if(!state.connected){queue(label);throw new UIError(t('Connect from the SSH tab first. This operation has not been executed.'));}}
async function action(label,fn,button) {if(button?.disabled)return; if(button)button.disabled=true;try {await fn();}catch(e){report(e,label);}finally{if(button)button.disabled=false;update();}}
function on(id,label,fn) {$(id).addEventListener('click',e=>action(label,fn,e.currentTarget));}
function form(id,label,fn) {$(id).addEventListener('submit',e=>{e.preventDefault();action(label,fn,$(id).querySelector('[type=submit]'));});}
function update() {for(const id of ['host','port','serverPort','user','password','deployPort','probe'])$(id).disabled=state.busy; $('trust').disabled=state.busy||!state.fingerprint; for(const id of ['connect','deploy'])$(id).disabled=state.busy||!state.fingerprint||!$('trust').checked; $('disconnect').disabled=state.busy||!state.connected; $('terminal-open').disabled=state.terminalOpen; $('terminal-close').disabled=!state.terminalOpen;}
function resetTrust(){state.fingerprint='';$('trust').checked=false;put($('fingerprint'),t('No host key inspected'));update();}
function config(deploy=false,probe=false) {const c={host:$('host').value.trim(),port:Number($('port').value),serverPort:Number($(deploy?'deployPort':'serverPort').value),user:$('user').value.trim(),password:probe?'':$('password').value,fingerprint:state.fingerprint};if(!c.host||/[\s/\\]/.test(c.host)||!Number.isInteger(c.port)||c.port<1||c.port>65535)throw new UIError(t('Enter a hostname or IP without a scheme, and an SSH port from 1 to 65535.'));if(!probe){if(!Number.isInteger(c.serverPort)||c.serverPort<1024||c.serverPort>65535)throw new UIError(t('Instance port must be 1024–65535.'));if(!c.user||!c.password)throw new UIError(t('SSH user and password are required.'));if(!c.fingerprint||!$('trust').checked)throw new UIError(t('Independently verify the host fingerprint first.'));}return c;}
function applyStatus(s) {Object.assign(state,{connected:!!s.connected,terminalOpen:!!s.terminalOpen});put($('connection'),s.connected?`SSH · ${s.host}:${s.serverPort}`:t('Disconnected'));$('connection').classList.toggle('online',!!s.connected);put($('token'),s.tokenReady?t('Token ready · protected'):t('Token unavailable'));update();}
async function status(){applyStatus(await call('Status'));}
async function rememberDevice(c){try{await call('SaveDevice',{name:c.host,host:c.host,port:c.port,serverPort:c.serverPort,user:c.user,password:c.password,fingerprint:c.fingerprint,rdpPort:3389});}catch(e){report(e,t('Remember device'));}}
async function authenticate(deploy) {const c=config(deploy);if(deploy&&!confirm(t('Install or start HomeCloud on {0}, instance {1}? This writes the application and starts a background service.',c.host,c.serverPort)))return;state.busy=true;update();notice(deploy?t('Installing and checking service…'):t('Establishing verified SSH connection…'));try{const result=await call(deploy?'Deploy':'Connect',c);if(deploy){$('serverPort').value=c.serverPort;notice(joinText(t('Server response'),': ',redact(result)));}else{applyStatus(result);await rememberDevice(c);notice(result.message?joinText(t('Server response'),': ',redact(result.message)):t('Connected securely. Token stays in the native backend.'));}}finally{$('password').value='';c.password='';state.busy=false;update();}}
for(const id of ['host','port'])$(id).addEventListener('input',resetTrust);
$('trust').addEventListener('change',update);
on('probe',t('Read host key'),async()=>{const c=config(false,true);resetTrust();state.busy=true;update();try{const fp=await call('ProbeHostKey',c);if(!/^SHA256:/.test(fp))throw new UIError(t('Invalid fingerprint returned.'));state.fingerprint=fp;put($('fingerprint'),fp);notice(t('Host key retrieved but NOT trusted. Verify independently before checking the box.'));}finally{state.busy=false;update();}});
form('connect-form','Connect',()=>authenticate(false));form('deploy-form','Deploy',()=>authenticate(true));
on('disconnect','Disconnect',async()=>{await call('Disconnect');rdpStop(true);applyStatus({});$('password').value='';notice(t('Disconnected. Live snapshots may be stale.'));});
async function discover(){ put($('discovery-status'),t('Scanning mDNS…'));try{const list=await call('Discover');$('discovered').replaceChildren();for(const s of list||[]){const b=document.createElement('button');put(b,`${s.name} · ${s.host}:${s.port}`);b.onclick=()=>{if(state.busy)return;$('host').value=s.host;$('serverPort').value=s.port;resetTrust();notice(t('Discovery selected. Confirm SSH port and verify the host key before connecting.'));};$('discovered').append(b);}put($('discovery-status'),list?.length?t('{0} server(s) found. Manual fields below always override discovery.',list.length):t('No advertisements found. Enter a host manually.'));}catch(e){put($('discovery-status'),t('Discovery unavailable. Enter a host manually.'));throw e;}}
on('discover',t('Discovery'),discover);
function bytes(n){n=Number(n);if(!Number.isFinite(n)||n<0)return '—';const units=['B','KiB','MiB','GiB','TiB'];let i=0;while(n>=1024&&i<4){n/=1024;i++;}return `${n.toFixed(i?1:0)} ${units[i]}`;}
function cell(row,value){const td=document.createElement('td');put(td,value??'—');row.append(td);return td;}
function empty(id,count,text){const row=document.createElement('tr');cell(row,text).colSpan=count;$(id).replaceChildren(row);}
async function dashboard(){const [s,ps]=await Promise.all([request('/api/v1/stats'),request('/api/v1/processes')]);const cpu=Number(Array.isArray(s.cpu_percent)?s.cpu_percent[0]:s.cpu_percent);put($('cpu'),Number.isFinite(cpu)?`${cpu.toFixed(1)}%`:'—');$('cpu-bar').value=cpu||0;const m=s.memory||{};put($('memory'),bytes(m.used));put($('memory-total'),bytes(m.total));$('memory-bar').value=Number(m.usedPercent)||0;$('processes').replaceChildren();for(const p of ps||[]){const r=document.createElement('tr');cell(r,p.pid);cell(r,p.name);cell(r,bytes(p.memory_bytes));$('processes').append(r);}if(!ps?.length)empty('processes',3,t('No processes returned.'));put($('updated'),t('Snapshot · {0}',new Date().toLocaleTimeString()));}
on('refresh-dashboard',t('Dashboard'),dashboard);
async function containers(){const list=await request('/api/v1/containers');$('container-list').replaceChildren();for(const c of list||[]){const r=document.createElement('tr');const name=(c.Names||[]).join(', ')||c.Id?.slice(0,12);const title=cell(r,name);const sub=document.createElement('small');put(sub,c.Image);title.append(sub);cell(r,joinText(t(c.State||''),' · ',c.Status||''));const actions=cell(r,'');const group=document.createElement('div');group.className='toolbar';for(const verb of ['start','stop','restart']){const b=document.createElement('button');put(b,t(verb));b.onclick=()=>action(t('Container {0}',t(verb)),async()=>{if(verb!=='start'&&!confirm(t('{0} {1}? Running work may be interrupted.',t(verb),name)))return;await request('/api/v1/containers/control',JSON.stringify({container_id:c.Id,action:verb}));await containers();notice(t('Container {0} completed.',t(verb)));},b);group.append(b);}actions.append(group);$('container-list').append(r);}if(!list?.length)empty('container-list',3,t('No containers on this server.'));}
on('refresh-containers',t('Containers'),containers);
function safePath(value){if(/[\x00-\x1f\\]/.test(value)||value.split('/').includes('..'))throw new UIError(t('Use a relative path without traversal or backslashes.'));return value;}
function iconFor(f){if(f.isDir)return '📁';const ext=(f.name.split('.').pop()||'').toLowerCase();const map={jpg:'🖼️',jpeg:'🖼️',png:'🖼️',gif:'🖼️',bmp:'🖼️',webp:'🖼️',svg:'🖼️',mp4:'🎬',mkv:'🎬',avi:'🎬',mov:'🎬',mp3:'🎵',wav:'🎵',ogg:'🎵',flac:'🎵',aac:'🎵',pdf:'📄',doc:'📝',docx:'📝',txt:'📝',md:'📝',xls:'📊',xlsx:'📊',csv:'📊',ppt:'📊',pptx:'📊',zip:'📦',tar:'📦',gz:'📦',rar:'📦','7z':'📦',exe:'⚙️',msi:'⚙️',sh:'⚙️',py:'🐍',js:'📜',html:'🌐',css:'🎨',json:'📋',xml:'📋',iso:'💿',img:'💿',ascii:'🔧'};return map[ext]||'📄';}
function updateBreadcrumb(){const parts=$('file-path').value.split('/').filter(Boolean);const bc=$('breadcrumb');bc.replaceChildren();const root=document.createElement('button');put(root,'🏠 /');root.onclick=async()=>{$('file-path').value='/';await files();};bc.append(root);let accumulated='';for(const part of parts){const sep=document.createElement('span');sep.className='sep';sep.textContent=' / ';bc.append(sep);accumulated+='/'+part;const b=document.createElement('button');put(b,part);b.onclick=async()=>{$('file-path').value=accumulated;await files();};bc.append(b);}}
function renderFileGrid(list){const grid=$('file-grid');grid.replaceChildren();for(const f of list||[]){const item=document.createElement('div');item.className='file-item'+(f.isDir?' dir':'');item.innerHTML=`<div class="icon">${iconFor(f)}</div><div class="name" title="${f.name}">${f.name}</div><div class="meta">${f.isDir?t('Folder'):bytes(f.size)}</div>`;item.onclick=async()=>{if(f.isDir){$('file-path').value=f.path;await files();}else{requireConnection('DownloadFile');await call('DownloadFile',safePath(f.path));notice(t('Download dialog completed.'));}};item.ondblclick=async()=>{if(f.isDir){$('file-path').value=f.path;await files();}};grid.append(item);}if(!list?.length){const empty=document.createElement('p');empty.className='hint';empty.style.textAlign='center';empty.style.padding='30px';put(empty,t('This folder is empty.'));grid.append(empty);}}
function renderFileTable(list){$('file-list').replaceChildren();for(const f of list||[]){const r=document.createElement('tr');cell(r,`${f.isDir?'▸ ':''}${f.name}`);cell(r,f.isDir?t('Folder'):bytes(f.size));cell(r,f.modified||'—');const td=cell(r,'');const b=document.createElement('button');put(b,f.isDir?t('Open'):t('Download'));b.onclick=()=>action(f.isDir?t('Open'):t('Download'),async()=>{if(f.isDir){$('file-path').value=f.path;await files();}else{requireConnection('DownloadFile');await call('DownloadFile',safePath(f.path));notice(t('Download dialog completed.'));}},b);td.append(b);$('file-list').append(r);}if(!list?.length)empty('file-list',4,t('This folder is empty.'));}
async function files(){requireConnection('ListFiles');let path=safePath($('file-path').value.trim()||'/');let list=await call('ListFiles',path);if(path==='/'){try{const home=await call('HomeDir');if(home&&home!=='/'){const homeList=await call('ListFiles',home);path=home;list=homeList;$('file-path').value=home;}}catch{}}state.currentFiles=list||[];updateBreadcrumb();if(state.fileView==='grid'){$('file-grid').hidden=false;$('file-table-wrap').hidden=true;renderFileGrid(list);}else{$('file-grid').hidden=true;$('file-table-wrap').hidden=false;renderFileTable(list);}}
async function download(path){requireConnection('DownloadFile');const saved=await call('DownloadFile',safePath(path));notice(saved?t('Downloaded to {0}',saved):t('Download dialog completed.'));}
function renderSearchResults(list){state.currentFiles=list||[];updateBreadcrumb();$('file-grid').hidden=false;$('file-table-wrap').hidden=true;const grid=$('file-grid');grid.replaceChildren();for(const f of list||[]){const item=document.createElement('div');item.className='file-item'+(f.isDir?' dir':'');item.innerHTML=`<div class="icon">${iconFor(f)}</div><div class="name" title="${f.path}">${f.name}</div><div class="meta">${f.path}</div>`;item.onclick=async()=>{if(f.isDir){$('file-path').value=f.path;await files();}else{await download(f.path);}};grid.append(item);}if(!list?.length){const empty=document.createElement('p');empty.className='hint';empty.style.textAlign='center';empty.style.padding='30px';put(empty,t('No matches.'));grid.append(empty);}}
on('view-toggle',t('Toggle view'),async()=>{state.fileView=state.fileView==='grid'?'table':'grid';$('view-toggle').querySelector('span').textContent=state.fileView==='grid'?t('Grid view'):t('Table view');await files();});
form('path-form',t('Files'),files);on('parent',t('Parent folder'),async()=>{const parts=$('file-path').value.split('/').filter(Boolean);parts.pop();$('file-path').value='/'+parts.join('/');await files();});on('file-home',t('Home folder'),async()=>{requireConnection('HomeDir');const home=await call('HomeDir');if(home){$('file-path').value=home;}await files();});form('search-form',t('Search'),async()=>{requireConnection('SearchFiles');const query=$('file-search').value.trim();if(!query)throw new UIError(t('Enter a search term.'));const dir=safePath($('file-path').value.trim()||'/');put($('search-status'),t('Searching…'));const list=await call('SearchFiles',dir,query);put($('search-status'),t('Search results: {0}',list?.length||0));renderSearchResults(list);});on('upload',t('Upload'),async()=>{requireConnection('UploadFile');let dir=safePath($('file-path').value.trim()||'/');if(dir==='/'){try{const home=await call('HomeDir');if(home&&home!=='/'){dir=home;$('file-path').value=home;}}catch{}}const choice=await new Promise(resolve=>{const box=document.createElement('div');box.className='modal';const card=document.createElement('div');card.className='modal-card';card.setAttribute('role','dialog');card.setAttribute('aria-modal','true');const lead=document.createElement('p');put(lead,t('Upload files into {0}? Existing files with the same name may be overwritten.',dir));const where=document.createElement('code');where.textContent=dir;const note=document.createElement('p');put(note,t('Files will be saved to {0}.',dir));const bar=document.createElement('div');bar.className='toolbar';const filesBtn=document.createElement('button');filesBtn.type='button';put(filesBtn,t('Upload files'));const folderBtn=document.createElement('button');folderBtn.type='button';put(folderBtn,t('Upload a folder'));const cancelBtn=document.createElement('button');cancelBtn.type='button';put(cancelBtn,t('Cancel'));bar.append(filesBtn,folderBtn,cancelBtn);card.append(lead,where,note,bar);box.append(card);const done=v=>{box.remove();resolve(v);};filesBtn.onclick=()=>done('files');folderBtn.onclick=()=>done('folder');cancelBtn.onclick=()=>done('');box.addEventListener('click',e=>{if(e.target===box)done('');});document.body.append(box);});if(!choice)return;const uploaded=await call(choice==='folder'?'UploadFolder':'UploadFile',dir);$('file-path').value=dir;await files();notice(uploaded?.length?t('Uploaded {0} file(s) to {1}.',uploaded.length,dir):t('Nothing was uploaded.'));});
form('mount-form',t('Drive mapping'),async()=>{const letter=$('drive').value.toUpperCase();if(!/^[D-Z]$/.test(letter))throw new UIError(t('Choose a drive letter D–Z.'));try{put($('mount-status'),joinText(t('Server response'),': ',await call('MountDrive',letter)));}catch(e){put($('mount-status'),joinText(t('Drive mapping'),': ',normalize(e)));queue(t('MountDrive · requires a supported authenticated mapping solution'));throw e;}});
// Bounded mirror of the most recent terminal output; kept for diagnostics and tests only.
let terminalTail='';

// Remote Desktop (mstsc.exe) + saved devices
function rdpRequest(){return {host:$('rdp-host').value.trim(),port:Number($('rdp-port').value),user:$('rdp-user').value.trim(),password:$('rdp-password').value,width:Number($('rdp-width').value),height:Number($('rdp-height').value),fullscreen:$('rdp-fullscreen').checked,remember:$('rdp-remember').checked};}
async function refreshDevices(){try{const list=await call('ListDevices');state.devices=list||[];}catch(e){state.devices=[];}const select=$('device-list');select.replaceChildren();if(!state.devices.length){const opt=document.createElement('option');opt.value='';put(opt,t('No saved devices.'));select.append(opt);return;}for(const d of state.devices){const opt=document.createElement('option');opt.value=d.id;put(opt,`${d.name||d.host} · ${d.host}:${d.rdpPort||3389}`);select.append(opt);}}
on('rdp-check',t('Check port'),async()=>{put($('rdp-status'),t('Checking RDP port…'));const s=await call('RDPCheck',rdpRequest());put($('rdp-status'),s.reachable?t('RDP port reachable.'):t('RDP port unreachable.'));if(!s.reachable)notice(s.message||'',true);});
on('rdp-probe',t('Probe xrdp over SSH'),async()=>{requireConnection('RDPProbe');put($('rdp-status'),t('Probing xrdp over SSH…'));const info=await call('RDPProbe',Number($('rdp-port').value)||3389);put($('rdp-status'),info.status==='ok'?t('xrdp ready.'):t('xrdp unavailable.'));if(info.message)notice(info.message,info.status!=='ok');});
on('rdp-connect',t('Connect via RDP'),async()=>{put($('rdp-status'),t('Launching Remote Desktop…'));const msg=await call('RDPConnect',rdpRequest());$('rdp-password').value='';put($('rdp-status'),msg);notice(t('RDP session launched.'));});
on('rdp-forget',t('Forget credentials'),async()=>{const msg=await call('RDPForget');put($('rdp-status'),msg);notice(t('RDP credentials removed.'));});
function applyDevice(d){$('host').value=d.host||'';$('port').value=d.port||22;$('serverPort').value=d.serverPort||8081;$('user').value=d.user||'';$('password').value=d.password||'';$('rdp-host').value=d.host||'';$('rdp-port').value=d.rdpPort||3389;$('rdp-user').value=d.user||'';$('rdp-password').value=d.password||'';$('device-name').value=d.name||'';if(d.fingerprint){state.fingerprint=d.fingerprint;put($('fingerprint'),d.fingerprint);$('trust').checked=true;}else{resetTrust();}update();}
on('device-load',t('Load'),async()=>{const id=$('device-list').value;if(!id)throw new UIError(t('Select a saved device first.'));applyDevice(await call('GetDevice',id));put($('device-status'),t('Device loaded.'));});
on('device-save',t('Save'),async()=>{const host=$('rdp-host').value.trim();const user=$('rdp-user').value.trim();if(!host||!user)throw new UIError(t('Enter the RDP host and user to save.'));const id=$('device-list').value;const existing=(state.devices||[]).find(d=>d.id===id);const same=existing&&existing.host===host;const saved=await call('SaveDevice',{id,name:$('device-name').value.trim(),host,port:(same&&existing.port)||Number($('port').value)||22,serverPort:(same&&existing.serverPort)||Number($('serverPort').value)||8081,user,password:$('rdp-password').value,rdpPort:Number($('rdp-port').value)||3389,fingerprint:(same&&existing.fingerprint)||''});await refreshDevices();$('device-list').value=saved.id||'';put($('device-status'),t('Device saved.'));});
on('device-delete',t('Delete'),async()=>{const id=$('device-list').value;if(!id)throw new UIError(t('Select a saved device first.'));if(!confirm(t('Delete this saved device?')))return;await call('DeleteDevice',id);await refreshDevices();put($('device-status'),t('Device deleted.'));});
async function rdpTab(){await refreshDevices();await rdpInfo();}
function rdpStop(silent){
  if(state.rdpTimer){clearInterval(state.rdpTimer);state.rdpTimer=null;}
  if(state.rdpAuto){
    state.rdpAuto=false;
    const btn=$('rdp-auto');
    if(btn){const s=btn.querySelector('span');if(s)s.textContent=t('Auto refresh');else btn.textContent=t('Auto refresh');}
    if(!silent)notice(t('Auto refresh stopped'));
  }
}
async function rdpCapture(){
  requireConnection('CaptureScreen');
  put($('rdp-shot-status'),t('Capturing screen…'));
  try{
    const b64=await call('CaptureScreen');
    const img=$('rdp-screen');
    img.src='data:image/png;base64,'+b64;
    img.style.display='block';
    $('rdp-placeholder').style.display='none';
    put($('rdp-shot-status'),joinText(t('Screen captured'),' · ',new Date().toLocaleTimeString()));
  }catch(e){
    put($('rdp-shot-status'),joinText(t('Screen capture failed'),': ',normalize(e)));
    throw e;
  }
}
on('rdp-capture',t('Capture screen'),rdpCapture);
on('rdp-auto',t('Auto refresh'),async()=>{
  state.rdpAuto=!state.rdpAuto;
  if(state.rdpAuto){
    $('rdp-auto').querySelector('span').textContent=t('Stop auto refresh');
    await rdpCapture();
    state.rdpTimer=setInterval(()=>{
      if(state.connected&&!state.busy){
        rdpCapture().catch(()=>{});
      }
    },3000);
    notice(t('Auto refresh started'));
  }else{
    rdpStop();
  }
});
async function rdpInfo(){
  if(!state.connected)return;
  try{
    const info=await call('GetScreenInfo');
    if(info&&info.dimensions){
      put($('rdp-info'),t('Screen dimensions: {0}',info.dimensions));
    }
  }catch(e){
    // Screen info is optional, ignore errors
  }
}
let term, termFit, termReady = false;
function ensureTerminal(){
  if(termReady)return true;
  if(typeof Terminal!=='function'||typeof FitAddon!=='function'||typeof FitAddon.FitAddon!=='function')return false;
  term=new Terminal({cursorBlink:true,scrollback:5000,fontFamily:'Consolas, "Cascadia Mono", ui-monospace, monospace',fontSize:14,theme:{background:'#0c1118',foreground:'#dee6f2',cursor:'#83d9ca',selectionBackground:'#254039'}});
  termFit=new FitAddon.FitAddon();
  term.loadAddon(termFit);
  term.open($('terminal-output'));
  term.onData(data=>{if(!state.terminalOpen)return;call('TerminalInput',data).catch(e=>report(e,t('Terminal input')));});
  term.attachCustomKeyEventHandler(event=>{
    if(event.type!=='keydown'||event.key!=='Enter'||event.ctrlKey||event.altKey||event.metaKey||event.shiftKey)return true;
    const line=(term.buffer?.active?.getLine(term.buffer.active.cursorY+term.buffer.active.baseY)?.translateToString(true)||'').trim();
    if(/\b(rm|rmdir|mkfs|dd|shutdown|reboot|halt|poweroff)\b/.test(line)&&!confirm(t('This command may destroy data or interrupt services. Execute it?')))return false;
    return true;
  });
  termReady=true;
  return true;
}
function terminalOutput(data){const text=String(data??'');terminalTail=(terminalTail+text).slice(-150000);if(!ensureTerminal())return;term.write(text);}
async function resizeTerminal(){if(!state.terminalOpen||!ensureTerminal()||!bridge()?.ResizeTerminal)return;termFit.fit();const cols=Math.max(2,Math.min(500,term.cols)),rows=Math.max(2,Math.min(200,term.rows));await call('ResizeTerminal',cols,rows);}
on('terminal-open',t('Open terminal'),async()=>{requireConnection('OpenTerminal');if(!window.runtime?.EventsOn){queue(t('terminal events'));throw new UIError(t('Native Wails event bridge unavailable. Terminal cannot display output.'));}if(!ensureTerminal())throw new UIError(t('Terminal emulator failed to load.'));term.reset();await call('OpenTerminal');state.terminalOpen=true;notice(t('[SSH session opened]'));await resizeTerminal();term.focus();update();});
on('terminal-close',t('Close terminal'),async()=>{await call('CloseTerminal');state.terminalOpen=false;update();});
on('clear-errors',t('Clear errors'),()=>{state.errors=[];put($('errors'),t('No errors.'));put($('error-count'),'0');});on('clear-queue',t('Clear reminders'),()=>{state.queue=[];$('queue').replaceChildren();});
const tabs=[...document.querySelectorAll('[data-tab]')];
function selectTab(name){state.tab=name;for(const button of tabs){const selected=button.dataset.tab===name;button.setAttribute('aria-selected',String(selected));button.tabIndex=selected?0:-1;$(button.dataset.tab).hidden=!selected;}if(name!=='rdp')rdpStop(true);if(name==='rdp'){action(t('Remote Desktop'),rdpTab);return;}if(state.connected){const load={dashboard,containers,files,terminal:resizeTerminal}[name];if(load)action(t({dashboard:'Dashboard',containers:'Containers',files:'Files',terminal:'Terminal'}[name]),load);}}
for(const button of tabs){button.onclick=()=>selectTab(button.dataset.tab);button.addEventListener('keydown',e=>{let i=tabs.indexOf(button);if(e.key==='ArrowRight')i=(i+1)%tabs.length;else if(e.key==='ArrowLeft')i=(i+tabs.length-1)%tabs.length;else if(e.key==='Home')i=0;else if(e.key==='End')i=tabs.length-1;else return;e.preventDefault();selectTab(tabs[i].dataset.tab);tabs[i].focus();});}
let resizeTimer;window.addEventListener('resize',()=>{clearTimeout(resizeTimer);resizeTimer=setTimeout(()=>action(t('Terminal resize'),resizeTerminal),200);});
window.addEventListener('unhandledrejection',e=>{e.preventDefault();report(e.reason,t('Unexpected failure'));});
let unsubscribers=[];
async function restoreDevice(){await refreshDevices();const id=(state.devices||[])[0]?.id;if(!id)return;applyDevice(await call('GetDevice',id));if($('device-list'))$('device-list').value=id;notice(t('Saved device restored.'));}
async function init(){selectTab('ssh');update();if(!bridge()){put($('runtime-state'),t('Static preview · native bridge unavailable'));put($('discovery-status'),t('Open the native app for mDNS discovery; manual fields are available.'));return;}put($('runtime-state'),t('Native Wails bridge'));if(window.runtime?.EventsOn){unsubscribers.push(window.runtime.EventsOn('terminal:output',terminalOutput),window.runtime.EventsOn('terminal:closed',message=>{state.terminalOpen=false;notice(joinText(t('[Session closed]'),' ',redact(message||'')));update();}),window.runtime.EventsOn('connection:closed',message=>{state.connected=false;state.terminalOpen=false;applyStatus({});report(message,t('Connection closed'));}));}await action(t('Connection status'),status);await action(t('Saved devices'),restoreDevice);await action(t('Discovery'),discover);}
window.addEventListener('beforeunload',()=>{$('password').value='';for(const off of unsubscribers)if(typeof off==='function')off();});
initLanguage();
init();
