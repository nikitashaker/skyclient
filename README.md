skyclient v 1.0
Docs (RU)

1. Подключение skytable
Подключение к базе данных осуществляется с помощью функции NewConnection:

func NewConnection(addr, username, password string) (*Client, error)

В качестве параметров передаются адрес подключения, имя пользователя БД, пароль пользователя БД в формате string
Возвращаемые значения: указатель на структуру Client, error - предполагаемая ошибка
Далее полученное значение будет использоваться при вызове других функций из библиотеки skytable

2. Выполнение запросов
В протоколе skytable указано два вида запросов: Simple query и Pipeline

1) Чтобы выполнить одиночный запрос (simple query) используется функция Query:

func (c *Client) Query(query string, params ...interface{}) ([]interface{}, error)

В качестве параметров передаются сам запрос в формате string, параметры любого типа
Возвращаемые значения: слайс значений любого типа (полученные значения), error - предполагаемая ошибка
BlueQL допускает только литералы в качестве параметров, поэтому, параметризировать запрос нужно следующим образом:

Пример использования функции выглядит следующим образом: 
responce, err := client.Query("SELECT * FROM myspace.mymodel WHERE username = ?", "nick")

Также возможно передавать запрос без параметров

2) Чтобы запарсить значения, полученные при запросе, используется функция QueryParse:

func (c *Client) QueryParse(query string, out interface{}, params ...interface{}) error

В качестве параметров передаются сам запрос в формате string, указатель на стукрутуру любого типа, параметры любого типа
Возвращаемые значения: error - предполагаемая ошибка
(!) Важно: запрос, обрабатываемый функцией QueryParse обязательно должен возвращать определенные значения. Зарпос, вернувший null, приведет к ошибке

Пример использования функции выглядит следующим образом: 
type User struct {
    Username string
	Surname string
    Age int
}

var u User
err := client.QueryParse("SELECT username, surname, age FROM myspace.mymodel WHERE username = ?", &u "nick")

3) Чтобы выполнить множественный запрос (pipeline) использутеся функция Pipeline:

func (c *Client) Pipeline(req ...Request) ([][]interface{}, error)

В качестве параметров передаются структуры типа Request, которые содержат Query - запрос типа string, Params = слайс значений любого типа

type Request struct {
	Query  string
	Params []interface{}
}

Возвращаемые значения: слайс слайсов значений любого типа (полученные значения), error - предполагаемая ошибка
Пример использования функции выглядит следующим образом:
responce, err := client.Pipeline(
    skyclient.Request{
        Query:  "SELECT * FROM apps.mymodel1 WHERE username = ?",
        Params: []interface{}{"nick"},
    },
    skyclient.Request{
        Query:  "SELECT * FROM apps.mymodel2 WHERE username = ?",
        Params: []interface{}{"john"},
    },
)

3. Закрытие подключения
Закрытие подключения осуществляется с помощью функции Close:

func (c *Client) Close() error

Возвращаемые значения: error - предполагаемая ошибка

4. (!) Важные уточнения
Есть определенные правила при использовании skyclient, без которых будет некорректная работа протокола

1) Использовать float64 вместо float32 (float32 не подходит для сериализации и получения типа, т.к. сервер не отправляет корректные байт типа данных для float32)
2) Всегда оставлять поле типа boolean в самом конце при создании модели (сервер не отправляет разделитель (\n) после получения значений)
3) Не использовать в запросах ключевое слово ALL. Протокол с ним работает некорректно