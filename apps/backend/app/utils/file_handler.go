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

// ensureDir memastikan bahwa direktori dari sebuah path file itu ada.
// Jika tidak ada, ia akan membuatnya.
func ensureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			return fmt.Errorf("gagal membuat direktori %s: %w", dirName, err)
		}
	}
	return nil
}

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
			// Ensure direktori untuk file log ada sebelum menulis
			logPath := "storage/logs/failed_deletes.log"
			_ = ensureDir(logPath)
			_ = os.WriteFile(logPath, []byte(fmt.Sprintf("%s: %v\n", path, err)), 0644)
		}
	}()
}

// MoveImageIfExists memindahkan file dari tmp ke finalDir dengan aman
func MoveImageIfExists(fileName string, finalDir string, errorMsg string) error {
	if fileName == "" {
		return nil
	}

	if !strings.HasPrefix(finalDir, "storage/") {
		finalDir = filepath.Join("storage", finalDir)
	}

	tmpPath := filepath.Join("storage", "uploads", "tmp", "images", fileName)
	finalPath := filepath.Join(finalDir, fileName)

	// Periksa apakah file tmp benar-benar ada
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		fmt.Printf("File %s tidak ada di folder sumber, skip operasi\n", tmpPath)
		return nil
	}

	// Pastikan folder tujuan akhir ADA sebelum melakukan apa pun
	if err := ensureDir(finalPath); err != nil {
		return fmt.Errorf("%s: %w", errorMsg, err)
	}

	// Jika file sudah ada di tujuan, skip pemindahan
	if _, err := os.Stat(finalPath); err == nil {
		fmt.Printf("File %s sudah ada di %s, skip operasi\n", fileName, finalDir)
		return nil
	}

	time.Sleep(100 * time.Millisecond)

	// Coba rename langsung
	if err := os.Rename(tmpPath, finalPath); err != nil {
		fmt.Printf("Gagal rename %s ke %s: %v, mencoba fallback\n", tmpPath, finalPath, err)

		// Fallback: Salin dan Hapus (mengatasi error "cross-device link" di Docker/Volume)
		srcFile, err := os.Open(tmpPath)
		if err != nil {
			return fmt.Errorf("%s: gagal membaca file sumber: %w", errorMsg, err)
		}
		defer srcFile.Close()

		dstFile, err := os.Create(finalPath)
		if err != nil {
			return fmt.Errorf("%s: gagal membuat file tujuan: %w", errorMsg, err)
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			dstFile.Close()
			return fmt.Errorf("%s: gagal menyalin ke tujuan: %w", errorMsg, err)
		}
		dstFile.Close()

		// Berhasil di-copy, sekarang hapus sumbernya
		asyncRemove(tmpPath)
		fmt.Printf("Berhasil memindahkan file %s ke %s menggunakan fallback (Copy & Delete)\n", tmpPath, finalPath)
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

	tmpPath := filepath.Join("storage", "uploads", "tmp", "images", fileName)
	objectKey := fmt.Sprintf("%s/%s", strings.TrimSuffix(finalDir, "/"), fileName)

	// Pastikan file lokal ada
	if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
		fmt.Printf("File lokal %s tidak ditemukan, skip upload ke R2\n", tmpPath)
		return nil
	}

	srcFile, err := os.Open(tmpPath)
	if err != nil {
		return fmt.Errorf("%s: gagal membuka file sumber: %w", errorMsg, err)
	}
	defer srcFile.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, srcFile); err != nil {
		return fmt.Errorf("%s: gagal membaca file sumber: %w", errorMsg, err)
	}

	ext := filepath.Ext(fileName)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = http.DetectContentType(buf.Bytes())
	}

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

	fmt.Printf("Mengupload %s ke bucket %s ...\n", objectKey, bucketName)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ACL:         "public-read",
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("%s: gagal upload ke Cloudflare R2: %w", errorMsg, err)
	}

	asyncRemove(tmpPath)
	fmt.Printf("Berhasil upload file %s ke R2 (Content-Type: %s) dan menghapus lokal tmp\n", tmpPath, contentType)

	return nil
}
