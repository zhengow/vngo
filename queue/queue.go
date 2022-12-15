package queue

type CommonChannel interface {
    send(msg interface{})
    recv()
}
