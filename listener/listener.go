package listener

type Listener interface {
    Listen() chan map[string]int
}
