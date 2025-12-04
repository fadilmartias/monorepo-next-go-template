package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "github.com/fadilmartias/dilz_code/apps/backend/config"
)

// safeRemove menghapus file dengan retry, logging, dan timeout maksimum
func safeRemove(path string) error {
	var err error
	const maxRetries = 20
	const retryDelay = 500 * time.Millisecond
	const maxTimeout = 15 * time.Second

	timeout := time.After(maxTimeout)
	for i := 0; i < maxRetries; i++ {
		select {
		case <-timeout:
			return fmt.Errorf("gagal menghapus file %s: timeout setelah %v", path, maxTimeout)
		default:
			if !isFileLocked(path) {
				err = os.Remove(path)
				if err == nil || os.IsNotExist(err) {
					fmt.Printf("Berhasil menghapus file: %s\n", path)
					return nil
				}
			}
			fmt.Printf("Gagal menghapus %s (percobaan %d): %v\n", path, i+1, err)
			time.Sleep(retryDelay)
		}
	}
	return fmt.Errorf("gagal menghapus file %s setelah %d percobaan: %w", path, maxRetries, err)
}

// isFileLocked memeriksa apakah file terkunci
func isFileLocked(path string) bool {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_EXCL, 0644)
	if err != nil {
		fmt.Printf("File %s terkunci atau tidak bisa dibuka: %v\n", path, err)
		return true
	}
	file.Close()
	return false
}

// asyncRemove menjalankan penghapusan file secara asinkronus dengan delay
func asyncRemove(path string) {
	go func() {
		time.Sleep(5 * time.Second) // Jeda lebih lama untuk memastikan file dilepaskan
		if err := safeRemove(path); err != nil {
			fmt.Printf("Gagal menghapus file secara asinkronus %s: %v\n", path, err)
			// Optional: simpan path ke file log untuk cleanup manual
			_ = os.WriteFile("failed_deletes.log", []byte(fmt.Sprintf("%s: %v\n", path, err)), 0644)
		}
	}()
}

// checkFilePermissions memeriksa apakah file memiliki izin yang cukup
func checkFilePermissions(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("gagal memeriksa izin file %s: %w", path, err)
	}
	file.Close()
	return nil
}

