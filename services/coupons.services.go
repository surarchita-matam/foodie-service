package services

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"foodie-service/models" 
)

type codeWithFile struct {
	code string
	file string
}

type CouponService struct {
	models *models.BaseModel
}

var couponService *CouponService

func NewCouponService(models *models.BaseModel) *CouponService {
	if couponService != nil {
		return couponService
	}

	couponService = &CouponService{models: models}
	return couponService
}

var (
	filePaths = []string{"couponbase1.gz", "couponbase2.gz", "couponbase3.gz"}
	once      sync.Once
	loadErr   error
)

const (
	batchSize = 20_000_000
)

// Init initializes the database connection and loads coupons if needed
func (cs *CouponService) Init(ctx context.Context) error {
	once.Do(func() {
		if err := cs.models.Coupons.InitCollection(ctx); err != nil {
			loadErr = fmt.Errorf("failed to initialize collection: %v", err)
			return
		}

		exists, err := cs.models.Coupons.CollectionExists(ctx)
		if err != nil {
			loadErr = fmt.Errorf("failed to check collection: %v", err)
			return
		}

		if !exists {
			err = cs.loadCouponsToDB(ctx)
			if err != nil {
				loadErr = fmt.Errorf("failed to load coupons: %v", err)
				return
			}
		} else {
			fmt.Println("Coupons collection already exists and has data. Skipping load.")
		}
	})

	return loadErr
}

func (cs *CouponService) loadCouponsToDB(ctx context.Context) error {
	fmt.Printf("Starting batch processing of coupons (batch size: %d records per file)...\n", batchSize)
	startTime := time.Now()

	totalProcessed := 0
	batchNumber := 1

	for {
		fmt.Printf("\nProcessing Batch #%d (offset: %d records)...\n", batchNumber, totalProcessed)
		batchStartTime := time.Now()

		codeChan := make(chan codeWithFile, 200000)
		readerErrors := make(chan error, len(filePaths))

		// Start file readers for this batch
		var wgReaders sync.WaitGroup
		for _, file := range filePaths {
			wgReaders.Add(1)
			go func(filename string) {
				defer wgReaders.Done()
				fmt.Printf("Processing file: %s (batch offset: %d, limit: %d)\n",
					filename, totalProcessed, batchSize)

				if err := readFileToChannelWithOffset(ctx, filename, codeChan, totalProcessed, batchSize); err != nil {
					readerErrors <- fmt.Errorf("error processing file %s: %v", filename, err)
					return
				}
				fmt.Printf("Completed reading batch from file %s\n", filename)
			}(file)
		}

		go func() {
			wgReaders.Wait()
			close(codeChan)
			close(readerErrors)
		}()

		// Process this batch
		codesWithFiles := make(map[string][]string)
		batchReadCount := 0
		readProgressTime := time.Now()

		for cf := range codeChan {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			files := codesWithFiles[cf.code]
			found := false
			for _, f := range files {
				if f == cf.file {
					found = true
					break
				}
			}
			if !found {
				codesWithFiles[cf.code] = append(files, cf.file)
			}

			batchReadCount++
			if batchReadCount%1000000 == 0 && time.Since(readProgressTime) > 5*time.Second {
				fmt.Printf("Collected %d codes in current batch...\n", batchReadCount)
				readProgressTime = time.Now()
			}
		}

		// Check for reader errors
		for err := range readerErrors {
			if err != nil {
				return err
			}
		}

		fmt.Printf("Batch #%d: Total codes collected: %d\n", batchNumber, batchReadCount)
		fmt.Printf("Batch #%d: Unique codes in batch: %d\n", batchNumber, len(codesWithFiles))

		appearanceCounts := make(map[int]int)
		for _, fileList := range codesWithFiles {
			appearanceCounts[len(fileList)]++
		}
		fmt.Printf("Batch #%d distribution of code appearances:\n", batchNumber)
		for appearances, count := range appearanceCounts {
			fmt.Printf("  Codes appearing in %d file(s): %d\n", appearances, count)
		}


		var couponsToInsert []models.Coupon
		for code, fileList := range codesWithFiles {
			if len(fileList) >= 2 {
				couponsToInsert = append(couponsToInsert, models.Coupon{
					Code:        code,
					FileList:    fileList,
					Appearances: len(fileList),
				})
			}
		}

		if len(couponsToInsert) > 0 {
			fmt.Printf("Batch #%d: Found %d codes appearing in multiple files. Inserting to DB...\n",
				batchNumber, len(couponsToInsert))
			if err := cs.models.Coupons.OptimizedBulkInsert(ctx, couponsToInsert); err != nil {
				return fmt.Errorf("failed to insert batch #%d: %v", batchNumber, err)
			}
		} else {
			fmt.Printf("Batch #%d: No codes found appearing in multiple files.\n", batchNumber)
		}

		batchDuration := time.Since(batchStartTime)
		fmt.Printf("Batch #%d completed in %v\n", batchNumber, batchDuration)

		if batchReadCount == 0 {
			fmt.Println("No more records to process. Finishing...")
			break
		}

		totalProcessed += batchSize
		batchNumber++
	}

	totalDuration := time.Since(startTime)
	fmt.Printf("\nTotal coupon processing completed in %v\n", totalDuration)
	return nil
}

// readFileToChannelWithOffset reads a file starting from a specific offset and sends codes to a channel
func readFileToChannelWithOffset(ctx context.Context, filename string, codes chan<- codeWithFile, offset, limit int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)
	buf := make([]byte, 1024*1024) 
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanWords)

	skipped := 0
	for skipped < offset && scanner.Scan() {
		skipped++
	}

	readCount := 0
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if readCount >= limit {
			break
		}

		word := strings.TrimSpace(scanner.Text())
		if len(word) >= 8 && len(word) <= 10 {
			codes <- codeWithFile{
				code: word,
				file: filename,
			}
			readCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file %s: %v", filename, err)
	}

	return nil
}

func (cs *CouponService) ValidateCode(code string) (bool, float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	valid, err := cs.models.Coupons.ValidateCoupon(ctx, code)
	if err != nil {
		fmt.Printf("Error validating code: %v\n", err)
		return false, 0
	}

	if valid {
		return true, 0.10
	}
	return false, 0
}

func (cs *CouponService) FetchCoupons() ([]models.Coupon, error) {

	coupons, err := cs.models.Coupons.FetchCoupons()
	if err != nil {
		return nil, err
	}
	return coupons, nil
}
