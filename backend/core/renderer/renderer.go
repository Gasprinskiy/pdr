package renderer

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"pdr/backend/core/document"
	"pdr/backend/pkg/render_pool"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
)

type RendererUsecase struct {
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
	renderPool *render_pool.Pool,
	documentsRepo DocumentsRepo,
	//
	cryptoTool Crypto,
	//
	dpi, totalJobs int,
	outDir string,
) *RendererUsecase {
	return &RendererUsecase{
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
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	fileSize := stat.Size()
	fileBytes := make([]byte, fileSize)
	_, err = file.Read(fileBytes)
	if err != nil {
		return err
	}

	pageCount, updateDate, err := u.getFileInfo(fileBytes)
	if err != nil {
		return err
	}

	doc := document.Document{
		ID:         u.cryptoTool.GenerateHexID(),
		UpdateDate: updateDate,
		PageCount:  pageCount,
		FilePath:   param.DocumentFullPath,
		Name:       param.DocumentFullPath,
		Size:       fileSize,
	}

	if err := u.documentsRepo.CreateNewDocument(ctx, doc); err != nil {
		return err
	}

	jobsChan := make(chan int, u.totalJobs)
	errChan := make(chan error)
	results := make(chan pageRenderInfo, u.totalJobs)

	for range u.totalJobs {
		go u.renderPages(u.outDir, &fileBytes, jobsChan, results, errChan)
	}

	go func() {
		for i := 0; i < pageCount; i++ {
			jobsChan <- i
		}
		close(jobsChan)
	}()

	go func() {
		for res := range results {
			page := document.DocumentPage{
				FilePath: res.FilePath,
				Index:    res.Index,
				DocID:    doc.ID,
			}

			if err := u.documentsRepo.CreateNewPage(ctx, page); err != nil {
				errChan <- err
				break
			}
		}
		close(errChan)
	}()

	if err := <-errChan; err != nil {
		return err
	}

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
	errChan chan<- error,
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
			errChan <- err
			return
		}
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

		filePath := fmt.Sprintf("%s/%d.png", outDir, index)
		f, err = os.Create(filePath)
		if err != nil {
			return
		}

		err = png.Encode(f, pageRender.Result.Image)
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
