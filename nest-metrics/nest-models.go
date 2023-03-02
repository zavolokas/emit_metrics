package nestmetrics

type NestDeviceResponse struct {
	Devices []NestDevice `json:"devices"`
}

type NestDevice struct {
	Name            string                            `json:"name"`
	Type            string                            `json:"type"`
	Traits          map[string]map[string]interface{} `json:"traits"`
	ParentRelations []NestRelation                    `json:"parentRelations"`
}

type NestRelation struct {
	DisplayName string `json:"displayName"`
}
