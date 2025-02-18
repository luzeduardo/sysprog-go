package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	mainC4P73RemoveSymlink()
	// mainC4P72CreateSymlink()
	// mainC4P68TraverssingDir()
	// mainC4P66Filepath()
	// mainC4P62FileStats()
	// mainC3P54CliStdStreams()
	// mainC3P45ManipulatingProcesses()
	// mainC1P35SignalingChannel()
	// mainC1P30BufferedChannels()
	// mainC1ChannelWaitGroup()
	// mainC1Channel()
	// mainC1P22Mutexes()
}

func mainC4P73RemoveSymlink() {
	filePath := "shortcut_to_doc.txt"
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error removing the file: %v\n", err)
		return
	}
	fmt.Printf("File removed: %s\n", filePath)
}

func mainC4P72CreateSymlink() {
	sourcePath := "doc.txt"
	symPath := "shortcut_to_doc.txt"

	err := os.Symlink(sourcePath, symPath)
	if err != nil {
		fmt.Printf("Error creating symlink: %v\n", err)
		return
	}

	fmt.Printf("Symlink created: %s -> %s\n", symPath, sourcePath)
}

func mainC4P68TraverssingDir() {
	mainC3P54CliStdStreams()
}

func mainC4P66Filepath() {
	dir := "/home/user"
	file := "document.txt"
	//joining paths
	fullPath := filepath.Join(dir, file)
	fmt.Println(fullPath)
	//cleaning paths
	uncleanPath := "/home/user/../user2"
	clearedPath := filepath.Clean(uncleanPath)
	fmt.Println(clearedPath)

	//splitting paths
	joinedPath := "/home/user/documents/document.txt"
	dirSplitted, fileSplitted := filepath.Split(joinedPath)
	fmt.Printf("Dir: %s - file: %s", dirSplitted, fileSplitted)
}

func mainC4P62FileStats() {
	info, err := os.Stat("main.go")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
			return
		} else {
			panic(err)
		}
	}
	fmt.Printf("File name: %s\n", info.Name())
	fmt.Printf("Size: %v KB\n", info.Size()/1024)
	fmt.Printf("Permissions: %s\n", info.Mode())
	permissionString := fmt.Sprintf("%o", info.Mode().Perm())
	fmt.Printf("Permission string: %s\n", permissionString)
	fmt.Printf("Last modified: %v", info.ModTime())
}

type Option func(*CliConfig) error

func WithErrStream(errStream io.Writer) Option {
	return func(c *CliConfig) error {
		c.ErrStream = errStream
		return nil
	}
}

func WithOutStream(outStream io.Writer) Option {
	return func(c *CliConfig) error {
		c.OutStream = outStream
		return nil
	}
}

type CliConfig struct {
	ErrStream, OutStream io.Writer
	OutputFile           string
}

func NewCliConfig(optsFunc ...Option) (CliConfig, error) {
	c := CliConfig{
		OutputFile: "",
		ErrStream:  os.Stderr,
		OutStream:  os.Stdout,
	}

	for _, optFn := range optsFunc {
		if err := optFn(&c); err != nil {
			return CliConfig{}, err
		}
	}
	return c, nil
}

func traverseDirs(rootDir []string, cfg CliConfig) {
	var outputWriter io.Writer
	if cfg.OutputFile != "" {
		outputFile, err := os.Create(cfg.OutputFile)
		if err != nil {
			fmt.Fprintf(cfg.ErrStream, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer outputFile.Close()
		outputWriter = io.MultiWriter(cfg.OutStream, outputFile)
	} else {
		outputWriter = cfg.OutStream
	}

	for _, directory := range rootDir {
		err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(cfg.ErrStream, "Error accessing the path %q: %v\n", path, err)
			}

			info, _ := d.Info()
			if info.Mode()&os.ModeSymlink != 0 {
				target, err := os.Readlink(path)
				if err != nil {
					fmt.Fprintf(cfg.ErrStream, "Error reading the symlink %q: %v\n", path, err)
				} else {
					_, err := os.Stat(target) //check if target of the symlink exists
					if err != nil {
						if os.IsNotExist(err) {
							fmt.Fprintf(outputWriter, "Broken symlink found: %s -> %s\n", path, target)
						} else {
							fmt.Fprintf(cfg.ErrStream, "Error reading the symlink target %s: %v\n", target, err)
						}
					}
				}
			}

			if path == ".git" {
				return filepath.SkipDir
			}

			if d.IsDir() {
				fmt.Fprintf(outputWriter, "%s\n", path)
			}
			return nil
		})

		if err != nil {
			fmt.Fprintf(cfg.ErrStream, "Error walking the path %q: %v\n", directory, err)
			continue
		}
	}
}

