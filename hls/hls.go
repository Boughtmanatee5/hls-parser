package hls

const (
	// hls tags
	tagPrefix            = "#EXT"
	headerTag            = "#EXTM3U"
	playlistTag          = "#EXT-X-PLAYLIST-TYPE"
	targetDurationTag    = "#EXT-X-TARGETDURATION"
	versionTag           = "#EXT-X-VERSION"
	mediaSequenceTag     = "#EXT-X-MEDIA-SEQUENCE"
	informationTag       = "#EXTINF"
	streamInformationTag = "#EXT-X-STREAM-INF"
	endListTag           = "#EXT-X-ENDLIST"

	//stream attributes
	bandwidth  = "BANDWIDTH"
	resolution = "RESOLUTION"
	codecs     = "CODECS"
)

type Segment struct {
	Duration float64
	URL      string
}

type Variant struct {
	AverageBandwidth int64
	Bandwidth        int64
	FrameRate        float64
	HDCPLevel        int64
	Resolution       int64
	VideoRange       string
	Codecs           []string
}

type Manifest struct {
	PlaylistType   string
	TargetDuration float64
	Version        int64
	MediaSequence  int64
	Segments       []*Segment
}
