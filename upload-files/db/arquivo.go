package db

type Arquivo struct {
	ID       uint   `json:"ID"`
	ImagemID uint   `json:"ImagemID"`
	Tamanho  string `json:"Tamanho"`
	Path     string `json:"Path"`
}

func InsertArquivo(arq *Arquivo) (uint, error) {
	stmt, err := db.Prepare("INSERT INTO public.arquivo(imagem_id, tamanho, path) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int
	rows, err := stmt.Query(arq.ImagemID, arq.Tamanho, arq.Path)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return uint(id), nil
}

func GetFileByUUIDAndSize(uuid string, size string) (*Arquivo, error) {
	file := new(Arquivo)

	stmt, err := db.Prepare(`
		SELECT 
			a.id, 
			a.imagem_id, 
			a.tamanho, 
			a.path 
		FROM 
			public.arquivo a INNER JOIN public.imagem i ON a.imagem_id = i.id 
		WHERE 
			i.uuid = $1 AND 
			a.tamanho = $2 
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(uuid, size)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&file.ID, &file.ImagemID, &file.Tamanho, &file.Path)
		if err != nil {
			return nil, err
		}
	}

	return file, nil

}
