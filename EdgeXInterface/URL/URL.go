package URL

const (
	EdgeXAddr = "http://localhost"
	CoreMetadataPort = "48081"
	CoreCommandPort = "48082"
)

var (
	CoreMetadata string = EdgeXAddr+":"+CoreMetadataPort
	DeviceRegistration string = CoreMetadata+"/api/v1/device"
	CoreCommand string = EdgeXAddr+":"+CoreCommandPort
	SetResource string = CoreCommand+"/api/v1/device/name"
)
