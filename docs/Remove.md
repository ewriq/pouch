### `pouch.Remove` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Remove` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, belirtilen bir Docker konteynerini sistemden kalıcı olarak kaldırmaktır. Bu işlem, arka planda `docker rm` komutunu çalıştırarak gerçekleştirilir. Ayrıca, çalışan bir konteyneri zorla kaldırma seçeneği sunarak esneklik sağlar.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// Remove, belirtilen ID'ye sahip bir Docker konteynerini kaldırır.
// 'force' parametresi true ise, çalışan bir konteyneri zorla kaldırır (-f).
// Başarılı olduğunda kaldırılan konteynerin ID'sini döndürür.
func Remove(id string, force bool) (string, error) {
	// "docker rm" komutunun temel argümanlarını hazırla.
	args := []string{"rm"}
	
	// Eğer force parametresi true ise, komuta "-f" (force) bayrağını ekle.
	if force {
		args = append(args, "-f")
	}

	// Kaldırılacak konteynerin ID'sini argümanların sonuna ekle.
	args = append(args, id)

	// Docker komutunu oluştur.
	cmd := exec.Command("docker", args...)

	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem Go hatasını hem de Docker'ın çıktısını içeren
		// detaylı bir hata mesajı döndür.
		return "", fmt.Errorf("kaldırma hatası: %v\nçıktı: %s", err, out)
	}
	
	// Docker rm başarılı olduğunda kaldırılan konteynerin ID'sini döndürür.
	// Bu çıktıyı temizleyip geri ver.
	return strings.TrimSpace(string(out)), nil
}
```

### Fonksiyon Detayları

1.  **Dinamik Argüman Oluşturma**:
    Fonksiyon, `args := []string{"rm"}` ile `docker rm` komutunun temel argümanlarını oluşturarak işe başlar. Ardından `if force { ... }` bloğu ile `force` parametresinin durumunu kontrol eder. Eğer `force` parametresi `true` olarak ayarlanmışsa, `-f` bayrağı `args` dilimine eklenir. Bu bayrak, `docker rm` komutuna çalışan bir konteyneri bile durdurup kaldırması talimatını verir.

2.  **Komutun Tamamlanması**:
    `args = append(args, id)` satırıyla, kaldırılacak konteynerin ID'si veya adı argüman listesinin sonuna eklenir. Böylece, çalıştırılacak tam komut (`docker rm [-f] [id]`) dinamik olarak oluşturulmuş olur.

3.  **Komutu Çalıştırma ve Hata Yönetimi**:
    `cmd.CombinedOutput()` ile oluşturulan komut çalıştırılır. `docker rm` komutu, çalışan bir konteyneri `force` bayrağı olmadan kaldırmaya çalışırsa veya belirtilen ID'ye sahip bir konteyner bulamazsa hata verir. `CombinedOutput` sayesinde hem bu hata mesajları hem de başarılı durumda Docker'ın döndürdüğü çıktı yakalanır. Bir hata (`err != nil`) durumunda, sorunun kaynağını belirlemeye yardımcı olmak için hem Go hatasını hem de Docker'ın çıktısını içeren zengin bir hata mesajı döndürülür.

4.  **Başarılı Çıktının Döndürülmesi**:
    `docker rm` komutu başarıyla çalıştığında, kaldırdığı konteynerin ID'sini veya adını standart çıktıya yazar. Fonksiyon bu çıktıyı (`out`) yakalar, `strings.TrimSpace` ile olası baştaki ve sondaki boşlukları temizler ve bu temizlenmiş metni (kaldırılan konteynerin ID'si) `string` olarak geri döndürür.

### Parametreler

*   `id (string)`: Kaldırılacak olan Docker konteynerinin kimliği (ID) veya adı.
*   `force (bool)`: Konteynerin zorla kaldırılıp kaldırılmayacağını belirten bir boolean değeri.
    *   `true`: Konteyner çalışıyorsa bile durdurulur ve kaldırılır (`docker rm -f`).
    *   `false`: Konteyner çalışıyorsa komut hata verir ve konteyner kaldırılmaz.

### Dönüş Değeri

*   `string`: İşlem başarılı olursa, kaldırılan konteynerin ID'si veya adı.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner bulunamazsa, `force` `false` iken konteyner çalışıyorsa veya başka bir Docker hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

*   Bu fonksiyonun çalışabilmesi için sistemde **Docker'ın kurulu** ve `docker` komutunun `PATH` içinde erişilebilir olması gerekir.

### Kullanım Örneği

`old-container` adında durdurulmuş bir konteyneri ve `live-container` adında çalışan bir konteyneri kaldırma senaryoları:

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	// Senaryo 1: Durdurulmuş bir konteyneri kaldırma (force=false)
	stoppedContainerID := "old-container"
	removedID, err := pouch.Remove(stoppedContainerID, false)
	if err != nil {
		log.Printf("'%s' konteyneri kaldırılırken hata oluştu: %v", stoppedContainerID, err)
	} else {
		log.Printf("Konteyner '%s' başarıyla kaldırıldı.", removedID)
	}

	// Senaryo 2: Çalışan bir konteyneri zorla kaldırma (force=true)
	runningContainerID := "live-container"
	removedIDForced, err := pouch.Remove(runningContainerID, true)
	if err != nil {
		log.Printf("'%s' konteyneri zorla kaldırılırken hata oluştu: %v", runningContainerID, err)
	} else {
		log.Printf("Çalışan konteyner '%s' başarıyla zorla kaldırıldı.", removedIDForced)
	}
}
```