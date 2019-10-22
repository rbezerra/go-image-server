package db

type Imagem struct {
	ID        uint
	UUID      string
	Descricao string
}

func ListAllImages() ([]*Imagem, error) {
	rows, err := db.Query("SELECT * FROM imagem")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imgs := make([]*Imagem, 0)
	for rows.Next() {
		img := new(Imagem)
		err := rows.Scan(&img.ID, &img.UUID, &img.Descricao)
		if err != nil {
			return nil, err
		}

		imgs = append(imgs, img)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return imgs, nil
}

func InsertImage(img *Imagem) error {
	stmt, err := db.Prepare("INSERT INTO public.imagem (uuid) VALUES ($1)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Query(img.UUID); err != nil {
		return err
	}

	return nil
}
