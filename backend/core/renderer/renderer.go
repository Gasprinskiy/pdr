package renderer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"pdr/backend/core/document"
	"pdr/backend/core/shared"
	"pdr/pkg/render_pool"
	"pdr/pkg/z_logger"
	"sync"

	"github.com/chai2010/webp"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

type RendererUsecase struct {
	log           z_logger.Logger
	renderPool    *render_pool.Pool
	documentsRepo DocumentsRepo
	//
	cryptoTool Crypto
	//
	dpi       int
	totalJobs int
	outDir    string
}

func NewRendererUsecase(
	log z_logger.Logger,
	renderPool *render_pool.Pool,
	documentsRepo DocumentsRepo,
	//
	cryptoTool Crypto,
	//
	dpi, totalJobs int,
	outDir string,
) *RendererUsecase {
	return &RendererUsecase{
		log:           log,
		renderPool:    renderPool,
		documentsRepo: documentsRepo,
		//
		cryptoTool: cryptoTool,
		//
		dpi:       dpi,
		totalJobs: totalJobs,
		outDir:    outDir,
	}
}

func (u *RendererUsecase) RenderPDFDocumentPages(ctx context.Context, param RenderPDFDocumentPagesParam) error {
	file, err := os.Open(param.DocumentFullPath)
	if err != nil {
		u.log.Error("could not open file", err)
		return shared.ErrFileRead
	}

	stat, err := file.Stat()
	if err != nil {
		u.log.Error("could not get file stats", err)
		return shared.ErrFileRead
	}

	fileSize := stat.Size()
	fileBytes := make([]byte, fileSize)
	_, err = file.Read(fileBytes)
	if err != nil {
		u.log.Error("could not read file", err)
		return shared.ErrFileRead
	}

	pageCount, updateDate, err := u.getFileInfo(fileBytes)
	if err != nil {
		u.log.Error("could not get file info", err)
		return shared.ErrFileRead
	}

	doc := document.Document{
		ID:         u.cryptoTool.GenerateHexID(),
		UpdateDate: updateDate,
		PageCount:  pageCount,
		FilePath:   param.DocumentFullPath,
		Name:       param.TempName(),
		Size:       fileSize,
	}

	if err := u.documentsRepo.CreateNewDocument(ctx, doc); err != nil {
		u.log.Error("could not create new document in repo", err)
		return shared.ErrLocalStorage
	}

	imagesOutPath := filepath.Join(u.outDir, doc.ID)

	if err := os.Mkdir(imagesOutPath, 0755); err != nil {
		u.log.Error("could not create directory", err)
		return shared.ErrFileRead
	}

	jobsChan := make(chan int, u.totalJobs)
	errChan := make(chan error)
	results := make(chan pageRenderInfo, u.totalJobs)

	var once sync.Once
	closeResultsChan := func() {
		once.Do(func() { close(results) })
	}

	for range u.totalJobs {
		go u.renderPages(imagesOutPath, &fileBytes, jobsChan, results, closeResultsChan)
	}

	go func() {
		for i := 0; i < pageCount; i++ {
			jobsChan <- i
		}
		close(jobsChan)
	}()

	go func() {
		var count int

		for res := range results {
			if res.Err != nil {
				errChan <- res.Err
				break
			}

			page := document.DocumentPage{
				FilePath: res.FilePath,
				Index:    res.Index,
				DocID:    doc.ID,
			}

			if err := u.documentsRepo.CreateNewPage(ctx, page); err != nil {
				errChan <- err
				break
			}

			count += 1
			param.OnUpdate(OnUpdatePayload{
				Count: count,
				OutOf: pageCount,
			})
		}
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		close(results)

		if err := os.RemoveAll(imagesOutPath); err != nil {
			u.log.Error("could not remove images folder", err)
		}
		u.log.Error("error while render or create page", err)
		return ErrWhileRender
	}
	close(results)

	return nil
}

func (u *RendererUsecase) getFileInfo(fileByles []byte) (int, string, error) {
	instance, err := u.renderPool.Instance()
	if err != nil {
		return 0, "", err
	}
	defer instance.Close()

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &fileByles,
	})
	if err != nil {
		return 0, "", err
	}

	result, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return 0, "", err
	}

	modDate, err := instance.FPDF_GetMetaText(&requests.FPDF_GetMetaText{
		Document: doc.Document,
		Tag:      "ModDate",
	})

	return result.PageCount, modDate.Value, nil
}

func (u *RendererUsecase) renderPages(
	outDir string,
	fileBytes *[]byte,
	jobsChan <-chan int,
	results chan<- pageRenderInfo,
	done func(),
) {
	var (
		err        error
		index      int
		instance   pdfium.Pdfium
		doc        *responses.OpenDocument
		pageRender *responses.RenderPageInDPI
		f          *os.File
	)

	defer func() {
		if err != nil {
			results <- pageRenderInfo{
				Err: err,
			}
			return
		}
		done()
	}()

	instance, err = u.renderPool.Instance()
	if err != nil {
		return
	}
	defer instance.Close()

	doc, err = instance.OpenDocument(&requests.OpenDocument{
		File: fileBytes,
	})
	if err != nil {
		return
	}

	for index = range jobsChan {
		pageRender, err = instance.RenderPageInDPI(&requests.RenderPageInDPI{
			DPI: u.dpi,
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    index,
				},
			},
		})
		if err != nil {
			return
		}

		filePath := fmt.Sprintf("%s/%d.webp", outDir, index)
		f, err = os.Create(filePath)
		if err != nil {
			return
		}

		err = webp.Encode(f, pageRender.Result.Image, &webp.Options{
			Lossless: true,
			Quality:  95,
		})
		if err != nil {
			f.Close()
			return
		}
		f.Close()

		results <- pageRenderInfo{
			Index:    index,
			FilePath: filePath,
		}
	}
}
