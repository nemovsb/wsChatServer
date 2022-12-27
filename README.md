# wsChatServer

## Задание:

Сделать чат сервер. Принимает подключения по websocket, 
Сообщения отправляются в формате json. Первым сообщением 
должен быть хендшейк для подключения к чат каналу. Сообщения 
отправленные клиентом в чат канал должны быть разосланы 
всем клиентам подключенным к этому чат каналу(broadcast). 
Чат каналов может быть несколько и создаются при запросе 
на подключение к каналу первого клиента. 
Когда у канала не остается клиентов то канал должен быть закрыт и удален

### Формат сообщения:

{
    "client_name":"ser",
    "chanel_name":"ch",
    "text":"123"
}

### Команды клиента (JS):
### Подключение:

var ws = new WebSocket("ws://127.0.0.1:8081");
var ws2 = new WebSocket("ws://127.0.0.1:8081");
var ws3 = new WebSocket("ws://127.0.0.1:8081");

### Вывод входящих сообщений:

ws.onmessage = e => console.log(e.data);
ws2.onmessage = e => console.log(e.data);
ws3.onmessage = e => console.log(e.data);

### Сообщения:

ws.send(`{"client_name":"ser","chanel_name":"ch","text":"123"}`); 
ws2.send(`{"client_name":"ser2","chanel_name":"ch","text":"123456789"}`);
ws3.send(`{"client_name":"ser3","chanel_name":"ch","text":"111111111111111111"}`);


ws3.send(`{"client_name":"ser3","chanel_name":"ch22222","text":"999"}`);

### Закрыть соединения:

ws.close();
ws2.close();
ws3.close();