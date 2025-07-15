
### `pouch.DeleteFile` Fonksiyonu Açıklaması

Bu Go kodu, `pouch` paketi içinde yer alan `DeleteFile` adında bir fonksiyon tanımlar. Fonksiyonun temel amacı, çalışan bir Docker konteynerinin içindeki belirli bir dosyayı silmektir. Bu işlemi, arka planda `docker exec` komutunu kullanarak, konteyner içerisinde `rm -f` komutunu çalıştırarak gerçekleştirir.

#### Kod Bloğu

```go
package pouch

import (
	"fmt"
	"os/exec"
)

// DeleteFile, belirtilen Docker konteynerinin içindeki bir dosyayı siler.
// Bu işlem, "docker exec [containerID] rm -f [filePath]" komutunu çalıştırır.
func DeleteFile(containerID, filePath string) error {
	// Komutu oluştur: docker exec [ID] rm -f [dosya-yolu]
	cmd := exec.Command("docker", "exec", containerID, "rm", "-f", filePath)
	
	// Komutu çalıştır ve standart çıktı ile standart hatayı birleştir.
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Hata durumunda, hem orijinal hatayı hem de komutun çıktısını
		// içeren açıklayıcı bir hata mesajı döndür.
		return fmt.Errorf("dosya silme hatası: %v, çıktı: %s", err, string(out))
	}
	
	// İşlem başarılıysa nil (hata yok) döndür.
	return nil
}
```

### Fonksiyon Detayları

Fonksiyonun çalışma mantığı şu adımlara ayrılabilir:

1.  **Komut Oluşturma**:
    `exec.Command("docker", "exec", containerID, "rm", "-f", filePath)` satırı, terminalde çalıştırılacak olan `docker exec [containerID] rm -f [filePath]` komutunu programatik olarak oluşturur.
    *   `docker exec`: Belirtilen bir konteynerin içinde bir komut çalıştırmayı sağlar.
    *   `rm -f`: Standart bir Linux/Unix komutudur. `-f` (force) bayrağı sayesinde, dosya mevcut olmasa bile komut hata vermez ve silme onayı istemeden işlemi gerçekleştirir.

2.  **Komutu Çalıştırma ve Çıktıyı Yakalama**:
    `cmd.CombinedOutput()` fonksiyonu, oluşturulan komutu çalıştırır. Olası bir hata durumunda, Docker'ın döndürdüğü hata mesajını (`stderr`) ve standart çıktıyı (`stdout`) yakalayarak hata ayıklamayı kolaylaştırır.

3.  **Hata Kontrolü**:
    `if err != nil` bloğu, komutun çalıştırılması sırasında bir hata olup olmadığını kontrol eder. `docker exec` komutu genellikle konteyner çalışmıyorsa veya belirtilen dosya bir dizinse ve `-r` bayrağı kullanılmadıysa hata verir. Fonksiyon, bu tür hataları yakalar ve `fmt.Errorf` ile detaylı bir hata mesajı oluşturarak geri döndürür.

4.  **Başarılı Durum**:
    Eğer komut hatasız bir şekilde tamamlanırsa (dosya silindi veya zaten mevcut değildi), fonksiyon `nil` değeri döndürerek işlemin başarılı olduğunu belirtir.

### Parametreler

*   `containerID (string)`: Dosyanın silineceği Docker konteynerinin kimliği (ID) veya adı.
*   `filePath (string)`: Konteynerin **içindeki** silinecek dosyanın mutlak yolu (örneğin, `/app/config.json`).

### Dönüş Değeri

*   `error`: Fonksiyon bir `error` değeri döndürür.
    *   İşlem başarılı olursa bu değer `nil` olur.
    *   Konteynere erişilemezse veya başka bir `docker exec` hatası oluşursa, komutun neden başarısız olduğuna dair bilgileri içeren bir hata nesnesi döndürülür.

### Önemli Notlar

*   Bu fonksiyonun çalışabilmesi için, kodu çalıştıran sistemde **Docker'ın yüklü olması** ve `docker` komutunun sistemin `PATH` değişkeninde bulunması gerekmektedir.
*   Fonksiyonun çağrıldığı sırada hedef `containerID`'ye sahip konteynerin **çalışıyor durumda** olması zorunludur.

### Kullanım Örneği

Diyelim ki `my-web-server` adlı bir konteynerin içindeki `/var/log/nginx/access.log` dosyasını silmek istiyoruz.

```go
package main

import (
	"log"
	// 'pouch' paketini projenize göre import etmeniz gerekir.
)

func main() {
	containerID := "my-web-server"
	fileToDelete := "/var/log/nginx/access.log"

	// pouch.DeleteFile fonksiyonunu çağır
	err := pouch.DeleteFile(containerID, fileToDelete)
	if err != nil {
		log.Fatalf("Konteyner içindeki dosya silinemedi: %v", err)
	}

	log.Printf("'%s' konteynerindeki '%s' dosyası başarıyla silindi.", containerID, fileToDelete)
}
```