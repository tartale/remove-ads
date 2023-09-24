package rmads

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/tartale/go/pkg/logz"
	"github.com/tartale/go/pkg/mathx"
	"github.com/tartale/go/pkg/slicez"
	"github.com/tartale/remove-ads/pkg/config"
)

// Invert will, given a list of non-overlapping
// continuously increasing segments, returns a list of segments that
// captures the inversion of the input list. For example,
// if the input segment is a list of time frames to cut out
// of a video, then the inversion would be a list of
// time frames to keep:
//
//	input:    {5:10}, {15:20}, {30:40}
//	inverted: {0:5},  {10:15}, {20:30}, {40:end}
func (s Segments) Invert(endOffset time.Duration) Segments {

	var inverted Segments

	if len(s) == 0 {
		inverted = append(inverted, Segment{StartOffset: 0, EndOffset: endOffset})
		return inverted
	}

	for i, segment := range s {
		if i == 0 && segment.StartOffset == 0 {
			continue
		}
		if i == 0 {
			inverted = append(inverted, Segment{
				StartOffset: 0 * time.Second,
				EndOffset:   segment.StartOffset,
			})
			continue
		}
		inverted = append(inverted, Segment{
			StartOffset: s[i-1].EndOffset,
			EndOffset:   s[i].StartOffset,
		})
	}

	lastSegment := slicez.MustGetLast(s)
	if lastSegment.EndOffset != endOffset {
		inverted = append(inverted, Segment{
			StartOffset: lastSegment.EndOffset,
			EndOffset:   endOffset,
		})
	}

	return inverted
}

// ffmpeg -y -i example.mp4 -vf select='between(t\,10\,20)+between(t\,30\,40\),setpts=N/FRAME_RATE/TB' -af aselect='between(t\,10\,20)+between(t\,30\,40\),asetpts=N/SR/TB'  test_output.mp4
// ffmpeg -y -i intTestTransportStream.ts -vf select='between(t\\,0\\,413)+between(t\\,607\\,970)+between(t\\,1166\\,1679)+between(t\\,1903\\,1980),setpts=N/FRAME_RATE/TB' -af aselect='between(t\\,0\\,413)+between(t\\,607\\,970)+between(t\\,1166\\,1679)+between(t\\,1903\\,1980),asetpts=N/SR/TB'"
func (s Segments) Remove(ctx context.Context, inputFilePath, outputFilePath string) error {

	var logger = logz.Logger()
	ffmpegCmd, err := s.makeRemoveCommand(inputFilePath, outputFilePath)
	if err != nil {
		return err
	}
	logger.Debugf("ffmpeg command: %s\n", ffmpegCmd.String())

	return nil
}

func (s Segments) makeRemoveCommand(inputFilePath, outputFilePath string) (*exec.Cmd, error) {

	if inputFilePath == "" {
		inputFilePath = "-"
	}
	if outputFilePath == "" {
		outputFilePath = "-"
	}

	var timeSelect []string
	for _, segment := range s {
		startOffset := mathx.Floor(segment.StartOffset.Seconds())
		endOffset := mathx.Ceil(segment.EndOffset.Seconds())
		timeSelect = append(timeSelect, fmt.Sprintf(`between(t\,%d\,%d)`, startOffset, endOffset))
	}

	timeSelectArg := strings.Join(timeSelect, "+")
	videoSelectArg := fmt.Sprintf("select='%s,setpts=N/FRAME_RATE/TB'", timeSelectArg)
	audioSelectArg := fmt.Sprintf("aselect='%s,asetpts=N/SR/TB'", timeSelectArg)
	ffmpegCmd := exec.Command(config.Values.FFmpegFilePath, "-y", "-i", inputFilePath,
		"-vf", videoSelectArg, "-af", audioSelectArg, outputFilePath)

	return ffmpegCmd, nil
}

/*

<VideoReDoProject Version="3">
<Filename>/Users/Tom/Projects/kmttg-plus/go/./output/download/LocalNews/LocalNews.ts</Filename><CutList>
<InputPIDList><VideoStreamPID>97</VideoStreamPID>
<AudioStreamPID>100</AudioStreamPID><SubtitlePID1>0</SubtitlePID1></InputPIDList>
<Cut><CutTimeStart>0</CutTimeStart> <CutTimeEnd>99766333</CutTimeEnd> </Cut>
<Cut><CutTimeStart>4956295781</CutTimeStart> <CutTimeEnd>7053955958</CutTimeEnd> </Cut>
<Cut><CutTimeStart>11669324333</CutTimeStart> <CutTimeEnd>13268922333</CutTimeEnd> </Cut>
<Cut><CutTimeStart>17732453185</CutTimeStart> <CutTimeEnd>17986217201</CutTimeEnd> </Cut>
</CutList>
<SceneList>
<SceneMarker Sequence="0" Timecode="0:00:09.97">99766333</SceneMarker>
<SceneMarker Sequence="1" Timecode="0:08:15.58">4955889201</SceneMarker>
<SceneMarker Sequence="2" Timecode="0:08:30.17">5101763275</SceneMarker>
<SceneMarker Sequence="3" Timecode="0:09:15.78">5557885268</SceneMarker>
<SceneMarker Sequence="4" Timecode="0:09:21.58">5615889666</SceneMarker>
<SceneMarker Sequence="5" Timecode="0:10:15.73">6157330614</SceneMarker>
<SceneMarker Sequence="6" Timecode="0:10:45.44">6454447806</SceneMarker>
<SceneMarker Sequence="7" Timecode="0:11:15.84">6758418332</SceneMarker>
<SceneMarker Sequence="8" Timecode="0:11:45.39">7053955958</SceneMarker>
<SceneMarker Sequence="9" Timecode="0:19:26.89">11668990667</SceneMarker>
<SceneMarker Sequence="10" Timecode="0:20:56.23">12562301711</SceneMarker>
<SceneMarker Sequence="11" Timecode="0:21:11.33">12713360843</SceneMarker>
<SceneMarker Sequence="12" Timecode="0:21:56.64">13166486667</SceneMarker>
<SceneMarker Sequence="13" Timecode="0:22:06.89">13268922333</SceneMarker>
<SceneMarker Sequence="14" Timecode="0:24:50.96">14909624730</SceneMarker>
<SceneMarker Sequence="15" Timecode="0:29:33.20">17732090428</SceneMarker>
<SceneMarker Sequence="16" Timecode="0:29:34.93">17749396984</SceneMarker>
<SceneMarker Sequence="17" Timecode="0:29:58.62">17986217201</SceneMarker>
</SceneList>
</VideoReDoProject>

*/
