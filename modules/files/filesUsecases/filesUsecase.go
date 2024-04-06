package filesUsecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/PHURINTOR/phurinshop/config"
	"github.com/PHURINTOR/phurinshop/modules/files"
)

// ======================================= Interface =========================================
type IFilesUsecase interface {
	UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error)
	DeleteFileGCP(req []*files.DeleteFileReq) error
}

// ======================================= Struct ============================================
type filesUsecase struct {
	cfg config.IConfig
}

type filesPub struct {
	bucket      string
	destination string
	file        *files.FileRes
}

// ======================================= Constructor =======================================
func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

// ----------Flow
//    upload ---> make public  ---> upload worker (func) ---> output

// ========================================Missing Function ==============================
// ------------------------- pre Function Upload File to GCP
// ข้อมูลที่รับเข้ามาต้อง make public เสียก่อน

// --------------- make public ---------
// makePublic gives all users read access to an object.
func (f *filesPub) makePublic(ctx context.Context, client *storage.Client) error {

	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("cvcv.Set: %w", err)
	}
	fmt.Printf("Blob %v is now publicly accessible.\n", f.destination)
	return nil
}

func (u *filesUsecase) uploadWorkers(ctx context.Context, client *storage.Client, jobs <-chan *files.FileReq, reuslt chan<- *files.FileRes, errs chan<- error) {
	// สิ่งที่ต้องใช้ใน GCP Files Upload คือ
	// *** object = jobsCh
	// *** bucket = config path and detail  to bucket

	for job := range jobs {
		container, err := job.File.Open() // .File.Open()  ---> io.util  = byte
		if err != nil {
			errs <- err
			return
		}

		// Return output byte type
		b, err := ioutil.ReadAll(container)
		if err != nil {
			errs <- err
			return
		}

		//***** buf = object files.Req to byte ก่อน
		buf := bytes.NewBuffer(b)

		wc := client.Bucket(u.cfg.App().Gcpbucket()).Object(job.Destination).NewWriter(ctx)

		// Upload an object with storage.Writer.
		if _, err = io.Copy(wc, buf); err != nil {
			errs <- fmt.Errorf("io.Copy: %w", err)
			return
		}

		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errs <- fmt.Errorf("Writer.Close: %w", err)
			return
		}
		fmt.Printf("%v uploaded to %v. \n", job.FileName, job.Extension)
		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().Gcpbucket(), job.Destination),
			},
			bucket:      u.cfg.App().Gcpbucket(),
			destination: job.Destination,
		}

		if err := newFile.makePublic(ctx, client); err != nil {
			errs <- err
			return
		}
		errs <- nil
		reuslt <- newFile.file

		// ด้วยความที่ goroutine รับเรื่อยๆ หากไม่มีส่งเข้าไปก็จะวนทำงานไปเรื่อยๆ ไม่จบ  ต้องให้ทุกครั้งที่รับ  ตัดจำนวน buffer ออกไปด้วย
	}

}

// --------------- Delete File Fuction pre to pool worker
func (u *filesUsecase) deleteFileWorker(ctx context.Context, client *storage.Client, jobs <-chan *files.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		o := client.Bucket(u.cfg.App().Gcpbucket()).Object(job.Destination)

		attrs, err := o.Attrs(ctx)
		if err != nil {
			errs <- fmt.Errorf("object.Attrs: %v", err)
			return
		}
		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errs <- fmt.Errorf("Object(%q).Delete: %v", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v delete. \n", job.Destination)
		errs <- nil
	}
}

//------------------------- Upload File to GCP -------

func (u *filesUsecase) UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error) {

	// 1. Create Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	//2. GCP  open connect to storage Bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	//3. Upload แบบ PoolWorker
	//------------------------------- Worker pool-------------------------------

	//	    		  INPUT         |   		Chan		   |	   OUTPUT
	// req(array files.FileReq) --> |jobsCh -+-+-+ resultsCh   | --->  res(array *files.FileRes)

	// 3.1 inital
	jobsCh := make(chan *files.FileReq, len(req)) //len(req) = buffer channal
	resultsCh := make(chan *files.FileRes, len(req))
	errCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0) //ไว้รอรับจาก resultsCh (output)

	// 3.2 Assign Req(array) to JobsCh(Chan) = (input) เพราะ req ไม่ได้เป็นตัวแปรชนิด Chan ** การจะให้ Chan สื่อสารกัน ต้องเป็น Chan ทั้งคู่
	for _, r := range req { // r = req
		jobsCh <- r
	}
	close(jobsCh)

	// 3.3 Worker on JobsCH
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		//******worker = function upload
		go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errCh)
	}

	// 3.4 output --> Result
	for a := 0; a < len(req); a++ {
		err := <-errCh // err เป็น reciver รับค่าจาก chan errCh ที่เรียกใช้ใน Function 3.3
		if err != nil {
			return nil, err
		}

		// ไม่มี error
		result := <-resultsCh
		res = append(res, result) //res (array)
	}
	return res, nil

	//่ jobsCh ทำงานไปเรื่อย  ----> ถ้ามี error ก็จะใส่ไปเรื่อยๆ  ---> result 3.4 ก็จะทำงานไปเรื่อยๆ เช่นกัน
}

// ------------------------- Delete File to GCP -------
func (u *filesUsecase) DeleteFileGCP(req []*files.DeleteFileReq) error {

	// 1. New Context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	//2. GCP  open connect to storage Bucket
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// 3.1 inital
	jobsCh := make(chan *files.DeleteFileReq, len(req)) //len(req) = buffer channal
	errCh := make(chan error, len(req))

	// 3.2 Assign Req(array) to JobsCh(Chan) = (input) เพราะ req ไม่ได้เป็นตัวแปรชนิด Chan ** การจะให้ Chan สื่อสารกัน ต้องเป็น Chan ทั้งคู่
	for _, r := range req { // r = req
		jobsCh <- r
	}
	close(jobsCh)

	// 3.3 Worker on JobsCH
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		//******worker = function upload
		go u.deleteFileWorker(ctx, client, jobsCh, errCh)
	}

	// 3.4 output --> Result
	for a := 0; a < len(req); a++ {
		err := <-errCh // err เป็น reciver รับค่าจาก chan errCh ที่เรียกใช้ใน Function 3.3
		return err

	}
	return nil

}