func mainC3P54CliStdStreams() {
	var outputFileName string
	flag.StringVar(&outputFileName, "f", "", "Output file (default: stdout)")
	flag.Parse()

	rootDir := os.Args[1:]
	if len(rootDir) == 0 {
		fmt.Fprintln(os.Stderr, "No rootDir provided")
		os.Exit(1)
	}
	cfg, err := NewCliConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
		os.Exit(1)
	}
	traverseDirs(rootDir, cfg)
}

func mainC3P45ManipulatingProcesses() {
	cmd := exec.Command("ls", "-l")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	//get the process id
	pid := os.Getpid()
	fmt.Println("Process ID: ", pid)
}

func mainC1P35SignalingChannel() {
	var wg sync.WaitGroup
	signalChannel := make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 1: Waiting for signal")
		<-signalChannel
		fmt.Println("Goroutine 1: Received signal")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 2 is about to send a signal")
		signalChannel <- true
		fmt.Println("Goroutine 2: Sent signal")
	}()

	wg.Wait()
	fmt.Println("All goroutines have finished their jobs")
}

func mainC1P30BufferedChannels() {
	clownChannel := make(chan int, 3)
	clowns := 5

	go func() {
		defer close(clownChannel)
		for clownID := range clownChannel {
			ballon := fmt.Sprintf("Ballon %d", clownID)
			fmt.Printf("Driver: Drove the car with %s inside \n", ballon)

			time.Sleep(time.Millisecond * 500)
			fmt.Printf("Driver: Clown finished with %s, the car is ready for more!\n", ballon)
		}
	}()

	var wg sync.WaitGroup

	for clown := 1; clown <= clowns; clown++ {
		wg.Add(1)

		go func(clownID int) {
			defer wg.Done()
			ballon := fmt.Sprintf("Ballon %d", clownID)
			fmt.Printf("Clown %d: Hopped into the car with %s\n", clownID, ballon)
			select {
			case clownChannel <- clownID:
				fmt.Printf("Clown %d: Finished with %s\n", clownID, ballon)
			default:
				fmt.Printf("Clown %d: Ops the car is full, can't fit %s!\n", clownID, ballon)
			}
		}(clown)

		fmt.Println("Circus car ride is over")
	}

	wg.Wait()
	fmt.Println("Circus car ride is over")
}

func mainC1ChannelWaitGroup() {
	balls := make(chan string)
	// create a WaitGroup to wait for the goroutines to finish
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		throwBalls("red", balls)
	}()

	go func() {
		defer wg.Done()
		throwBalls("blue", balls)
	}()

	go func() {
		wg.Wait()
		// close the channel after the goroutines are done
		// so the range loop in the main goroutine can finish
		close(balls)
	}()

	for color := range balls {
		fmt.Printf("%s ball received\n", color)
	}
}

func mainC1Channel() {
	balls := make(chan string)
	go throwBalls("red", balls)
	fmt.Println(<-balls, "received")
}

func throwBalls(color string, balls chan string) {
	fmt.Printf("throwing the %s ball\n", color)
	balls <- color
}

func mainC1P22Mutexes() {
	m := sync.Mutex{}
	fmt.Println("Total Items Packed: ", PackItems(&m, 0))
}

func PackItems(m *sync.Mutex, totalItems int) int {
	const workers = 2
	const itemsPerWorker = 1000

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < itemsPerWorker; j++ {
				m.Lock()
				itemsPacked := totalItems
				itemsPacked++
				totalItems = itemsPacked
				m.Unlock()
			}
		}(i)
	}

	wg.Wait()
	return totalItems
}
