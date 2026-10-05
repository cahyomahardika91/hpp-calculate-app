package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
)

type ExportStatic struct {
}

// Signature The name and signature of the console command.
func (r *ExportStatic) Signature() string {
	return "export:static"
}

// Description The console command description.
func (r *ExportStatic) Description() string {
	return "Ekspor resources/views/hpp.tmpl menjadi file statis mandiri di public/index.html dan docs/index.html"
}

// Extend The console command extend.
func (r *ExportStatic) Extend() command.Extend {
	return command.Extend{Category: "export"}
}

// Handle Execute the console command.
func (r *ExportStatic) Handle(ctx console.Context) error {
	tmplPath := filepath.Join("resources", "views", "hpp.tmpl")
	data, err := os.ReadFile(tmplPath)
	if err != nil {
		ctx.Error("Gagal membaca file " + tmplPath + ": " + err.Error())
		return err
	}

	content := string(data)

	// Hapus BOM jika ada
	content = strings.TrimPrefix(content, "\ufeff")

	// Ganti variabel framework Go dengan versi statis
	content = strings.ReplaceAll(content, "{{ .version }}", support.Version)

	// Pastikan simbol aman dengan entitas HTML
	content = strings.ReplaceAll(content, "©", "&copy;")
	content = strings.ReplaceAll(content, "Â©", "&copy;")

	targets := []string{
		filepath.Join("public", "index.html"),
		filepath.Join("docs", "index.html"),
	}

	for _, target := range targets {
		dir := filepath.Dir(target)
		if err := os.MkdirAll(dir, 0755); err != nil {
			ctx.Error("Gagal membuat direktori " + dir + ": " + err.Error())
			return err
		}

		// Tulis file UTF-8 bersih tanpa BOM
		buf := bytes.NewBufferString(content)
		if err := os.WriteFile(target, buf.Bytes(), 0644); err != nil {
			ctx.Error("Gagal menulis file " + target + ": " + err.Error())
			return err
		}
		ctx.Info("Berhasil dibuat: " + target)
	}

	ctx.Success("Ekspor file statis selesai! Siap di-deploy ke GitHub Pages atau Vercel.")
	return nil
}
