package stringmethods

func Charactercount(s string) int {
    count := 0
    for _, c := range s {
        if c != ' ' {
            count++
        }
    }
    return count
}