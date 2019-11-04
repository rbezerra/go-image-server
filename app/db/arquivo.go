package db

import "database/sql"

type Arquivo struct {
	ID       uint   `json:"ID"`
	UUID     string `json:"UUID"`
	Path     string `json:"Path"`
	ImagemID uint   `json:"ImagemID"`
	Tamanho  string `json:"Tamanho"`
	Original bool   `json:"Original"`
}

func InsertArquivo(arq *Arquivo) (uint, error) {
	stmt, err := db.Prepare("INSERT INTO public.arquivo(imagem_id, tamanho, path, original) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int
	rows, err := stmt.Query(arq.ImagemID, arq.Tamanho, arq.Path, arq.Original)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return uint(id), nil
}

func GetFileByUUIDAndSize(uuid string, size string) (*Arquivo, error) {

	sqlWithSize := `
		SELECT 
			a.id, 
			a.imagem_id, 
			a.tamanho, 
			a.path,
			a.original
		FROM 
			public.arquivo a INNER JOIN public.imagem i ON a.imagem_id = i.id 
		WHERE 
			i.uuid = $1 AND 
			a.tamanho = $2 
	`

	sqlOriginalFile := `
		SELECT 
			a.id, 
			a.imagem_id, 
			a.tamanho, 
			a.path,
			a.original
		FROM 
			public.arquivo a INNER JOIN public.imagem i ON a.imagem_id = i.id 
		WHERE 
			i.uuid = $1 AND 
			a.original = true 
	`

	file := new(Arquivo)

	var sqlStmt string
	if size != "" {
		sqlStmt = sqlWithSize
	} else {
		sqlStmt = sqlOriginalFile
	}

	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if size != "" {
		rows, err = stmt.Query(uuid, size)
	} else {
		rows, err = stmt.Query(uuid)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&file.ID, &file.ImagemID, &file.Tamanho, &file.Path, &file.Original)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return file, nil

}
