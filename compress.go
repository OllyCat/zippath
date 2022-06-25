package zippath

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Compress(p string, n string) (err error) {
	// создаём файл архива
	f, err := os.Create(n)
	if err != nil {
		return fmt.Errorf("Could not create archive: %w", err)
	}
	defer f.Close()

	// writer для zip архива
	z := zip.NewWriter(f)
	defer z.Close()

	// счётчик запакованных файлов
	var count int

	// проход по содержимому папки
	err = filepath.Walk(n, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// игнорируем, если директорий
		if info.IsDir() {
			return nil
		}

		// если файл - открываем его на чтение
		rf, e := os.Open(path)
		if e != nil {
			return fmt.Errorf("Error open file: %w", e)
		}
		// закрываем по окончании
		defer rf.Close()

		// создаём файл в архиве
		zf, e := z.Create(path)
		if e != nil {
			return fmt.Errorf("Error archive file: %w", e)
		}

		// копируем содержимое файла в архив
		_, e = io.Copy(zf, rf)
		if e != nil {
			return fmt.Errorf("Error copy file: %w", e)
		}
		// если всё хорошо - увеличим счётчик файлов
		count++
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error create compressed file: %w", err)
	}

	return
}
