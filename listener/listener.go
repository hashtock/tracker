package listener

type Listener interface {
    Listen() chan map[string]int
    Stop()

    SetTags(tags []string)
    Tags() []string
}
