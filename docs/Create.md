
### `pouch.Create` Fonksiyonu ve `CreateOptions` Yapısı Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Create` fonksiyonunu ve bu fonksiyon tarafından kullanılan `CreateOptions` yapısını (struct) tanımlar. Bu kodun temel amacı, `docker create` komutunu programatik olarak ve esnek bir şekilde çalıştırmak için bir arayüz sağlamaktır. Kullanıcı, `CreateOptions` yapısını doldurarak yeni bir Docker konteyneri oluşturmak için gereken tüm ayarları belirleyebilir.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// CreateOptions, bir Docker konteyneri oluştururken kullanılacak 
// yapılandırma seçeneklerini içeren bir yapıdır.
type CreateOptions struct {
	Name        string            // --name: Konteynerin adı
	Image       string            // Kullanılacak Docker imajı (zorunlu)
	Port        string            // -p: Port yönlendirmesi (örn: "8080:80")
	HostDataDir string            // -v: Yerel dizini konteynere bağlama (volume)
	Network     string            // --network: Konteynerin bağlanacağı ağ
	Hostname    string            // --hostname: Konteynerin hostname'i
	UserUIDGID  string            // --user: Konteyner içinde çalışacak kullanıcı (UID:GID)
	MemoryLimit string            // --memory: Hafıza limiti (örn: "512m")
	EntryPoint  string            // --entrypoint: İmajın varsayılan entrypoint'ini geçersiz kılma
	CPULimit    float64           // --cpus: CPU limiti (örn: 1.5)
	EnvVars     map[string]string // -e: Ortam değişkenleri
	Labels      map[string]string // --label: Etiketler
}

