package num2word

func Converter(firstNum, secondNum int) (string, string) {
	var num2WordList = []string{"нуль", "один", "два", "три", "четыре", "пять",
		"шесть", "семь", "восемь", "девять", "десять",
		"одиннадцать", "двеннадцать", "триннадцать", "четырнадцать", "пятнадцать",
		"шестнадцать", "семнадцать", "восемнадцать", "двевятнадцать", "двадцать"}
	fisrtNum2Word := num2WordList[firstNum]
	secondNum2Word := num2WordList[secondNum]

	return fisrtNum2Word, secondNum2Word
}
