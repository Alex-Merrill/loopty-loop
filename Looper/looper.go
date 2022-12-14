package looper

import (
	"fmt"
	"math"
	"sync"
	"time"

	vidio "github.com/AlexEidt/Vidio"
)

type Video struct {
	frames  [][]float64
	width   int
	height  int
	fps     float64
	bitrate int
}

type Looper struct {
	vid                      Video
	minDuration, maxDuration int
	startFrame               int
	endFrame                 int
}

func NewLoop(path string, minDuration, maxDuration int) Looper {
	// TODO: Probably wanna resize the video to smaller size so its faster to process? or give an option to do so.
	// maybe use grayscale?
	fmt.Println("Reading video...")
	startTime := time.Now()
	video := readVideo(path)
	fmt.Printf("Done! Time to read video: %v\n", time.Since(startTime))

	return Looper{
		vid:         video,
		minDuration: minDuration,
		maxDuration: maxDuration,
		startFrame:  0,
		endFrame:    len(video.frames) - 1,
	}
}

func (l *Looper) Start() (bool, error) {
	fmt.Println("Reading frames...")
	startTime := time.Now()
	diffs := l.getAllFrameDiffs()
	fmt.Printf("Time to calc diffs of all eligible frames: %v\n", time.Since(startTime))
	fmt.Printf("diffs: %v\n", diffs[:1])

	startTime = time.Now()
	start, end := l.getBestLoop(diffs)
	fmt.Printf("Time to calc best start/end frames: %v\n", time.Since(startTime))
	fmt.Printf("Found best start/end at frames: %d, %d\n", start, end)

	var err error
	err = nil

	return true, err
}

func (l *Looper) getBestLoop(diffs [][]float64) (int, int) {
    // convert diffs slice to be slice of slices of structs
    // this will allow us to sort the array by difference value
    // to find the best loop from a given frame in O(f) where f is 
    // number of frames, for total runtime of O(f^2), instead of having to
    // brute force find best loop from a given frame in O(f^2) for total O(f^3)
    // TODO: Convert to tuple array way earlier when initially process data to avoid
    // doing this work again. Then we can just sort it here.
    type frames struct {
        frame int
        difference float64
    }
    var tupleDiffs [][]frames
    for row := 0; row < len(diffs); row++ {
        currFrames := make([]frames, len(diffs[row]))
		for col := 0; col < len(diffs[0]); col++ {
            frameToAdd := frames{
                frame: col,
                difference: diffs[row][col],
            } 
            currFrames[col] = frameToAdd
		}
        sort.Slice(currFrames, func(i,j int) bool {
            return currFrames[i].difference < currFrames[j].difference
        })
        tupleDiffs = append(tupleDiffs, currFrames)
    }


	var currStart int
	var currEnd int
    var currDifference float64
	var currDuration float64
	for row := 0; row < len(diffs); row++ {
		for col := 0; col < len(diffs[0]); col++ {
            if diffs[currStart]
		}
	}

	return start, end
}

func (l *Looper) getAllFrameDiffs() [][]float64 {
	frameDiffs := make([][]float64, len(l.vid.frames))

	// only calc frame differences if there are enough frames after it to satisfy min duration
	fps := int(l.vid.fps)
	lastFrameEligible := len(l.vid.frames) - (fps * l.minDuration)

	for i := 0; i < lastFrameEligible; i++ {
		startTime := time.Now()
		fmt.Printf("Starting frame diff for frame %d...\n", i)
		frameDiffs[i] = l.getFrameDiffs(i)
		duration := time.Since(startTime)
		fmt.Printf("Frame %d done! Took %v\n", i, duration)
	}

	return frameDiffs
}

// returns slice of avg pixel diferences for frame idx and all frames after it
func (l *Looper) getFrameDiffs(idx int) []float64 {
	diffs := make([]float64, len(l.vid.frames))

	lastFrameEligible := idx + (int(l.vid.fps) * l.maxDuration)

	lenToUse := int(math.Min(float64(len(l.vid.frames)), float64(lastFrameEligible)))

	var wg sync.WaitGroup
	wg.Add(lenToUse - idx - 1)
	for i := idx + 1; i < lenToUse; i++ {
		go func(i int) {
			defer wg.Done()
			diffs[i] = getFramePixelDiff(l.vid.frames[idx], l.vid.frames[i])
		}(i)
	}
	wg.Wait()

	return diffs
}

// returns average pixel difference between frames f1 and f2
func getFramePixelDiff(f1, f2 []float64) float64 {
	totalDiff := 0.0
	for i := 0; i < len(f1); i += 3 {
		r1 := f1[i]
		g1 := f1[i+1]
		b1 := f1[i+2]
		r2 := f2[i]
		g2 := f2[i+1]
		b2 := f2[i+2]

		rBar := (r1 + r2) / 2
		dR := r1 - r2
		dB := b1 - b2
		dG := g1 - g2

		p1 := (2 + rBar/256) * math.Pow(dR, 2)
		p2 := 4 * math.Pow(dG, 2)
		p3 := (2 + (255-rBar)/256) * math.Pow(dB, 2)

		dC := math.Sqrt(p1 + p2 + p3)

		totalDiff += dC
	}

	totalDiff /= float64(len(f1))

	return totalDiff
}

func writeImageFromFrame(filename string, f []float64, w, h int) {
	fB := make([]byte, len(f))
	for i, float := range f {
		fB[i] = byte(float)
	}
	vidio.Write(filename, w, h, fB)
}

func readVideo(path string) Video {
	video, err := vidio.NewVideo(path)
	checkErr(err)

	var frames [][]float64
	for video.Read() {
		frame := video.FrameBuffer()
		newFrame := make([]float64, len(frame))
		for i, b := range frame {
			newFrame[i] = float64(b)
		}
		frames = append(frames, newFrame)
	}

	return Video{
		frames:  frames,
		width:   video.Width(),
		height:  video.Height(),
		fps:     video.FPS(),
		bitrate: video.Bitrate(),
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
