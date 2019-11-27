package web

import "fmt"

func getColor(n float64) string {
	switch {
	case n < 5:
		return RED
	case 5 <= n && n < 8:
		return YELLOW
	case 8 <= n:
		return GREEN
	}
	return ""
}

func percent(n float64) string {
	return fmt.Sprintf("%.2f", n)
}
