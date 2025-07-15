### `pouch.Restart` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `Restart` adında bir fonksiyonu tanımlar. Fonksiyonun temel amacı, durmuş veya çalışır durumda olan belirli bir Docker konteynerini yeniden başlatmaktır. Bu işlem, arka planda `docker restart` komutunu çalıştırarak gerçekleştirilir ve başarılı olduğunda yeniden başlatılan konteynerin kimliğini döndürür.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
	"strings"
)

// Restart, belirtilen ID'ye sahip bir Docker konteynerini yeniden başlatır.
// Başarılı olursa, yeniden başlatılan konteynerin ID'sini döndürür.
func Restart(id string) (string, error) {
	// "docker restart [id]" komutunu oluştur.
	cmd := exec.Command("docker", "restart", id)

	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return "", fmt.Errorf("restart hatası: %v\nçıktı: %s", err, out)
	}
	
	// docker restart başarılı olduğunda konteynerin ID'sini döndürür.
	// Bu çıktıyı temizleyip geri ver.
	return strings.TrimSpace(string(out)), nil
}
```

### Fonksiyon Detayları

1.  **Komut Oluşturma**:
    `exec.Command("docker", "restart", id)` satırı, terminalde `docker restart [containerID]` komutunu çalıştırmak için bir komut nesnesi oluşturur. Bu komut, hedef konteyner durmuş ise onu başlatır; çalışıyorsa önce durdurur sonra tekrar başlatır.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. `docker restart` komutu başarılı olduğunda yeniden başlattığı konteynerin ID'sini standart çıktıya yazar. Hata durumunda ise (örneğin konteyner bulunamadığında) hata mesajını standart hataya yazar. `CombinedOutput` bu iki akışı da yakalayarak hem başarılı sonucu hem de hata detaylarını tek bir `out` değişkeninde toplar.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun sıfırdan farklı bir çıkış koduyla sonlanıp sonlanmadığını kontrol eder. Eğer belirtilen ID'ye sahip bir konteyner yoksa `docker restart` komutu hata verir. Fonksiyon bu hatayı yakalar ve `fmt.Errorf` ile hem Go seviyesindeki hatayı (`err`) hem de Docker'dan gelen spesifik hata mesajını (`out`) içeren detaylı bir hata mesajı oluşturur.

4.  **Başarılı Çıktının İşlenmesi**:
    Komut başarıyla tamamlandığında, `docker restart` çıktısı olarak verilen konteyner ID'sini içeren `out` byte dizisi, `string(out)` ile metne dönüştürülür. Genellikle bu çıktının sonunda bir satır sonu karakteri bulunur. `strings.TrimSpace`, bu gereksiz boşluğu temizleyerek sadece saf konteyner ID'sinin döndürülmesini sağlar.

### Parametreler

*   `id (string)`: Yeniden başlatılacak olan Docker konteynerinin kimliği (ID) veya adı.

### Dönüş Değeri

*   `string`: İşlem başarılı olursa, yeniden başlatılan konteynerin kimliği veya adı.
*   `error`:
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Eğer konteyner bulunamazsa veya başka bir Docker hatası oluşursa, detaylı bilgi içeren bir hata nesnesi döndürülür.

### Bağımlılıklar

*   Bu fonksiyonun çalışabilmesi için, kodu yürüten sistemde **Docker'ın yüklü** olması ve `docker` komutunun sistemin `PATH` değişkeni üzerinden erişilebilir olması gerekmektedir.

### Kullanım Örneği

`my-web-server` adındaki bir konteyneri yeniden başlatmak için `Restart` fonksiyonunun nasıl kullanılacağını gösteren bir örnek:

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerToRestart := "my-web-server"

	log.Printf("'%s' konteyneri yeniden başlatılıyor...", containerToRestart)

	// pouch.Restart fonksiyonunu çağır
	restartedID, err := pouch.Restart(containerToRestart)
	if err != nil {
		log.Fatalf("Konteyner yeniden başlatılamadı: %v", err)
	}

	log.Printf("Konteyner ('%s') başarıyla yeniden başlatıldı.", restartedID)
}
```