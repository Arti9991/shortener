package exitchecker

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// функция TestMyAnalyzer применяет тестируемый анализатор ExitAnalyzer
// к пакетам из папки testdata и проверяет ожидания
func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ExitAnalyzer, "./...")
}
