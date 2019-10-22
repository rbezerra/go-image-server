package db

type Imagem struct {
	ID        uint
	UUID      string
	Descricao string
}

func ListAllImages() ([]*Imagem, error) {
	rows, err := db.Query("SELECT id, uuid, descricao FROM imagem")
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

func InsertImage(img *Imagem) (uint, error) {
	stmt, err := db.Prepare("INSERT INTO public.imagem (uuid) VALUES ($1) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var id int
	rows, err := stmt.Query(img.UUID)
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

func GetImageByUUID(UUID string) (*Imagem, error) {
	stmt, err := db.Prepare("SELECT id, uuid, descricao FROM public.imagem WHERE uuid = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(UUID)
	if err != nil {
		return nil, err
	}

	img := new(Imagem)
	for rows.Next() {
		err := rows.Scan(&img.ID, &img.UUID, &img.Descricao)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return img, nil

}
