package hls

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) (*Manifest, error) {
	manifest := &Manifest{}
	scanner := bufio.NewScanner(r)

	var segments []*Segment
	var segmentDuration float64
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if i == 0 {
			if line != headerTag {
				return nil, fmt.Errorf("invalid hls manifest missing %s header", headerTag)
			}
			continue
		}

		// there's a smarter way to figure this out
		if strings.Contains(line, tagPrefix) {
			keyValuePair := strings.Split(line, ":")
			key := keyValuePair[0]
			switch key {
			case playlistTag:
				value := keyValuePair[1]
				manifest.PlaylistType = value
			case targetDurationTag:
				value := keyValuePair[1]
				duration, parseDurationErr := strconv.ParseFloat(value, 64)
				if parseDurationErr != nil {
					return nil, fmt.Errorf("error parsing target duration %w", parseDurationErr)
				}
				manifest.TargetDuration = duration
			case versionTag:
				value := keyValuePair[1]
				version, parseVersionErr := strconv.ParseInt(value, 10, 64)
				if parseVersionErr != nil {
					return nil, fmt.Errorf("error parsing version %w", parseVersionErr)
				}
				manifest.Version = version
			case mediaSequenceTag:
				value := keyValuePair[1]
				mediaSequence, parseMediaSequenceErr := strconv.ParseInt(value, 10, 64)
				if parseMediaSequenceErr != nil {
					return nil, fmt.Errorf("eror parsing media sequence %w", parseMediaSequenceErr)
				}
				manifest.MediaSequence = mediaSequence
			case streamInformationTag:
				infoString := keyValuePair[1]
				attributes := strings.Split(infoString, ",")
				variant := &Variant{}
				for _, attribute := range attributes {
					attributeKeyValue := strings.Split(attribute, "=")
					if len(attributeKeyValue) > 2 {
						continue
					}
					switch attributeKeyValue[0] {
					case bandwidth:
						variantBandwidth, parseErr := strconv.ParseInt(attributeKeyValue[1], 10, 64)
						if parseErr != nil {
							return nil, fmt.Errorf("error parsing hls variant bandwidth")
						}
						variant.Bandwidth = variantBandwidth
					}
				}

			case informationTag:
				// strip the comma before parsing
				value := strings.ReplaceAll(keyValuePair[1], ",", "")
				duration, parseDurationErr := strconv.ParseFloat(value, 64)
				if parseDurationErr != nil {
					return nil, fmt.Errorf("error parsing segment duration %w", parseDurationErr)
				}
				segmentDuration = duration
			default:
				// unrecognized tag
			}
			continue
		}

		if segmentDuration == 0 {
			return nil, fmt.Errorf("segment at index %d was not preceeded by an EXTINF tag", i)
		}

		segment := &Segment{
			Duration: segmentDuration,
		}
		url, urlParseErr := url.Parse(line)
		if urlParseErr != nil {
			return nil, fmt.Errorf("failed to parse segment url at index %d with error %w", i, urlParseErr)
		}
		segment.URL = url.String()
		segments = append(segments, segment)
	}
	manifest.Segments = segments

	return manifest, nil
}
