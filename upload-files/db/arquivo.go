package db

type Arquivo struct {
	ID       uint
	ImagemID uint
	Tamanho  string
	Path     string
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