// Create, verilen CreateOptions'a göre yeni bir Docker konteyneri oluşturur.
// Başarılı olursa oluşturulan konteynerin ID'sini, aksi takdirde bir hata döndürür.
func Create(opt CreateOptions) (string, error) {
	// "docker create" komutunun temel argümanlarını oluştur.
	args := []string{"create"}

	// CreateOptions'taki alanların dolu olup olmamasına göre argümanları dinamik olarak ekle.
	if opt.Name != "" {
		args = append(args, "--name", opt.Name)
	}

	for k, v := range opt.EnvVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	if opt.Port != "" {
		args = append(args, "-p", opt.Port)
	}

	if opt.HostDataDir != "" {
		args = append(args, "-v", opt.HostDataDir)
	}

	// Her zaman yeniden başlayacak şekilde ayarla.
	args = append(args, "--restart", "always")

	if opt.Network != "" {
		args = append(args, "--network", opt.Network)
	}
    
	// Diğer opsiyonel parametreler...
	if opt.Hostname != "" {
		args = append(args, "--hostname", opt.Hostname)
	}

	if opt.MemoryLimit != "" {
		args = append(args, "--memory", opt.MemoryLimit)
	}

	if opt.CPULimit > 0 {
		args = append(args, "--cpus", fmt.Sprintf("%.1f", opt.CPULimit))
	}

	for k, v := range opt.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}

	if opt.UserUIDGID != "" {
		args = append(args, "--user", opt.UserUIDGID)
	}

	// İnteraktif bir TTY oturumu için gerekli flag'ler.
	args = append(args, "-i", "-t")

	if opt.EntryPoint != "" {
		args = append(args, "--entrypoint", opt.EntryPoint)
	}

	// İmaj adı zorunludur. Kontrol et.
	if opt.Image == "" {
		return "", fmt.Errorf("imaj adı belirtilmelidir")
	}
	args = append(args, opt.Image)

	// Komutu oluştur ve çalıştır.
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("konteyner oluşturma hatası: %v\nçıktı: %s", err, string(out))
	}

	// Docker çıktısındaki boşlukları temizleyerek sadece konteyner ID'sini döndür.
	return strings.TrimSpace(string(out)), nil
}
```

### Yapı ve Fonksiyon Detayları

#### `CreateOptions` Yapısı

Bu yapı, "options pattern" (seçenekler deseni) adı verilen bir tasarım desenini kullanır. Fonksiyona çok sayıda parametre geçmek yerine, tüm bu parametreleri tek bir yapı içinde toplar. Bu, kodun okunabilirliğini ve yönetilebilirliğini artırır. Her bir alan, `docker create` komutunun bir flag'ine karşılık gelir.

#### `Create` Fonksiyonu

Bu fonksiyon, `docker create` komutunu adım adım oluşturur ve çalıştırır:

1.  **Argüman Hazırlığı**:
    `args := []string{"create"}` satırıyla bir `string` dilimi (slice) başlatılır. Bu dilim, `docker` komutuna verilecek tüm argümanları içerecektir.

2.  **Dinamik Argüman Ekleme**:
    Fonksiyon, `CreateOptions` yapısındaki her bir alanı kontrol eder. Eğer bir alan (örneğin `opt.Name`) boş değilse, ilgili Docker flag'i (`--name`) ve değeri `args` dilimine eklenir. `EnvVars` ve `Labels` gibi `map` türündeki alanlar için bir döngü kullanılarak her bir anahtar-değer çifti `-e` veya `--label` flag'i ile eklenir.

3.  **Sabit Argümanlar**:
    Bazı argümanlar koşulsuz olarak eklenir:
    *   `--restart always`: Konteynerin herhangi bir nedenle durması durumunda Docker tarafından otomatik olarak yeniden başlatılmasını sağlar.
    *   `-i -t`: Genellikle interaktif bir oturum için kullanılır (`-it`). Konteynerin standart girdisini (stdin) açık tutar ve bir pseudo-TTY ayırır. Bu, konteynere daha sonra `attach` olmak için faydalıdır.

4.  **Zorunlu Alan Kontrolü**:
    `if opt.Image == ""` kontrolü, bir konteyner oluşturmak için imaj adının mutlak bir gereklilik olduğunu garanti eder. Eğer imaj adı belirtilmemişse, fonksiyon bir hata döndürerek işlemi durdurur.

5.  **Komutun Çalıştırılması**:
    `exec.Command("docker", args...)` ile tüm argümanlar kullanılarak bir komut nesnesi oluşturulur. `cmd.CombinedOutput()` ile komut çalıştırılır ve hem standart çıktı hem de standart hata tek bir değişkende toplanır. Bu, hata ayıklamayı kolaylaştırır çünkü hata durumunda Docker'ın verdiği mesajı da görebiliriz.

6.  **Çıktının İşlenmesi**:
    `docker create` komutu başarılı olduğunda, standart çıktıya yeni oluşturulan konteynerin tam ID'sini yazar ve ardından bir satır sonu karakteri ekler. `strings.TrimSpace(string(out))` fonksiyonu, bu ID'nin başındaki ve sonundaki tüm boşlukları (satır sonu dahil) temizleyerek saf konteyner ID'sini döndürür.

### Parametreler ve Dönüş Değeri

*   **Parametre**:
    *   `opt (CreateOptions)`: Konteyneri oluşturmak için kullanılacak tüm ayarları içeren yapı.
*   **Dönüş Değeri**:
    *   `string`: Başarılı olursa, oluşturulan konteynerin ID'si.
    *   `error`: İşlem sırasında bir hata oluşursa, detaylı bir hata nesnesi.

### Kullanım Örneği

Aşağıda, `redis` imajını kullanarak `my-redis-db` adında bir konteyner oluşturma örneği verilmiştir.

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	options := pouch.CreateOptions{
		Name:        "my-redis-db",
		Image:       "redis:alpine",
		Port:        "6379:6379",
		MemoryLimit: "256m",
		EnvVars: map[string]string{
			"REDIS_REPLICATION_MODE": "master",
		},
		Labels: map[string]string{
			"project": "awesome-app",
			"owner":   "dev-team",
		},
	}

	containerID, err := pouch.Create(options)
	if err != nil {
		log.Fatalf("Konteyner oluşturulamadı: %v", err)
	}

	log.Printf("Konteyner başarıyla oluşturuldu. ID: %s\n", containerID)
}
```