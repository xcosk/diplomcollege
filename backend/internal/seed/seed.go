package seed

import (
	"database/sql"
	"encoding/json"
)

type seedCourse struct {
	Level       string
	Title       string
	Description string
	Lessons     []seedLesson
}

type seedLesson struct {
	Title   string
	Content string
	Quiz    []seedQuestion
}

type seedQuestion struct {
	Question     string
	Options      []string
	CorrectIndex int
	Explanation  string
}

func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM courses").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	courses := seedCourses()
	for _, c := range courses {
		res, err := db.Exec("INSERT INTO courses(level, title, description) VALUES(?,?,?)", c.Level, c.Title, c.Description)
		if err != nil {
			return err
		}
		courseID, _ := res.LastInsertId()
		for i, l := range c.Lessons {
			resLesson, err := db.Exec("INSERT INTO lessons(course_id, title, content, order_index) VALUES(?,?,?,?)", courseID, l.Title, l.Content, i+1)
			if err != nil {
				return err
			}
			lessonID, _ := resLesson.LastInsertId()
			for _, q := range l.Quiz {
				options, _ := json.Marshal(q.Options)
				_, err = db.Exec("INSERT INTO lesson_quiz_questions(lesson_id, question, options_json, correct_index, explanation) VALUES(?,?,?,?,?)", lessonID, q.Question, string(options), q.CorrectIndex, q.Explanation)
				if err != nil {
					return err
				}
			}
		}
	}

	for _, q := range seedPlacementQuestions() {
		options, _ := json.Marshal(q.Options)
		_, err := db.Exec("INSERT INTO placement_questions(question, options_json, correct_index) VALUES(?,?,?)", q.Question, string(options), q.CorrectIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedCourses() []seedCourse {
	courses := []seedCourse{
		{
			Level:       "base",
			Title:       "Go: Базовый уровень",
			Description: "Старт с нуля: установка, синтаксис, типы, управляющие конструкции и первые структуры данных.",
			Lessons: []seedLesson{
				{
					Title: "Старт и первая программа",
					Content: `В этом уроке мы установим Go и напишем первую программу. Это фундамент для всех последующих тем.

Шаг 1. Установка и проверка.
- Скачайте Go с официального сайта.
- Проверьте установку командой: go version.

Шаг 2. Первая программа.
- Создайте файл main.go.
- Укажите пакет main.
- Напишите функцию main.

Мини‑практика:
1) Напишите программу, которая печатает ваше имя.
2) Добавьте вывод текущей даты (через time.Now()).

Важно:
- Все начинается с пакета main.
- Точка входа — функция main.
- Для вывода используйте fmt.Println.

Дополнение:
- Go компилируется в один бинарный файл.
- Быстрый запуск: go run main.go.
- Сборка: go build.
`,
					Quiz: []seedQuestion{
						{Question: "Какая функция является точкой входа?", Options: []string{"start()", "main()", "run()", "init()"}, CorrectIndex: 1, Explanation: "Точка входа — функция main() в пакете main."},
						{Question: "Какой пакет нужен для вывода в консоль?", Options: []string{"io", "fmt", "print", "log"}, CorrectIndex: 1, Explanation: "fmt предоставляет функции Print/Println."},
						{Question: "Правильно ли имя файла main.go?", Options: []string{"Да, это стандартное имя", "Нет, нужно app.go", "Нужно program.go", "Имя файла не важно"}, CorrectIndex: 0, Explanation: "Имя файла произвольное, но main.go часто используется."},
					},
				},
				{
					Title: "Переменные и типы",
					Content: `Разберемся с базовыми типами и объявлением переменных.

Объявление:
- var age int = 25
- name := "Go" (короткая запись)

Базовые типы:
- int, int64 для целых чисел
- float64 для дробей
- string для строк
- bool для логики

Рекомендации:
- Используйте короткую запись внутри функций.
- Явно указывайте тип, если это повышает читаемость.

Мини‑практика:
1) Объявите 3 переменные разных типов.
2) Выведите их значения в консоль.
3) Измените одно значение и снова выведите.
`,
					Quiz: []seedQuestion{
						{Question: "Как объявить переменную с выводом типа?", Options: []string{"name := \"Go\"", "var name string", "let name = \"Go\"", "const name = \"Go\""}, CorrectIndex: 0, Explanation: "Короткая запись := выводит тип автоматически."},
						{Question: "Какой тип подходит для дробных чисел?", Options: []string{"int", "float64", "string", "bool"}, CorrectIndex: 1, Explanation: "float64 — стандартный тип для дробных чисел."},
						{Question: "Можно ли использовать := вне функции?", Options: []string{"Да", "Нет"}, CorrectIndex: 1, Explanation: "Короткая запись допустима только внутри функций."},
					},
				},
				{
					Title: "Условия и циклы",
					Content: `Освоим if, for и switch.

if:
- if x > 10 { ... }
- else if / else

for:
- Единственный цикл в Go
- for i := 0; i < 10; i++ {}
- for i < 10 {} (как while)

switch:
- switch по значениям
- не требует break

Мини‑практика:
1) Цикл от 1 до 10.
2) Отметьте четные числа.
3) Используйте switch для вывода "odd/even".
`,
					Quiz: []seedQuestion{
						{Question: "Какая конструкция используется для цикла?", Options: []string{"loop", "while", "for", "repeat"}, CorrectIndex: 2, Explanation: "В Go есть только for."},
						{Question: "Нужны ли скобки вокруг условия if?", Options: []string{"Да", "Нет"}, CorrectIndex: 1, Explanation: "Скобки не обязательны."},
						{Question: "switch в Go требует break?", Options: []string{"Да, всегда", "Нет, не требует"}, CorrectIndex: 1, Explanation: "В Go switch не проваливается по умолчанию."},
					},
				},
				{
					Title: "Функции и пакеты",
					Content: `Функции — главный строительный блок Go.

Синтаксис:
- func add(a int, b int) int { return a + b }
- Можно возвращать несколько значений.

Пакеты:
- Организуют код.
- Импортируются через import "package".

Мини‑практика:
1) Функция Sum(a, b) и Mul(a, b).
2) Вызовите их из main.
3) Верните два значения из одной функции.
`,
					Quiz: []seedQuestion{
						{Question: "Можно ли вернуть два значения из функции?", Options: []string{"Да", "Нет"}, CorrectIndex: 0, Explanation: "Go поддерживает множественные возвращаемые значения."},
						{Question: "Как называется корневой пакет программы?", Options: []string{"root", "main", "app", "core"}, CorrectIndex: 1, Explanation: "Основной пакет — main."},
						{Question: "Какой синтаксис у функции?", Options: []string{"func name()", "function name()", "def name()", "fn name()"}, CorrectIndex: 0, Explanation: "Go использует ключевое слово func."},
					},
				},
				{
					Title: "Срезы, карты и структуры",
					Content: `Основные структуры данных Go.

Срезы (slice):
- Динамические массивы
- append для добавления

Карты (map):
- map[keyType]valueType
- Быстрый доступ по ключу

Структуры (struct):
- Группировка полей
- Создание своих типов

Мини‑практика:
1) Создайте struct User {Name, Age}.
2) Сделайте map пользователей.
3) Добавьте и выведите значения.
`,
					Quiz: []seedQuestion{
						{Question: "Как создать срез?", Options: []string{"var s []int", "var s map", "s := map[]", "s := struct{}"}, CorrectIndex: 0, Explanation: "Срез объявляется как []тип."},
						{Question: "map в Go — это", Options: []string{"список", "массив", "ключ‑значение", "очередь"}, CorrectIndex: 2, Explanation: "map хранит пары ключ‑значение."},
						{Question: "Как объявить struct?", Options: []string{"type User struct { ... }", "struct User { ... }", "class User { ... }", "record User { ... }"}, CorrectIndex: 0, Explanation: "struct объявляется через type."},
					},
				},
				{
					Title: "Указатели и ссылки",
					Content: `Указатели позволяют работать с адресами в памяти.

- &x возвращает адрес переменной
- *p — разыменование
- Используются для изменения значений и экономии памяти

Мини‑практика:
1) Напишите функцию inc(x *int).
2) Передайте указатель на переменную.
3) Проверьте, что значение изменилось.
`,
					Quiz: []seedQuestion{
						{Question: "Что делает оператор &?", Options: []string{"Разыменовывает", "Берет адрес", "Создает копию", "Удаляет ссылку"}, CorrectIndex: 1, Explanation: "& возвращает адрес переменной."},
						{Question: "Что значит *p в Go?", Options: []string{"Адрес p", "Значение по адресу p", "Создание указателя", "Ошибка компиляции"}, CorrectIndex: 1, Explanation: "*p — разыменование."},
						{Question: "Зачем нужны указатели?", Options: []string{"Для UI", "Для ссылочной передачи и экономии", "Для форматирования", "Не нужны"}, CorrectIndex: 1, Explanation: "Указатели передают данные без копирования."},
					},
				},
				{
					Title: "Пакеты, модули и зависимости",
					Content: `Организация проекта и зависимостей.

- Инициализация: go mod init
- Добавление зависимостей: go get
- Разделение по пакетам

Мини‑практика:
1) Создайте пакет utils.
2) Вынесите туда функцию.
3) Подключите в main.
`,
					Quiz: []seedQuestion{
						{Question: "Команда для инициализации модуля?", Options: []string{"go init", "go mod init", "go start", "go module"}, CorrectIndex: 1, Explanation: "Используется go mod init."},
						{Question: "Где хранится список зависимостей?", Options: []string{"main.go", "go.mod", "README", "go.sum"}, CorrectIndex: 1, Explanation: "go.mod содержит зависимости."},
						{Question: "Как подключить пакет из проекта?", Options: []string{"import \"./utils\"", "import \"project/utils\"", "import \"utils\"", "import \"./\""}, CorrectIndex: 1, Explanation: "Импорт по модульному пути."},
					},
				},
			},
		},
		{
			Level:       "mid",
			Title:       "Go: Средний уровень",
			Description: "Интерфейсы, ошибки, конкурентность, тестирование и работа с данными.",
			Lessons: []seedLesson{
				{
					Title: "Структуры, методы, интерфейсы",
					Content: `Углубляемся в типы и поведение.

- Методы привязываются к типам: func (u User) FullName() string
- Интерфейсы описывают поведение
- Любой тип реализует интерфейс автоматически

Мини‑практика:
1) Создайте интерфейс Stringer.
2) Реализуйте его для структуры User.
3) Выведите результат.
`,
					Quiz: []seedQuestion{
						{Question: "Что описывает интерфейс?", Options: []string{"Данные", "Поведение", "Память", "Потоки"}, CorrectIndex: 1, Explanation: "Интерфейс описывает набор методов."},
						{Question: "Нужно ли явно писать implements?", Options: []string{"Да", "Нет"}, CorrectIndex: 1, Explanation: "Реализация интерфейса не требует ключевого слова."},
						{Question: "Где объявляются методы?", Options: []string{"Внутри struct", "Снаружи, через receiver", "В пакете main", "В комментариях"}, CorrectIndex: 1, Explanation: "Метод имеет receiver перед именем функции."},
					},
				},
				{
					Title: "Ошибки и обработка",
					Content: `Ошибки — часть дизайна Go.

- Возвращайте error как последнее значение.
- Проверяйте err сразу после вызова.
- Создавайте ошибки через errors.New или fmt.Errorf.

Мини‑практика:
1) Сделайте функцию деления.
2) Возвращайте ошибку при делении на ноль.
3) Проверьте в main.
`,
					Quiz: []seedQuestion{
						{Question: "Где обычно находится error в возвращаемых значениях?", Options: []string{"Первым", "Последним", "В середине", "Не важно"}, CorrectIndex: 1, Explanation: "По соглашению error — последнее значение."},
						{Question: "Как создать новую ошибку?", Options: []string{"errors.New()", "error()", "fmt.Error()", "make(error)"}, CorrectIndex: 0, Explanation: "Используйте errors.New или fmt.Errorf."},
						{Question: "Что делать с err?", Options: []string{"Игнорировать", "Проверить сразу", "Проверить в конце", "Проверять редко"}, CorrectIndex: 1, Explanation: "Go‑код проверяет err сразу."},
					},
				},
				{
					Title: "Горутины и каналы",
					Content: `Конкурентность — сильная сторона Go.

- go f() запускает горутину
- каналы передают данные между горутинами
- select слушает несколько каналов

Мини‑практика:
1) Запустите 2 горутины.
2) Соберите результаты через канал.
3) Объедините их.
`,
					Quiz: []seedQuestion{
						{Question: "Как запускается горутина?", Options: []string{"thread f()", "go f()", "async f()", "spawn f()"}, CorrectIndex: 1, Explanation: "go f() запускает горутину."},
						{Question: "Что такое канал?", Options: []string{"Файл", "Очередь для данных", "Пакет", "Модуль"}, CorrectIndex: 1, Explanation: "Канал передает данные между горутинами."},
						{Question: "select нужен для", Options: []string{"циклов", "ветвления", "ожидания каналов", "форматирования"}, CorrectIndex: 2, Explanation: "select ожидает готовый канал."},
					},
				},
				{
					Title: "Тестирование",
					Content: `Тесты пишутся в файлах *_test.go.

- Используйте пакет testing.
- Запуск: go test ./...
- Табличные тесты — стандартная практика.

Мини‑практика:
1) Напишите тест для функции суммы.
2) Используйте таблицу входов.
`,
					Quiz: []seedQuestion{
						{Question: "Как называется пакет для тестов?", Options: []string{"test", "testing", "assert", "check"}, CorrectIndex: 1, Explanation: "Стандартный пакет testing."},
						{Question: "Какой суффикс у тестовых файлов?", Options: []string{"_spec.go", "_test.go", ".test.go", "_tests.go"}, CorrectIndex: 1, Explanation: "Файлы называются *_test.go."},
						{Question: "Как запускать тесты?", Options: []string{"go run", "go test", "go build", "go vet"}, CorrectIndex: 1, Explanation: "Запуск тестов: go test."},
					},
				},
				{
					Title: "Файлы и JSON",
					Content: `Практика работы с данными.

- os и io для чтения/записи
- encoding/json для сериализации
- struct tags для формата

Мини‑практика:
1) Сохраните структуру в JSON файл.
2) Прочитайте обратно.
`,
					Quiz: []seedQuestion{
						{Question: "Какой пакет кодирует JSON?", Options: []string{"encoding/json", "json", "fmt", "encoding"}, CorrectIndex: 0, Explanation: "encoding/json — стандартный пакет."},
						{Question: "Теги полей в JSON пишутся", Options: []string{"в комментариях", "в struct tag", "в отдельном файле", "не нужны"}, CorrectIndex: 1, Explanation: "Теги указываются в struct tag."},
						{Question: "Что делает json.Marshal?", Options: []string{"Декодирует", "Кодирует в JSON", "Пишет в файл", "Открывает канал"}, CorrectIndex: 1, Explanation: "Marshal кодирует в JSON."},
					},
				},
				{
					Title: "HTTP клиент и сервер",
					Content: `Создаем базовые HTTP сервисы.

- net/http для клиентских и серверных запросов
- Handler обрабатывает запросы
- Статус‑коды и JSON ответы

Мини‑практика:
1) Поднимите сервер /ping.
2) Верните {"ok":true}.
3) Проверьте через браузер.
`,
					Quiz: []seedQuestion{
						{Question: "Какой пакет для HTTP в Go?", Options: []string{"net/http", "http", "net/client", "server"}, CorrectIndex: 0, Explanation: "Используется net/http."},
						{Question: "Какой статус‑код для успешного GET?", Options: []string{"200", "201", "204", "400"}, CorrectIndex: 0, Explanation: "200 OK."},
						{Question: "Что такое handler?", Options: []string{"Файл", "Функция обработки запроса", "База данных", "Логгер"}, CorrectIndex: 1, Explanation: "Handler обрабатывает запрос."},
					},
				},
				{
					Title: "Работа с контекстом",
					Content: `Контекст управляет временем жизни операций.

- context.Background() базовый контекст
- context.WithTimeout для дедлайна
- ctx.Done() для отмены

Мини‑практика:
1) Добавьте таймаут 1 сек.
2) Отмените запрос по истечению времени.
`,
					Quiz: []seedQuestion{
						{Question: "Где обычно начинается context?", Options: []string{"context.Start()", "context.Background()", "context.New()", "context.Base()"}, CorrectIndex: 1, Explanation: "Context начинается с Background."},
						{Question: "Зачем нужен ctx.Done()?", Options: []string{"Для логов", "Для отмены", "Для JSON", "Для UI"}, CorrectIndex: 1, Explanation: "ctx.Done() сигнализирует отмену."},
						{Question: "Что делает WithTimeout?", Options: []string{"Останавливает GC", "Задает дедлайн", "Создает горутину", "Ломает тесты"}, CorrectIndex: 1, Explanation: "WithTimeout задает дедлайн."},
					},
				},
			},
		},
		{
			Level:       "pro",
			Title:       "Go: Профессиональный уровень",
			Description: "Архитектура, контекст, производительность, базы данных и production‑подходы.",
			Lessons: []seedLesson{
				{
					Title: "Архитектура и модули",
					Content: `Профессиональный подход начинается с структуры проекта.

- Пакеты по доменным зонам
- Разделение логики и инфраструктуры
- Go modules для зависимостей

Мини‑практика:
1) Спроектируйте структуру api/domain/storage.
2) Опишите ответственность каждого пакета.
`,
					Quiz: []seedQuestion{
						{Question: "Что такое Go module?", Options: []string{"Папка проекта", "Система управления зависимостями", "Пакет тестов", "Логгер"}, CorrectIndex: 1, Explanation: "Go modules управляют зависимостями."},
						{Question: "Зачем разделять домен и инфраструктуру?", Options: []string{"Для красоты", "Для тестируемости и гибкости", "Для ускорения", "Не нужно"}, CorrectIndex: 1, Explanation: "Разделение улучшает поддержку и тестирование."},
						{Question: "Где хранить бизнес‑логику?", Options: []string{"В handlers", "В domain/usecase", "В cmd", "В tests"}, CorrectIndex: 1, Explanation: "Бизнес‑логика — в доменном слое."},
					},
				},
				{
					Title: "Контекст, таймауты и отмена",
					Content: `Контекст управляет временем жизни операций.

- context.Context передается по цепочке
- WithTimeout / WithCancel
- Уважайте ctx.Done()

Мини‑практика:
1) Оберните запрос БД в WithTimeout.
2) Завершите по дедлайну.
`,
					Quiz: []seedQuestion{
						{Question: "Что делает context.WithTimeout?", Options: []string{"Логирует", "Создает контекст с дедлайном", "Запускает горутину", "Открывает файл"}, CorrectIndex: 1, Explanation: "WithTimeout задает дедлайн."},
						{Question: "Как проверять отмену?", Options: []string{"ctx.Done()", "ctx.Stop()", "ctx.Cancel()", "ctx.Wait()"}, CorrectIndex: 0, Explanation: "Отмена через ctx.Done()."},
						{Question: "Контекст нужно", Options: []string{"создавать глобально", "передавать в функции", "хранить в переменной", "игнорировать"}, CorrectIndex: 1, Explanation: "Context передается по цепочке вызовов."},
					},
				},
				{
					Title: "Производительность и профилирование",
					Content: `Оптимизация начинается с измерений.

- pprof для профилей CPU/памяти
- Избегайте лишних аллокаций
- Проверяйте горячие пути

Мини‑практика:
1) Подключите pprof.
2) Снимите профиль CPU.
`,
					Quiz: []seedQuestion{
						{Question: "Что такое pprof?", Options: []string{"Линтер", "Профайлер", "Фреймворк", "Тест‑раннер"}, CorrectIndex: 1, Explanation: "pprof — профайлер для Go."},
						{Question: "С чего начинается оптимизация?", Options: []string{"С догадок", "С измерений", "С удаления кода", "С переписывания"}, CorrectIndex: 1, Explanation: "Сначала измерения."},
						{Question: "Что ухудшает производительность?", Options: []string{"Лишние аллокации", "Тесты", "Интерфейсы", "Ошибки"}, CorrectIndex: 0, Explanation: "Лишние аллокации повышают нагрузку на GC."},
					},
				},
				{
					Title: "Работа с базой данных",
					Content: `Профессиональная работа с БД включает транзакции.

- database/sql — стандарт
- Всегда закрывайте rows
- Используйте транзакции

Мини‑практика:
1) Выполните транзакцию с двумя запросами.
`,
					Quiz: []seedQuestion{
						{Question: "Какой пакет стандартный для БД?", Options: []string{"sql", "database/sql", "db", "storage"}, CorrectIndex: 1, Explanation: "database/sql — стандартный пакет."},
						{Question: "Зачем нужны транзакции?", Options: []string{"Для скорости", "Для атомарности", "Для логов", "Для UI"}, CorrectIndex: 1, Explanation: "Транзакции обеспечивают целостность."},
						{Question: "Что нужно делать с rows?", Options: []string{"Игнорировать", "Закрывать", "Логировать", "Сериализовать"}, CorrectIndex: 1, Explanation: "rows нужно закрывать."},
					},
				},
				{
					Title: "REST API и middleware",
					Content: `В production важны стандарты и слои.

- Разделяйте handlers и сервисы
- Middleware для логирования, авторизации, CORS
- Следите за статус‑кодами

Мини‑практика:
1) Добавьте middleware для API‑ключа.
`,
					Quiz: []seedQuestion{
						{Question: "Что делает middleware?", Options: []string{"Обрабатывает запросы вокруг handler", "Пишет в БД", "Рисует UI", "Собирает билды"}, CorrectIndex: 0, Explanation: "Middleware оборачивает обработчики."},
						{Question: "Какой статус код для успешного POST?", Options: []string{"200", "201", "204", "302"}, CorrectIndex: 1, Explanation: "201 Created используется после создания."},
						{Question: "CORS нужен для", Options: []string{"Безопасности файлов", "Доступа с другого домена", "Шифрования", "Логирования"}, CorrectIndex: 1, Explanation: "CORS разрешает запросы с другого домена."},
					},
				},
				{
					Title: "Кэширование и балансировка",
					Content: `Ускоряем сервисы и повышаем устойчивость.

- Кэш в памяти (map + mutex)
- TTL и инвалидация
- Балансировка нагрузки

Мини‑практика:
1) Реализуйте in‑memory cache с TTL.
`,
					Quiz: []seedQuestion{
						{Question: "Зачем нужен кэш?", Options: []string{"Для красоты", "Для ускорения", "Для логов", "Не нужен"}, CorrectIndex: 1, Explanation: "Кэш уменьшает время ответа."},
						{Question: "Что такое TTL?", Options: []string{"Время жизни записи", "Тип данных", "Шифрование", "Метод HTTP"}, CorrectIndex: 0, Explanation: "TTL — время жизни записи."},
						{Question: "Для чего балансировка?", Options: []string{"Для стилей", "Для распределения нагрузки", "Для JSON", "Для тестов"}, CorrectIndex: 1, Explanation: "Балансировка распределяет нагрузку."},
					},
				},
				{
					Title: "Наблюдаемость: логи и метрики",
					Content: `Production требует наблюдаемости.

- Структурированные логи
- Метрики Prometheus
- Трассировка запросов

Мини‑практика:
1) Добавьте логгер.
2) Запишите метрику счетчика.
`,
					Quiz: []seedQuestion{
						{Question: "Что такое метрика?", Options: []string{"Файл", "Числовой показатель", "Функция", "Тест"}, CorrectIndex: 1, Explanation: "Метрики — числовые показатели."},
						{Question: "Зачем нужна трассировка?", Options: []string{"Для UI", "Для отслеживания цепочки запросов", "Для JSON", "Не нужна"}, CorrectIndex: 1, Explanation: "Трассировка показывает путь запроса."},
						{Question: "Почему важны логи?", Options: []string{"Для отладки и мониторинга", "Для ускорения", "Для дизайна", "Не важны"}, CorrectIndex: 0, Explanation: "Логи помогают найти проблемы."},
					},
				},
			},
		},
	}
	for i := range courses {
		for j := range courses[i].Lessons {
			courses[i].Lessons[j].Content = enrichContent(courses[i].Lessons[j].Title, courses[i].Lessons[j].Content)
		}
	}
	return courses
}

func enrichContent(title, content string) string {
	return content + "\n\n### Разбор по шагам\n" +
		"1) Прочитайте теорию выше.\n" +
		"2) Повторите примеры кода.\n" +
		"3) Выполните практику и проверьте себя.\n" +
		"\n### Пример кода\n```go\n" + snippetByTitle(title) + "\n```\n" +
		"\n### Дополнительная практика\n" +
		"- Сделайте небольшое упражнение по теме.\n" +
		"- Расширьте пример собственными данными.\n" +
		"- Проверьте крайние случаи.\n" +
		"\n### Мини‑проект\n" +
		miniProjectByTitle(title) + "\n" +
		"\n### Словарь терминов\n" +
		glossaryByTitle(title) + "\n"
}

func snippetByTitle(title string) string {
	switch title {
	case "Старт и первая программа":
		return "package main\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nfunc main() {\n\tfmt.Println(\"Привет, Go!\")\n\tfmt.Println(time.Now().Format(\"2006-01-02\"))\n}"
	case "Переменные и типы":
		return "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tname := \"Go\"\n\tage := 3\n\trating := 4.9\n\tok := true\n\tfmt.Println(name, age, rating, ok)\n}"
	case "Условия и циклы":
		return "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfor i := 1; i <= 10; i++ {\n\t\tif i%2 == 0 {\n\t\t\tfmt.Println(i, \"even\")\n\t\t} else {\n\t\t\tfmt.Println(i, \"odd\")\n\t\t}\n\t}\n}"
	case "Функции и пакеты":
		return "package main\n\nimport \"fmt\"\n\nfunc sum(a, b int) int { return a + b }\nfunc mul(a, b int) int { return a * b }\n\nfunc main() {\n\tfmt.Println(sum(2,3), mul(2,3))\n}"
	case "Срезы, карты и структуры":
		return "package main\n\nimport \"fmt\"\n\ntype User struct {\n\tName string\n\tAge  int\n}\n\nfunc main() {\n\tusers := []User{{\"Ann\", 20}, {\"Bob\", 25}}\n\tindex := map[string]User{\"ann\": users[0]}\n\tfmt.Println(index[\"ann\"])\n}"
	case "Указатели и ссылки":
		return "package main\n\nimport \"fmt\"\n\nfunc inc(x *int) {\n\t*x = *x + 1\n}\n\nfunc main() {\n\tvalue := 10\n\tinc(&value)\n\tfmt.Println(value)\n}"
	case "Пакеты, модули и зависимости":
		return "package main\n\nimport (\n\t\"fmt\"\n)\n\nfunc main() {\n\tfmt.Println(\"go mod init example\")\n}"
	case "Структуры, методы, интерфейсы":
		return "package main\n\nimport \"fmt\"\n\ntype User struct {\n\tName string\n}\n\nfunc (u User) String() string {\n\treturn \"User: \" + u.Name\n}\n\nfunc main() {\n\tu := User{Name: \"Anna\"}\n\tfmt.Println(u.String())\n}"
	case "Ошибки и обработка":
		return "package main\n\nimport (\n\t\"errors\"\n\t\"fmt\"\n)\n\nfunc div(a, b int) (int, error) {\n\tif b == 0 {\n\t\treturn 0, errors.New(\"division by zero\")\n\t}\n\treturn a / b, nil\n}\n\nfunc main() {\n\tres, err := div(10, 2)\n\tfmt.Println(res, err)\n}"
	case "Горутины и каналы":
		return "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tch := make(chan string)\n\tgo func() { ch <- \"hello\" }()\n\tfmt.Println(<-ch)\n}"
	case "Тестирование":
		return "package main\n\nfunc Sum(a, b int) int { return a + b }\n\n// go test ./... (в отдельном файле *_test.go)"
	case "Файлы и JSON":
		return "package main\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n)\n\ntype User struct {\n\tName string `json:\"name\"`\n}\n\nfunc main() {\n\tu := User{Name: \"Go\"}\n\tb, _ := json.Marshal(u)\n\tfmt.Println(string(b))\n}"
	case "HTTP клиент и сервер":
		return "package main\n\nimport (\n\t\"fmt\"\n\t\"net/http\"\n)\n\nfunc main() {\n\thttp.HandleFunc(\"/ping\", func(w http.ResponseWriter, r *http.Request) {\n\t\tfmt.Fprint(w, \"ok\")\n\t})\n\thttp.ListenAndServe(\":8081\", nil)\n}"
	case "Работа с контекстом":
		return "package main\n\nimport (\n\t\"context\"\n\t\"time\"\n)\n\nfunc main() {\n\tctx, cancel := context.WithTimeout(context.Background(), time.Second)\n\tdefer cancel()\n\t<-ctx.Done()\n}"
	case "Архитектура и модули":
		return "// project/\n//   cmd/\n//   internal/\n//     api/\n//     domain/\n//     storage/"
	case "Контекст, таймауты и отмена":
		return "package main\n\nimport (\n\t\"context\"\n\t\"time\"\n)\n\nfunc main() {\n\tctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)\n\tdefer cancel()\n\t<-ctx.Done()\n}"
	case "Производительность и профилирование":
		return "// go tool pprof -http=:8080 cpu.out\n// import _ \"net/http/pprof\""
	case "Работа с базой данных":
		return "// tx, _ := db.Begin()\n// tx.Exec(\"INSERT ...\")\n// tx.Commit()"
	case "REST API и middleware":
		return "func middleware(next http.Handler) http.Handler {\n\treturn http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n\t\t// before\n\t\tnext.ServeHTTP(w, r)\n\t})\n}"
	case "Кэширование и балансировка":
		return "type Cache struct {\n\tdata map[string]string\n}\n\nfunc (c *Cache) Get(k string) (string, bool) {\n\tv, ok := c.data[k]\n\treturn v, ok\n}"
	case "Наблюдаемость: логи и метрики":
		return "// log.Println(\"request\", time.Now())\n// counter.Inc()"
	default:
		return "package main\n\nfunc main() {}"
	}
}

func miniProjectByTitle(title string) string {
	switch title {
	case "Старт и первая программа":
		return "Соберите \"Hello Go\" приложение: выводите приветствие, дату и аргументы командной строки."
	case "Переменные и типы":
		return "Сделайте калькулятор бюджета: доход, расход, остаток. Выведите итог."
	case "Условия и циклы":
		return "Напишите программу, которая выводит таблицу умножения 1–10."
	case "Функции и пакеты":
		return "Разделите калькулятор на пакет mathutils и основной пакет."
	case "Срезы, карты и структуры":
		return "Сделайте каталог пользователей с поиском по имени."
	case "Указатели и ссылки":
		return "Создайте структуру Counter и изменяйте её через указатели."
	case "Пакеты, модули и зависимости":
		return "Соберите простой CLI‑проект с несколькими пакетами."
	case "Структуры, методы, интерфейсы":
		return "Сделайте интерфейс Notifier и реализации Email/SMS."
	case "Ошибки и обработка":
		return "Сделайте функцию парсинга строки в число с валидацией."
	case "Горутины и каналы":
		return "Запустите 3 горутины и соберите результаты в одном канале."
	case "Тестирование":
		return "Покройте тестами небольшую библиотеку вычислений."
	case "Файлы и JSON":
		return "Сохраните список пользователей в JSON и загрузите обратно."
	case "HTTP клиент и сервер":
		return "Сделайте сервер /status и клиент, который проверяет его."
	case "Работа с контекстом":
		return "Оберните HTTP запрос в контекст с таймаутом."
	case "Архитектура и модули":
		return "Спроектируйте структуру проекта из 3 слоев: api, domain, storage."
	case "Контекст, таймауты и отмена":
		return "Добавьте таймауты ко всем запросам БД."
	case "Производительность и профилирование":
		return "Снимите профиль CPU и найдите самую тяжелую функцию."
	case "Работа с базой данных":
		return "Сделайте CRUD для таблицы users с транзакцией."
	case "REST API и middleware":
		return "Создайте middleware логирования запросов."
	case "Кэширование и балансировка":
		return "Добавьте кэш к медленной функции."
	case "Наблюдаемость: логи и метрики":
		return "Подключите логгер и счетчик запросов."
	default:
		return "Соберите небольшой пример, который объединяет ключевые идеи урока."
	}
}

func glossaryByTitle(title string) string {
	switch title {
	case "Старт и первая программа":
		return "- `main` — точка входа программы.\n- `package` — пространство имен.\n- `fmt` — пакет для форматированного вывода."
	case "Переменные и типы":
		return "- `var` — явное объявление переменной.\n- `:=` — короткое объявление.\n- `тип` — описание данных."
	case "Условия и циклы":
		return "- `if` — ветвление.\n- `for` — цикл.\n- `switch` — множественный выбор."
	case "Функции и пакеты":
		return "- `func` — объявление функции.\n- `return` — возврат значения.\n- `import` — подключение пакета."
	case "Срезы, карты и структуры":
		return "- `slice` — динамический массив.\n- `map` — ключ‑значение.\n- `struct` — пользовательский тип."
	case "Указатели и ссылки":
		return "- `&` — адрес переменной.\n- `*` — разыменование.\n- `nil` — отсутствие ссылки."
	case "Пакеты, модули и зависимости":
		return "- `go.mod` — описание модуля.\n- `go.sum` — контрольные суммы.\n- `module` — имя проекта."
	case "Структуры, методы, интерфейсы":
		return "- `method` — функция с receiver.\n- `interface` — контракт поведения.\n- `receiver` — объект метода."
	case "Ошибки и обработка":
		return "- `error` — интерфейс ошибки.\n- `nil` — нет ошибки.\n- `fmt.Errorf` — форматирование ошибки."
	case "Горутины и каналы":
		return "- `goroutine` — легковесный поток.\n- `channel` — передача данных.\n- `select` — ожидание каналов."
	case "Тестирование":
		return "- `testing` — пакет тестов.\n- `*_test.go` — файл тестов.\n- `табличные тесты` — тест‑кейсы."
	case "Файлы и JSON":
		return "- `json.Marshal` — кодирование.\n- `json.Unmarshal` — декодирование.\n- `tag` — метаданные поля."
	case "HTTP клиент и сервер":
		return "- `handler` — обработчик запроса.\n- `status code` — код ответа.\n- `endpoint` — URL сервиса."
	case "Работа с контекстом":
		return "- `context` — управление жизненным циклом.\n- `timeout` — ограничение времени.\n- `cancel` — отмена."
	case "Архитектура и модули":
		return "- `domain` — бизнес‑логика.\n- `api` — слой входа.\n- `storage` — слой данных."
	case "Контекст, таймауты и отмена":
		return "- `deadline` — крайний срок.\n- `cancel` — отмена.\n- `propagation` — передача контекста."
	case "Производительность и профилирование":
		return "- `pprof` — профайлер.\n- `allocation` — выделение памяти.\n- `hot path` — горячий участок."
	case "Работа с базой данных":
		return "- `transaction` — атомарность.\n- `rows` — результат запроса.\n- `scan` — чтение данных."
	case "REST API и middleware":
		return "- `middleware` — обертка запроса.\n- `CORS` — доступ доменов.\n- `status` — код ответа."
	case "Кэширование и балансировка":
		return "- `cache` — временное хранение.\n- `TTL` — срок жизни.\n- `load balancing` — балансировка."
	case "Наблюдаемость: логи и метрики":
		return "- `metrics` — показатели.\n- `logging` — журналы.\n- `tracing` — трассировка."
	default:
		return "- термин — краткое определение."
	}
}

func seedPlacementQuestions() []seedQuestion {
	return []seedQuestion{
		{Question: "Что такое пакет main?", Options: []string{"Пакет для тестов", "Пакет с точкой входа", "Пакет для JSON", "Пакет для логов"}, CorrectIndex: 1},
		{Question: "Как объявить переменную с выводом типа?", Options: []string{"let x = 1", "x := 1", "var x", "int x"}, CorrectIndex: 1},
		{Question: "Что такое интерфейс?", Options: []string{"Тип данных", "Набор методов", "Файл", "Пакет"}, CorrectIndex: 1},
		{Question: "Как запустить горутину?", Options: []string{"async f()", "go f()", "thread f()", "spawn f()"}, CorrectIndex: 1},
		{Question: "Что делает json.Marshal?", Options: []string{"Читает файл", "Кодирует в JSON", "Создает канал", "Считает байты"}, CorrectIndex: 1},
		{Question: "Где находится error в возвращаемых значениях?", Options: []string{"Первым", "Последним", "В середине", "Не важно"}, CorrectIndex: 1},
		{Question: "Что делает context.WithTimeout?", Options: []string{"Ставит дедлайн", "Логирует", "Создает файл", "Запускает тесты"}, CorrectIndex: 0},
		{Question: "Зачем нужны транзакции?", Options: []string{"Для красоты", "Для атомарности", "Для скорости", "Для UI"}, CorrectIndex: 1},
		{Question: "Какой командой запускаются тесты?", Options: []string{"go run", "go test", "go build", "go vet"}, CorrectIndex: 1},
		{Question: "Что такое pprof?", Options: []string{"Профайлер", "Линтер", "Фреймворк", "Библиотека UI"}, CorrectIndex: 0},
	}
}
