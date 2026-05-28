package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func SeedSetting(db *gorm.DB, count int) {
	log.Printf("Seeding %d setting...", count)
	var items []models.Setting

	const (
		tnc = `
<p>Dengan mengakses situs web ini, kami menganggap Anda menerima syarat dan ketentuan ini. Jangan lanjutkan menggunakan DilZ Topup jika Anda tidak setuju untuk mengambil semua syarat dan ketentuan yang tercantum di halaman ini.</p>

<h3>Lisensi</h3>
<p>Kecuali dinyatakan lain, DilZ Topup dan/atau pemberi lisensinya memiliki hak kekayaan intelektual untuk semua materi di DilZ Topup. Semua hak kekayaan intelektual dilindungi. Anda dapat mengakses ini dari DilZ Topup untuk penggunaan pribadi Anda sendiri dengan tunduk pada batasan yang ditetapkan dalam syarat dan ketentuan ini.</p>
<p>Anda tidak boleh:</p>
<ul>
    <li>Menerbitkan ulang materi dari DilZ Topup</li>
    <li>Menjual, menyewakan, atau mensublisensikan materi dari DilZ Topup</li>
    <li>Mereproduksi, menggandakan, atau menyalin materi dari DilZ Topup</li>
    <li>Mendistribusikan kembali konten dari DilZ Topup</li>
</ul>

<h3>Akun Pengguna</h3>
<p>Jika Anda membuat akun di situs web kami, Anda bertanggung jawab untuk menjaga keamanan akun Anda dan Anda sepenuhnya bertanggung jawab atas semua aktivitas yang terjadi di bawah akun dan tindakan lain yang diambil sehubungan dengannya. Anda harus segera memberitahu kami tentang penggunaan akun Anda yang tidak sah atau pelanggaran keamanan lainnya.</p>

<h3>Penafian</h3>
<p>Sejauh diizinkan oleh hukum yang berlaku, kami mengecualikan semua representasi, jaminan, dan ketentuan yang berkaitan dengan situs web kami dan penggunaan situs web ini. Tidak ada dalam penafian ini yang akan:</p>
<ul>
    <li>Membatasi atau mengecualikan tanggung jawab kami atau Anda atas kematian atau cedera pribadi;</li>
    <li>Membatasi atau mengecualikan tanggung jawab kami atau Anda atas penipuan atau pernyataan keliru yang menipu;</li>
    <li>Membatasi salah satu dari tanggung jawab kami atau Anda dengan cara apa pun yang tidak diizinkan berdasarkan hukum yang berlaku.</li>
</ul>

<h3>Perubahan Ketentuan</h3>
<p>Kami berhak untuk mengubah syarat dan ketentuan ini kapan saja. Dengan terus menggunakan situs web ini, Anda setuju untuk terikat oleh versi terbaru dari syarat dan ketentuan ini.</p>

<h3>Hubungi Kami</h3>
<p>Jika Anda memiliki pertanyaan tentang Syarat dan Ketentuan ini, silakan hubungi kami melalui email di <a href="mailto:support@dilztopup.com">support@dilztopup.com</a>.</p>
		`

		pp = `
<p>Terima kasih telah mengunjungi DilZ Topup. Kami menghargai privasi Anda dan berkomitmen untuk melindunginya. Kebijakan Privasi ini menjelaskan jenis informasi pribadi yang kami kumpulkan, bagaimana kami menggunakannya, dan pilihan yang Anda miliki terkait informasi Anda.</p>

<h3>Informasi yang Kami Kumpulkan</h3>
<p>Kami mengumpulkan informasi tentang Anda melalui berbagai cara ketika Anda menggunakan situs kami, termasuk:</p>
<ul>
    <li><strong>Informasi yang Anda Berikan Secara Langsung:</strong> Seperti nama, alamat email, nomor telepon, dan informasi lain yang Anda masukkan saat mendaftar, melakukan pemesanan, atau menghubungi kami.</li>
    <li><strong>Informasi yang Dikumpulkan Secara Otomatis:</strong> Kami secara otomatis mencatat informasi tentang Anda dan komputer Anda. Contohnya, kami mencatat alamat IP, sistem operasi, jenis browser, dan waktu kunjungan Anda.</li>
    <li><strong>Cookies:</strong> Kami dapat menggunakan "cookies" untuk menyimpan preferensi pengunjung dan merekam informasi sesi.</li>
</ul>

<h3>Bagaimana Kami Menggunakan Informasi Anda</h3>
<p>Kami menggunakan informasi yang kami kumpulkan untuk:</p>
<ul>
    <li>Menjalankan, memelihara, dan meningkatkan situs, produk, dan layanan kami.</li>
    <li>Memproses dan mengirimkan entri dan hadiah kontes.</li>
    <li>Menanggapi komentar dan pertanyaan Anda serta menyediakan layanan pelanggan.</li>
    <li>Mengirimkan informasi, termasuk konfirmasi, faktur, pemberitahuan teknis, pembaruan, dan pesan dukungan.</li>
</ul>

<h3>Keamanan Data</h3>
<p>Kami mengambil langkah-langkah yang wajar untuk membantu melindungi informasi pribadi Anda dalam upaya untuk mencegah kehilangan, penyalahgunaan, dan akses tidak sah, pengungkapan, perubahan, dan perusakan.</p>

<h3>Hak Anda</h3>
<p>Anda memiliki hak untuk mengakses, memperbaiki, atau menghapus informasi pribadi Anda yang kami miliki. Anda dapat melakukannya dengan menghubungi kami melalui informasi kontak di bawah ini.</p>

<h3>Perubahan pada Kebijakan Privasi Ini</h3>
<p>Kami dapat mengubah kebijakan privasi ini dari waktu ke waktu. Jika kami melakukan perubahan, kami akan memberi tahu Anda dengan merevisi tanggal "Terakhir diperbarui" di atas kebijakan ini.</p>

<h3>Hubungi Kami</h3>
<p>Jika Anda memiliki pertanyaan tentang kebijakan privasi ini, silakan hubungi kami di:</p>
<p>DilZ Topup<br>Email: <a href="mailto:support@dilztopup.com">support@dilztopup.com</a></p>
		`
	)

	// Optional: predefined example
	sample := models.Setting{
		BaseModel: models.BaseModel{ID: "ST1"},
		Key:       "privacy-policy",
		Value:     pp,
		Type:      "text",
	}
	items = append(items, sample)

	sample = models.Setting{
		BaseModel: models.BaseModel{ID: "ST2"},
		Key:       "terms-and-conditions",
		Value:     tnc,
		Type:      "text",
	}
	items = append(items, sample)

	sample = models.Setting{
		BaseModel: models.BaseModel{ID: "ST3"},
		Key:       "ppn",
		Value:     "0.11",
		Type:      "number",
	}
	items = append(items, sample)

	sample = models.Setting{
		BaseModel: models.BaseModel{ID: "ST4"},
		Key:       "digiflazz_balance",
		Value:     "143000",
		Type:      "number",
	}
	items = append(items, sample)

	result := db.CreateInBatches(&items, 100)
	if result.Error != nil {
		log.Printf("Could not seed setting: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d setting successfully.\n", count)
	}
}
