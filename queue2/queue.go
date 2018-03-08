package queue2

const (
	namespace = "STEAM2_"
)

func RunConsumers() {

	//go func() {
	//
	//	for _ := range closeChannel {
	//		fmt.Println("Reconnecting")
	//		err := connect()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//		}
	//	}
	//
	//}()

	go receiveChanges()
	//go changeConsumer()
	//go packageConsumer()
	//go playerConsumer()
}
