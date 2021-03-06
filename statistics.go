package main

import (
    "fmt"
    "log"
    "net/http"
    "sort"
    "strconv"
    "strings"
    "math"
)

const (
    pageTop    = `<!DOCTYPE HTML><html><head>
<style>.error{color:#FF0000;}</style></head><title>Statistics</title>
<body><h3>Statistics</h3>
<p>Computes basic statistics for a given list of numbers</p>`
    form       = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />
<input type="text" name="numbers" size="30"><br />
<input type="submit" value="Calculate">
</form>`
    pageBottom = `</body></html>`
    anError    = `<p class="error">%s</p>`
)

type statistics struct {
    numbers []float64
    mean    float64
    median  float64
    mode    float64
    std  float64
    variance    float64
}

func main() {
    http.HandleFunc("/", homePage)
    fmt.Printf("Starting server.\n Listening on :9001")
    if err := http.ListenAndServe(":9001", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func homePage(writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm() // Must be called before writing response
    fmt.Fprint(writer, pageTop, form)
    if err != nil {
        fmt.Fprintf(writer, anError, err)
    } else {
        if numbers, message, ok := processRequest(request); ok {
            stats := getStats(numbers)
            fmt.Fprint(writer, formatStats(stats))
        } else if message != "" {
            fmt.Fprintf(writer, anError, message)
        }
    }
    fmt.Fprint(writer, pageBottom)
}

func processRequest(request *http.Request) ([]float64, string, bool) {
    var numbers []float64
    if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
        text := strings.Replace(slice[0], ",", " ", -1)
        for _, field := range strings.Fields(text) {
            if x, err := strconv.ParseFloat(field, 64); err != nil {
                return numbers, "'" + field + "' is invalid", false
            } else {
                numbers = append(numbers, x)
            }
        }
    }
    if len(numbers) == 0 {
        return numbers, "", false // no data first time form is shown
    }
    return numbers, "", true
}

func formatStats(stats statistics) string {
    return fmt.Sprintf(`<table border="1">
<tr><th colspan="2">Results</th></tr>
<tr><td>Numbers</td><td>%v</td></tr>
<tr><td>Count</td><td>%d</td></tr>
<tr><td>Arithmetic Mean</td><td>%f</td></tr>
<tr><td>Median</td><td>%f</td></tr>
<tr><td>Mode</td><td>%f</td></tr>
<tr><td>Variance</td><td>%f</td></tr>
<tr><td>Std Deviation</td><td>%f</td></tr>
</table>`, stats.numbers, len(stats.numbers), stats.mean, stats.median,
         stats.mode, stats.variance, stats.std)
}

func getStats(numbers []float64) (stats statistics) {
    stats.numbers = numbers
    sort.Float64s(stats.numbers)
    stats.mean = mean(numbers)
    stats.median = median(numbers)
    stats.mode = mode(numbers)
    stats.variance = variance(numbers)
    stats.std = std(numbers)
    return stats
}

func sum(numbers []float64) (total float64) {
    for _, x := range numbers {
        total += x
    }
    return total
}

func mean(numbers []float64) float64 {
    return sum(numbers) / float64(len(numbers))
}

func median(numbers []float64) float64 {
    middle := len(numbers) / 2
    result := numbers[middle]
    if len(numbers)%2 == 0 {
        result = (result + numbers[middle-1]) / 2
    }
    return result
}

func mode(numbers []float64) float64 {
    counter := make(map[float64]int)
    for _, i := range numbers {
        counter[i] += 1
    }
    max := 0.0
    max_count := 0
    for key, value := range counter {
        if value > max_count {
            max = key
            max_count = value
        }
    }
    return max
 }

func variance(numbers []float64) float64 {
    mean_val := mean(numbers)
    cov := 0.0
    for _, i := range numbers {
        cov += math.Pow((i - mean_val), 2)
    }    
    return cov / (float64(len(numbers)) - 1.0)
 }

 func std(numbers []float64) float64 {
    return math.Sqrt(variance(numbers))
 }


