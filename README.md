приложение-чата, работает в двух режимах работы - клиент чата и сервер, обслуживающий чат.</br>
Клиент перед отправкой сообщения добавлять к нему префикс - имя клиента и двоеточие.</br>
Ввод сообщений в клиенте должен осуществляться через prompt, а отправка по переносу строки (нажатие enter). </br>
Режим работы приложения определяться в файле конфигурации. </br>

**Сервер** обрабатывает подключения и для каждого подключения при получении от него сообщения - </br> 
отправляет его всем клиентам (шаблон fan-out).</br>
