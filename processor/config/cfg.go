package config

var (
	ProcessorAddress       = "127.0.0.1:8080"
	StorageAddress         = "127.0.0.1:8081"
	ProcessorHealthAddress = "127.0.0.1:8888"
)

func init() {
	//load file, get envs, load kubernetes config map...
}
