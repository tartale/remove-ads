package rmads

type TivoClipMetadata struct {
	ClipMetadata []struct {
		Segment []struct {
			StartOffset string `json:"startOffset"`
			EndOffset   string `json:"endOffset"`
		} `json:"segment"`
	} `json:"clipMetadata"`
}