// MoveImageIfExists memindahkan file dari tmp ke finalDir dengan aman
func MoveImageIfExists(fileName string, finalDir string, errorMsg string) error {
	if fileName == "" {
		return nil
	}

	if !strings.HasPrefix(finalDir, "public/") {
		finalDir = fmt.Sprintf("public/%s", finalDir)
	}

	tmpPath := filepath.Join("public", "uploads", "tmp", "images", fileName)
	finalPath := filepath.Join(finalDir, fileName)

	// Jika file tidak ada di tmpPath, skip operasi
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		fmt.Printf("File %s tidak ada di folder sumber, operasi dipindahkan\n", tmpPath)
		return nil
	}

	// Pastikan folder tujuan ada
	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		return fmt.Errorf("%s: gagal membuat folder tujuan: %w", errorMsg, err)
	}

	// Jika file sudah ada di finalDir, skip operasi
	if _, err := os.Stat(finalPath); !os.IsNotExist(err) {
		fmt.Printf("File %s sudah ada di %s, operasi dipindahkan\n", fileName, finalDir)
		return nil
	}

	// Periksa izin file sumber
	// if err := checkFilePermissions(tmpPath); err != nil {
	// 	return fmt.Errorf("%s: %w", errorMsg, err)
	// }

	// Kalau file sumber tidak ada → langsung lanjut (return nil)
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		fmt.Printf("File sumber %s tidak ditemukan, skip pemindahan\n", tmpPath)
		return nil
	}

	// Beri jeda singkat untuk memastikan file tidak sedang digunakan
	time.Sleep(100 * time.Millisecond)

	// Coba rename langsung
	if err := os.Rename(tmpPath, finalPath); err != nil {
		fmt.Printf("Gagal rename %s ke %s: %v, mencoba fallback\n", tmpPath, finalPath, err)

		// Fallback: gunakan file sementara untuk menyalin
		tmpFile, err := os.CreateTemp("", "temp-*.webp")
		if err != nil {
			return fmt.Errorf("%s: gagal membuat file sementara: %w", errorMsg, err)
		}
		tmpFilePath := tmpFile.Name()
		defer tmpFile.Close()

		// Baca file sumber
		srcFile, err := os.Open(tmpPath)
		if err != nil {
			return fmt.Errorf("%s: gagal membaca file sumber: %w", errorMsg, err)
		}
		defer srcFile.Close()

		// Salin ke file sementara
		if _, err := io.Copy(tmpFile, srcFile); err != nil {
			return fmt.Errorf("%s: gagal menyalin ke file sementara: %w", errorMsg, err)
		}

		// Tutup file sementara
		tmpFile.Close()

		// Salin dari file sementara ke tujuan
		if err := os.Rename(tmpFilePath, finalPath); err != nil {
			// Jika rename gagal, coba salin manual
			data, err := os.ReadFile(tmpFilePath)
			if err != nil {
				return fmt.Errorf("%s: gagal membaca file sementara: %w", errorMsg, err)
			}
			if err := os.WriteFile(finalPath, data, 0644); err != nil {
				return fmt.Errorf("%s: gagal menulis file tujuan: %w", errorMsg, err)
			}
		}

		// Hapus file sementara
		if err := safeRemove(tmpFilePath); err != nil {
			fmt.Printf("Peringatan: gagal menghapus file sementara %s: %v\n", tmpFilePath, err)
		}

		// Coba hapus file sumber secara asinkronus
		asyncRemove(tmpPath)
		fmt.Printf("Berhasil memindahkan file %s ke %s menggunakan fallback\n", tmpPath, finalPath)
	} else {
		fmt.Printf("Berhasil memindahkan file %s ke %s menggunakan rename\n", tmpPath, finalPath)
	}

	return nil
}

func UploadImageToR2IfExists(fileName string, finalDir string, errorMsg string) error {
	if fileName == "" {
		return nil
	}

	cloudflareR2Config := appConfig.LoadCloudflareR2Config()

	bucketName := cloudflareR2Config.BucketName
	r2AccountID := cloudflareR2Config.AccountID
	r2AccessKey := cloudflareR2Config.AccessKey
	r2SecretKey := cloudflareR2Config.SecretKey

	tmpPath := filepath.Join("public", "uploads", "tmp", "images", fileName)
	objectKey := fmt.Sprintf("%s/%s", strings.TrimSuffix(finalDir, "/"), fileName)

	// Pastikan file ada
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		fmt.Printf("File %s tidak ditemukan, skip upload\n", tmpPath)
		return nil
	}

	// Buka file sumber
	srcFile, err := os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("%s: gagal membuka file sumber: %w", errorMsg, err)
	}
	defer srcFile.Close()

	// Baca isi file ke buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, srcFile); err != nil {
		return fmt.Errorf("%s: gagal membaca file sumber: %w", errorMsg, err)
	}

	// Deteksi MIME type
	ext := filepath.Ext(fileName)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = http.DetectContentType(buf.Bytes())
	}

	// Setup konfigurasi ke R2
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r2AccountID)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2AccessKey, r2SecretKey, "")),
		config.WithRegion("auto"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: "auto",
			}, nil
		})),
	)
	if err != nil {
		return fmt.Errorf("%s: gagal memuat konfigurasi R2: %w", errorMsg, err)
	}

	client := s3.NewFromConfig(cfg)

	// Upload ke R2
	fmt.Printf("Mengupload %s ke bucket %s ...\n", objectKey, bucketName)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ACL:         "public-read", // agar bisa diakses via CDN
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("%s: gagal upload ke Cloudflare R2: %w", errorMsg, err)
	}

	// Upload sukses, hapus file tmp secara async
	asyncRemove(tmpPath)
	fmt.Printf("Berhasil upload file %s ke R2 (Content-Type: %s) dan menghapus lokal tmp\n", tmpPath, contentType)

	return nil
}
