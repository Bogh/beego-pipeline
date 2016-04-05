package pipeline

type Processor struct {
	// type of asset
	t       string
	outputs Outputs
}

func NewProcessor(t string, outputs Outputs) *Processor {
	return &Processor{t, outputs}
}

// func (p *Processor) Process() error {
// 	// compile then compress

// 	// get list of files for each output and
// 	// normalize paths and check if they exist, otherwise issue an ignore
// 	for _, output := range p.outputs {
// 		err := p.Compress(output)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	return nil
// }

// // Accepts an io.Writer and returns an io.Reader
// func (p *Processor) Compress(output Output) error {
// 	compressor, ok := compressors[p.t]
// 	if !ok {
// 		beego.Debug(compressors)
// 		return fmt.Errorf("Compressor not found for type: %s", p.t)
// 	}

// 	// normalized paths
// 	sources, _ := output.Paths()

// 	files := make([]io.Reader, len(sources))
// 	for i, path := range sources {
// 		f, _ := os.Open(path)
// 		files[i] = io.Reader(f)
// 		defer f.Close()
// 	}
// 	done := make(chan bool)

// 	r, err := compressor.Compress(done, io.MultiReader(files...))
// 	if err != nil {
// 		beego.Error(err)
// 		return err
// 	}

// 	// output File

// 	go func() {
// 		defer r.Close()
// 		oFile, _ := os.OpenFile(output.NOutput(), os.O_WRONLY|os.O_CREATE, 0644)
// 		defer oFile.Close()
// 		n, _ := io.Copy(oFile, r)
// 		beego.Debug("Bytes written to file: ", n)
// 		<-done
// 	}()
// 	return nil
// }
