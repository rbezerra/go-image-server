package handlers

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"../db"
	"../utils"
	"github.com/nfnt/resize"
)

const uploadPath = "./temp-images"

func UploadFile(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")

	if err != nil {
		utils.RenderError(w, "INVALID_FILE", http.StatusBadRequest)
		fmt.Println("INVALID_FILE")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		utils.RenderError(w, "INVALID_FILE", http.StatusInternalServerError)
		fmt.Println("INVALID_FILE")
		fmt.Println(err)
	}

	filetype := http.DetectContentType(fileBytes)
	fileName := utils.CreateUUIDFileName(handler.Filename)
	switch filetype {
	case "image/jpeg", "image/jpg", "image/gif", "image/png":
	default:
		utils.RenderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
		fmt.Println("INVALID_FILE_TYPE")
		fmt.Println(err)
		return
	}

	fileEndings, err := mime.ExtensionsByType(filetype)
	if err != nil {
		utils.RenderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
		fmt.Println("CANT_READ_FILE_TYPE")
		fmt.Println(err)
		return
	}

	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	fmt.Printf("FileType: %s, File: %s\n", filetype, newPath)

	newFile, err := os.Create(newPath)
	if err != nil {
		utils.RenderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		fmt.Println("CANT_WRITE_FILE")
		fmt.Println(err)
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		utils.RenderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
		fmt.Println("CANT_WRITE_FILE")
		fmt.Println(err)
		return
	}

	//salvar referencia da imagem no banco
	img := new(db.Imagem)
	img.UUID = fileName
	imgID, err := db.InsertImage(img)
	if err != nil {
		utils.RenderError(w, "CANT_SAVE_IMAGE_INFO_ON DATABASE", http.StatusInternalServerError)
		fmt.Println("CANT_SAVE_IMAGE_INFO_ON DATABASE")
		fmt.Println(err)
		return
	}

	//salvar referência do arquivo original
	arq := new(db.Arquivo)
	arq.Tamanho = "original"
	arq.Path = newPath
	arq.ImagemID = imgID
	_, err = db.InsertArquivo(arq)
	if err != nil {
		utils.RenderError(w, "CANT_SAVE_FILE_INFO_ON DATABASE", http.StatusInternalServerError)
		fmt.Println("CANT_SAVE_FILE_INFO_ON DATABASE")
		fmt.Println(err)
		return
	}

	if imagesCreated, err := createStandardImages(newPath, fileBytes, fileName); err != nil || imagesCreated == 0 {
		utils.RenderError(w, "CANT_CREATE_IMAGES", http.StatusInternalServerError)
		fmt.Println("CANT_CREATE_IMAGES")
		fmt.Println(err)
	}

	w.Write([]byte(fileName))

}

func ListImages(w http.ResponseWriter, r *http.Request) {
	imgs, err := db.ListAllImages()

	if err != nil {
		utils.RenderError(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	for _, img := range imgs {
		fmt.Fprintln(w, img.ID, img.UUID, img.Descricao)
	}
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	tamanho := mux.Vars(r)["tamanho"]

	file, err := db.GetFileByUUIDAndSize(uuid, tamanho)
	if err != nil {
		utils.RenderError(w, http.StatusText(500), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	//se não for encotrado o arquivo deve ser gerado
	if file == nil {
		utils.RenderError(w, http.StatusText(404), http.StatusNotFound)
	} else {
		json.NewEncoder(w).Encode(file)
	}

}

func createStandardImages(originalFilePath string, originalFileBytes []byte, fileName string) (int, error) {
	standardSizes := []string{"150x112", "100x75", "75x75", "75x100", "180x240", "375x500", "768x1024"}
	imagesCreated := 0
	filetype := http.DetectContentType(originalFileBytes)

	file, err := os.Open(originalFilePath)
	if err != nil {
		return imagesCreated, err
	}

	var img image.Image
	if filetype == "image/png" {
		img, err = png.Decode(file)
		if err != nil {
			return imagesCreated, err
		}
	}

	if filetype == "image/jpg" || filetype == "image/jpeg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			return imagesCreated, err
		}
	}

	file.Close()

	DBImg, err := db.GetImageByUUID(fileName)
	if err != nil {
		return imagesCreated, err
	}

	for _, size := range standardSizes {
		s := strings.Split(size, "x")
		h, _ := strconv.ParseUint(s[0], 10, 32)
		w, _ := strconv.ParseUint(s[1], 10, 32)
		height := uint(h)
		width := uint(w)
		newImg := resize.Resize(width, height, img, resize.NearestNeighbor)

		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			return imagesCreated, err
		}

		newPath := filepath.Join(uploadPath, fileName+"-"+size+fileEndings[0])
		outFile, err := os.Create(newPath)
		if err != nil {
			return imagesCreated, err
		}
		defer outFile.Close()

		if filetype == "image/png" {
			png.Encode(outFile, newImg)
		}

		if filetype == "image/jpeg" || filetype == "image/jpg" {
			jpeg.Encode(outFile, newImg, nil)
		}

		imagesCreated++

		//salvar referencia na tabela de arquivo no banco
		arq := new(db.Arquivo)
		arq.ImagemID = DBImg.ID
		arq.Path = newPath
		arq.Tamanho = size
		_, err = db.InsertArquivo(arq)
		if err != nil {
			fmt.Println("CANT_SAVE_FILE_INFO_ON DATABASE")
			fmt.Println(err)
			return imagesCreated, nil
		}
	}

	return imagesCreated, nil
}