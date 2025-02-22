package spentcalories

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	MInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// parsePackage принимает строку с данными формата
// "количество шагов, вид активности, продолжительность прогулки"
// функция парсит строку, переводит данные в соответсвтующие типы и возвращает эти значения
// В случае ошибки - возвращает ее.
func parseTraining(data string) (int, string, time.Duration, error) {

	input := strings.Split(data, ",")
	if len(input) != 3 {
		return 0, "", 0, errors.New("Неверный формат вхоных данных")
	}

	steps, err := strconv.Atoi(input[0])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка при попытке получить целое число шагов")
	}

	activity := input[1]

	activityDuration, err := time.ParseDuration(input[2])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка при попытке получить продолжительность активности")
	}

	return steps, activity, activityDuration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / MInKm

}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	return distance(steps) / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, activity, activityDuration, err := parseTraining(data)

	if err != nil {
		return "Ошибка при попытке получить информацию о тренировке"
	}

	calories := 0.0
	switch activity {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, activityDuration)
	case "Бег":
		calories = RunningSpentCalories(steps, weight, activityDuration)
	default:
		return "Неизвестный тип тренировки"
	}

	return fmt.Sprintf("Тип тренировки: %s\n Длительность: %.2f ч.\n Дистанция: %.2f км.\n Скорость: %.2f км/ч\n Сожгли калорий: %.2f\n",
		activity, activityDuration.Hours(), distance(steps), meanSpeed(steps, activityDuration), calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed(steps, duration)) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	if height <= 0 {
		return 0
	}
	return (((walkingCaloriesWeightMultiplier * weight) + (math.Pow(meanSpeed(steps, duration), 2.0)/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH)

}
