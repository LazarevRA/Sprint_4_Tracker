package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength     = 0.65   // длина шага в метрах
	convKilometers = 1000.0 // метров в километре
)

// parsePackage принимает строку с данными формата "количество шагов, продолжительность прогулки"
// функция парсит строку, переводит данные в соответсвтующие типы и возвращает эти значения
// В случае ошибки - возвращает ее.
func parsePackage(data string) (int, time.Duration, error) {

	input := strings.Split(data, ",")

	if len(input) != 2 {
		return 0, time.Duration(0), errors.New("Неверный формат входных данных")
	}

	steps, err := strconv.Atoi(input[0])
	if err != nil {
		return 0, time.Duration(0), errors.New("Ошибка при попытке получить целое число шагов")
	}

	walkDuration, err := time.ParseDuration(input[1])
	if err != nil {
		return 0, time.Duration(0), errors.New("Ошибка при попытке получить продолжительность прогулки")
	}

	return steps, walkDuration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {

	var convKilometers = 1000.0

	steps, activityDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	distance := StepLength * float64(steps) / convKilometers
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, activityDuration)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
